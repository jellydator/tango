package indc

import (
	"errors"

	"github.com/shopspring/decimal"
	"github.com/swithek/chartype"
)

var (
	// ErrInvalidDataPointCount is returned when insufficient amount of data points is
	// provided.
	ErrInvalidDataPointCount = errors.New("insufficient amount of data points")

	// ErrInvalidLength is returned when provided length is less than 1.
	ErrInvalidLength = errors.New("length cannot be less than 1")

	// ErrMANotSet is returned when ma field is nil.
	ErrMANotSet = errors.New("ma value not set")
)

// resize cuts given array based on length to use for
// calculations.
func resize(dd []decimal.Decimal, l int) ([]decimal.Decimal, error) {
	if l < 1 {
		return nil, ErrInvalidLength
	}

	if l > len(dd) {
		return nil, ErrInvalidDataPointCount
	}

	return dd[len(dd)-l:], nil
}

// resizeCandles cuts given array based on length to use for
// calculations.
func resizeCandles(cc []chartype.Candle, l int) ([]chartype.Candle, error) {
	if l < 1 {
		return nil, ErrInvalidLength
	}

	if l > len(cc) || l < 1 {
		return nil, ErrInvalidDataPointCount
	}

	return cc[len(cc)-l:], nil
}
