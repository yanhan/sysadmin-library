package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.opentelemetry.io/contrib/detectors/aws/ec2"
	"go.opentelemetry.io/contrib/detectors/aws/eks"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	_name      = "fib"
	_paramName = "n"
	_port      = 8831
)

var (
	_app                   *App
	_otlpCollectorHostPort string
)

type App struct {
	r io.Reader
	l *log.Logger
}

func NewApp(r io.Reader, l *log.Logger) *App {
	return &App{r: r, l: l}
}

func (a *App) CheckN(ctx context.Context, s string) (uint, error) {
	_, span := otel.Tracer(_name).Start(ctx, "CheckN")
	defer span.End()

	var n uint
	_, err := fmt.Sscanf(s, "%d\n", &n)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}

	// Store n as a string to not overflow an int64
	nStr := strconv.FormatUint(uint64(n), 10)
	span.SetAttributes(attribute.String("request.n", nStr))
	return n, err
}

func (a *App) Compute(ctx context.Context, n uint) (uint64, error) {
	var span trace.Span
	ctx, span = otel.Tracer(_name).Start(ctx, "Compute")
	defer span.End()

	f, err := func(ctx context.Context) (uint64, error) {
		_, span := otel.Tracer(_name).Start(ctx, "Fibonacci")
		defer span.End()
		f, err := Fibonacci(n)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		return f, err
	}(ctx)
	return f, err
}

func newExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint(_otlpCollectorHostPort), otlptracegrpc.WithDialOption(grpc.WithBlock()))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %+v", err)
	}
	return traceExporter, err
}

func newResource() *resource.Resource {
	ctx := context.Background()
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("xrayfib"),
			semconv.ServiceVersionKey.String("v0.2.12"),
			attribute.String("environment", "demo"),
		),
	)
	eksResourceDetector := eks.NewResourceDetector()
	eksResource, _ := eksResourceDetector.Detect(ctx)
	r, _ = resource.Merge(r, eksResource)
	ec2ResourceDetector := ec2.NewResourceDetector()
	ec2Resource, _ := ec2ResourceDetector.Detect(ctx)
	r, _ = resource.Merge(r, ec2Resource)
	return r
}

func initTracerProvider(l *log.Logger) (func(context.Context) error, error) {
	ctx := context.Background()
	dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(dialCtx, _otlpCollectorHostPort, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector at %s: %+v", _otlpCollectorHostPort, err)
	}
	defer conn.Close()
	l.Printf("Connected to gRPC server %q\n", _otlpCollectorHostPort)
	exp, err := newExporter(ctx)
	if err != nil {
		l.Fatal(err)
	}
	l.Printf("Created new gRPC trace exporter: %+v\n", exp)
	idg := xray.NewIDGenerator()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithIDGenerator(idg),
		sdktrace.WithResource(newResource()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})
	return tp.Shutdown, nil
}

func Fibonacci(n uint) (uint64, error) {
	if n <= 1 {
		return uint64(n), nil
	}

	if n > 93 {
		return 0, fmt.Errorf("unsupported fibonacci number %d: too large", n)
	}

	var n2, n1 uint64 = 0, 1
	for i := uint(2); i < n; i++ {
		n2, n1 = n1, n1+n2
	}
	return n2 + n1, nil
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func fibHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	newCtx, span := otel.Tracer(_name).Start(ctx, "fibHandler")
	defer span.End()
	query := r.URL.Query()
	values, ok := query[_paramName]
	if !ok || len(values) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Query parameter '%s' is missing", _paramName)
		return
	}
	n, err := _app.CheckN(newCtx, values[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_app.l.Printf("_app.Poll(newCtx, %d): %+v\n", n, err)
		fmt.Fprintf(w, "Bad request, please try with valid input")
		return
	}
	result, err := _app.Compute(newCtx, n)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_app.l.Printf("_app.Compute(newCtx, %d): %+v\n", n, err)
		fmt.Fprintf(w, "Internal Server Error")
		return
	}
	w.WriteHeader(http.StatusOK)
	response := fmt.Sprintf("Fib(%d) = %d", n, result)
	_app.l.Println(response)
	fmt.Fprint(w, response)
}

func main() {
	l := log.New(os.Stdout, "", 0)
	ctx := context.Background()

	_otlpCollectorHostPort = os.Getenv("OTLP_COLLECTOR_HOST_PORT")
	if len(_otlpCollectorHostPort) <= 0 {
		l.Fatal("Environment variable OTLP_COLLECTOR_HOST_PORT is empty")
	}
	l.Printf("Setting _otlpCollectorHostPort to %q\n", _otlpCollectorHostPort)

	shutdown, err := initTracerProvider(l)
	if err != nil {
		l.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			l.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	_app = NewApp(os.Stdin, l)
	http.HandleFunc("/healthz", healthCheckHandler)
	// We need to wrap our handlers with otelhttp.NewHandler to get the
	// URL, Method, Response code, Client IP, etc to show up on AWS XRay.
	// The value of the string argument does not seem to matter
	fibH := otelhttp.NewHandler(http.HandlerFunc(fibHandler), "fib")
	http.Handle("/fib", fibH)
	l.Printf("Listening at port %d\n", _port)
	l.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", _port), nil))
}
