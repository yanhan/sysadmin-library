package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v7"
)

const (
	envVarRedisAuthToken = "REDIS_AUTH_TOKEN"
	envVarRedisEndpoint  = "REDIS_ENDPOINT"
)

func main() {
	authToken := os.Getenv(envVarRedisAuthToken)
	if authToken == "" {
		log.Printf("Env var '%s' must be specified\n", envVarRedisAuthToken)
		os.Exit(1)
	}
	redisUrl := os.Getenv(envVarRedisEndpoint)
	if redisUrl == "" {
		log.Printf("Env var '%s' must be specified\n", envVarRedisEndpoint)
		os.Exit(1)
	}
	systemCertPool, err := x509.SystemCertPool()
	if err != nil {
		log.Printf("x509.SystemCertPool: %v\n", err)
		os.Exit(1)
	}
	// This works as long as it's launched on Amazon Linux.
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: authToken,
		DB:       0,
		TLSConfig: &tls.Config{
			RootCAs: systemCertPool,
		},
	})
	defer func() {
		err := rdb.Close()
		if err != nil {
			log.Printf("redis.Cient.Close: %v\n", err)
		}
	}()
	pong, err := rdb.Ping().Result()
	log.Println(pong, err)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if len(key) == 0 {
			_, err := w.Write([]byte("Missing 'key'"))
			if err != nil {
				log.Printf("ResponseWriter.Write: %v\n", err)
				w.WriteHeader(500)
			} else {
				w.WriteHeader(400)
			}
		} else {
			result, err := rdb.Get(key).Result()
			if err != nil {
				log.Printf("redis.Client.Get(%s); err = %v\n", key, err)
				w.WriteHeader(500)
				w.Write([]byte("Cannot find given key"))
			} else {
				log.Printf("redis.Client.Get: %s\n", result)
				w.Write([]byte(result))
			}
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
