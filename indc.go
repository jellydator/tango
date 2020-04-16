package indc

import (
	"encoding/json"
	"math"

	"github.com/shopspring/decimal"
)

// Aroon holds all the neccesary information needed to calculate aroon.
type Aroon struct {
	// Trend configures which aroon trend to measure (it can either
	// be up or down).
	Trend string `json:"trend"`

	// Length specifies how many data points should be used
	// in calculations.
	Length int `json:"length"`
}

// InitAroon verifies provided values and
// initializes aroon indicator.
func InitAroon(t string, lh int) (Aroon, error) {
	a := Aroon{Trend: t, Length: lh}

	if err := a.Validate(); err != nil {
		return Aroon{}, err
	}

	return a, nil
}

// Validate checks all Aroon settings stored in func receiver to
// make sure that they're matching their requirements.
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

		if a.Trend == "up" && v.LessThanOrEqual(dd[i]) ||
			a.Trend == "down" && !v.LessThan(dd[i]) {
			v = dd[i]
			p = decimal.NewFromInt(int64(a.Length - i - 1))
		}
	}

	return decimal.NewFromInt(int64(a.Length)).Sub(p).
		Mul(decimal.NewFromInt(100)).Div(decimal.NewFromInt(int64(a.Length))), nil
}

// Count determines the total amount of data points needed for Aroon
// calculation by using settings stored in the receiver.
func (a Aroon) Count() int {
	return a.Length
}

// CCI holds all the neccesary information needed to calculate commodity
// channel index.
type CCI struct {
	// Origin configures what calculations to use when computing CCI value.
	Origin Source `json:"origin"`
}

// InitCCI verifies provided values and
// initializes commodity channel index indicator.
func InitCCI(o Source) (CCI, error) {
	c := CCI{Origin: o}

	if err := c.Validate(); err != nil {
		return CCI{}, err
	}

	return c, nil
}

// Validate checks all CCI settings stored in func receiver to make sure that
// they're matching their requirements.
func (c CCI) Validate() error {
	if err := c.Origin.Validate(); err != nil {
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

	m, err := c.Origin.Calc(dd)
	if err != nil {
		return decimal.Zero, err
	}

	return dd[len(dd)-1].Sub(m).Div(decimal.NewFromFloat(0.015).
		Mul(meanDeviation(dd))), nil
}

// Count determines the total amount of data points needed for CCI
// calculation by using settings stored in the receiver.
func (c CCI) Count() int {
	return c.Origin.Count()
}

// DEMA holds all the neccesary information needed to calculate
// double exponential moving average.
type DEMA struct {
	// Length specifies how many data points should be used
	// in calculations.
	Length int `json:"length"`
}

// InitDEMA verifies provided values and
// initializes double exponential moving average indicator.
func InitDEMA(lh int) (DEMA, error) {
	d := DEMA{Length: lh}

	if err := d.Validate(); err != nil {
		return DEMA{}, err
	}

	return d, nil
}

// Validate checks all DEMA settings stored in func receiver to
// make sure that they're matching their requirements.
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

	v := make([]decimal.Decimal, d.Length)

	s := SMA{Length: d.Length}
	v[0], err = s.Calc(dd[:d.Length])
	if err != nil {
		return decimal.Zero, err
	}

	e := EMA{Length: d.Length}

	for i := d.Length; i < len(dd); i++ {
		v[i-d.Length+1] = e.CalcNext(v[i-d.Length], dd[i])
	}

	r := v[0]

	for i := 0; i < len(v); i++ {
		r = e.CalcNext(r, v[i])
	}

	return r, nil
}

// Count determines the total amount of data points needed for DEMA
// calculation by using settings stored in the receiver.
func (d DEMA) Count() int {
	return d.Length*2 - 1
}

// EMA holds all the neccesary information needed to calculate exponential
// moving average.
type EMA struct {
	// Length specifies how many data points should be used
	// in calculations.
	Length int `json:"length"`
}

// InitEMA verifies provided values and
// initializes exponential moving average indicator.
func InitEMA(lh int) (EMA, error) {
	e := EMA{Length: lh}

	if err := e.Validate(); err != nil {
		return EMA{}, err
	}

	return e, nil
}

// Validate checks all EMA settings stored in func receiver to make sure that
// they're matching their requirements.
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

	return r, nil
}

// CalcNext calculates sequential EMA value by using previous ema.
func (e EMA) CalcNext(l, n decimal.Decimal) decimal.Decimal {
	m := e.multiplier()
	return n.Mul(m).Add(l.Mul(decimal.NewFromInt(1).Sub(m)))
}

// multiplier calculates EMA multiplier value by using settings stored
// in the func receiver.
func (e EMA) multiplier() decimal.Decimal {
	return decimal.NewFromFloat(2.0 / float64(e.Length+1))
}

// Count determines the total amount of data points needed for EMA
// calculation by using settings stored in the receiver.
func (e EMA) Count() int {
	return e.Length*2 - 1
}

// HMA holds all the neccesary information needed to calculate
// hull moving average.
type HMA struct {
	// WMA configures base moving average.
	WMA WMA `json:"wma"`
}

// InitHMA verifies provided values and
// initializes hull moving average indicator.
func InitHMA(w WMA) (HMA, error) {
	h := HMA{WMA: w}

	if err := h.Validate(); err != nil {
		return HMA{}, err
	}

	return h, nil
}

// Validate checks all HMA settings stored in func receiver to make sure that
// they're matching their requirements.
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
		r1, err := w1.Calc(dd[:len(dd)-l+i+1])
		if err != nil {
			return decimal.Zero, nil
		}

		r2, err := w2.Calc(dd[:len(dd)-l+i+1])
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

// MACD holds all the neccesary information needed to calculate
// difference between two source indicators.
type MACD struct {
	// Origin1 configures what calculations to use when computing first
	// macd value.
	Origin1 Source `json:"origin1"`

	// Origin2 configures what calculations to use when computing second
	// macd value.
	Origin2 Source `json:"origin2"`
}

// InitMACD verifies provided values and
// initializes MACD indicator.
func InitMACD(o1, o2 Source) (MACD, error) {
	m := MACD{Origin1: o1, Origin2: o2}

	if err := m.Validate(); err != nil {
		return MACD{}, err
	}

	return m, nil
}

// Validate checks all MACD settings stored in func receiver
// to make sure that they're matching their requirements.
func (m MACD) Validate() error {
	if err := m.Origin1.Validate(); err != nil {
		return err
	}

	if err := m.Origin2.Validate(); err != nil {
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

	r1, err := m.Origin1.Calc(dd)
	if err != nil {
		return decimal.Zero, err
	}

	r2, err := m.Origin2.Calc(dd)
	if err != nil {
		return decimal.Zero, err
	}

	r := r1.Sub(r2)

	return r, nil
}

// Count determines the total amount of data points needed for MACD
// calculation by using settings stored in the receiver.
func (m MACD) Count() int {
	c1 := m.Origin1.Count()
	c2 := m.Origin2.Count()

	if c1 > c2 {
		return c1
	}

	return c2
}

// ROC holds all the neccesary information needed to calculate rate
// of change.
type ROC struct {
	// Length specifies how many data points should be used
	// in calculations.
	Length int `json:"length"`
}

// InitROC verifies provided values and
// initializes rate of change indicator.
func InitROC(lh int) (ROC, error) {
	r := ROC{Length: lh}

	if err := r.Validate(); err != nil {
		return ROC{}, err
	}

	return r, nil
}

// Validate checks all ROC settings stored in func receiver to make sure that
// they're matching their requirements.
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

	n := dd[len(dd)-1]
	l := dd[0]

	return n.Sub(l).Div(l).Mul(decimal.NewFromInt(100)), nil
}

// Count determines the total amount of data points needed for ROC
// calculation by using settings stored in the receiver.
func (r ROC) Count() int {
	return r.Length
}

// RSI holds all the neccesary information needed to calculate relative
// strength index.
type RSI struct {
	// Length specifies how many data points should be used
	// in calculations.
	Length int `json:"length"`
}

// InitRSI verifies provided values and
// initializes relative strength index indicator.
func InitRSI(lh int) (RSI, error) {
	r := RSI{Length: lh}

	if err := r.Validate(); err != nil {
		return RSI{}, err
	}

	return r, nil
}

// Validate checks all RSI settings stored in func receiver to make sure that
// they're matching their requirements.
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

	return decimal.NewFromInt(100).Sub(decimal.NewFromInt(100).
		Div(decimal.NewFromInt(1).Add(ag.Div(al)))), nil
}

// Count determines the total amount of data points needed for RSI
// calculation by using settings stored in the receiver.
func (r RSI) Count() int {
	return r.Length
}

// SMA holds all the neccesary information needed to calculate simple
// moving average.
type SMA struct {
	// Length specifies how many data points should be used
	// in calculations.
	Length int `json:"length"`
}

// InitSMA verifies provided values and
// initializes simple moving average indicator.
func InitSMA(lh int) (SMA, error) {
	s := SMA{Length: lh}

	if err := s.Validate(); err != nil {
		return SMA{}, err
	}

	return s, nil
}

// Validate checks all SMA settings stored in func receiver to make sure that
// they're matching their requirements.
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
	// Length specifies how many data points should be used
	// in calculations.
	Length int `json:"length"`
}

// InitStoch verifies provided values and
// initializes stochastic indicator.
func InitStoch(lh int) (Stoch, error) {
	s := Stoch{Length: lh}

	if err := s.Validate(); err != nil {
		return Stoch{}, err
	}

	return s, nil
}

// Validate checks all stochastic settings stored in func receiver to make
// sure that they're matching their requirements.
func (s Stoch) Validate() error {
	if s.Length < 1 {
		return ErrInvalidLength
	}
	return nil
}

// Calc calculates stochastic value by using settings stored in
// the func receiver.
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
	// Length specifies how many data points should be used
	// in calculations.
	Length int `json:"length"`
}

// InitWMA verifies provided values and
// initializes weighted moving average indicator.
func InitWMA(lh int) (WMA, error) {
	w := WMA{Length: lh}

	if err := w.Validate(); err != nil {
		return WMA{}, err
	}

	return w, nil
}

// Validate checks all WMA settings stored in func receiver to make sure that
// they're matching their requirements.
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

// Indicator is an interface that every indicator should implement.
type Indicator interface {
	// Validate should check whether the configuration options are
	// of a valid format.
	Validate() error

	// Calc should calculate and return indicator's value.
	Calc(dd []decimal.Decimal) (decimal.Decimal, error)

	// Count shoul determines the total amount of data points needed
	// for indicator's calculation.
	Count() int
}

// Source is a wrapper type allowing a more convenient work with the
// indicator interface.
type Source struct {
	Indicator
}

// InitSource verifies provided indicator and
// initializes source.
func InitSource(i Indicator) (Source, error) {
	if err := i.Validate(); err != nil {
		return Source{}, err
	}

	return Source{i}, nil
}

// Validate checks all Source values stored in func receiver to make sure
// that they're matching provided requirements.
func (s Source) Validate() error {
	if err := s.Validate(); err != nil {
		return err
	}

	return nil
}

// UnmarshalJSON parse JSON into an indicator source.
func (s *Source) UnmarshalJSON(d []byte) error {
	var id struct {
		Name string `json:"name"`
	}

	if err := json.Unmarshal(d, &id.Name); err != nil {
		return err
	}

	ind, err := newIndicator(id.Name)
	if err != nil {
		return ErrInvalidSrcName
	}

	if err := json.Unmarshal(d, &ind); err != nil {
		return err
	}

	s.Indicator = ind

	return nil
}

// MarshalJSON converts source data into JSON.
func (s *Source) MarshalJSON() ([]byte, error) {
	name, err := indicatorName(s.Indicator)
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		Indicator
		Name string `json:"name"`
	}{Indicator: s.Indicator, Name: name})
}
