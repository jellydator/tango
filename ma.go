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
	Calc(dd []decimal.Decimal) (decimal.Decimal, error)

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
func ValidateSMA(l int) error {
	s := SMA{Length: l}
	return s.Validate()
}

// CalcSMA calculates SMA value by using settings passed as parameters.
func CalcSMA(dd []decimal.Decimal, l int) (decimal.Decimal, error) {
	s := SMA{Length: l}
	return s.Calc(dd)
}

// CountSMA determines the total amount of data points needed for SMA
// calculation by using settings passed as parameters.
func CountSMA(l int) int {
	s := SMA{Length: l}
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
func ValidateEMA(l int) error {
	e := EMA{Length: l}
	return e.Validate()
}

// CalcEMA calculates EMA value by using settings passed as parameters.
func CalcEMA(dd []decimal.Decimal, l int) (decimal.Decimal, error) {
	e := EMA{Length: l}
	return e.Calc(dd)
}

// CountEMA determines the total amount of data points needed for EMA
// calculation by using settings passed as parameters.
func CountEMA(l int) int {
	e := EMA{Length: l}
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
func ValidateWMA(l int) error {
	w := WMA{Length: l}
	return w.Validate()
}

// CalcWMA calculates WMA value by using settings passed as parameters.
func CalcWMA(dd []decimal.Decimal, l int) (decimal.Decimal, error) {
	w := WMA{Length: l}
	return w.Calc(dd)
}

// CountWMA determines the total amount of data points needed for WMA
// calculation by using settings passed as parameters.
func CountWMA(l int) int {
	w := WMA{Length: l}
	return w.Count()
}
