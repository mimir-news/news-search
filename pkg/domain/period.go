package domain

import (
	"errors"
	"strconv"
	"time"
)

// Common errors that can occur in parsing a period.
var (
	ErrInvalidPeriod           = errors.New("invalid period")
	ErrInvalidPeriodQuantifier = errors.New("invalid period quantifier, must be a positive integer")
	ErrInvalidPeriodUnit       = errors.New("invalid period unit, must be D, W, M or Y")
)

// Base durations.
const (
	Day   = 24 * time.Hour
	Week  = 7 * Day
	Month = 30 * Day
	Year  = 365 * Day
)

var durationMap = map[string]time.Duration{
	"D": Day,
	"W": Week,
	"M": Month,
	"Y": Year,
}

// GetDuration parses a period string int a duration.
func GetDuration(periodStr string) (time.Duration, error) {
	p, err := parsePeriod(periodStr)
	if err != nil {
		return Day, err
	}

	d, err := translateDuration(p.unit)
	if err != nil {
		return Day, err
	}

	duration := time.Duration(p.quantifier) * d
	return duration, nil
}

type period struct {
	quantifier int
	unit       string
}

func translateDuration(unit string) (time.Duration, error) {
	d, ok := durationMap[unit]
	if !ok {
		return Day, ErrInvalidPeriodUnit
	}
	return d, nil
}

func parsePeriod(periodStr string) (period, error) {
	strlen := len(periodStr)
	if strlen < 2 {
		return period{}, ErrInvalidPeriod
	}

	q, err := strconv.Atoi(periodStr[:strlen-1])
	if err != nil || q < 1 {
		return period{}, ErrInvalidPeriodQuantifier
	}

	unit := periodStr[strlen-1:]
	return period{quantifier: q, unit: unit}, nil
}
