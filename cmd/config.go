package main

import (
	"log"
	"os"

	"github.com/mimir-news/pkg/httputil/auth"

	"github.com/mimir-news/news-search/pkg/domain"
	"github.com/mimir-news/pkg/dbutil"
)

// Service metadata.
const (
	ServiceName    = "news-search"
	ServiceVersion = "1.1"
)

const (
	defaultSearchLimit = 5
	maxSearchLimit     = 10
	defaultPeriod      = domain.Day
)

var (
	unsecuredRoutes = []string{"/health"}
)

type config struct {
	db             dbutil.Config
	port           string
	JWTCredentials auth.JWTCredentials
}

func getConfig() config {
	jwtCredentials := getJWTCredentials(mustGetenv("JWT_CREDENTIALS_FILE"))

	return config{
		db:             dbutil.MustGetConfig("DB"),
		port:           mustGetenv("SERVICE_PORT"),
		JWTCredentials: jwtCredentials,
	}
}

func getJWTCredentials(filename string) auth.JWTCredentials {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	credentials, err := auth.ReadJWTCredentials(f)
	if err != nil {
		log.Fatal(err)
	}

	return credentials
}

func mustGetenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("No value for key: %s\n", key)
	}

	return val
}
