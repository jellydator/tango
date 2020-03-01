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
func ValidateRSI(l int) error {
	r := RSI{Length: l}
	return r.Validate()
}

// CalcRSI calculates RSI value by using settings passed as parameters.
func CalcRSI(dd []decimal.Decimal, l int) (decimal.Decimal, error) {
	r := RSI{Length: l}
	return r.Calc(dd)
}

// CountRSI determines the total amount of data points needed for RSI
// calculation by using settings passed as parameters.
func CountRSI(l int) int {
	r := RSI{Length: l}
	return r.Count()
}

// Stoch holds all the neccesary information needed to calculate stochastic
// oscillator.
type Stoch struct {
	// Length specifies how many data points should be used.
	Length int `json:"length"`
}

// Validate checks all stochastic settings stored in func receiver to make sure that
// they're meeting each of their own requirements.
func (s Stoch) Validate() error {
	if s.Length < 1 {
		return ErrInvalidLength
	}
	return nil
}

// Calc calculates stochastic value by using settings stored in the func receiver.
func (s Stoch) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
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

// Count determines the total amount of data points needed for stochastic
// calculation by using settings stored in the receiver.
func (s Stoch) Count() int {
	return s.Length
}

// ValidateStoch checks all settings passed as parameters to make sure that
// they're meeting each of their own requirements.
func ValidateStoch(l int) error {
	s := Stoch{Length: l}
	return s.Validate()
}

// CalcStoch calculates stochastic value by using settings passed as parameters.
func CalcStoch(dd []decimal.Decimal, l int) (decimal.Decimal, error) {
	s := Stoch{Length: l}
	return s.Calc(dd)
}

// CountStoch determines the total amount of data points needed for stochastic
// calculation by using settings passed as parameters.
func CountStoch(l int) int {
	s := Stoch{Length: l}
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
func ValidateROC(l int) error {
	r := ROC{Length: l}
	return r.Validate()
}

// CalcROC calculates ROC value by using settings passed as parameters.
func CalcROC(dd []decimal.Decimal, l int) (decimal.Decimal, error) {
	r := ROC{Length: l}
	return r.Calc(dd)
}

// CountROC determines the total amount of data points needed for ROC
// calculation by using settings passed as parameters.
func CountROC(l int) int {
	r := ROC{Length: l}
	return r.Count()
}

// CCI holds all the neccesary information needed to calculate commodity
// channel index.
type CCI struct {
	// MA configures moving average.
	MA MA `json:"ma"`
}

// Validate checks all CCI settings stored in func receiver to make sure that
// they're meeting each of their own requirements.
func (c CCI) Validate() error {
	if c.MA == nil {
		return ErrMANotSet
	}

	if err := c.MA.Validate(); err != nil {
		return err
	}
	return nil
}

// Calc calculates CCI value by using settings stored in the func receiver.
func (c CCI) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, c.Count())
	if err != nil {
		return decimal.Zero, err
	}

	ma, err := c.MA.Calc(dd)
	if err != nil {
		return decimal.Zero, err
	}

	return dd[len(dd)-1].Sub(ma).Div(decimal.NewFromFloat(0.015).Mul(meanDeviation(dd))).Round(8), nil
}

// Count determines the total amount of data points needed for CCI
// calculation by using settings stored in the receiver.
func (c CCI) Count() int {
	return c.MA.Count()
}

// ValidateCCI checks all settings passed as parameters to make sure that
// they're meeting each of their own requirements.
func ValidateCCI(MA MA) error {
	c := CCI{MA: MA}
	return c.Validate()
}

// CalcCCI calculates CCI value by using settings passed as parameters.
func CalcCCI(dd []decimal.Decimal, MA MA) (decimal.Decimal, error) {
	c := CCI{MA: MA}
	return c.Calc(dd)
}

// CountCCI determines the total amount of data points needed for CCI
// calculation by using settings passed as parameters.
func CountCCI(MA MA) int {
	c := CCI{MA: MA}
	return c.MA.Count()
}
