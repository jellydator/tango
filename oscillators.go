package indc

import (
	"github.com/shopspring/decimal"
	"github.com/swithek/chartype"
)

// RSI holds all the neccesary information needed to calculate relative
// strength index.
type RSI struct {
	// Length specifies how many candles should be used.
	Length int `json:"length"`

	// Offset specifies how many latest candles should be skipped.
	Offset int `json:"offset"`

	// Src specifies which price field of the candle should be used.
	Src chartype.CandleField `json:"src"`
}

// Validate checks all RSI settings stored in func receiver to make sure that
// they're meeting each of their own requirements.
func (r RSI) Validate() error {
	if r.Length < 1 {
		return ErrInvalidLength
	}

	if r.Offset < 0 {
		return ErrInvalidOffset
	}

	if err := r.Src.Validate(); err != nil {
		return err
	}

	return nil
}

// Calc calculates RSI value by using settings stored in the func receiver.
func (r RSI) Calc(cc []chartype.Candle) (decimal.Decimal, error) {
	if r.CandleCount() > len(cc) {
		return decimal.Zero, ErrInvalidCandleCount
	}

	ag := decimal.Zero
	al := decimal.Zero

	for i := len(cc) - r.CandleCount() + 1; i < len(cc)-r.CandleCount()+r.Length; i++ {
		if r.Src.Extract(cc[i]).Sub(r.Src.Extract(cc[i-1])).LessThan(decimal.Zero) {
			al = al.Add(r.Src.Extract(cc[i]).Sub(r.Src.Extract(cc[i-1])).Abs())
		} else {
			ag = ag.Add(r.Src.Extract(cc[i]).Sub(r.Src.Extract(cc[i-1])))
		}
	}

	ag = ag.Div(decimal.NewFromInt(int64(r.Length)))
	al = al.Div(decimal.NewFromInt(int64(r.Length)))

	return decimal.NewFromInt(100).Sub(decimal.NewFromInt(100).Div(decimal.NewFromInt(1).Add(ag.Div(al)))).Round(8), nil
}

// CandleCount determines the total amount of candles needed for RSI
// calculation by using settings stored in the receiver.
func (r RSI) CandleCount() int {
	return r.Length + r.Offset
}

// ValidateRSI checks all settings passed as parameters to make sure that
// they're meeting each of their own requirements.
func ValidateRSI(len, off int, src chartype.CandleField) error {
	r := RSI{Length: len, Offset: off, Src: src}
	return r.Validate()
}

// CalcRSI calculates RSI value by using settings passed as parameters.
func CalcRSI(cc []chartype.Candle, len, off int, src chartype.CandleField) (decimal.Decimal, error) {
	r := RSI{Length: len, Offset: off, Src: src}
	return r.Calc(cc)
}

// CandleCountRSI determines the total amount of candles needed for RSI
// calculation by using settings passed as parameters.
func CandleCountRSI(len, off int) int {
	r := RSI{Length: len, Offset: off}
	return r.CandleCount()
}

// ROC holds all the neccesary information needed to calculate rate
// of change.
type ROC struct {
	// Length specifies how many candles should be used.
	Length int `json:"length"`

	// Offset specifies how many latest candles should be skipped.
	Offset int `json:"offset"`

	// Src specifies which price field of the candle should be used.
	Src chartype.CandleField `json:"src"`
}

// Validate checks all ROC settings stored in func receiver to make sure that
// they're meeting each of their own requirements.
func (r ROC) Validate() error {
	if r.Length < 1 {
		return ErrInvalidLength
	}

	if r.Offset < 0 {
		return ErrInvalidOffset
	}

	if err := r.Src.Validate(); err != nil {
		return err
	}

	return nil
}

// Calc calculates ROC value by using settings stored in the func receiver.
func (r ROC) Calc(cc []chartype.Candle) (decimal.Decimal, error) {
	return decimal.Zero, nil
}

// CandleCount determines the total amount of candles needed for ROC
// calculation by using settings stored in the receiver.
func (r ROC) CandleCount() int {
	return r.Length + r.Offset
}

// ValidateROC checks all settings passed as parameters to make sure that
// they're meeting each of their own requirements.
func ValidateROC(len, off int, src chartype.CandleField) error {
	r := ROC{Length: len, Offset: off, Src: src}
	return r.Validate()
}

// CalcROC calculates ROC value by using settings passed as parameters.
func CalcROC(cc []chartype.Candle, len, off int, src chartype.CandleField) (decimal.Decimal, error) {
	r := ROC{Length: len, Offset: off, Src: src}
	return r.Calc(cc)
}

// CandleCountROC determines the total amount of candles needed for ROC
// calculation by using settings passed as parameters.
func CandleCountROC(len, off int) int {
	r := ROC{Length: len, Offset: off}
	return r.CandleCount()
}
