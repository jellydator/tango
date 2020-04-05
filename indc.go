package indc

import (
	"math"

	"github.com/shopspring/decimal"
)

// Aroon holds all the neccesary information needed to calculate aroon.
type Aroon struct {
	// Trend configures which aroon trend to measure (it can either be up or down).
	Trend string `json:"trend"`

	// Length specifies how many data points should be used.
	Length int `json:"length"`
}

// Validate checks all Aroon settings stored in func receiver to make sure that
// they're meeting each of their own requirements.
func (a Aroon) Validate() error {
	if a.Trend != "down" && a.Trend != "up" {
		return ErrInvalidType
	}

	if a.Length < 1 {
		return ErrInvalidLength
	}

	return nil
}

// Calc calculates Aroon value by using settings stored in the func receiver.
func (a Aroon) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, a.Count())
	if err != nil {
		return decimal.Zero, err
	}

	v := decimal.Zero
	p := decimal.Zero

	for i := 0; i < len(dd); i++ {
		if v.Equal(decimal.Zero) {
			v = dd[i]
		}
		if a.Trend == "up" && v.LessThanOrEqual(dd[i]) || a.Trend == "down" && !v.LessThan(dd[i]) {
			v = dd[i]
			p = decimal.NewFromInt(int64(a.Length - i - 1))
		}
	}

	return decimal.NewFromInt(int64(a.Length)).Sub(p).Mul(decimal.NewFromInt(100)).Div(decimal.NewFromInt(int64(a.Length))), nil
}

// Count determines the total amount of data points needed for Aroon
// calculation by using settings stored in the receiver.
func (a Aroon) Count() int {
	return a.Length
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

// DEMA holds all the neccesary information needed to calculate double exponential
// moving average.
type DEMA struct {
	// Length specifies how many data points should be used.
	Length int `json:"length"`
}

// Validate checks all DEMA settings stored in func receiver to make sure that
// they're meeting each of their own requirements.
func (d DEMA) Validate() error {
	if d.Length < 1 {
		return ErrInvalidLength
	}
	return nil
}

// Calc calculates DEMA value by using settings stored in the func receiver.
func (d DEMA) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, d.Count())
	if err != nil {
		return decimal.Zero, err
	}

	r := make([]decimal.Decimal, d.Length)

	sma := SMA{Length: d.Length}
	r[0], err = sma.Calc(dd[:d.Length])
	if err != nil {
		return decimal.Zero, err
	}

	e := EMA{Length: d.Length}

	for i := d.Length; i < len(dd); i++ {
		r[i-d.Length+1] = e.CalcNext(r[i-d.Length], dd[i])
	}

	v := r[0]

	for i := 0; i < len(r); i++ {
		v = e.CalcNext(v, r[i])
	}

	return v.Round(8), nil
}

// Count determines the total amount of data points needed for DEMA
// calculation by using settings stored in the receiver.
func (d DEMA) Count() int {
	return d.Length*2 - 1
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

	s := SMA{Length: e.Length}
	r, err := s.Calc(dd[:e.Length])
	if err != nil {
		return decimal.Zero, err
	}

	for i := e.Length; i < len(dd); i++ {
		r = e.CalcNext(r, dd[i])
	}

	return r.Round(8), nil
}

// CalcNext calculates sequential EMA value by using previous ema.
func (e EMA) CalcNext(l, n decimal.Decimal) decimal.Decimal {
	mul := e.multiplier()
	return n.Mul(mul).Add(l.Mul(decimal.NewFromInt(1).Sub(mul)))
}

// multiplier calculates EMA multiplier value by using settings stored in the func receiver.
func (e EMA) multiplier() decimal.Decimal {
	return decimal.NewFromFloat(2.0 / float64(e.Length+1))
}

// Count determines the total amount of data points needed for EMA
// calculation by using settings stored in the receiver.
func (e EMA) Count() int {
	return e.Length*2 - 1
}

// HMA holds all the neccesary information needed to calculate hull moving average.
type HMA struct {
	// WMA configures base moving average.
	WMA WMA `json:"wma"`
}

// Validate checks all HMA settings stored in func receiver to make sure that
// they're meeting each of their own requirements.
func (h HMA) Validate() error {
	if h.WMA == (WMA{}) {
		return ErrMANotSet
	}
	if h.WMA.Length < 1 {
		return ErrInvalidLength
	}
	return nil
}

// Calc calculates HMA value by using settings stored in the func receiver.
func (h HMA) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, h.Count())
	if err != nil {
		return decimal.Zero, err
	}

	l := int(math.Sqrt(float64(h.WMA.Count())))

	w1 := WMA{Length: h.WMA.Count() / 2}
	w2 := h.WMA
	w3 := WMA{Length: l}

	v := make([]decimal.Decimal, l)

	for i := 0; i < l; i++ {
		r1, err := w1.Calc(dd[:len(dd)-l+i + 1])
		if err != nil {
			return decimal.Zero, nil
		}

		r2, err := w2.Calc(dd[:len(dd)-l+i + 1])
		if err != nil {
			return decimal.Zero, nil
		}

		v[i] = r1.Mul(decimal.NewFromInt(2)).Sub(r2)
	}

	r, err := w3.Calc(v)
	if err != nil {
		return decimal.Zero, err
	}
	return r, nil
}

// Count determines the total amount of data points needed for HMA
// calculation by using settings stored in the receiver.
func (h HMA) Count() int {
	return h.WMA.Count()*2 - 1
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
func (m MACD) Validate() error {
	if m.MA1 == nil || m.MA2 == nil {
		return ErrMANotSet
	}

	if err := m.MA1.Validate(); err != nil {
		return err
	}

	if err := m.MA2.Validate(); err != nil {
		return err
	}

	return nil
}

// Calc calculates MACD value by using settings stored in the func receiver.
func (m MACD) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, m.Count())
	if err != nil {
		return decimal.Zero, err
	}

	r1, err := m.MA1.Calc(dd)
	if err != nil {
		return decimal.Zero, err
	}

	r2, err := m.MA2.Calc(dd)
	if err != nil {
		return decimal.Zero, err
	}

	r := r1.Sub(r2)

	return r, nil
}

// Count determines the total amount of data points needed for MACD
// calculation by using settings stored in the receiver.
func (m MACD) Count() int {
	c1 := m.MA1.Count()
	c2 := m.MA2.Count()

	if c1 > c2 {
		return c1
	}

	return c2
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

	r := decimal.Zero

	for i := 0; i < len(dd); i++ {
		r = r.Add(dd[i])
	}

	return r.Div(decimal.NewFromInt(int64(s.Length))), nil
}

// Count determines the total amount of data points needed for SMA
// calculation by using settings stored in the receiver.
func (s SMA) Count() int {
	return s.Length
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

	r := decimal.Zero

	wi := decimal.NewFromFloat(float64(w.Length*(w.Length+1)) / 2.0)

	for i := 0; i < len(dd); i++ {
		r = r.Add(dd[i].Mul(decimal.NewFromInt(int64(i + 1)).Div(wi)))
	}

	return r, nil
}

// Count determines the total amount of data points needed for WMA
// calculation by using settings stored in the receiver.
func (w WMA) Count() int {
	return w.Length
}

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
