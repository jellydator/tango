package indc

import (
	"errors"

	"github.com/shopspring/decimal"
	"github.com/swithek/chartype"
)

var (
	// ErrCandlesCount is returned when insufficient amount of candles is
	// provided.
	ErrCandlesCount = errors.New("insufficient amount of candles")
)

// SMA holds all the neccesary information needed to calculate a simple
// moving average.
type SMA struct {
	// Length specifies how many candles should be used.
	Length int `json:"length"`

	// Offset specifies how many latest candles should be skipped.
	Offset int `json:"offset"`

	// Src specifies which price field of the candle should be used.
	Src chartype.CandleField `json:"src"`
}

// Calc calculates SMA value by using settings stored in the func receiver.
func (s SMA) Calc(cc []chartype.Candle) (decimal.Decimal, error) {
	if s.CandlesCount() > len(cc) {
		return decimal.Zero, ErrCandlesCount
	}

	res := decimal.Zero

	for i := len(cc) - 1 - s.Offset; i >= len(cc)-s.CandlesCount(); i-- {
		res = res.Add(s.Src.Extract(cc[i]))
	}

	return res.Div(decimal.NewFromInt(int64(s.Length))), nil
}

// CandlesCount determines the total amount of candles needed for SMA
// calculation by using settings stored in the receiver.
func (s SMA) CandlesCount() int {
	return s.Length + s.Offset
}

// CalcSMA calculates SMA value by using settings passed as parameters.
func CalcSMA(cc []chartype.Candle, l, off int, src chartype.CandleField) (decimal.Decimal, error) {
	s := SMA{Length: l, Offset: off, Src: src}
	return s.Calc(cc)
}

// CandlesCountSMA determines the total amount of candles needed for SMA
// calculation by using settings passed as parameters.
func CandlesCountSMA(l, off int) int {
	s := SMA{Length: l, Offset: off}
	return s.CandlesCount()
}
