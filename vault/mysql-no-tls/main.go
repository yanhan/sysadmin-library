package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	if dbUser == "" {
		log.Fatal("Please supply DB_USER environment variable.")
	}
	if dbPassword == "" {
		log.Fatal("Please supply DB_PASSWORD environment variable.")
	}
	if dbHost == "" {
		log.Fatal("Please supply DB_HOST environment variable.")
	}
	dbName := "apps_galore"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to MySQL: %v\n", err)
	}
	defer db.Close()
	http.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {
		err := db.Ping()
		if err != nil {
			log.Fatalln(err)
		} else {
			w.Write([]byte("OK"))
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// vim:noet
