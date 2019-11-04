package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Reference:
// https://www.alexedwards.net/blog/serverless-api-with-go-and-aws-lambda

func genericErrorResponse(statusCode int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: "error",
	}
}

func genericOkResponse(statusCode int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: "ok",
	}
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	data, ok := req.QueryStringParameters["data"]
	if !ok {
		err := errors.New("'data' not found in request parameters.")
		log.Println(err)
		return genericErrorResponse(http.StatusBadRequest), err
	}
	log.Printf("data = %s\n", data)
	r := events.APIGatewayProxyResponse{
		StatusCode: http.StatusAccepted,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: "Yea",
	}
	return r, nil
}

func main() {
	lambda.Start(HandleRequest)
}
