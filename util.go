package indc

import (
	"errors"

	"github.com/shopspring/decimal"
)

var (
	// ErrInvalidCandleCount is returned when insufficient amount of candles is
	// provided.
	ErrInvalidCandleCount = errors.New("insufficient amount of candles")

	// ErrInvalidLength is returned when provided length is less than 1.
	ErrInvalidLength = errors.New("length cannot be less than 1")

	// ErrMANotSet is returned when ma field is nil.
	ErrMANotSet = errors.New("macd ma value not set")
)

func Resize(dd []decimal.Decimal, l, offset int) ([]decimal.Decimal, error) {
	if l+offset > len(dd) {
		return decimal.Zero, ErrInvalidCandleCount
	}
	return decimal.Zero, nil
}
