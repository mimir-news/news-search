package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mimir-news/news-search/pkg/repository"
	"github.com/mimir-news/news-search/pkg/service"
	"github.com/mimir-news/pkg/httputil/auth"
	"github.com/mimir-news/pkg/id"
	"github.com/mimir-news/pkg/schema/news"
	"github.com/stretchr/testify/assert"
)

func TestHandleGetStockNews(t *testing.T) {
	assert := assert.New(t)

	userID := id.New()
	symbol := "TEST"

	conf := getTestConfig()
	server := newServer(getTestEnv(conf, nil), conf)
	token := getTestToken(conf, userID)

	baseRoute := fmt.Sprintf("/v1/news/%s", symbol)
	testInvalidNewsParameters(token, baseRoute, server.Handler, t)

	expectedArtices := []news.Article{
		news.Article{ID: id.New()},
		news.Article{ID: id.New()},
	}
	articeRepo := &repository.MockArticleRepo{
		FindArticlesArticles: expectedArtices,
	}
	server = newServer(getTestEnv(conf, articeRepo), conf)

	req := createTestRequest(token, baseRoute+"?period=2D")
	res := performTestRequest(server.Handler, req)

	assert.Equal(http.StatusOK, res.Code)
	assert.Equal(1, articeRepo.FindArticlesInvocations)
	assert.Equal(symbol, articeRepo.FindArticlesArgSymbol)
	assert.Equal(defaultSearchLimit, articeRepo.FindArticlesArgLimit)

	articles := make([]news.Article, 0)
	err := json.NewDecoder(res.Body).Decode(&articles)
	assert.NoError(err)
	assert.Equal(len(expectedArtices), len(articles))
	for i, ea := range expectedArtices {
		assert.Equal(ea.ID, articles[i].ID, i)
	}

	articeRepo.UnsetArgs()
	req = createTestRequest(token, baseRoute+"?limit=6")
	res = performTestRequest(server.Handler, req)

	assert.Equal(http.StatusOK, res.Code)
	assert.Equal(1, articeRepo.FindArticlesInvocations)
	assert.Equal(symbol, articeRepo.FindArticlesArgSymbol)
	assert.Equal(6, articeRepo.FindArticlesArgLimit)

	err = json.NewDecoder(res.Body).Decode(&articles)
	assert.NoError(err)
	assert.Equal(len(expectedArtices), len(articles))
	for i, ea := range expectedArtices {
		assert.Equal(ea.ID, articles[i].ID, i)
	}
}

func TestHandleGetNews(t *testing.T) {
	assert := assert.New(t)

	userID := id.New()

	conf := getTestConfig()
	server := newServer(getTestEnv(conf, nil), conf)
	token := getTestToken(conf, userID)

	baseRoute := "/v1/news"
	testInvalidNewsParameters(token, baseRoute, server.Handler, t)

	expectedArtices := []news.Article{
		news.Article{ID: id.New()},
		news.Article{ID: id.New()},
	}
	articeRepo := &repository.MockArticleRepo{
		FindAllArticlesArticles: expectedArtices,
	}
	server = newServer(getTestEnv(conf, articeRepo), conf)

	req := createTestRequest(token, baseRoute+"?period=2D")
	res := performTestRequest(server.Handler, req)

	assert.Equal(http.StatusOK, res.Code)
	assert.Equal(1, articeRepo.FindAllArticlesInvocations)
	assert.Equal(defaultSearchLimit, articeRepo.FindAllArticlesArgLimit)

	articles := make([]news.Article, 0)
	err := json.NewDecoder(res.Body).Decode(&articles)
	assert.NoError(err)
	assert.Equal(len(expectedArtices), len(articles))
	for i, ea := range expectedArtices {
		assert.Equal(ea.ID, articles[i].ID, i)
	}

	articeRepo.UnsetArgs()
	req = createTestRequest(token, baseRoute+"?limit=6")
	res = performTestRequest(server.Handler, req)

	assert.Equal(http.StatusOK, res.Code)
	assert.Equal(1, articeRepo.FindAllArticlesInvocations)
	assert.Equal(6, articeRepo.FindAllArticlesArgLimit)

	err = json.NewDecoder(res.Body).Decode(&articles)
	assert.NoError(err)
	assert.Equal(len(expectedArtices), len(articles))
	for i, ea := range expectedArtices {
		assert.Equal(ea.ID, articles[i].ID, i)
	}
}

func testInvalidNewsParameters(token, baseRoute string, handler http.Handler, t *testing.T) {
	assert := assert.New(t)

	req := createTestRequest(token, baseRoute+"?period=0W&limit=3")
	res := performTestRequest(handler, req)

	assert.Equal(http.StatusBadRequest, res.Code)

	req = createTestRequest(token, baseRoute+"?period=D&limit=3")
	res = performTestRequest(handler, req)

	assert.Equal(http.StatusBadRequest, res.Code)

	req = createTestRequest(token, baseRoute+"?period=1D&limit=0")
	res = performTestRequest(handler, req)

	assert.Equal(http.StatusBadRequest, res.Code)

	req = createTestRequest(token, baseRoute+"?period=1D&limit=50")
	res = performTestRequest(handler, req)

	assert.Equal(http.StatusBadRequest, res.Code)
}

func performTestRequest(r http.Handler, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func getTestEnv(cfg config, articleRepo repository.ArticleRepo) *env {
	return &env{
		newsSvc: service.NewNewsService(articleRepo),
	}
}

func getTestConfig() config {
	return config{
		JWTCredentials: auth.JWTCredentials{
			Issuer: "news-search-test-issuer",
			Secret: "my-secret",
		},
	}
}

func getTestSigner(cfg config) auth.Signer {
	return auth.NewSigner(cfg.JWTCredentials, 24*time.Hour)
}

func getTestToken(cfg config, userID string) string {
	signer := getTestSigner(cfg)

	token, err := signer.Sign(id.New(), auth.User{ID: userID, Role: auth.UserRole})
	if err != nil {
		log.Fatal(err)
	}

	return token
}

func createTestRequest(token, route string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		log.Fatal(err)
	}

	if token != "" {
		bearerToken := auth.AuthTokenPrefix + token
		req.Header.Set(auth.AuthHeaderKey, bearerToken)
	}

	return req
}
