package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"database/sql"
	"github.com/go-sql-driver/mysql"
)

const (
	rdsCombinedCaBundleDefaultPath = "/home/ec2-user/rds-combined-ca-bundle.pem"
)

func main() {
	rdsTlsCaCertPath := os.Getenv("RDS_TLS_CA_CERT_PATH")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	if rdsTlsCaCertPath == "" {
		rdsTlsCaCertPath = rdsCombinedCaBundleDefaultPath
	}
	if dbUser == "" {
		log.Fatal("Please supply DB_USER environment variable.")
	}
	if dbPassword == "" {
		log.Fatal("Please supply DB_PASSWORD environment variable.")
	}
	if dbHost == "" {
		log.Fatal("Please supply DB_HOST environment variable.")
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
	dbName := "acedb"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=custom", dbUser, dbPassword, dbHost, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to MySQL: %v\n", err)
	}
	defer db.Close()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		stmtOut, err := db.Prepare("SELECT val FROM hearts WHERE id = ?")
		if err != nil {
			msg := fmt.Sprintf("Error running SELECT statement: %v", err)
			log.Print(msg)
			w.WriteHeader(500)
			w.Write([]byte(msg))
		} else {
			defer stmtOut.Close()
			var v string
			err = stmtOut.QueryRow(13).Scan(&v)
			if err != nil {
				msg := fmt.Sprintf("Error QueryRow: %v", err)
				log.Print(msg)
				w.WriteHeader(500)
				w.Write([]byte(msg))
			} else {
				w.Write([]byte(fmt.Sprintf("Got %s", v)))
			}
		}
	})
	http.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// vim:noet
