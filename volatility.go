package indc

import "github.com/shopspring/decimal"

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
