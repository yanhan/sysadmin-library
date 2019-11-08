package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"database/sql"
	"github.com/go-sql-driver/mysql"
)

const (
	rdsTlsCaCertPath = "/etc/ssl/certs/rds-combined-ca-bundle.pem"
)

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	if dbUser == "" {
		log.Fatal("Please supply DB_USER environment variable.")
	}
	if dbPassword == "" {
		log.Fatal("Please supply DB_PASSWORD environment variable.")
	}
	if dbHost == "" {
		log.Fatal("Please supply DB_HOST environment variable.")
	}
	if dbName == "" {
		log.Fatal("Please supply DB_NAME environment variable.")
	}
	pemFileData, err := ioutil.ReadFile(rdsTlsCaCertPath)
	if err != nil {
		log.Fatalf("Error opening RDS combined ca bundle: %v\n", err)
	}
	rootCaCertPool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("Error getting system cert pool: %v\n", err)
	}
	ok := rootCaCertPool.AppendCertsFromPEM(pemFileData)
	if !ok {
		log.Fatal("Failed to append certs from RDS combined ca bundle")
	}
	mysql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs: rootCaCertPool,
	})
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=custom", dbUser, dbPassword, dbHost, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to MySQL: %v\n", err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT title FROM books LIMIT 1")
	if err != nil {
		msg := fmt.Sprintf("Error running SELECT statement: %v", err)
		log.Print(msg)
	} else {
		if rows.Next() {
			title := ""
			err := rows.Scan(&title)
			if err != nil {
				log.Printf("rows.Scan() error = %v\n", err)
				return
			}
			cols, err := rows.Columns()
			if err != nil {
				log.Printf("rows.Columns() error = %v\n", err)
			} else {
				log.Printf("Columns = %s\n", strings.Join(cols, ","))
				log.Printf("Book title = \"%s\"\n", title)
			}
		} else {
			log.Println("No results.")
		}
	}
}

// vim:noet
