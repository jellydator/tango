package indc

import (
	"errors"

	"github.com/shopspring/decimal"
	"github.com/swithek/chartype"
)

var (
	// ErrInvalidCandleCount is returned when insufficient amount of candles is
	// provided.
	ErrInvalidCandleCount = errors.New("insufficient amount of candles")

	// ErrInvalidLength is returned when provided length is less than 1.
	ErrInvalidLength = errors.New("length cannot be less than 1")

	// ErrMANotSet is returned when ma field is nil.
	ErrMANotSet = errors.New("ma value not set")
)

// resize cuts given array based on length to use for
// calculations.
func resize(dd []decimal.Decimal, l int) ([]decimal.Decimal, error) {
	if l > len(dd) || l < 1 {
		return nil, ErrInvalidCandleCount
	}
	return dd[len(dd)-l:], nil
}

// resizeCandles cuts given array based on length to use for
// calculations.
func resizeCandles(cc []chartype.Candle, l int) ([]chartype.Candle, error) {
	if l > len(cc) || l < 1 {
		return nil, ErrInvalidCandleCount
	}
	return cc[len(cc)-l:], nil
}
