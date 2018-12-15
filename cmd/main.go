package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mimir-news/pkg/httputil"
	"github.com/mimir-news/pkg/httputil/auth"
)

func main() {
	conf := getConfig()
	e := setupEnv(conf)
	defer e.close()
	server := newServer(e, conf)

	log.Printf("Starting %s on port: %s\n", ServiceName, conf.port)
	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func newServer(e *env, conf config) *http.Server {
	r := newRouter(e, conf)

	r.GET("/v1/news/:symbol", e.handleGetStockNews)
	r.GET("/v1/news", e.throwNotImplemented)

	return &http.Server{
		Addr:    ":" + conf.port,
		Handler: r,
	}
}

func newRouter(e *env, conf config) *gin.Engine {
	authOpts := auth.NewOptions(conf.tokenSecret, conf.tokenVerificationKey, unsecuredRoutes...)
	r := httputil.NewRouter(ServiceName, ServiceVersion, e.healthCheck)
	r.Use(auth.RequireToken(authOpts))

	return r
}

func (e *env) healthCheck() error {
	return e.db.Ping()
}
