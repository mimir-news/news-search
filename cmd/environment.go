package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mimir-news/pkg/httputil"
)

type env struct {
	db *sql.DB
}

func setupEnv(conf config) *env {
	db, err := conf.db.ConnectPostgres()
	if err != nil {
		log.Fatal(err)
	}

	return &env{
		db: db,
	}
}

func (e *env) close() {
	err := e.db.Close()
	if err != nil {
		log.Println(err)
	}
}

func (e *env) throwNotImplemented(c *gin.Context) {
	c.Error(httputil.NewError("Not implemented", http.StatusNotImplemented))
}
