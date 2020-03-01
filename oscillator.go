package indc

import (
	"github.com/shopspring/decimal"
	"github.com/swithek/chartype"
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

	for i := 0; i < len(dd); i++ {
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

	// Length specifies how many candles should be used.
	Length int `json:"length"`

	// Src specifies which price field of the candle should be used.
	Src chartype.CandleField `json:"src"`
}

// Validate checks all STOCH settings stored in func receiver to make sure that
// they're meeting each of their own requirements.
func (s STOCH) Validate() error {
	if s.Length < 1 {
		return ErrInvalidLength
	}

	if err := s.Src.Validate(); err != nil {
		return err
	}

	return nil
}

// Calc calculates STOCH value by using settings stored in the func receiver.
func (s STOCH) Calc(cc []chartype.Candle) (decimal.Decimal, error) {
	if s.CandleCount() > len(cc) {
		return decimal.Zero, ErrInvalidCandleCount
	}

	l := s.Src.Extract(cc[len(cc)-s.CandleCount()])
	h := s.Src.Extract(cc[len(cc)-s.CandleCount()])

	for i := len(cc) - s.CandleCount() + 1; i < len(cc)-s.CandleCount()+s.Length; i++ {
		if s.Src.Extract(cc[i]).LessThan(l) {
			l = s.Src.Extract(cc[i])
		}
		if s.Src.Extract(cc[i]).GreaterThan(h) {
			h = s.Src.Extract(cc[i])
		}
	}

	return s.Src.Extract(cc[len(cc)-1]).Sub(l).Div(h.Sub(l)).Mul(decimal.NewFromInt(100)), nil
}

// CandleCount determines the total amount of candles needed for STOCH
// calculation by using settings stored in the receiver.
func (s STOCH) CandleCount() int {
	return s.Length
}

// ValidateSTOCH checks all settings passed as parameters to make sure that
// they're meeting each of their own requirements.
func ValidateSTOCH(len int, src chartype.CandleField) error {
	s := STOCH{Length: len, Src: src}
	return s.Validate()
}

// CalcSTOCH calculates STOCH value by using settings passed as parameters.
func CalcSTOCH(cc []chartype.Candle, len int, src chartype.CandleField) (decimal.Decimal, error) {
	s := STOCH{Length: len, Src: src}
	return s.Calc(cc)
}

// CandleCountSTOCH determines the total amount of candles needed for STOCH
// calculation by using settings passed as parameters.
func CandleCountSTOCH(len int) int {
	s := STOCH{Length: len}
	return s.CandleCount()
}

// ROC holds all the neccesary information needed to calculate rate
// of change.
type ROC struct {
	// Length specifies how many candles should be used.
	Length int `json:"length"`

	// Src specifies which price field of the candle should be used.
	Src chartype.CandleField `json:"src"`
}

// Validate checks all ROC settings stored in func receiver to make sure that
// they're meeting each of their own requirements.
func (r ROC) Validate() error {
	if r.Length < 1 {
		return ErrInvalidLength
	}

	if err := r.Src.Validate(); err != nil {
		return err
	}

	return nil
}

// Calc calculates ROC value by using settings stored in the func receiver.
func (r ROC) Calc(cc []chartype.Candle) (decimal.Decimal, error) {
	if r.CandleCount() > len(cc) {
		return decimal.Zero, ErrInvalidCandleCount
	}

	l := r.Src.Extract(cc[len(cc)-1])
	s := r.Src.Extract(cc[len(cc)-r.CandleCount()])
	return l.Sub(s).Div(s).Mul(decimal.NewFromInt(100)).Round(8), nil
}

// CandleCount determines the total amount of candles needed for ROC
// calculation by using settings stored in the receiver.
func (r ROC) CandleCount() int {
	return r.Length
}

// ValidateROC checks all settings passed as parameters to make sure that
// they're meeting each of their own requirements.
func ValidateROC(len int, src chartype.CandleField) error {
	r := ROC{Length: len, Src: src}
	return r.Validate()
}

// CalcROC calculates ROC value by using settings passed as parameters.
func CalcROC(cc []chartype.Candle, len int, src chartype.CandleField) (decimal.Decimal, error) {
	r := ROC{Length: len, Src: src}
	return r.Calc(cc)
}

// CandleCountROC determines the total amount of candles needed for ROC
// calculation by using settings passed as parameters.
func CandleCountROC(len int) int {
	r := ROC{Length: len}
	return r.CandleCount()
}
