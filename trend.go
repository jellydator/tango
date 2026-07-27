package tango

import "github.com/shopspring/decimal"

// ADX holds all the necessary information needed to calculate average
// directional index.
// The zero value is not usable.
type ADX struct {
	// valid specifies whether ADX parameters were validated.
	valid bool

	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewADX validates provided configuration options and
// creates new ADX indicator instance.
func NewADX(length int) (ADX, error) {
	adx := ADX{
		length: length,
	}

	if err := adx.validate(); err != nil {
		return ADX{}, err
	}

	return adx, nil
}

// validate checks whether the indicator has valid configuration properties.
func (adx *ADX) validate() error {
	if adx.length < 1 {
		return ErrInvalidLength
	}

	adx.valid = true

	return nil
}

// Calc calculates ADX along with both directional index values from the
// provided candles ordered from the oldest to the newest. The
// directional movements and true ranges are smoothed with the Wilder
// smoothing method, the first length directional index values seed a
// simple average which is then smoothed the same way.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/a/adx.asp.
// All credits are due to J. Welles Wilder Jr. who developed ADX indicator.
func (adx ADX) Calc(cc []Candle) (
	value decimal.Decimal,
	plusDI decimal.Decimal,
	minusDI decimal.Decimal,
	err error,
) {

	if !adx.valid {
		return decimal.Zero, decimal.Zero, decimal.Zero, ErrInvalidIndicator
	}

	if len(cc) != adx.Count() {
		return decimal.Zero, decimal.Zero, decimal.Zero, ErrInvalidDataSize
	}

	pdm, mdm := directionalMovements(cc)
	trs := trueRanges(cc)
	length := decimal.NewFromInt(int64(adx.length))

	spdm := decimal.Sum(pdm[0], pdm[1:adx.length]...)
	smdm := decimal.Sum(mdm[0], mdm[1:adx.length]...)
	str := decimal.Sum(trs[0], trs[1:adx.length]...)

	dxs := make([]decimal.Decimal, 0, 2*adx.length+1)

	plusDI, minusDI = directionalIndexes(spdm, smdm, str)
	dxs = append(dxs, directionalStrength(plusDI, minusDI))

	for i := adx.length; i < len(trs); i++ {
		spdm = spdm.Sub(spdm.Div(length)).Add(pdm[i])
		smdm = smdm.Sub(smdm.Div(length)).Add(mdm[i])
		str = str.Sub(str.Div(length)).Add(trs[i])

		plusDI, minusDI = directionalIndexes(spdm, smdm, str)
		dxs = append(dxs, directionalStrength(plusDI, minusDI))
	}

	value = Average(dxs[:adx.length])

	for _, dx := range dxs[adx.length:] {
		value = value.Mul(length.Sub(_one)).Add(dx).Div(length)
	}

	return value, plusDI, minusDI, nil
}

// Count determines the total amount of data points needed for ADX
// calculation.
func (adx ADX) Count() int {
	return adx.length*3 + 1
}

// directionalMovements calculates the positive and the negative
// directional movement of every candle that has a preceding candle. The
// candles must be ordered from the oldest to the newest.
func directionalMovements(cc []Candle) (pdm, mdm []decimal.Decimal) {
	pdm = make([]decimal.Decimal, 0, len(cc)-1)
	mdm = make([]decimal.Decimal, 0, len(cc)-1)

	for i := 1; i < len(cc); i++ {
		up := cc[i].High.Sub(cc[i-1].High)
		down := cc[i-1].Low.Sub(cc[i].Low)

		if up.GreaterThan(down) && up.IsPositive() {
			pdm = append(pdm, up)
		} else {
			pdm = append(pdm, decimal.Zero)
		}

		if down.GreaterThan(up) && down.IsPositive() {
			mdm = append(mdm, down)
		} else {
			mdm = append(mdm, decimal.Zero)
		}
	}

	return pdm, mdm
}

// directionalIndexes calculates both directional index values from the
// smoothed directional movements and the smoothed true range. Both
// values are zero when the smoothed true range is zero.
func directionalIndexes(spdm, smdm, str decimal.Decimal) (plusDI, minusDI decimal.Decimal) {
	if str.IsZero() {
		return decimal.Zero, decimal.Zero
	}

	return spdm.Mul(_hundred).Div(str), smdm.Mul(_hundred).Div(str)
}

// directionalStrength calculates the directional index strength from
// both directional index values. The strength is zero when both values
// are zero.
func directionalStrength(plusDI, minusDI decimal.Decimal) decimal.Decimal {
	dnm := plusDI.Add(minusDI)
	if dnm.IsZero() {
		return decimal.Zero
	}

	return plusDI.Sub(minusDI).Abs().Mul(_hundred).Div(dnm)
}

// SuperTrend holds all the necessary information needed to calculate
// the SuperTrend indicator.
// The zero value is not usable.
type SuperTrend struct {
	// valid specifies whether SuperTrend parameters were validated.
	valid bool

	// multiplier specifies by how much the average true range is
	// scaled before it is applied to the median price.
	multiplier decimal.Decimal

	// atr specifies the average true range used for the band offsets.
	atr ATR
}

// NewSuperTrend validates provided configuration options and
// creates new SuperTrend indicator instance.
func NewSuperTrend(multiplier decimal.Decimal, length int) (SuperTrend, error) {
	atr, err := NewATR(length)
	if err != nil {
		return SuperTrend{}, err
	}

	st := SuperTrend{
		multiplier: multiplier,
		atr:        atr,
	}

	if err := st.validate(); err != nil {
		return SuperTrend{}, err
	}

	return st, nil
}

// validate checks whether the indicator has valid configuration properties.
func (st *SuperTrend) validate() error {
	if st.multiplier.Cmp(decimal.Zero) <= 0 {
		return ErrInvalidMultiplier
	}

	st.valid = true

	return nil
}

// Calc calculates the SuperTrend value and its trend from the provided
// candles ordered from the oldest to the newest. The band tightening
// logic is iterated over the most recent length candles, starting in
// the uptrend state when the close of the first iterated candle is
// above its median price and in the downtrend state otherwise.
// Calculation is based on the widely used SuperTrend formula.
// All credits are due to Olivier Seban who developed the indicator.
func (st SuperTrend) Calc(cc []Candle) (value decimal.Decimal, trend Trend, err error) {
	if !st.valid {
		return decimal.Zero, 0, ErrInvalidIndicator
	}

	if len(cc) != st.Count() {
		return decimal.Zero, 0, ErrInvalidDataSize
	}

	var upper, lower decimal.Decimal

	trend = TrendDown

	for i := 0; i < st.atr.length; i++ {
		window := cc[i : i+st.atr.Count()]
		last := window[len(window)-1]

		atr, err := st.atr.Calc(window)
		if err != nil {
			// unlikely to happen
			return decimal.Zero, 0, err
		}

		median := last.High.Add(last.Low).Div(decimal.NewFromInt(2))
		offset := atr.Mul(st.multiplier)

		basicUpper := median.Add(offset)
		basicLower := median.Sub(offset)

		if i == 0 {
			upper = basicUpper
			lower = basicLower

			if last.Close.GreaterThan(median) {
				trend = TrendUp
			}

			continue
		}

		prevClose := window[len(window)-2].Close

		if basicUpper.LessThan(upper) || prevClose.GreaterThan(upper) {
			upper = basicUpper
		}

		if basicLower.GreaterThan(lower) || prevClose.LessThan(lower) {
			lower = basicLower
		}

		trend = superTrendShift(trend, last.Close, upper, lower)
	}

	if trend == TrendUp {
		return lower, trend, nil
	}

	return upper, trend, nil
}

// superTrendShift determines the SuperTrend trend state after the
// candle with the provided close price when the final bands are in
// effect.
func superTrendShift(trend Trend, closePrice, upper, lower decimal.Decimal) Trend {
	if trend == TrendUp && closePrice.LessThan(lower) {
		return TrendDown
	}

	if trend == TrendDown && closePrice.GreaterThan(upper) {
		return TrendUp
	}

	return trend
}

// Count determines the total amount of data points needed for
// SuperTrend calculation.
func (st SuperTrend) Count() int {
	return st.atr.Count() + st.atr.length - 1
}

// PSAR holds all the necessary information needed to calculate the
// parabolic stop and reverse indicator.
// The zero value is not usable.
type PSAR struct {
	// valid specifies whether PSAR parameters were validated.
	valid bool

	// step specifies by how much the acceleration factor grows with
	// every new extreme point.
	step decimal.Decimal

	// limit specifies the maximum value of the acceleration factor.
	limit decimal.Decimal

	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewPSAR validates provided configuration options and
// creates new PSAR indicator instance.
func NewPSAR(step, limit decimal.Decimal, length int) (PSAR, error) {
	psar := PSAR{
		step:   step,
		limit:  limit,
		length: length,
	}

	if err := psar.validate(); err != nil {
		return PSAR{}, err
	}

	return psar, nil
}

// validate checks whether the indicator has valid configuration properties.
func (psar *PSAR) validate() error {
	if psar.length < 2 {
		return ErrInvalidLength
	}

	if psar.step.Cmp(decimal.Zero) <= 0 || psar.step.GreaterThan(psar.limit) {
		return ErrInvalidStep
	}

	psar.valid = true

	return nil
}

// psarState holds the intermediate state of the PSAR calculation.
type psarState struct {
	// up specifies whether the current trend is an uptrend.
	up bool

	// sar specifies the current stop and reverse value.
	sar decimal.Decimal

	// ep specifies the current extreme point of the trend.
	ep decimal.Decimal

	// af specifies the current acceleration factor.
	af decimal.Decimal
}

// Calc calculates the PSAR value and its trend from the provided
// candles ordered from the oldest to the newest. The iteration starts
// in the uptrend state when the second candle closes at or above the
// close of the first candle and in the downtrend state otherwise.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/p/parabolicindicator.asp.
// All credits are due to J. Welles Wilder Jr. who developed PSAR indicator.
func (psar PSAR) Calc(cc []Candle) (value decimal.Decimal, trend Trend, err error) {
	if !psar.valid {
		return decimal.Zero, 0, ErrInvalidIndicator
	}

	if len(cc) != psar.Count() {
		return decimal.Zero, 0, ErrInvalidDataSize
	}

	st := psarState{
		up:  cc[1].Close.GreaterThanOrEqual(cc[0].Close),
		af:  psar.step,
		sar: cc[0].High,
		ep:  cc[1].Low,
	}

	if st.up {
		st.sar = cc[0].Low
		st.ep = cc[1].High
	}

	for i := 2; i < len(cc); i++ {
		st.sar = st.sar.Add(st.af.Mul(st.ep.Sub(st.sar)))

		if st.up {
			st.sar = decimal.Min(st.sar, cc[i-1].Low, cc[i-2].Low)
			st = psar.riseStep(st, cc[i])

			continue
		}

		st.sar = decimal.Max(st.sar, cc[i-1].High, cc[i-2].High)
		st = psar.fallStep(st, cc[i])
	}

	if st.up {
		return st.sar, TrendUp, nil
	}

	return st.sar, TrendDown, nil
}

// riseStep advances the PSAR state with the next candle while in an
// uptrend, reversing the trend when the candle low crosses the stop and
// reverse value.
func (psar PSAR) riseStep(st psarState, c Candle) psarState {
	if c.Low.LessThan(st.sar) {
		return psarState{
			up:  false,
			sar: st.ep,
			ep:  c.Low,
			af:  psar.step,
		}
	}

	if c.High.GreaterThan(st.ep) {
		st.ep = c.High
		st.af = decimal.Min(st.af.Add(psar.step), psar.limit)
	}

	return st
}

// fallStep advances the PSAR state with the next candle while in a
// downtrend, reversing the trend when the candle high crosses the stop
// and reverse value.
func (psar PSAR) fallStep(st psarState, c Candle) psarState {
	if c.High.GreaterThan(st.sar) {
		return psarState{
			up:  true,
			sar: st.ep,
			ep:  c.High,
			af:  psar.step,
		}
	}

	if c.Low.LessThan(st.ep) {
		st.ep = c.Low
		st.af = decimal.Min(st.af.Add(psar.step), psar.limit)
	}

	return st
}

// Count determines the total amount of data points needed for PSAR
// calculation.
func (psar PSAR) Count() int {
	return psar.length
}

// Ichimoku holds all the necessary information needed to calculate the
// Ichimoku cloud values.
// The zero value is not usable.
type Ichimoku struct {
	// valid specifies whether Ichimoku parameters were validated.
	valid bool

	// conversion specifies how many data points should be used for
	// the conversion line.
	conversion int

	// base specifies how many data points should be used for the
	// base line.
	base int

	// span specifies how many data points should be used for the
	// leading span B line.
	span int
}

// NewIchimoku validates provided configuration options and
// creates new Ichimoku indicator instance.
func NewIchimoku(conversion, base, span int) (Ichimoku, error) {
	ich := Ichimoku{
		conversion: conversion,
		base:       base,
		span:       span,
	}

	if err := ich.validate(); err != nil {
		return Ichimoku{}, err
	}

	return ich, nil
}

// validate checks whether the indicator has valid configuration properties.
func (ich *Ichimoku) validate() error {
	if ich.conversion < 1 || ich.base < 1 || ich.span < 1 {
		return ErrInvalidLength
	}

	ich.valid = true

	return nil
}

// Calc calculates all Ichimoku cloud line values from the provided
// candles ordered from the oldest to the newest. Every line is the
// midpoint between the highest high and the lowest low of its window,
// except for the leading span A line which is the midpoint between the
// conversion and the base lines. The leading span values are not
// displaced forward.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/i/ichimoku-cloud.asp.
// All credits are due to Goichi Hosoda who developed the indicator.
func (ich Ichimoku) Calc(cc []Candle) (
	conversion decimal.Decimal,
	base decimal.Decimal,
	spanA decimal.Decimal,
	spanB decimal.Decimal,
	err error,
) {

	if !ich.valid {
		return decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero, ErrInvalidIndicator
	}

	if len(cc) != ich.Count() {
		return decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero, ErrInvalidDataSize
	}

	conversion = rangeMidpoint(cc[len(cc)-ich.conversion:])
	base = rangeMidpoint(cc[len(cc)-ich.base:])
	spanA = conversion.Add(base).Div(decimal.NewFromInt(2))
	spanB = rangeMidpoint(cc[len(cc)-ich.span:])

	return conversion, base, spanA, spanB, nil
}

// Count determines the total amount of data points needed for Ichimoku
// calculation.
func (ich Ichimoku) Count() int {
	return max(ich.conversion, ich.base, ich.span)
}

// rangeMidpoint calculates the midpoint between the highest high and
// the lowest low of the provided candles.
func rangeMidpoint(cc []Candle) decimal.Decimal {
	high := cc[0].High
	low := cc[0].Low

	for _, c := range cc[1:] {
		if c.High.GreaterThan(high) {
			high = c.High
		}

		if c.Low.LessThan(low) {
			low = c.Low
		}
	}

	return high.Add(low).Div(decimal.NewFromInt(2))
}
