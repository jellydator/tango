package indc

import (
	"github.com/shopspring/decimal"
)

// RSI holds all the neccesary information needed to calculate relative
// strength index.
type RSI struct {
	// Length specifies how many data points should be used.
	Length int `json:"length"`
}

// Validate checks all RSI settings stored in func receiver to make sure that
// they're meeting each of their own requirements.
func (r RSI) Validate() error {
	if r.Length < 1 {
		return ErrInvalidLength
	}

	return nil
}

// Calc calculates RSI value by using settings stored in the func receiver.
func (r RSI) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, r.Count())
	if err != nil {
		return decimal.Zero, err
	}

	ag := decimal.Zero
	al := decimal.Zero

	for i := 1; i < len(dd); i++ {
		if dd[i].Sub(dd[i-1]).LessThan(decimal.Zero) {
			al = al.Add(dd[i].Sub(dd[i-1]).Abs())
		} else {
			ag = ag.Add(dd[i].Sub(dd[i-1]))
		}
	}

	ag = ag.Div(decimal.NewFromInt(int64(r.Length)))
	al = al.Div(decimal.NewFromInt(int64(r.Length)))

	return decimal.NewFromInt(100).Sub(decimal.NewFromInt(100).Div(decimal.NewFromInt(1).Add(ag.Div(al)))).Round(8), nil
}

// Count determines the total amount of data points needed for RSI
// calculation by using settings stored in the receiver.
func (r RSI) Count() int {
	return r.Length
}

// ValidateRSI checks all settings passed as parameters to make sure that
// they're meeting each of their own requirements.
func ValidateRSI(len int) error {
	r := RSI{Length: len}
	return r.Validate()
}

// CalcRSI calculates RSI value by using settings passed as parameters.
func CalcRSI(dd []decimal.Decimal, len int) (decimal.Decimal, error) {
	r := RSI{Length: len}
	return r.Calc(dd)
}

// CountRSI determines the total amount of data points needed for RSI
// calculation by using settings passed as parameters.
func CountRSI(len int) int {
	r := RSI{Length: len}
	return r.Count()
}

// STOCH holds all the neccesary information needed to calculate stochastic
// oscillator.
type STOCH struct {
	// Length specifies how many data points should be used.
	Length int `json:"length"`
}

// Validate checks all STOCH settings stored in func receiver to make sure that
// they're meeting each of their own requirements.
func (s STOCH) Validate() error {
	if s.Length < 1 {
		return ErrInvalidLength
	}

	return nil
}

// Calc calculates STOCH value by using settings stored in the func receiver.
func (s STOCH) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, s.Count())
	if err != nil {
		return decimal.Zero, err
	}

	l := dd[0]
	h := dd[0]

	for i := 0; i < len(dd); i++ {
		if dd[i].LessThan(l) {
			l = dd[i]
		}
		if dd[i].GreaterThan(h) {
			h = dd[i]
		}
	}

	return dd[len(dd)-1].Sub(l).Div(h.Sub(l)).Mul(decimal.NewFromInt(100)), nil
}

// Count determines the total amount of data points needed for STOCH
// calculation by using settings stored in the receiver.
func (s STOCH) Count() int {
	return s.Length
}

// ValidateSTOCH checks all settings passed as parameters to make sure that
// they're meeting each of their own requirements.
func ValidateSTOCH(len int) error {
	s := STOCH{Length: len}
	return s.Validate()
}

// CalcSTOCH calculates STOCH value by using settings passed as parameters.
func CalcSTOCH(dd []decimal.Decimal, len int) (decimal.Decimal, error) {
	s := STOCH{Length: len}
	return s.Calc(dd)
}

// CountSTOCH determines the total amount of data points needed for STOCH
// calculation by using settings passed as parameters.
func CountSTOCH(len int) int {
	s := STOCH{Length: len}
	return s.Count()
}

// ROC holds all the neccesary information needed to calculate rate
// of change.
type ROC struct {
	// Length specifies how many data points should be used.
	Length int `json:"length"`
}

// Validate checks all ROC settings stored in func receiver to make sure that
// they're meeting each of their own requirements.
func (r ROC) Validate() error {
	if r.Length < 1 {
		return ErrInvalidLength
	}

	return nil
}

// Calc calculates ROC value by using settings stored in the func receiver.
func (r ROC) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, r.Count())
	if err != nil {
		return decimal.Zero, err
	}

	l := dd[len(dd)-1]
	s := dd[0]
	return l.Sub(s).Div(s).Mul(decimal.NewFromInt(100)).Round(8), nil
}

// Count determines the total amount of data points needed for ROC
// calculation by using settings stored in the receiver.
func (r ROC) Count() int {
	return r.Length
}

// ValidateROC checks all settings passed as parameters to make sure that
// they're meeting each of their own requirements.
func ValidateROC(len int) error {
	r := ROC{Length: len}
	return r.Validate()
}

// CalcROC calculates ROC value by using settings passed as parameters.
func CalcROC(dd []decimal.Decimal, len int) (decimal.Decimal, error) {
	r := ROC{Length: len}
	return r.Calc(dd)
}

// CountROC determines the total amount of data points needed for ROC
// calculation by using settings passed as parameters.
func CountROC(len int) int {
	r := ROC{Length: len}
	return r.Count()
}
