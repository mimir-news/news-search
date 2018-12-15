package repository

import (
	"time"

	"github.com/mimir-news/pkg/schema/news"
)

// ArticleRepo interface for interacting with persisted articles.
type ArticleRepo interface {
	FindArticles(symbol string, after time.Time, limit int) ([]news.Article, error)
	FindAllArticles(after time.Time, limit int) ([]news.Article, error)
}

// MockArticleRepo mock implementation of ArticleRepo.
type MockArticleRepo struct {
	FindArticlesArgSymbol   string
	FindArticlesArgAfter    time.Time
	FindArticlesArgLimit    int
	FindArticlesArticles    []news.Article
	FindArticlesErr         error
	FindArticlesInvocations int

	FindAllArticlesArgAfter    time.Time
	FindAllArticlesArgLimit    int
	FindAllArticlesArticles    []news.Article
	FindAllArticlesErr         error
	FindAllArticlesInvocations int
}

// UnsetArgs sets all MockArticleRepo fields to their default value.
func (ar *MockArticleRepo) UnsetArgs() {
	var defaultTime time.Time

	ar.FindArticlesArgSymbol = ""
	ar.FindArticlesArgAfter = defaultTime
	ar.FindArticlesArgLimit = 0
	ar.FindArticlesInvocations = 0

	ar.FindAllArticlesArgAfter = defaultTime
	ar.FindAllArticlesArgLimit = 0
	ar.FindAllArticlesInvocations = 0
}

// FindArticles mock implementation.
func (ar *MockArticleRepo) FindArticles(symbol string, after time.Time, limit int) ([]news.Article, error) {
	ar.FindArticlesArgSymbol = symbol
	ar.FindArticlesArgAfter = after
	ar.FindArticlesArgLimit = limit
	ar.FindArticlesInvocations++

	return ar.FindArticlesArticles, ar.FindArticlesErr
}

// FindAllArticles mock implementation.
func (ar *MockArticleRepo) FindAllArticles(after time.Time, limit int) ([]news.Article, error) {
	ar.FindAllArticlesArgAfter = after
	ar.FindAllArticlesArgLimit = limit
	ar.FindAllArticlesInvocations++

	return ar.FindAllArticlesArticles, ar.FindAllArticlesErr
}
