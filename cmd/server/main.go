// Package main
package main

import (
	"database/sql"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Siddharth-Chhatbar/reliable-email-delivery/internal/api"
	_ "github.com/lib/pq"
)

func main() {
	requiredEnvironment := []string{
		"DB_HOST",
		"DB_PORT",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_DB",
	}
	var missingEnvironment []string
	for _, name := range requiredEnvironment {
		if strings.TrimSpace(os.Getenv(name)) == "" {
			missingEnvironment = append(missingEnvironment, name)
		}
	}
	if len(missingEnvironment) > 0 {
		slog.Error("Missing required environment variables", "variables", missingEnvironment)
		os.Exit(1)
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		slog.Error("Invalid database port", "variable", "DB_PORT")
		os.Exit(1)
	}

	connectionURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   net.JoinHostPort(host, port),
		Path:   dbname,
	}
	query := connectionURL.Query()
	query.Set("sslmode", "disable")
	connectionURL.RawQuery = query.Encode()

	db, err := sql.Open("postgres", connectionURL.String())
	if err != nil {
		slog.Error("Error connecting to DB", "error", err)
		os.Exit(1)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(30 * time.Minute)
	defer db.Close()

	server := &http.Server{
		Addr:        ":8080",
		Handler:     api.NewHandler(db),
		ReadTimeout: 5 * time.Second,
	}

	slog.Info("Starting http server", "address", server.Addr)

	if err := server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		slog.Error("HTTP server stopped unexpectedly", "error", err)
		os.Exit(1)
	}
}
