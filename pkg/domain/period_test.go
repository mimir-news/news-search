package domain_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/mimir-news/news-search/pkg/domain"
	"github.com/stretchr/testify/assert"
)

func TestGetDuration(t *testing.T) {
	type testCase struct {
		str      string
		duration time.Duration
		err      error
	}

	testCases := []testCase{
		testCase{
			str:      "1D",
			duration: domain.Day,
			err:      nil,
		},
		testCase{
			str:      "12D",
			duration: 12 * domain.Day,
			err:      nil,
		},
		testCase{
			str:      "2W",
			duration: 2 * domain.Week,
			err:      nil,
		},
		testCase{
			str:      "5M",
			duration: 5 * domain.Month,
			err:      nil,
		},
		testCase{
			str:      "3Y",
			duration: 3 * domain.Year,
			err:      nil,
		},
		testCase{
			str: "Y",
			err: domain.ErrInvalidPeriod,
		},
		testCase{
			str: "XY",
			err: domain.ErrInvalidPeriodQuantifier,
		},
		testCase{
			str: "1X",
			err: domain.ErrInvalidPeriodUnit,
		},
	}

	assert := assert.New(t)
	for i, tc := range testCases {
		d, err := domain.GetDuration(tc.str)
		assert.Equal(tc.err, err, fmt.Sprintf("%d - Wrong err for str: %s", i, tc.str))

		if tc.err == nil {
			assert.Equal(tc.duration, d, fmt.Sprintf("%d - Wrong duration for str: %s", i, tc.str))
		}
	}
}
