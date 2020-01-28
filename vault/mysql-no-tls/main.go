package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DbCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbSecretsFile := os.Getenv("DB_SECRETS_FILE")
	if dbHost == "" {
		log.Fatal("Please supply DB_HOST environment variable.")
	}
	if dbSecretsFile == "" {
		log.Fatal("Please supply DB_SECRETS_FILE environment variable.")
	}
	data, err := ioutil.ReadFile(dbSecretsFile)
	if err != nil {
		log.Fatalf("Error reading DB_SECRETS_FILE %s", dbSecretsFile)
	}
	var dbCreds DbCredentials
	err = json.Unmarshal(data, &dbCreds)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON for DB credentials: %v", err)
	}
	dbName := "apps_galore"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbCreds.Username, dbCreds.Password, dbHost, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to MySQL: %v\n", err)
	}
	defer db.Close()
	http.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Database username = %s\n", dbCreds.Username)
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
