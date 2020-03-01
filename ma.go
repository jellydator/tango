package indc

import (
	"github.com/shopspring/decimal"
)

// MA interface holds all the placeholder functions required that every
// moving average has to have.
type MA interface {
	// Validate makes sure that the moving average is valid.
	Validate() error

	// Calc calculates moving average value by using settings stored in the func receiver.
	Calc(cc []decimal.Decimal) (decimal.Decimal, error)

	// Count determines the total amount of data points needed for moving averages
	// calculation by using settings stored in the receiver.
	Count() int
}

// SMA holds all the neccesary information needed to calculate simple
// moving average.
type SMA struct {
	// Length specifies how many data points should be used.
	Length int `json:"length"`
}

// Validate checks all SMA settings stored in func receiver to make sure that
// they're meeting each of their own requirements.
func (s SMA) Validate() error {
	if s.Length < 1 {
		return ErrInvalidLength
	}

	return nil
}

// Calc calculates SMA value by using settings stored in the func receiver.
func (s SMA) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, s.Count())
	if err != nil {
		return decimal.Zero, err
	}

	res := decimal.Zero

	for i := 0; i < len(dd); i++ {
		res = res.Add(dd[i])
	}

	return res.Div(decimal.NewFromInt(int64(s.Length))), nil
}

// Count determines the total amount of data points needed for SMA
// calculation by using settings stored in the receiver.
func (s SMA) Count() int {
	return s.Length
}

// ValidateSMA checks all settings passed as parameters to make sure that
// they're meeting each of their own requirements.
func ValidateSMA(len int) error {
	s := SMA{Length: len}
	return s.Validate()
}

// CalcSMA calculates SMA value by using settings passed as parameters.
func CalcSMA(dd []decimal.Decimal, len int) (decimal.Decimal, error) {
	s := SMA{Length: len}
	return s.Calc(dd)
}

// CountSMA determines the total amount of data points needed for SMA
// calculation by using settings passed as parameters.
func CountSMA(len int) int {
	s := SMA{Length: len}
	return s.Count()
}

// EMA holds all the neccesary information needed to calculate exponential
// moving average.
type EMA struct {
	// Length specifies how many data points should be used.
	Length int `json:"length"`
}

// Validate checks all EMA settings stored in func receiver to make sure that
// they're meeting each of their own requirements.
func (e EMA) Validate() error {
	if e.Length < 1 {
		return ErrInvalidLength
	}

	return nil
}

// Calc calculates EMA value by using settings stored in the func receiver.
func (e EMA) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, e.Count())

	if err != nil {
		return decimal.Zero, err
	}

	res, err := CalcSMA(dd[len(dd)-e.Length:], e.Length)

	if err != nil {
		return decimal.Zero, err
	}

	mul := e.multiplier()

	for i := e.Length; i < len(dd); i++ {
		res = dd[i].Mul(mul).Add(res.Mul(decimal.NewFromInt(1).Sub(mul)))
	}

	return res, nil
}

// multiplier calculates EMA multiplier value by using settings stored in the func receiver.
func (e EMA) multiplier() decimal.Decimal {
	return decimal.NewFromFloat(2.0 / float64(e.Length+1))
}

// Count determines the total amount of data points needed for EMA
// calculation by using settings stored in the receiver.
func (e EMA) Count() int {
	return e.Length * 2
}

// ValidateEMA checks all settings passed as parameters to make sure that
// they're meeting each of their own requirements.
func ValidateEMA(len int) error {
	e := EMA{Length: len}
	return e.Validate()
}

// CalcEMA calculates EMA value by using settings passed as parameters.
func CalcEMA(dd []decimal.Decimal, len int) (decimal.Decimal, error) {
	e := EMA{Length: len}
	return e.Calc(dd)
}

// CountEMA determines the total amount of data points needed for EMA
// calculation by using settings passed as parameters.
func CountEMA(len int) int {
	e := EMA{Length: len}
	return e.Count()
}

// WMA holds all the neccesary information needed to calculate weighted
// moving average.
type WMA struct {
	// Length specifies how many data points should be used.
	Length int `json:"length"`
}

// Validate checks all WMA settings stored in func receiver to make sure that
// they're meeting each of their own requirements.
func (w WMA) Validate() error {
	if w.Length < 1 {
		return ErrInvalidLength
	}

	return nil
}

// Calc calculates WMA value by using settings stored in the func receiver.
func (w WMA) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, w.Count())

	if err != nil {
		return decimal.Zero, err
	}

	res := decimal.Zero

	weight := decimal.NewFromFloat(float64(w.Length*(w.Length+1)) / 2.0)

	for i := 0; i < len(dd); i++ {
		res = res.Add(dd[i].Mul(decimal.NewFromInt(int64(i + 1)).Div(weight)))
	}

	return res, nil
}

// Count determines the total amount of data points needed for WMA
// calculation by using settings stored in the receiver.
func (w WMA) Count() int {
	return w.Length
}

// ValidateWMA checks all settings passed as parameters to make sure that
// they're meeting each of their own requirements.
func ValidateWMA(len int) error {
	w := WMA{Length: len}
	return w.Validate()
}

// CalcWMA calculates WMA value by using settings passed as parameters.
func CalcWMA(dd []decimal.Decimal, len int) (decimal.Decimal, error) {
	w := WMA{Length: len}
	return w.Calc(dd)
}

// CountWMA determines the total amount of data points needed for WMA
// calculation by using settings passed as parameters.
func CountWMA(len int) int {
	w := WMA{Length: len}
	return w.Count()
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
func (macd MACD) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, macd.Count())
	if err != nil {
		return decimal.Zero, err
	}

	res1, err := macd.MA1.Calc(dd)
	if err != nil {
		return decimal.Zero, err
	}

	res2, err := macd.MA2.Calc(dd)
	if err != nil {
		return decimal.Zero, err
	}

	res := res1.Sub(res2)

	return res, nil
}

// Count determines the total amount of data points needed for MACD
// calculation by using settings stored in the receiver.
func (macd MACD) Count() int {
	c1 := macd.MA1.Count()
	c2 := macd.MA2.Count()
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
func CalcMACD(dd []decimal.Decimal, ma1, ma2 MA) (decimal.Decimal, error) {
	macd := MACD{MA1: ma1, MA2: ma2}
	return macd.Calc(dd)
}

// CountMACD determines the total amount of data points needed for MACD
// calculation by using settings passed as parameters.
func CountMACD(ma1, ma2 MA) int {
	macd := MACD{MA1: ma1, MA2: ma2}
	return macd.Count()
}
