package indc

import "github.com/shopspring/decimal"

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
