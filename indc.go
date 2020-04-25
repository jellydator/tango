package indc

import (
	"encoding/json"
	"errors"
	"math"

	"github.com/shopspring/decimal"
)

// Names of all available indicators.
const (
	NameAroon = "aroon"
	NameCCI   = "cci"
	NameDEMA  = "dema"
	NameEMA   = "ema"
	NameMACD  = "macd"
	NameHMA   = "hma"
	NameROC   = "roc"
	NameRSI   = "rsi"
	NameSMA   = "sma"
	NameSRSI  = "srsi"
	NameStoch = "stoch"
	NameWMA   = "wma"
)

// Indicator is an interface that every indicator should implement.
type Indicator interface {
	// Calc should calculate and return indicator's value.
	Calc(dd []decimal.Decimal) (decimal.Decimal, error)

	// Count determines the total amount of data points needed
	// for indicator's calculation.
	Count() int

	// namedMarshalJSON helps indicators using other indicators to identify
	// their names.
	namedMarshalJSON() ([]byte, error)
}

// Aroon holds all the neccesary information needed to calculate aroon.
type Aroon struct {
	// trend configures which aroon trend to measure (it can either
	// be up or down).
	trend String

	// length specifies how many data points should be used
	// in calculations.
	length int
}

// NewAroon verifies provided parameters and
// creates aroon indicator.
func NewAroon(trend String, length int) (Aroon, error) {
	a := Aroon{trend: trend, length: length}

	if err := a.validate(); err != nil {
		return Aroon{}, err
	}

	return a, nil
}

// validate checks all Aroon parameters stored in func receiver to
// make sure that they're matching their requirements.
func (a Aroon) validate() error {
	if a.trend != "down" && a.trend != "up" {
		return errors.New("invalid trend")
	}

	if a.length < 1 {
		return ErrInvalidLength
	}

	return nil
}

// Calc calculates Aroon by using parameters stored in the func receiver.
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

// Count determines the total amount of data points needed for Aroon
// calculation by using parameters stored in the receiver.
func (a Aroon) Count() int {
	return a.length
}

// UnmarshalJSON parse JSON into an indicator source.
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

// MarshalJSON converts Aroon data into byte array.
func (a Aroon) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		T String `json:"trend"`
		L int    `json:"length"`
	}{
		T: a.trend, L: a.length,
	})
}

// namedMarshalJSON converts Aroon data with its name into byte array.
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

// CCI holds all the neccesary information needed to calculate commodity
// channel index.
type CCI struct {
	// source configures what calculations to use when computing CCI.
	source Indicator
}

// NewCCI verifies provided parameters and
// creates commodity channel index indicator.
func NewCCI(source Indicator) (CCI, error) {
	c := CCI{source: source}

	if err := c.validate(); err != nil {
		return CCI{}, err
	}

	return c, nil
}

// validate checks all CCI parameters stored in func receiver to make sure that
// they're matching their requirements.
func (c CCI) validate() error {
	if c.source == nil {
		return ErrInvalidSource
	}

	return nil
}

// Calc calculates CCI by using parameters stored in the func receiver.
func (c CCI) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, c.Count())
	if err != nil {
		return decimal.Zero, err
	}

	m, err := c.source.Calc(dd)
	if err != nil {
		return decimal.Zero, err
	}

	return dd[len(dd)-1].Sub(m).Div(decimal.NewFromFloat(0.015).
		Mul(meanDeviation(dd))), nil
}

// Count determines the total amount of data points needed for CCI
// calculation by using parameters stored in the receiver.
func (c CCI) Count() int {
	return c.source.Count()
}

// UnmarshalJSON parse JSON into an indicator source.
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

// MarshalJSON converts CCI data into byte array.
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

// namedMarshalJSON converts CCI data with its name into byte array.
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

// DEMA holds all the neccesary information needed to calculate
// double exponential moving average.
type DEMA struct {
	// length specifies how many data points should be used
	// in calculations.
	length int
}

// NewDEMA verifies provided parameters and
// creates double exponential moving average indicator.
func NewDEMA(length int) (DEMA, error) {
	d := DEMA{length: length}

	if err := d.validate(); err != nil {
		return DEMA{}, err
	}

	return d, nil
}

// Validate checks all DEMA parameters stored in func receiver to
// make sure that they're matching their requirements.
func (dm DEMA) validate() error {
	if dm.length < 1 {
		return ErrInvalidLength
	}
	return nil
}

// Calc calculates DEMA by using parameters stored in the func receiver.
func (dm DEMA) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, dm.Count())
	if err != nil {
		return decimal.Zero, err
	}

	v := make([]decimal.Decimal, dm.length)

	s := SMA{length: dm.length}
	v[0], _ = s.Calc(dd[:dm.length])

	e := EMA{length: dm.length}

	for i := dm.length; i < len(dd); i++ {
		v[i-dm.length+1] = e.CalcNext(v[i-dm.length], dd[i])
	}

	r := v[0]

	for i := 0; i < len(v); i++ {
		r = e.CalcNext(r, v[i])
	}

	return r, nil
}

// Count determines the total amount of data points needed for DEMA
// calculation by using parameters stored in the receiver.
func (dm DEMA) Count() int {
	return dm.length*2 - 1
}

// UnmarshalJSON parse JSON into an indicator source.
func (dm *DEMA) UnmarshalJSON(d []byte) error {
	var i struct {
		L int `json:"length"`
	}

	if err := json.Unmarshal(d, &i); err != nil {
		return err
	}

	dm.length = i.L

	if err := dm.validate(); err != nil {
		return err
	}

	return nil
}

// MarshalJSON converts DEMA data into byte array.
func (dm DEMA) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		L int `json:"length"`
	}{
		L: dm.length,
	})
}

// namedMarshalJSON converts DEMA data with its name into byte array.
func (dm DEMA) namedMarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		N String `json:"name"`
		L int    `json:"length"`
	}{
		N: NameDEMA,
		L: dm.length,
	})
}

// EMA holds all the neccesary information needed to calculate exponential
// moving average.
type EMA struct {
	// length specifies how many data points should be used
	// in calculations.
	length int
}

// NewEMA verifies provided parameters and
// creates exponential moving average indicator.
func NewEMA(length int) (EMA, error) {
	e := EMA{length: length}

	if err := e.validate(); err != nil {
		return EMA{}, err
	}

	return e, nil
}

// Validate checks all EMA parameters stored in func receiver to make sure that
// they're matching their requirements.
func (e EMA) validate() error {
	if e.length < 1 {
		return ErrInvalidLength
	}
	return nil
}

// Calc calculates EMA by using parameters stored in the func receiver.
func (e EMA) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, e.Count())
	if err != nil {
		return decimal.Zero, err
	}

	s := SMA{length: e.length}
	r, _ := s.Calc(dd[:e.length])

	for i := e.length; i < len(dd); i++ {
		r = e.CalcNext(r, dd[i])
	}

	return r, nil
}

// CalcNext calculates sequential EMA by using previous ema.
func (e EMA) CalcNext(l, n decimal.Decimal) decimal.Decimal {
	m := e.multiplier()
	return n.Mul(m).Add(l.Mul(decimal.NewFromInt(1).Sub(m)))
}

// multiplier calculates EMA multiplier by using parameters stored
// in the func receiver.
func (e EMA) multiplier() decimal.Decimal {
	return decimal.NewFromFloat(2.0 / float64(e.length+1))
}

// Count determines the total amount of data points needed for EMA
// calculation by using parameters stored in the receiver.
func (e EMA) Count() int {
	return e.length*2 - 1
}

// UnmarshalJSON parse JSON into an indicator source.
func (e *EMA) UnmarshalJSON(d []byte) error {
	var i struct {
		L int `json:"length"`
	}

	if err := json.Unmarshal(d, &i); err != nil {
		return err
	}

	e.length = i.L

	if err := e.validate(); err != nil {
		return err
	}

	return nil
}

// MarshalJSON converts EMA data into byte array.
func (e EMA) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		L int `json:"length"`
	}{
		L: e.length,
	})
}

// namedMarshalJSON converts EMA data with its name into byte array.
func (e EMA) namedMarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		N String `json:"name"`
		L int    `json:"length"`
	}{
		N: NameEMA,
		L: e.length,
	})
}

// HMA holds all the neccesary information needed to calculate
// hull moving average.
type HMA struct {
	// wma configures base moving average.
	wma WMA
}

// NewHMA verifies provided parameters and
// creates hull moving average indicator.
func NewHMA(w WMA) (HMA, error) {
	h := HMA{wma: w}

	if err := h.validate(); err != nil {
		return HMA{}, err
	}

	return h, nil
}

// validate checks all HMA parameters stored in func receiver to make sure that
// they're matching their requirements.
func (h HMA) validate() error {
	if h.wma == (WMA{}) {
		return errors.New("invalid wma")
	}

	if h.wma.length < 1 {
		return ErrInvalidLength
	}

	return nil
}

// Calc calculates HMA by using parameters stored in the func receiver.
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

// Count determines the total amount of data points needed for HMA
// calculation by using parameters stored in the receiver.
func (h HMA) Count() int {
	return h.wma.Count()*2 - 1
}

// UnmarshalJSON parse JSON into an indicator source.
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

// MarshalJSON converts HMA data into byte array.
func (h HMA) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		W WMA `json:"wma"`
	}{
		W: h.wma,
	})
}

// namedMarshalJSON converts HMA data with its name into byte array.
func (h HMA) namedMarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		N   String `json:"name"`
		WMA WMA    `json:"wma"`
	}{
		N:   NameHMA,
		WMA: h.wma,
	})
}

// MACD holds all the neccesary information needed to calculate
// difference between two source indicators.
type MACD struct {
	// source1 configures what calculations to use when computing first
	// macd.
	source1 Indicator

	// source2 configures what calculations to use when computing second
	// macd.
	source2 Indicator
}

// NewMACD verifies provided parameters and
// creates MACD indicator.
func NewMACD(source1, source2 Indicator) (MACD, error) {
	m := MACD{source1: source1, source2: source2}

	if err := m.validate(); err != nil {
		return MACD{}, err
	}

	return m, nil
}

// validate checks all MACD parameters stored in func receiver
// to make sure that they're matching their requirements.
func (m MACD) validate() error {
	if m.source1 == nil || m.source2 == nil {
		return ErrInvalidSource
	}

	return nil
}

// Calc calculates MACD by using parameters stored in the func receiver.
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

// Count determines the total amount of data points needed for MACD
// calculation by using parameters stored in the receiver.
func (m MACD) Count() int {
	c1 := m.source1.Count()
	c2 := m.source2.Count()

	if c1 > c2 {
		return c1
	}

	return c2
}

// UnmarshalJSON parse JSON into an indicator source.
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

// MarshalJSON converts MACD data into byte array.
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

// namedMarshalJSON converts MACD data with its name into byte array.
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

// ROC holds all the neccesary information needed to calculate rate
// of change.
type ROC struct {
	// length specifies how many data points should be used
	// in calculations.
	length int
}

// NewROC verifies provided parameters and
// creates rate of change indicator.
func NewROC(length int) (ROC, error) {
	r := ROC{length: length}

	if err := r.validate(); err != nil {
		return ROC{}, err
	}

	return r, nil
}

// Validate checks all ROC parameters stored in func receiver to make sure that
// they're matching their requirements.
func (r ROC) validate() error {
	if r.length < 1 {
		return ErrInvalidLength
	}
	return nil
}

// Calc calculates ROC by using parameters stored in the func receiver.
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
// calculation by using parameters stored in the receiver.
func (r ROC) Count() int {
	return r.length
}

// UnmarshalJSON parse JSON into an indicator source.
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

// MarshalJSON converts ROC data into byte array.
func (r ROC) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		L int `json:"length"`
	}{
		L: r.length,
	})
}

// namedMarshalJSON converts ROC data with its name into byte array.
func (r ROC) namedMarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		N String `json:"name"`
		L int    `json:"length"`
	}{
		N: NameROC,
		L: r.length,
	})
}

// RSI holds all the neccesary information needed to calculate relative
// strength index.
type RSI struct {
	// length specifies how many data points should be used
	// in calculations.
	length int
}

// NewRSI verifies provided parameters and
// creates relative strength index indicator.
func NewRSI(length int) (RSI, error) {
	r := RSI{length: length}

	if err := r.validate(); err != nil {
		return RSI{}, err
	}

	return r, nil
}

// Validate checks all RSI parameters stored in func receiver to make sure that
// they're matching their requirements.
func (r RSI) validate() error {
	if r.length < 1 {
		return ErrInvalidLength
	}
	return nil
}

// Calc calculates RSI by using parameters stored in the func receiver.
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

	ag = ag.Div(decimal.NewFromInt(int64(r.length)))
	al = al.Div(decimal.NewFromInt(int64(r.length)))

	return decimal.NewFromInt(100).Sub(decimal.NewFromInt(100).
		Div(decimal.NewFromInt(1).Add(ag.Div(al)))), nil
}

// Count determines the total amount of data points needed for RSI
// calculation by using parameters stored in the receiver.
func (r RSI) Count() int {
	return r.length
}

// UnmarshalJSON parse JSON into an indicator source.
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

// MarshalJSON converts RSI data into byte array.
func (r RSI) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		L int `json:"length"`
	}{
		L: r.length,
	})
}

// namedMarshalJSON converts RSI data with its name into byte array.
func (r RSI) namedMarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		N String `json:"name"`
		L int    `json:"length"`
	}{
		N: NameRSI,
		L: r.length,
	})
}

// SMA holds all the neccesary information needed to calculate simple
// moving average.
type SMA struct {
	// length specifies how many data points should be used
	// in calculations.
	length int
}

// NewSMA verifies provided parameters and
// creates simple moving average indicator.
func NewSMA(length int) (SMA, error) {
	s := SMA{length: length}

	if err := s.validate(); err != nil {
		return SMA{}, err
	}

	return s, nil
}

// validate checks all SMA parameters stored in func receiver to make sure that
// they're matching their requirements.
func (s SMA) validate() error {
	if s.length < 1 {
		return ErrInvalidLength
	}
	return nil
}

// Calc calculates SMA by using parameters stored in the func receiver.
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

// Count determines the total amount of data points needed for SMA
// calculation by using parameters stored in the receiver.
func (s SMA) Count() int {
	return s.length
}

// UnmarshalJSON parse JSON into an indicator source.
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

// MarshalJSON converts SMA data into byte array.
func (s SMA) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		L int `json:"length"`
	}{
		L: s.length,
	})
}

// namedMarshalJSON converts SMA data with its name into byte array.
func (s SMA) namedMarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		N String `json:"name"`
		L int    `json:"length"`
	}{
		N: NameSMA,
		L: s.length,
	})
}

// SRSI holds all the neccesary information needed to calculate stoch
// relative strenght index.
type SRSI struct {
	// rsi configures relative strengh index.
	rsi RSI
}

// NewSRSI verifies provided parameters and
// creates hull moving average indicator.
func NewSRSI(r RSI) (SRSI, error) {
	s := SRSI{rsi: r}

	if err := s.validate(); err != nil {
		return SRSI{}, err
	}

	return s, nil
}

// validate checks all SRSI parameters stored in func receiver to make sure that
// they're matching their requirements.
func (s SRSI) validate() error {
	if s.rsi == (RSI{}) {
		return errors.New("invalid rsi")
	}

	if s.rsi.length < 1 {
		return ErrInvalidLength
	}

	return nil
}

// Calc calculates SRSI by using parameters stored in the func receiver.
func (s SRSI) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	dd, err := resize(dd, s.Count())
	if err != nil {
		return decimal.Zero, err
	}

	c, err := s.rsi.Calc(dd)
	if err != nil {
		return decimal.Zero, err
	}

	l := c
	h := c

	for i := 1; i < len(dd); i++ {
		r, err := s.rsi.Calc(dd[:s.rsi.Count()+i])
		if err != nil {
			return decimal.Zero, err
		}

		if r.LessThan(l) {
			l = dd[i]
		}

		if r.GreaterThan(h) {
			h = dd[i]
		}
	}

	return c.Sub(l).Div(h.Sub(l)), nil
}

// Count determines the total amount of data points needed for SRSI
// calculation by using parameters stored in the receiver.
func (s SRSI) Count() int {
	return s.rsi.Count()*2 - 1
}

// UnmarshalJSON parse JSON into an indicator source.
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

// MarshalJSON converts SRSI data into byte array.
func (s SRSI) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		R RSI `json:"rsi"`
	}{
		R: s.rsi,
	})
}

// namedMarshalJSON converts SRSI data with its name into byte array.
func (s SRSI) namedMarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		N String `json:"name"`
		R RSI    `json:"rsi"`
	}{
		N: NameHMA,
		R: s.rsi,
	})
}

// Stoch holds all the neccesary information needed to calculate stochastic
// oscillator.
type Stoch struct {
	// length specifies how many data points should be used
	// in calculations.
	length int
}

// NewStoch verifies provided parameters and
// creates stochastic indicator.
func NewStoch(length int) (Stoch, error) {
	s := Stoch{length: length}

	if err := s.validate(); err != nil {
		return Stoch{}, err
	}

	return s, nil
}

// Validate checks all stochastic parameters stored in func receiver to make
// sure that they're matching their requirements.
func (s Stoch) validate() error {
	if s.length < 1 {
		return ErrInvalidLength
	}
	return nil
}

// Calc calculates stochastic by using parameters stored in
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
// calculation by using parameters stored in the receiver.
func (s Stoch) Count() int {
	return s.length
}

// UnmarshalJSON parse JSON into an indicator source.
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

// MarshalJSON converts Stoch data into byte array.
func (s Stoch) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		L int `json:"length"`
	}{
		L: s.length,
	})
}

// namedMarshalJSON converts Stoch data with its name into byte array.
func (s Stoch) namedMarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		N String `json:"name"`
		L int    `json:"length"`
	}{
		N: NameStoch,
		L: s.length,
	})
}

// WMA holds all the neccesary information needed to calculate weighted
// moving average.
type WMA struct {
	// length specifies how many data points should be used
	// in calculations.
	length int
}

// NewWMA verifies provided parameters and
// creates weighted moving average indicator.
func NewWMA(length int) (WMA, error) {
	w := WMA{length: length}

	if err := w.validate(); err != nil {
		return WMA{}, err
	}

	return w, nil
}

// Validate checks all WMA parameters stored in func receiver to make sure that
// they're matching their requirements.
func (w WMA) validate() error {
	if w.length < 1 {
		return ErrInvalidLength
	}
	return nil
}

// Calc calculates WMA by using parameters stored in the func receiver.
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

// Count determines the total amount of data points needed for WMA
// calculation by using parameters stored in the receiver.
func (w WMA) Count() int {
	return w.length
}

// UnmarshalJSON parse JSON into an indicator source.
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

// MarshalJSON converts WMA data into byte array.
func (w WMA) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		L int `json:"length"`
	}{
		L: w.length,
	})
}

// namedMarshalJSON converts WMA data with its name into byte array.
func (w WMA) namedMarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		N String `json:"name"`
		L int    `json:"length"`
	}{
		N: NameWMA,
		L: w.length,
	})
}
