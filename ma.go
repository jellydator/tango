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
)

// SMA holds all the neccesary information needed to calculate simple
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
	if s.CandleCount() > len(cc) {
		return decimal.Zero, ErrInvalidCandleCount
	}

	res := decimal.Zero

	for i := len(cc) - 1 - s.Offset; i >= len(cc)-s.CandleCount(); i-- {
		res = res.Add(s.Src.Extract(cc[i]))
	}

	return res.Div(decimal.NewFromInt(int64(s.Length))), nil
}

// CandleCount determines the total amount of candles needed for SMA
// calculation by using settings stored in the receiver.
func (s SMA) CandleCount() int {
	return s.Length + s.Offset
}

// CalcSMA calculates SMA value by using settings passed as parameters.
func CalcSMA(cc []chartype.Candle, l, off int, src chartype.CandleField) (decimal.Decimal, error) {
	s := SMA{Length: l, Offset: off, Src: src}
	return s.Calc(cc)
}

// CandleCountSMA determines the total amount of candles needed for SMA
// calculation by using settings passed as parameters.
func CandleCountSMA(l, off int) int {
	s := SMA{Length: l, Offset: off}
	return s.CandleCount()
}

// EMA holds all the neccesary information needed to calculate exponential
// moving average.
type EMA struct {
	// Length specifies how many candles should be used.
	Length int `json:"length"`

	// Offset specifies how many latest candles should be skipped.
	Offset int `json:"offset"`

	// Src specifies which price field of the candle should be used.
	Src chartype.CandleField `json:"src"`
}

// Calc calculates EMA value by using settings stored in the func receiver.
func (e EMA) Calc(cc []chartype.Candle) (decimal.Decimal, error) {
	if e.CandleCount() > len(cc) {
		return decimal.Zero, ErrInvalidCandleCount
	}

	res, err := CalcSMA(cc, e.Length, e.Offset, e.Src)

	if err != nil {
		return decimal.Zero, err
	}

	mul := e.Multiplier()

	for i := len(cc) - e.CandleCount(); i <= len(cc)-1-e.Offset; i++ {
		res = e.Src.Extract(cc[i]).Mul(mul).Add(res.Mul(decimal.NewFromInt(1).Sub(mul)))
	}

	return res, nil
}

// Multiplier calculates EMA multiplier value by using settings stored in the func receiver.
func (e EMA) Multiplier() decimal.Decimal {
	return decimal.NewFromFloat(2.0 / float64(e.Length+1))
}

// CandleCount determines the total amount of candles needed for EMA
// calculation by using settings stored in the receiver.
func (e EMA) CandleCount() int {
	return e.Length + e.Offset
}

// CalcEMA calculates EMA value by using settings passed as parameters.
func CalcEMA(cc []chartype.Candle, l, off int, src chartype.CandleField) (decimal.Decimal, error) {
	e := EMA{Length: l, Offset: off, Src: src}
	return e.Calc(cc)
}

// CandleCountEMA determines the total amount of candles needed for EMA
// calculation by using settings passed as parameters.
func CandleCountEMA(l, off int) int {
	e := EMA{Length: l, Offset: off}
	return e.CandleCount()
}
