// Package indc provides types and functions to calculate values of various
// market indicators.
package indc

import (
	"math"

	"github.com/shopspring/decimal"
)

// BB holds all the necessary information needed to calculate Bollinger Bands.
// The zero value is not usable.
type BB struct {
	// valid specifies whether BB paremeters were validated.
	valid bool

	// band specifies which bollinger band to calculate.
	band Band

	// stdDev specifies how to adjust standard deviation.
	stdDev decimal.Decimal

	// sma specifies SMA indicator configuration.
	sma SMA
}

// NewBB validates provided configuration options and creates
// new BB indicator.
func NewBB(band Band, stdDev decimal.Decimal, length int) (BB, error) {
	sma, err := NewSMA(length)
	if err != nil {
		return BB{}, err
	}

	bb := BB{
		band:   band,
		stdDev: stdDev,
		sma:    sma,
	}

	if err := bb.validate(); err != nil {
		return BB{}, err
	}

	return bb, nil
}

// validate checks whether the indicator has valid configuration properties.
func (bb *BB) validate() error {
	if err := bb.band.Validate(); err != nil {
		return err
	}

	bb.valid = true

	return nil
}

// Calc calculates BB from the provided data points slice.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/b/bollingerbands.asp.
// All credits are due to John Bollinger who developed BB indicator.
func (bb BB) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	if !bb.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != bb.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	res, err := bb.sma.Calc(dd)
	if err != nil {
		// unlikely to happen
		return decimal.Zero, err
	}

	sdev := sdev(dd).Mul(bb.stdDev)

	switch bb.band {
	case BandUpper:
		return res.Add(sdev), nil
	case BandLower:
		return res.Sub(sdev), nil
	default: // BB is validated, only BandWidth is left.
		return res.Add(sdev).Sub(res.Sub(sdev)).Div(res).Mul(_hundred), nil
	}
}

// Count determines the total amount of data points needed for BB
// calculation.
func (bb BB) Count() int {
	return bb.sma.Count()
}

// DEMA holds all the necessary information needed to calculate
// double exponential moving average.
// The zero value is not usable.
type DEMA struct {
	// valid specifies whether DEMA paremeters were validated.
	valid bool

	// ema specifies what ema should be used for dema calculations.
	ema EMA
}

// NewDEMA validates provided configuration options and creates
// new DEMA indicator.
func NewDEMA(length int) (DEMA, error) {
	ema, err := NewEMA(length)
	if err != nil {
		return DEMA{}, err
	}

	return DEMA{
		valid: true,
		ema:   ema,
	}, nil
}

// Calc calculates DEMA from the provided data points slice.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/d/double-exponential-moving-average.asp.
// All credits are due to Patrick Mulloy who developed DEMA indicator.
func (dema DEMA) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	if !dema.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != dema.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	pres := make([]decimal.Decimal, dema.ema.sma.length)

	var err error

	pres[0], err = dema.ema.sma.Calc(dd[:dema.ema.sma.length])
	if err != nil {
		// unlikely to happen
		return decimal.Zero, err
	}

	for i := dema.ema.sma.length; i < len(dd); i++ {
		pres[i-dema.ema.sma.length+1], err = dema.ema.CalcNext(pres[i-dema.ema.sma.length], dd[i])
		if err != nil {
			// unlikely to happen
			return decimal.Zero, err
		}
	}

	res := pres[0]

	for i := 0; i < len(pres); i++ {
		res, err = dema.ema.CalcNext(res, pres[i])
		if err != nil {
			// unlikely to happen
			return decimal.Zero, err
		}
	}

	return res, nil
}

// Count determines the total amount of data points needed for DEMA
// calculation.
func (dema DEMA) Count() int {
	return dema.ema.Count()
}

// EMA holds all the necessary information needed to calculate exponential
// moving average.
// The zero value is not usable.
type EMA struct {
	// valid specifies whether DEMA paremeters were validated.
	valid bool

	// sma specifies what sma should be used for ema calculations.
	sma SMA
}

// NewEMA validates provided configuration options and
// creates new EMA indicator.
func NewEMA(length int) (EMA, error) {
	sma, err := NewSMA(length)
	if err != nil {
		return EMA{}, err
	}

	return EMA{
		valid: true,
		sma:   sma,
	}, nil
}

// Calc calculates EMA from the provided data points slice.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/e/ema.asp.
func (ema EMA) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	if !ema.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != ema.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	res, err := ema.sma.Calc(dd[:ema.sma.length])
	if err != nil {
		// unlikely to happen
		return decimal.Zero, err
	}

	for i := ema.sma.length; i < len(dd); i++ {
		res, err = ema.CalcNext(res, dd[i])
		if err != nil {
			// unlikely to happen
			return decimal.Zero, err
		}
	}

	return res, nil
}

// CalcNext calculates sequential EMA by using previous EMA.
func (ema EMA) CalcNext(lres, dec decimal.Decimal) (decimal.Decimal, error) {
	if !ema.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	mtp := ema.multiplier()

	return dec.Mul(mtp).Add(lres.Mul(decimal.NewFromInt(1).Sub(mtp))), nil
}

// multiplier calculates EMA multiplier.
func (ema EMA) multiplier() decimal.Decimal {
	return decimal.NewFromInt(2).Div(decimal.NewFromInt(int64(ema.sma.length) + 1))
}

// Count determines the total amount of data points needed for EMA
// calculation.
func (ema EMA) Count() int {
	return ema.sma.length*2 - 1
}

// HMA holds all the necessary information needed to calculate
// hull moving average.
// The zero value is not usable.
type HMA struct {
	// valid specifies whether HMA paremeters were validated.
	valid bool

	// wma specifies the base moving average.
	wma WMA
}

// NewHMA validates provided configuration options and
// creates new HMA indicator.
func NewHMA(length int) (HMA, error) {
	wma, err := NewWMA(length)
	if err != nil {
		return HMA{}, err
	}

	return HMA{
		valid: true,
		wma:   wma,
	}, nil
}

// Calc calculates HMA from the provided data points slice.
// Calculation is based on formula provided by fidelity.
// https://www.fidelity.com/learning-center/trading-investing/technical-analysis/technical-indicator-guide/hull-moving-average.
// All credits are due to Alan Hull who developed HMA indicator.
func (h HMA) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	if !h.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != h.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	wma1 := WMA{length: h.wma.length / 2, valid: true}
	wma2 := WMA{length: int(math.Sqrt(float64(h.wma.length))), valid: true}

	res := make([]decimal.Decimal, wma2.length)

	for i := 0; i < wma2.length; i++ {
		res1, err := wma1.Calc(dd[i : wma1.length+i])
		if err != nil {
			// unlikely to happen
			return decimal.Zero, err
		}

		res2, err := h.wma.Calc(dd[i : h.wma.length+i])
		if err != nil {
			// unlikely to happen
			return decimal.Zero, err
		}

		res[i] = res1.Mul(decimal.NewFromInt(2)).Sub(res2)
	}

	return wma2.Calc(res)
}

// Count determines the total amount of data points needed for HMA
// calculation.
func (h HMA) Count() int {
	return int(math.Sqrt(float64(h.wma.length))) + h.wma.length - 1
}

// SMA holds all the necessary information needed to calculate simple
// moving average.
// The zero value is not usable.
type SMA struct {
	// valid specifies whether SMA paremeters were validated.
	valid bool

	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewSMA validates provided configuration options and
// creates new SMA indicator.
func NewSMA(length int) (SMA, error) {
	sma := SMA{
		length: length,
	}

	if err := sma.validate(); err != nil {
		return SMA{}, err
	}

	return sma, nil
}

// validate checks whether the indicator has valid configuration properties.
func (sma *SMA) validate() error {
	if sma.length < 1 {
		return ErrInvalidLength
	}

	sma.valid = true

	return nil
}

// Calc calculates SMA from the provided data points slice.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/s/sma.asp.
func (sma SMA) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	if !sma.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != sma.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	res := decimal.Zero
	for i := 0; i < len(dd); i++ {
		res = res.Add(dd[i])
	}

	return res.Div(decimal.NewFromInt(int64(sma.length))), nil
}

// Count determines the total amount of data points needed for SMA
// calculation.
func (sma SMA) Count() int {
	return sma.length
}

// WMA holds all the necessary information needed to calculate weighted
// moving average.
// The zero value is not usable.
type WMA struct {
	// valid specifies whether WMA paremeters were validated.
	valid bool

	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewWMA validates provided configuration options and
// creates new WMA indicator.
func NewWMA(length int) (WMA, error) {
	wma := WMA{
		length: length,
	}

	if err := wma.validate(); err != nil {
		return WMA{}, err
	}

	return wma, nil
}

// validate checks whether the indicator has valid configuration properties.
func (wma *WMA) validate() error {
	if wma.length < 1 {
		return ErrInvalidLength
	}

	wma.valid = true

	return nil
}

// Calc calculates WMA from the provided data points slice.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/articles/technical/060401.asp.
func (wma WMA) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	if !wma.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(dd) != wma.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	res := decimal.Zero

	weight := decimal.NewFromInt(int64(wma.length * (wma.length + 1))).Div(decimal.NewFromInt(2))

	for i := 0; i < len(dd); i++ {
		res = res.Add(dd[i].Mul(decimal.NewFromInt(int64(i + 1)).Div(weight)))
	}

	return res, nil
}

// Count determines the total amount of data points needed for WMA
// calculation.
func (wma WMA) Count() int {
	return wma.length
}
