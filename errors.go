package indc

import (
	"errors"
)

var (
	// ErrInvalidCandleCount is returned when insufficient amount of candles is
	// provided.
	ErrInvalidCandleCount = errors.New("insufficient amount of candles")

	// ErrInvalidLength is returned when provided length is less than 1.
	ErrInvalidLength = errors.New("length cannot be less than 1")

	// ErrInvalidOffset is returned when provided offset is less than 0.
	ErrInvalidOffset = errors.New("offset cannot be less than 0")

	// ErrMANotSet is returned when ma field is nil.
	ErrMANotSet = errors.New("macd ma value not set")
)
