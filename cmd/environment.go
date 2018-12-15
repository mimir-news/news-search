package main

import (
	"database/sql"
	"log"

	"github.com/mimir-news/news-search/pkg/service"
)

type env struct {
	db      *sql.DB
	newsSvc service.NewsService
}

func setupEnv(conf config) *env {
	db, err := conf.db.ConnectPostgres()
	if err != nil {
		log.Fatal(err)
	}

	return &env{
		db:      db,
		newsSvc: service.NewNewsService(nil),
	}
}

func (e *env) close() {
	err := e.db.Close()
	if err != nil {
		log.Println(err)
	}
}
