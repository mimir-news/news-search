package service

import (
	"time"

	"github.com/mimir-news/news-search/pkg/repository"
	"github.com/mimir-news/pkg/schema/news"
)

// NewsService service for retrieving and ranking news articles.
type NewsService interface {
	GetStockNews(symbol string, since time.Duration, limit int) ([]news.Article, error)
	GetNews(since time.Duration, limit int) ([]news.Article, error)
}

// NewNewsService creates a NewsService using the default implementaiton.
func NewNewsService(articleRepo repository.ArticleRepo) NewsService {
	return &newsSvc{
		articleRepo: articleRepo,
	}
}

// newsSvc default NewsService implementation.
type newsSvc struct {
	articleRepo repository.ArticleRepo
}

// GetStockNews gets the highest ranked stock new after a given time.
func (ns *newsSvc) GetStockNews(symbol string, since time.Duration, limit int) ([]news.Article, error) {
	after := calcAfter(since)
	articles, err := ns.articleRepo.FindArticles(symbol, after, limit)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (ns *newsSvc) GetNews(since time.Duration, limit int) ([]news.Article, error) {
	return nil, nil
}

func calcAfter(since time.Duration) time.Time {
	now := time.Now().UTC()
	return now.Add(-since)
}
