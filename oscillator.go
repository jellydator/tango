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

// STOCH holds all the neccesary information needed to calculate stochastic
// oscillator.
type STOCH struct {

	// Length specifies how many candles should be used.
	Length int `json:"length"`

	// Offset specifies how many latest candles should be skipped.
	Offset int `json:"offset"`

	// Src specifies which price field of the candle should be used.
	Src chartype.CandleField `json:"src"`
}

// Validate checks all STOCH settings stored in func receiver to make sure that
// they're meeting each of their own requirements.
func (s STOCH) Validate() error {
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

	return s.Src.Extract(cc[len(cc)-s.Offset-1]).Sub(l).Div(h.Sub(l)).Mul(decimal.NewFromInt(100)), nil
}

// CandleCount determines the total amount of candles needed for STOCH
// calculation by using settings stored in the receiver.
func (s STOCH) CandleCount() int {
	return s.Length + s.Offset
}

// ValidateSTOCH checks all settings passed as parameters to make sure that
// they're meeting each of their own requirements.
func ValidateSTOCH(len, off int, src chartype.CandleField) error {
	s := STOCH{Length: len, Offset: off, Src: src}
	return s.Validate()
}

// CalcSTOCH calculates STOCH value by using settings passed as parameters.
func CalcSTOCH(cc []chartype.Candle, len, off int, src chartype.CandleField) (decimal.Decimal, error) {
	s := STOCH{Length: len, Offset: off, Src: src}
	return s.Calc(cc)
}

// CandleCountSTOCH determines the total amount of candles needed for STOCH
// calculation by using settings passed as parameters.
func CandleCountSTOCH(len, off int) int {
	s := STOCH{Length: len, Offset: off}
	return s.CandleCount()
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
	if r.CandleCount() > len(cc) {
		return decimal.Zero, ErrInvalidCandleCount
	}

	l := r.Src.Extract(cc[len(cc)-r.Offset-1])
	s := r.Src.Extract(cc[len(cc)-r.CandleCount()])
	return l.Sub(s).Div(s).Mul(decimal.NewFromInt(100)).Round(8), nil
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
