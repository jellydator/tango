package indc

import (
	"encoding/json"
	"errors"
	"math"

	"github.com/shopspring/decimal"
)

// Indicator is an interface that every indicator should implement.
//go:generate moq -out ./indicator_mock_test.go . Indicator
type Indicator interface {
	// Calc should calculate indicator's value.
	Calc(dd []decimal.Decimal) (decimal.Decimal, error)

	// Count should determine the total amount of data points needed
	// for indicator's calculation.
	Count() int

	// namedMarshalJSON converts indicator and its name to JSON.
	namedMarshalJSON() ([]byte, error)
}

// NameAroon returns Aroon indicator name.
const NameAroon = "aroon"

// Aroon holds all the necessary information needed to calculate Aroon.
// The zero value is not usable.
type Aroon struct {
	// trend specifies which aroon trend to use during the
	// calculation process. Allowed values: up, down.
	trend String

	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewAroon validates provided configuration options and
// creates Aroon indicator.
func NewAroon(trend String, length int) (Aroon, error) {
	a := Aroon{trend: trend, length: length}

	if err := a.validate(); err != nil {
		return Aroon{}, err
	}

	return a, nil
}

// Length returns length configuration option.
func (a Aroon) Length() int {
	return a.length
}

// Trend returns trend configuration option.
func (a Aroon) Trend() String {
	return a.trend
}

// validate checks whether Aroon was configured properly or not.
func (a Aroon) validate() error {
	if a.trend != "down" && a.trend != "up" {
		return errors.New("invalid trend")
	}

	if a.length < 1 {
		return ErrInvalidLength
	}

	return nil
}

// Calc calculates Aroon from the provided data slice.
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

		if a.trend == "up" && v.LessThanOrEqual(dd[i]) ||
			a.trend == "down" && !v.LessThan(dd[i]) {

			v = dd[i]
			p = decimal.NewFromInt(int64(a.length - i - 1))
		}
	}

	return decimal.NewFromInt(int64(a.length)).Sub(p).
		Mul(decimal.NewFromInt(100)).Div(decimal.NewFromInt(int64(a.length))), nil
}

// Count determines the total amount of data needed for Aroon
// calculation.
func (a Aroon) Count() int {
	return a.length
}

// UnmarshalJSON parses JSON into Aroon structure.
func (a *Aroon) UnmarshalJSON(d []byte) error {
	var i struct {
		T String `json:"trend"`
		L int    `json:"length"`
	}

	if err := json.Unmarshal(d, &i); err != nil {
		return err
	}

	a.trend = i.T
	a.length = i.L

	if err := a.validate(); err != nil {
		return err
	}

	return nil
}

// MarshalJSON converts Aroon configuration data into JSON.
func (a Aroon) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		T String `json:"trend"`
		L int    `json:"length"`
	}{
		T: a.trend, L: a.length,
	})
}

// namedMarshalJSON converts Aroon configuration data with its
// name into JSON.
func (a Aroon) namedMarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		N String `json:"name"`
		T String `json:"trend"`
		L int    `json:"length"`
	}{
		N: NameAroon,
		T: a.trend,
		L: a.length,
	})
}

// NameCCI returns CCI indicator name.
const NameCCI = "cci"

// CCI holds all the necessary information needed to calculate commodity
// channel index.
// The zero value is not usable.
type CCI struct {
	// source specifies the base indicator to be used by the CCI.
	source Indicator
}

// NewCCI validates provided configuration options and creates commodity
// channel index indicator.
func NewCCI(source Indicator) (CCI, error) {
	c := CCI{source: source}

	if err := c.validate(); err != nil {
		return CCI{}, err
	}

	return c, nil
}

// Sub returns source configuration option.
func (c CCI) Sub() Indicator {
	return c.source
}

// validate checks whether CCI was configured properly or not.
func (c CCI) validate() error {
	if c.source == nil {
		return ErrInvalidSource
	}

	return nil
}

// Calc calculates CCI from the provided data slice.
func (c CCI) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, c.Count())
	if err != nil {
		return decimal.Zero, err
	}

	m, err := c.source.Calc(dd)
	if err != nil {
		return decimal.Zero, err
	}

	denom := decimal.NewFromFloat(0.015).Mul(meanDeviation(dd))

	if denom.Equal(decimal.Zero) {
		return decimal.Zero, nil
	}

	return dd[len(dd)-1].Sub(m).Div(denom), nil
}

// Count determines the total amount of data needed for CCI
// calculation.
func (c CCI) Count() int {
	return c.source.Count()
}

// UnmarshalJSON parses JSON into CCI structure.
func (c *CCI) UnmarshalJSON(d []byte) error {
	var i struct {
		S json.RawMessage `json:"source"`
	}

	if err := json.Unmarshal(d, &i); err != nil {
		return err
	}

	s, err := fromJSON(i.S)
	if err != nil {
		return err
	}

	c.source = s

	return nil
}

// MarshalJSON converts CCI configuration data into JSON.
func (c CCI) MarshalJSON() ([]byte, error) {
	s, err := c.source.namedMarshalJSON()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		S json.RawMessage `json:"source"`
	}{
		S: s,
	})
}

// namedMarshalJSON converts CCI configuration data with its
// name into JSON.
func (c CCI) namedMarshalJSON() ([]byte, error) {
	s, err := c.source.namedMarshalJSON()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		N String          `json:"name"`
		S json.RawMessage `json:"source"`
	}{
		N: NameCCI,
		S: s,
	})
}

// NameDEMA returns DEMA indicator name.
const NameDEMA = "dema"

// DEMA holds all the necessary information needed to calculate
// double exponential moving average.
// The zero value is not usable.
type DEMA struct {
	// ema specifies what ema should be used for dema calculations.
	ema EMA
}

// NewDEMA validates provided configuration options and creates double
// exponential moving average indicator.
func NewDEMA(ema EMA) (DEMA, error) {
	d := DEMA{ema: ema}

	if err := d.validate(); err != nil {
		return DEMA{}, err
	}

	return d, nil
}

// Length returns length configuration option.
func (dm DEMA) Length() int {
	return dm.ema.sma.Length()
}

// validate checks whether DEMA was configured properly or not.
func (dm DEMA) validate() error {
	if err := dm.ema.validate(); err != nil {
		return err
	}

	return nil
}

// Calc calculates DEMA from the provided data slice.
func (dm DEMA) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, dm.Count())
	if err != nil {
		return decimal.Zero, err
	}

	v := make([]decimal.Decimal, dm.Length())

	v[0], _ = dm.ema.sma.Calc(dd[:dm.Length()])

	for i := dm.Length(); i < len(dd); i++ {
		v[i-dm.Length()+1] = dm.ema.CalcNext(v[i-dm.Length()], dd[i])
	}

	r := v[0]

	for i := 0; i < len(v); i++ {
		r = dm.ema.CalcNext(r, v[i])
	}

	return r, nil
}

// Count determines the total amount of data needed for DEMA
// calculation.
func (dm DEMA) Count() int {
	return dm.ema.sma.Length()*2 - 1
}

// UnmarshalJSON parses JSON into DEMA structure.
func (dm *DEMA) UnmarshalJSON(d []byte) error {
	var i struct {
		EMA struct {
			SMA struct {
				L int `json:"length"`
			} `json:"sma"`
		} `json:"ema"`
	}

	if err := json.Unmarshal(d, &i); err != nil {
		return err
	}

	s, err := NewSMA(i.EMA.SMA.L)
	if err != nil {
		return err
	}

	e, err := NewEMA(s)
	if err != nil {
		return err
	}

	dm.ema = e

	if err := dm.validate(); err != nil {
		return err
	}

	return nil
}

// MarshalJSON converts DEMA configuration data into JSON.
func (dm DEMA) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		E EMA `json:"ema"`
	}{
		E: dm.ema,
	})
}

// namedMarshalJSON converts DEMA configuration data with its
// name into JSON.
func (dm DEMA) namedMarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		N String `json:"name"`
		E EMA    `json:"ema"`
	}{
		N: NameDEMA,
		E: dm.ema,
	})
}

// NameEMA returns EMA indicator name.
const NameEMA = "ema"

// EMA holds all the necessary information needed to calculate exponential
// moving average.
// The zero value is not usable.
type EMA struct {
	// sma specifies first EMA calculations SMA parameters.
	sma SMA
}

// NewEMA validates provided configuration options and
// creates exponential moving average indicator.
func NewEMA(sma SMA) (EMA, error) {
	e := EMA{sma: sma}

	if err := e.validate(); err != nil {
		return EMA{}, err
	}

	return e, nil
}

// Length returns length configuration option.
func (e EMA) Length() int {
	return e.sma.Count()
}

// validate checks whether EMA was configured properly or not.
func (e EMA) validate() error {
	if err := e.sma.validate(); err != nil {
		return err
	}

	return nil
}

// Calc calculates EMA from the provided data slice.
func (e EMA) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, e.Count())
	if err != nil {
		return decimal.Zero, err
	}

	r, _ := e.sma.Calc(dd[:e.Length()])

	for i := e.Length(); i < len(dd); i++ {
		r = e.CalcNext(r, dd[i])
	}

	return r, nil
}

// CalcNext calculates sequential EMA by using previous ema.
func (e EMA) CalcNext(l, n decimal.Decimal) decimal.Decimal {
	m := e.multiplier()
	return n.Mul(m).Add(l.Mul(decimal.NewFromInt(1).Sub(m)))
}

// multiplier calculates EMA multiplier.
func (e EMA) multiplier() decimal.Decimal {
	return decimal.NewFromFloat(2.0 / float64(e.Length()+1))
}

// Count determines the total amount of data needed for EMA
// calculation.
func (e EMA) Count() int {
	return e.Length()*2 - 1
}

// UnmarshalJSON parses JSON into EMA structure.
func (e *EMA) UnmarshalJSON(d []byte) error {
	var i struct {
		SMA struct {
			L int `json:"length"`
		} `json:"sma"`
	}

	if err := json.Unmarshal(d, &i); err != nil {
		return err
	}

	s, err := NewSMA(i.SMA.L)
	if err != nil {
		return err
	}

	e.sma = s

	if err := e.validate(); err != nil {
		return err
	}

	return nil
}

// MarshalJSON converts EMA configuration data into JSON.
func (e EMA) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		S SMA `json:"sma"`
	}{
		S: e.sma,
	})
}

// namedMarshalJSON converts EMA configuration data with its
// name into JSON.
func (e EMA) namedMarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		N String `json:"name"`
		S SMA    `json:"sma"`
	}{
		N: NameEMA,
		S: e.sma,
	})
}

// NameHMA returns HMA indicator name.
const NameHMA = "hma"

// HMA holds all the necessary information needed to calculate
// hull moving average.
// The zero value is not usable.
type HMA struct {
	// wma specifies the base moving average.
	wma WMA
}

// NewHMA validates provided configuration options and
// creates hull moving average indicator.
func NewHMA(w WMA) (HMA, error) {
	h := HMA{wma: w}

	if err := h.validate(); err != nil {
		return HMA{}, err
	}

	return h, nil
}

// WMA returns wma configuration option.
func (h HMA) WMA() WMA {
	return h.wma
}

// validate checks whether HMA was configured properly or not.
func (h HMA) validate() error {
	if h.wma == (WMA{}) {
		return errors.New("invalid wma")
	}

	if h.wma.length < 1 {
		return ErrInvalidLength
	}

	return nil
}

// Calc calculates HMA from the provided data slice.
func (h HMA) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, h.Count())
	if err != nil {
		return decimal.Zero, err
	}

	l := int(math.Sqrt(float64(h.wma.Count())))

	w1 := WMA{length: h.wma.Count() / 2}
	w2 := h.wma
	w3 := WMA{length: l}

	v := make([]decimal.Decimal, l)

	for i := 0; i < l; i++ {
		r1, _ := w1.Calc(dd[:len(dd)-l+i+1])

		r2, _ := w2.Calc(dd[:len(dd)-l+i+1])

		v[i] = r1.Mul(decimal.NewFromInt(2)).Sub(r2)
	}

	r, _ := w3.Calc(v)

	return r, nil
}

// Count determines the total amount of data needed for HMA
// calculation.
func (h HMA) Count() int {
	return h.wma.Count()*2 - 1
}

// UnmarshalJSON parses JSON into HMA structure.
func (h *HMA) UnmarshalJSON(d []byte) error {
	var i struct {
		WMA struct {
			L int `json:"length"`
		} `json:"wma"`
	}

	if err := json.Unmarshal(d, &i); err != nil {
		return err
	}

	w, err := NewWMA(i.WMA.L)
	if err != nil {
		return err
	}

	h.wma = w

	return nil
}

// MarshalJSON converts HMA configuration data into JSON.
func (h HMA) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		W WMA `json:"wma"`
	}{
		W: h.wma,
	})
}

// namedMarshalJSON converts HMA configuration data with its
// name into JSON.
func (h HMA) namedMarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		N String `json:"name"`
		W WMA    `json:"wma"`
	}{
		N: NameHMA,
		W: h.wma,
	})
}

// NameMACD returns MACD indicator name.
const NameMACD = "macd"

// MACD holds all the necessary information needed to calculate
// difference between two source indicators.
// The zero value is not usable.
type MACD struct {
	// source1 specifies the first base indicator.
	source1 Indicator

	// source2 specifies the second base indicator.
	source2 Indicator
}

// NewMACD validates provided configuration options and
// creates MACD indicator.
func NewMACD(source1, source2 Indicator) (MACD, error) {
	m := MACD{source1: source1, source2: source2}

	if err := m.validate(); err != nil {
		return MACD{}, err
	}

	return m, nil
}

// Sub1 returns source1 configuration option.
func (m MACD) Sub1() Indicator {
	return m.source1
}

// Sub2 returns source2 configuration option.
func (m MACD) Sub2() Indicator {
	return m.source2
}

// validate checks whether MACD was configured properly or not.
func (m MACD) validate() error {
	if m.source1 == nil || m.source2 == nil {
		return ErrInvalidSource
	}

	return nil
}

// Calc calculates MACD from the provided data slice.
func (m MACD) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, m.Count())
	if err != nil {
		return decimal.Zero, err
	}

	r1, err := m.source1.Calc(dd)
	if err != nil {
		return decimal.Zero, err
	}

	r2, err := m.source2.Calc(dd)
	if err != nil {
		return decimal.Zero, err
	}

	r := r1.Sub(r2)

	return r, nil
}

// Count determines the total amount of data needed for MACD
// calculation.
func (m MACD) Count() int {
	c1 := m.source1.Count()
	c2 := m.source2.Count()

	if c1 > c2 {
		return c1
	}

	return c2
}

// UnmarshalJSON parses JSON into MACD structure.
func (m *MACD) UnmarshalJSON(d []byte) error {
	var i struct {
		S1 json.RawMessage `json:"source1"`
		S2 json.RawMessage `json:"source2"`
	}

	if err := json.Unmarshal(d, &i); err != nil {
		return err
	}

	s1, err := fromJSON(i.S1)
	if err != nil {
		return err
	}

	m.source1 = s1

	s2, err := fromJSON(i.S2)
	if err != nil {
		return err
	}

	m.source2 = s2

	return nil
}

// MarshalJSON converts MACD configuration data into JSON.
func (m MACD) MarshalJSON() ([]byte, error) {
	s1, err := m.source1.namedMarshalJSON()
	if err != nil {
		return nil, err
	}

	s2, err := m.source2.namedMarshalJSON()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		S1 json.RawMessage `json:"source1"`
		S2 json.RawMessage `json:"source2"`
	}{
		S1: s1, S2: s2,
	})
}

// namedMarshalJSON converts MACD configuration data with its
// name into JSON.
func (m MACD) namedMarshalJSON() ([]byte, error) {
	s1, err := m.source1.namedMarshalJSON()
	if err != nil {
		return nil, err
	}

	s2, err := m.source2.namedMarshalJSON()
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		N  String          `json:"name"`
		S1 json.RawMessage `json:"source1"`
		S2 json.RawMessage `json:"source2"`
	}{
		N:  NameMACD,
		S1: s1,
		S2: s2,
	})
}

// NameROC returns ROC indicator name.
const NameROC = "roc"

// ROC holds all the necessary information needed to calculate rate
// of change.
// The zero value is not usable.
type ROC struct {
	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewROC validates provided configuration options and
// creates rate of change indicator.
func NewROC(length int) (ROC, error) {
	r := ROC{length: length}

	if err := r.validate(); err != nil {
		return ROC{}, err
	}

	return r, nil
}

// Length returns length configuration option.
func (r ROC) Length() int {
	return r.length
}

// validate checks whether ROC was configured properly or not.
func (r ROC) validate() error {
	if r.length < 1 {
		return ErrInvalidLength
	}

	return nil
}

// Calc calculates ROC from the provided data slice.
func (r ROC) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, r.Count())
	if err != nil {
		return decimal.Zero, err
	}

	n := dd[len(dd)-1]
	l := dd[len(dd)-r.Count()]

	if l.Equal(decimal.Zero) {
		return decimal.Zero, nil
	}

	return n.Sub(l).Div(l).Mul(decimal.NewFromInt(100)), nil
}

// Count determines the total amount of data needed for ROC
// calculation.
func (r ROC) Count() int {
	return r.length
}

// UnmarshalJSON parses JSON into ROC structure.
func (r *ROC) UnmarshalJSON(d []byte) error {
	var i struct {
		L int `json:"length"`
	}

	if err := json.Unmarshal(d, &i); err != nil {
		return err
	}

	r.length = i.L

	if err := r.validate(); err != nil {
		return err
	}

	return nil
}

// MarshalJSON converts ROC configuration data into JSON.
func (r ROC) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		L int `json:"length"`
	}{
		L: r.length,
	})
}

// namedMarshalJSON converts ROC configuration data with its
// name into JSON.
func (r ROC) namedMarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		N String `json:"name"`
		L int    `json:"length"`
	}{
		N: NameROC,
		L: r.length,
	})
}

// NameRSI returns RSI indicator name.
const NameRSI = "rsi"

// RSI holds all the necessary information needed to calculate relative
// strength index.
// The zero value is not usable.
type RSI struct {
	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewRSI validates provided configuration options and
// creates relative strength index indicator.
func NewRSI(length int) (RSI, error) {
	r := RSI{length: length}

	if err := r.validate(); err != nil {
		return RSI{}, err
	}

	return r, nil
}

// Length returns length configuration option.
func (r RSI) Length() int {
	return r.length
}

// validate checks whether RSI was configured properly or not.
func (r RSI) validate() error {
	if r.length < 1 {
		return ErrInvalidLength
	}

	return nil
}

// Calc calculates RSI from the provided data slice.
func (r RSI) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, r.Count())
	if err != nil {
		return decimal.Zero, err
	}

	ag := decimal.Zero
	al := decimal.Zero
	length := decimal.NewFromInt(int64(r.length))

	for i := 1; i < len(dd); i++ {
		if dd[i].Sub(dd[i-1]).LessThan(decimal.Zero) {
			al = al.Add(dd[i].Sub(dd[i-1]).Abs())
		} else {
			ag = ag.Add(dd[i].Sub(dd[i-1]))
		}
	}

	if ag == decimal.Zero {
		return decimal.NewFromInt(0), nil
	}

	if al == decimal.Zero {
		return decimal.NewFromInt(100), nil
	}

	ag = ag.Div(length)

	al = al.Div(length)

	return decimal.NewFromInt(100).Sub(decimal.NewFromInt(100).
		Div(decimal.NewFromInt(1).Add(ag.Div(al)))), nil
}

// Count determines the total amount of data needed for RSI
// calculation.
func (r RSI) Count() int {
	return r.length
}

// UnmarshalJSON parses JSON into RSI structure.
func (r *RSI) UnmarshalJSON(d []byte) error {
	var i struct {
		L int `json:"length"`
	}

	if err := json.Unmarshal(d, &i); err != nil {
		return err
	}

	r.length = i.L

	if err := r.validate(); err != nil {
		return err
	}

	return nil
}

// MarshalJSON converts RSI configuration data into JSON.
func (r RSI) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		L int `json:"length"`
	}{
		L: r.length,
	})
}

// namedMarshalJSON converts RSI configuration data with its
// name into JSON.
func (r RSI) namedMarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		N String `json:"name"`
		L int    `json:"length"`
	}{
		N: NameRSI,
		L: r.length,
	})
}

// NameSMA returns SMA indicator name.
const NameSMA = "sma"

// SMA holds all the necessary information needed to calculate simple
// moving average.
// The zero value is not usable.
type SMA struct {
	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewSMA validates provided configuration options and
// creates simple moving average indicator.
func NewSMA(length int) (SMA, error) {
	s := SMA{length: length}

	if err := s.validate(); err != nil {
		return SMA{}, err
	}

	return s, nil
}

// Length returns length configuration option.
func (s SMA) Length() int {
	return s.length
}

// validate checks whether SMA was configured properly or not.
func (s SMA) validate() error {
	if s.length < 1 {
		return ErrInvalidLength
	}

	return nil
}

// Calc calculates SMA from the provided data slice.
func (s SMA) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, s.Count())
	if err != nil {
		return decimal.Zero, err
	}

	r := decimal.Zero

	for i := 0; i < len(dd); i++ {
		r = r.Add(dd[i])
	}

	return r.Div(decimal.NewFromInt(int64(s.length))), nil
}

// Count determines the total amount of data needed for SMA
// calculation.
func (s SMA) Count() int {
	return s.length
}

// UnmarshalJSON parses JSON into SMA structure.
func (s *SMA) UnmarshalJSON(d []byte) error {
	var i struct {
		L int `json:"length"`
	}

	if err := json.Unmarshal(d, &i); err != nil {
		return err
	}

	s.length = i.L

	if err := s.validate(); err != nil {
		return err
	}

	return nil
}

// MarshalJSON converts SMA configuration data into JSON.
func (s SMA) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		L int `json:"length"`
	}{
		L: s.length,
	})
}

// namedMarshalJSON converts SMA configuration data with its
// name into JSON.
func (s SMA) namedMarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		N String `json:"name"`
		L int    `json:"length"`
	}{
		N: NameSMA,
		L: s.length,
	})
}

// NameSRSI returns SRSI indicator name.
const NameSRSI = "srsi"

// SRSI holds all the necessary information needed to calculate stoch
// relative strength index.
// The zero value is not usable.
type SRSI struct {
	// rsi specifies the base relative strength index.
	rsi RSI
}

// NewSRSI validates provided configuration options and
// creates stochastic relative strength index indicator.
func NewSRSI(r RSI) (SRSI, error) {
	s := SRSI{rsi: r}

	if err := s.validate(); err != nil {
		return SRSI{}, err
	}

	return s, nil
}

// RSI returns rsi configuration option.
func (s SRSI) RSI() RSI {
	return s.rsi
}

// validate checks whether SRSI was configured properly or not.
func (s SRSI) validate() error {
	if s.rsi == (RSI{}) {
		return errors.New("invalid rsi")
	}

	if s.rsi.length < 1 {
		return ErrInvalidLength
	}

	return nil
}

// Calc calculates SRSI from the provided data slice.
func (s SRSI) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	v, err := calcMultiple(dd, s.rsi.length, s.rsi)
	if err != nil {
		return decimal.Zero, err
	}

	c := v[0]
	h := v[0]
	l := v[0]

	for i := 1; i < len(v); i++ {
		if h.LessThan(v[i]) {
			h = v[i]
		}

		if l.GreaterThan(v[i]) {
			l = v[i]
		}
	}

	demin := h.Sub(l)
	if demin.Equal(decimal.Zero) {
		return decimal.Zero, nil
	}

	return c.Sub(l).Div(demin), nil
}

// Count determines the total amount of data needed for SRSI
// calculation.
func (s SRSI) Count() int {
	return s.rsi.Count()*2 - 1
}

// UnmarshalJSON parses JSON into SRSI structure.
func (s *SRSI) UnmarshalJSON(d []byte) error {
	var i struct {
		RSI struct {
			L int `json:"length"`
		} `json:"rsi"`
	}

	if err := json.Unmarshal(d, &i); err != nil {
		return err
	}

	r, err := NewRSI(i.RSI.L)
	if err != nil {
		return err
	}

	s.rsi = r

	return nil
}

// MarshalJSON converts SRSI configuration data into JSON.
func (s SRSI) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		R RSI `json:"rsi"`
	}{
		R: s.rsi,
	})
}

// namedMarshalJSON converts SRSI configuration data with its
// name into JSON.
func (s SRSI) namedMarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		N String `json:"name"`
		R RSI    `json:"rsi"`
	}{
		N: NameSRSI,
		R: s.rsi,
	})
}

// NameStoch returns Stoch  indicator name.
const NameStoch = "stoch"

// Stoch holds all the necessary information needed to calculate stochastic
// oscillator.
// The zero value is not usable.
type Stoch struct {
	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewStoch validates provided configuration options and
// creates stochastic indicator.
func NewStoch(length int) (Stoch, error) {
	s := Stoch{length: length}

	if err := s.validate(); err != nil {
		return Stoch{}, err
	}

	return s, nil
}

// Length returns length configuration option.
func (s Stoch) Length() int {
	return s.length
}

// validate checks whether Stoch was configured properly or not.
func (s Stoch) validate() error {
	if s.length < 1 {
		return ErrInvalidLength
	}

	return nil
}

// Calc calculates Stoch from the provided data slice.
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

	demin := h.Sub(l)
	if demin.Equal(decimal.Zero) {
		return decimal.Zero, nil
	}

	return dd[len(dd)-1].Sub(l).Div(demin).Mul(decimal.NewFromInt(100)), nil
}

// Count determines the total amount of data needed for Stoch
// calculation.
func (s Stoch) Count() int {
	return s.length
}

// UnmarshalJSON parses JSON into Stoch structure.
func (s *Stoch) UnmarshalJSON(d []byte) error {
	var i struct {
		L int `json:"length"`
	}

	if err := json.Unmarshal(d, &i); err != nil {
		return err
	}

	s.length = i.L

	if err := s.validate(); err != nil {
		return err
	}

	return nil
}

// MarshalJSON converts Stoch configuration data into JSON.
func (s Stoch) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		L int `json:"length"`
	}{
		L: s.length,
	})
}

// namedMarshalJSON converts Stoch configuration data with its
// name into JSON.
func (s Stoch) namedMarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		N String `json:"name"`
		L int    `json:"length"`
	}{
		N: NameStoch,
		L: s.length,
	})
}

// NameWMA returns WMA  indicator name.
const NameWMA = "wma"

// WMA holds all the necessary information needed to calculate weighted
// moving average.
// The zero value is not usable.
type WMA struct {
	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewWMA validates provided configuration options and
// creates weighted moving average indicator.
func NewWMA(length int) (WMA, error) {
	w := WMA{length: length}

	if err := w.validate(); err != nil {
		return WMA{}, err
	}

	return w, nil
}

// Length returns length configuration option.
func (w WMA) Length() int {
	return w.length
}

// validate checks whether WMA was configured properly or not.
func (w WMA) validate() error {
	if w.length < 1 {
		return ErrInvalidLength
	}

	return nil
}

// Calc calculates WMA from the provided data slice.
func (w WMA) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, w.Count())
	if err != nil {
		return decimal.Zero, err
	}

	r := decimal.Zero

	wi := decimal.NewFromFloat(float64(w.length*(w.length+1)) / 2.0)

	for i := 0; i < len(dd); i++ {
		r = r.Add(dd[i].Mul(decimal.NewFromInt(int64(i + 1)).Div(wi)))
	}

	return r, nil
}

// Count determines the total amount of data needed for WMA
// calculation.
func (w WMA) Count() int {
	return w.length
}

// UnmarshalJSON parses JSON into WMA structure.
func (w *WMA) UnmarshalJSON(d []byte) error {
	var i struct {
		L int `json:"length"`
	}

	if err := json.Unmarshal(d, &i); err != nil {
		return err
	}

	w.length = i.L

	if err := w.validate(); err != nil {
		return err
	}

	return nil
}

// MarshalJSON converts WMA configuration data into JSON.
func (w WMA) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		L int `json:"length"`
	}{
		L: w.length,
	})
}

// namedMarshalJSON converts WMA configuration data with its
// name into JSON.
func (w WMA) namedMarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		N String `json:"name"`
		L int    `json:"length"`
	}{
		N: NameWMA,
		L: w.length,
	})
}
