package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/mimir-news/news-search/pkg/domain"

	"github.com/mimir-news/pkg/httputil"

	"github.com/gin-gonic/gin"
)

func (e *env) handleGetStockNews(c *gin.Context) {
	period, limit, err := getPeriodAndLimit(c)
	if err != nil {
		c.Error(err)
		return
	}

	symbol := c.Param("symbol")
	articles, err := e.newsSvc.GetStockNews(symbol, period, limit)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, articles)
}

func getPeriodAndLimit(c *gin.Context) (time.Duration, int, error) {
	p, err := getPeriod(c)
	if err != nil {
		return defaultPeriod, 0, err
	}

	l, err := getLimit(c)
	if err != nil {
		return defaultPeriod, 0, err
	}

	return p, l, nil
}

func getPeriod(c *gin.Context) (time.Duration, error) {
	periodStr, ok := c.GetQuery("period")
	if !ok {
		return defaultPeriod, nil
	}

	d, err := domain.GetDuration(periodStr)
	if err != nil {
		return d, errBadRequest(err.Error())
	}

	return d, nil
}

func getLimit(c *gin.Context) (int, error) {
	val, ok := c.GetQuery("limit")
	if !ok {
		return defaultSearchLimit, nil
	}

	limit, err := strconv.Atoi(val)
	if err != nil {
		return 0, errBadRequest("Failed to parse limit as integer.")
	}

	if limit < 1 {
		return 0, errBadRequest("Limit must be at least 1.")
	}

	if limit > maxSearchLimit {
		return 0, errBadRequest("Search limit to high,")
	}

	return limit, nil
}

func errBadRequest(msg string) error {
	return httputil.NewError(msg, http.StatusBadRequest)
}
