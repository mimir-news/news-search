package repository

import (
	"database/sql"
	"strings"
	"time"

	"github.com/mimir-news/pkg/schema/news"
)

const (
	keywordDelimiter = ","
)

// ArticleRepo interface for interacting with persisted articles.
type ArticleRepo interface {
	FindArticles(symbol string, after time.Time, limit int) ([]news.Article, error)
	FindAllArticles(after time.Time, limit int) ([]news.Article, error)
}

// NewArticleRepo creates new ArticleRepo using default implementation.
func NewArticleRepo(db *sql.DB) ArticleRepo {
	return &pgArticleRepo{
		db: db,
	}
}

// pgArticleRepo postgres implementation of ArticleRepo.
type pgArticleRepo struct {
	db *sql.DB
}

const findStockArticlesQuery = `
	SELECT a.id, a.url, a.title, a.keywords, a.article_date FROM article a
		INNER JOIN article_cluster c ON c.lead_article_id = a.id
		WHERE c.symbol = $1 
		AND c.article_date >= $2
		ORDER BY c.score DESC
		LIMIT $3`

func (ar *pgArticleRepo) FindArticles(symbol string, after time.Time, limit int) ([]news.Article, error) {
	rows, err := ar.db.Query(findStockArticlesQuery, symbol, after, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return mapRowsToArticles(rows)
}

const findArticlesQuery = `
	SELECT a.id, a.url, a.title, a.keywords, a.article_date FROM article a
		INNER JOIN article_cluster c ON c.lead_article_id = a.id
		AND c.article_date >= $1
		ORDER BY c.score DESC
		LIMIT $2`

func (ar *pgArticleRepo) FindAllArticles(after time.Time, limit int) ([]news.Article, error) {
	rows, err := ar.db.Query(findArticlesQuery, after, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return mapRowsToArticles(rows)
}

func mapRowsToArticles(rows *sql.Rows) ([]news.Article, error) {
	articles := make([]news.Article, 0)
	for rows.Next() {
		var a news.Article
		var keywords sql.NullString
		err := rows.Scan(&a.ID, &a.URL, &a.Title, &keywords, &a.ArticleDate)
		if err != nil {
			return nil, err
		}

		a.Keywords = splitKeywords(keywords)
		articles = append(articles, a)
	}

	return articles, nil
}

func splitKeywords(keywords sql.NullString) []string {
	if !keywords.Valid {
		return make([]string, 0)
	}

	return strings.Split(keywords.String, keywordDelimiter)
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
