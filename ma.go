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

	// ErrInvalidOffset is returned when provided offset is less than 0.
	ErrInvalidOffset = errors.New("offset cannot be less than 0")

	// ErrMANotSet is returned when ma field is nil.
	ErrMANotSet = errors.New("macd ma value not set")
)

// MA interface holds all the placeholder functions required that every
// moving average has to have.
type MA interface {
	// Validate makes sure that the moving average is valid.
	Validate() error

	// Calc calculates moving average value by using settings stored in the func receiver.
	Calc(cc []chartype.Candle) (decimal.Decimal, error)

	// CandleCount determines the total amount of candles needed for moving averages
	// calculation by using settings stored in the receiver.
	CandleCount() int
}

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

// Validate checks all SMA settings stored in func receiver to make sure that
// they're meeting each of their own requirements.
func (s SMA) Validate() error {
	if s.Length < 1 {
		return ErrInvalidLength
	}

	if s.Offset < 0 {
		return ErrInvalidOffset
	}

	if err := s.Src.Validate(); err != nil {
		return err
	}

	return nil
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

// ValidateSMA checks all settings passed as parameters to make sure that
// they're meeting each of their own requirements.
func ValidateSMA(len, off int, src chartype.CandleField) error {
	s := SMA{Length: len, Offset: off, Src: src}
	return s.Validate()
}

// CalcSMA calculates SMA value by using settings passed as parameters.
func CalcSMA(cc []chartype.Candle, len, off int, src chartype.CandleField) (decimal.Decimal, error) {
	s := SMA{Length: len, Offset: off, Src: src}
	return s.Calc(cc)
}

// CandleCountSMA determines the total amount of candles needed for SMA
// calculation by using settings passed as parameters.
func CandleCountSMA(len, off int) int {
	s := SMA{Length: len, Offset: off}
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

// Validate checks all EMA settings stored in func receiver to make sure that
// they're meeting each of their own requirements.
func (e EMA) Validate() error {
	if e.Length < 1 {
		return ErrInvalidLength
	}

	if e.Offset < 0 {
		return ErrInvalidOffset
	}

	if err := e.Src.Validate(); err != nil {
		return err
	}

	return nil
}

// Calc calculates EMA value by using settings stored in the func receiver.
func (e EMA) Calc(cc []chartype.Candle) (decimal.Decimal, error) {
	if e.CandleCount() > len(cc) {
		return decimal.Zero, ErrInvalidCandleCount
	}

	res, err := CalcSMA(cc, e.Length, e.Offset+e.Length, e.Src)

	if err != nil {
		return decimal.Zero, err
	}

	mul := e.multiplier()

	for i := len(cc) - e.CandleCount() + e.Length; i < len(cc)-e.Offset; i++ {
		res = e.Src.Extract(cc[i]).Mul(mul).Add(res.Mul(decimal.NewFromInt(1).Sub(mul)))
	}

	return res, nil
}

// multiplier calculates EMA multiplier value by using settings stored in the func receiver.
func (e EMA) multiplier() decimal.Decimal {
	return decimal.NewFromFloat(2.0 / float64(e.Length+1))
}

// CandleCount determines the total amount of candles needed for EMA
// calculation by using settings stored in the receiver.
func (e EMA) CandleCount() int {
	return e.Length*2 + e.Offset
}

// ValidateEMA checks all settings passed as parameters to make sure that
// they're meeting each of their own requirements.
func ValidateEMA(len, off int, src chartype.CandleField) error {
	e := EMA{Length: len, Offset: off, Src: src}
	return e.Validate()
}

// CalcEMA calculates EMA value by using settings passed as parameters.
func CalcEMA(cc []chartype.Candle, len, off int, src chartype.CandleField) (decimal.Decimal, error) {
	e := EMA{Length: len, Offset: off, Src: src}
	return e.Calc(cc)
}

// CandleCountEMA determines the total amount of candles needed for EMA
// calculation by using settings passed as parameters.
func CandleCountEMA(len, off int) int {
	e := EMA{Length: len, Offset: off}
	return e.CandleCount()
}

// WMA holds all the neccesary information needed to calculate weighted
// moving average.
type WMA struct {
	// Length specifies how many candles should be used.
	Length int `json:"length"`

	// Offset specifies how many latest candles should be skipped.
	Offset int `json:"offset"`

	// Src specifies which price field of the candle should be used.
	Src chartype.CandleField `json:"src"`
}

// Validate checks all WMA settings stored in func receiver to make sure that
// they're meeting each of their own requirements.
func (w WMA) Validate() error {
	if w.Length < 1 {
		return ErrInvalidLength
	}

	if w.Offset < 0 {
		return ErrInvalidOffset
	}

	if err := w.Src.Validate(); err != nil {
		return err
	}

	return nil
}

// Calc calculates WMA value by using settings stored in the func receiver.
func (w WMA) Calc(cc []chartype.Candle) (decimal.Decimal, error) {
	if w.CandleCount() > len(cc) {
		return decimal.Zero, ErrInvalidCandleCount
	}

	res := decimal.Zero

	weight := decimal.NewFromFloat(float64(w.Length*(w.Length+1)) / 2.0)

	for i := len(cc) - w.CandleCount(); i < len(cc)-w.CandleCount()+w.Length; i++ {
		res = res.Add(w.Src.Extract(cc[i]).Mul(decimal.NewFromInt(int64(i - len(cc) + w.CandleCount() + 1)).Div(weight)))
	}

	return res, nil
}

// CandleCount determines the total amount of candles needed for WMA
// calculation by using settings stored in the receiver.
func (w WMA) CandleCount() int {
	return w.Length + w.Offset
}

// ValidateWMA checks all settings passed as parameters to make sure that
// they're meeting each of their own requirements.
func ValidateWMA(len, off int, src chartype.CandleField) error {
	w := WMA{Length: len, Offset: off, Src: src}
	return w.Validate()
}

// CalcWMA calculates WMA value by using settings passed as parameters.
func CalcWMA(cc []chartype.Candle, len, off int, src chartype.CandleField) (decimal.Decimal, error) {
	w := WMA{Length: len, Offset: off, Src: src}
	return w.Calc(cc)
}

// CandleCountWMA determines the total amount of candles needed for WMA
// calculation by using settings passed as parameters.
func CandleCountWMA(len, off int) int {
	w := WMA{Length: len, Offset: off}
	return w.CandleCount()
}

// MACD holds all the neccesary information needed to calculate moving averages
// convergence divergence.
type MACD struct {
	// MA1 configures first moving average.
	MA1 MA `json:"ma1"`

	// MA2 configures second moving average.
	MA2 MA `json:"ma2"`
}

// Validate checks all MACD settings stored in func receiver to make sure that
// they're meeting each of their own requirements.
func (macd MACD) Validate() error {
	if macd.MA1 == nil || macd.MA2 == nil {
		return ErrMANotSet
	}

	if err := macd.MA1.Validate(); err != nil {
		return err
	}

	if err := macd.MA2.Validate(); err != nil {
		return err
	}

	return nil
}

// Calc calculates MACD value by using settings stored in the func receiver.
func (macd MACD) Calc(cc []chartype.Candle) (decimal.Decimal, error) {
	res1, err := macd.MA1.Calc(cc)
	if err != nil {
		return decimal.Zero, err
	}

	res2, err := macd.MA2.Calc(cc)
	if err != nil {
		return decimal.Zero, err
	}

	res := res1.Sub(res2)

	return res, nil
}

// CandleCount determines the total amount of candles needed for MACD
// calculation by using settings stored in the receiver.
func (macd MACD) CandleCount() int {
	c1 := macd.MA1.CandleCount()
	c2 := macd.MA2.CandleCount()
	if c1 > c2 {
		return c1
	}
	return c2
}

// ValidateMACD checks all settings passed as parameters to make sure that
// they're meeting each of their own requirements.
func ValidateMACD(ma1, ma2 MA) error {
	macd := MACD{MA1: ma1, MA2: ma2}
	return macd.Validate()
}

// CalcMACD calculates MACD value by using settings passed as parameters.
func CalcMACD(cc []chartype.Candle, ma1, ma2 MA) (decimal.Decimal, error) {
	macd := MACD{MA1: ma1, MA2: ma2}
	return macd.Calc(cc)
}

// CandleCountMACD determines the total amount of candles needed for MACD
// calculation by using settings passed as parameters.
func CandleCountMACD(ma1, ma2 MA) int {
	macd := MACD{MA1: ma1, MA2: ma2}
	return macd.CandleCount()
}
