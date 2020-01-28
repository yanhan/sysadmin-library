package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DbCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var dbSecretsFile string
var db *sql.DB
var dbCreds DbCredentials

func connectToDb(dbHost string) {
	if db != nil {
		db.Close()
		db = nil
	}
	data, err := ioutil.ReadFile(dbSecretsFile)
	if err != nil {
		log.Fatalf("Error reading DB_SECRETS_FILE %s", dbSecretsFile)
	}
	err = json.Unmarshal(data, &dbCreds)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON for DB credentials: %v", err)
	}
	dbName := "apps_galore"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbCreds.Username, dbCreds.Password, dbHost, dbName)
	log.Printf("Connecting to DB as user %s\n", dbCreds.Username)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to MySQL: %v\n", err)
	}
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbSecretsFile = os.Getenv("DB_SECRETS_FILE")
	if dbHost == "" {
		log.Fatal("Please supply DB_HOST environment variable.")
	}
	if dbSecretsFile == "" {
		log.Fatal("Please supply DB_SECRETS_FILE environment variable.")
	}
	connectToDb(dbHost)
	defer db.Close()
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP)
	go func() {
		for _ = range ch {
			connectToDb(dbHost)
		}
	}()
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
