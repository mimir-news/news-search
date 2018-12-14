package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/mimir-news/news-search/pkg/domain"
	"github.com/mimir-news/pkg/dbutil"
)

// Service metadata.
const (
	ServiceName    = "news-search"
	ServiceVersion = "0.1"
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
	db                   dbutil.Config
	port                 string
	tokenSecret          string
	tokenVerificationKey string
}

func getConfig() config {
	tokenSecret := getSecret(mustGetenv("TOKEN_SECRETS_FILE"))

	return config{
		db:                   dbutil.MustGetConfig("DB"),
		tokenSecret:          tokenSecret.Secret,
		tokenVerificationKey: tokenSecret.Key,
		port:                 mustGetenv("SERVICE_PORT"),
	}
}

type secret struct {
	Secret string `json:"secret"`
	Key    string `json:"key"`
}

func getSecret(filename string) secret {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var s secret
	err = json.Unmarshal(content, &s)
	if err != nil {
		log.Fatal(err)
	}

	return s
}

func mustGetenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("No value for key: %s\n", key)
	}

	return val
}
