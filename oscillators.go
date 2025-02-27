// Package tango provides types and functions to calculate values of various
// market indicators.
package tango

import "github.com/shopspring/decimal"

// Aroon holds all the necessary information needed to calculate Aroon.
// The zero value is not usable.
type Aroon struct {
	// valid specifies whether Aroon paremeters were validated.
	valid bool

	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewAroon validates provided configuration options and
// creates new Aroon indicator instance.
func NewAroon(length int) (Aroon, error) {
	aroon := Aroon{
		length: length,
	}

	if err := aroon.validate(); err != nil {
		return Aroon{}, err
	}

	return aroon, nil
}

// validate checks whether the indicator has valid configuration properties.
func (aroon *Aroon) validate() error {
	if aroon.length < 1 {
		return ErrInvalidLength
	}

	aroon.valid = true

	return nil
}

// Calc calculates both Aroon trends from the provided data points slice.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/a/aroon.asp.
// All credits are due to Tushar Chande who developed Aroon indicator.
func (aroon Aroon) Calc(dd []decimal.Decimal) (
	uptrend decimal.Decimal,
	downtrend decimal.Decimal,
	err error,
) {

	if !aroon.valid {
		return decimal.Zero, decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != aroon.Count() {
		return decimal.Zero, decimal.Zero, ErrInvalidDataSize
	}

	min := dd[len(dd)-1]
	minIndex := decimal.NewFromInt(0)
	foundMin := false

	max := dd[len(dd)-1]
	maxIndex := decimal.NewFromInt(0)
	foundMax := false

	for i := len(dd) - 2; i >= 0 && (!foundMin || !foundMax); i-- {
		if !foundMin && min.GreaterThan(dd[i]) {
			min = dd[i]
			minIndex = decimal.NewFromInt(int64(aroon.length - i))
		} else if !min.Equal(dd[i]) {
			foundMin = true
		}

		if !foundMax && max.LessThan(dd[i]) {
			max = dd[i]
			maxIndex = decimal.NewFromInt(int64(aroon.length - i))
		} else if !max.Equal(dd[i]) {
			foundMax = true
		}
	}

	return aroon.calc(maxIndex), aroon.calc(minIndex), nil
}

// CalcTrend calculates specified Aroon trend from the provided data points slice.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/a/aroon.asp.
// All credits are due to Tushar Chande who developed Aroon indicator.
func (aroon Aroon) CalcTrend(dd []decimal.Decimal, trend Trend) (decimal.Decimal, error) {
	if err := trend.Validate(); err != nil {
		return decimal.Zero, err
	}

	uptrend, downtrend, err := aroon.Calc(dd)
	if err != nil {
		return decimal.Zero, err
	}

	if trend == TrendDown {
		return downtrend, nil
	}

	return uptrend, nil
}

func (aroon Aroon) calc(index decimal.Decimal) decimal.Decimal {
	return decimal.NewFromInt(int64(aroon.length)).Sub(index).
		Mul(_hundred).Div(decimal.NewFromInt(int64(aroon.length)))
}

// Count determines the total amount of data points needed for Aroon
// calculation.
func (aroon Aroon) Count() int {
	return aroon.length + 1
}

// CCI holds all the necessary information needed to calculate commodity
// channel index.
// The zero value is not usable.
type CCI struct {
	// valid specifies whether CCI paremeters were validated.
	valid bool

	// ma specifies moving average indicator configuration.
	ma MA
}

// NewCCI validates provided configuration options and creates
// new CCI indicator.
// If provided factor is zero, default value is going to be used (0.015f).
func NewCCI(mat MAType, length int) (CCI, error) {
	ma, err := NewMA(mat, length)
	if err != nil {
		return CCI{}, err
	}

	cci := CCI{
		ma:    ma,
		valid: true,
	}

	return cci, nil
}

// Calc calculates CCI from the provided data points slice.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/c/commoditychannelindex.asp.
// All credits are due to Donald Lambert who developed CCI indicator.
func (cci CCI) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	if !cci.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != cci.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	res, err := cci.ma.Calc(dd)
	if err != nil {
		return decimal.Zero, err
	}

	dnm := decimal.RequireFromString("0.015").Mul(MeanDeviation(dd))

	if dnm.Equal(decimal.Zero) {
		return decimal.Zero, nil
	}

	return dd[len(dd)-1].Sub(res).Div(dnm), nil
}

// Count determines the total amount of data points needed for CCI
// calculation.
func (cci CCI) Count() int {
	return cci.ma.Count()
}

// ROC holds all the necessary information needed to calculate rate
// of change.
// The zero value is not usable.
type ROC struct {
	// valid specifies whether ROC paremeters were validated.
	valid bool

	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewROC validates provided configuration options and
// creates new ROC indicator.
func NewROC(length int) (ROC, error) {
	roc := ROC{length: length}

	if err := roc.validate(); err != nil {
		return ROC{}, err
	}

	return roc, nil
}

// validate checks whether the indicator has valid configuration properties.
func (roc *ROC) validate() error {
	if roc.length < 1 {
		return ErrInvalidLength
	}

	roc.valid = true

	return nil
}

// Calc calculates ROC from the provided data points slice.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/p/pricerateofchange.asp.
func (roc ROC) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	if !roc.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != roc.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	curr := dd[0]
	last := dd[len(dd)-1]

	return curr.Div(last).Sub(_one).Mul(_hundred), nil
}

// Count determines the total amount of data points needed for ROC
// calculation.
func (roc ROC) Count() int {
	return roc.length
}

// RSI holds all the necessary information needed to calculate relative
// strength index.
// The zero value is not usable.
type RSI struct {
	// valid specifies whether RSI paremeters were validated.
	valid bool

	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewRSI validates provided configuration options and
// creates new RSI indicator.
func NewRSI(length int) (RSI, error) {
	rsi := RSI{
		length: length,
	}

	if err := rsi.validate(); err != nil {
		return RSI{}, err
	}

	return rsi, nil
}

// validate checks whether the indicator has valid configuration properties.
func (rsi *RSI) validate() error {
	if rsi.length < 1 {
		return ErrInvalidLength
	}

	rsi.valid = true

	return nil
}

// Calc calculates RSI from the provided data points slice.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/r/rsi.asp.
// All credits are due to J. Welles Wilder Jr. who developed RSI indicator.
func (rsi RSI) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	if !rsi.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != rsi.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	ag := decimal.Zero
	al := decimal.Zero
	length := decimal.NewFromInt(int64(rsi.length))

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
		return _hundred, nil
	}

	ag = ag.Div(length)

	al = al.Div(length)

	return _hundred.Sub(_hundred.Div(decimal.NewFromInt(1).Add(ag.Div(al)))), nil
}

// Count determines the total amount of data points needed for RSI
// calculation.
func (rsi RSI) Count() int {
	return rsi.length
}

// StochRSI holds all the necessary information needed to calculate stoch
// relative strength index.
// The zero value is not usable.
type StochRSI struct {
	// valid specifies whether StochRSI paremeters were validated.
	valid bool

	// rsi specifies the base relative strength index.
	rsi RSI
}

// NewStochRSI validates provided configuration options and
// creates new StochRSI indicator.
func NewStochRSI(length int) (StochRSI, error) {
	rsi, err := NewRSI(length)
	if err != nil {
		return StochRSI{}, err
	}

	return StochRSI{
		valid: true,
		rsi:   rsi,
	}, nil
}

// Calc calculates StochRSI from the provided data slice.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/s/stochrsi.asp.
func (s StochRSI) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	if !s.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != s.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	res := make([]decimal.Decimal, s.rsi.length)

	var err error
	for i := 0; i < s.rsi.length; i++ {
		res[i], err = s.rsi.Calc(dd[i : s.rsi.length+i])
		if err != nil {
			// unlikely to happen
			return decimal.Zero, err
		}
	}

	curr := res[0]
	max := res[0]
	min := res[0]

	for i := 1; i < len(res); i++ {
		if max.LessThan(res[i]) {
			max = res[i]
		}

		if min.GreaterThan(res[i]) {
			min = res[i]
		}
	}

	if max.Equal(min) {
		return decimal.Zero, nil
	}

	return curr.Sub(min).Div(max.Sub(min)), nil
}

// Count determines the total amount of data needed for StochRSI
// calculation.
func (s StochRSI) Count() int {
	return s.rsi.length*2 - 1
}

// Stoch holds all the necessary information needed to calculate stochastic
// oscillator.
// The zero value is not usable.
type Stoch struct {
	// valid specifies whether Stoch paremeters were validated.
	valid bool

	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewStoch validates provided configuration options and
// creates new Stoch indicator.
func NewStoch(length int) (Stoch, error) {
	stoch := Stoch{
		length: length,
	}

	if err := stoch.validate(); err != nil {
		return Stoch{}, err
	}

	return stoch, nil
}

// validate checks whether the indicator has valid configuration properties.
func (stoch *Stoch) validate() error {
	if stoch.length < 1 {
		return ErrInvalidLength
	}

	stoch.valid = true

	return nil
}

// Calc calculates Stoch from the provided data points slice.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/s/stochasticoscillator.asp.
func (stoch Stoch) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	if !stoch.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != stoch.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	low := dd[0]
	high := dd[0]

	for i := 0; i < len(dd); i++ {
		if dd[i].LessThan(low) {
			low = dd[i]
		}

		if dd[i].GreaterThan(high) {
			high = dd[i]
		}
	}

	dnm := high.Sub(low)
	if dnm.Equal(decimal.Zero) {
		return decimal.Zero, nil
	}

	return dd[len(dd)-1].Sub(low).Div(dnm).Mul(_hundred), nil
}

// Count determines the total amount of data points needed for Stoch
// calculation.
func (stoch Stoch) Count() int {
	return stoch.length
}
