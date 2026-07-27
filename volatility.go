package tango

import "github.com/shopspring/decimal"

// ATR holds all the necessary information needed to calculate average
// true range.
// The zero value is not usable.
type ATR struct {
	// valid specifies whether ATR parameters were validated.
	valid bool

	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewATR validates provided configuration options and
// creates new ATR indicator instance.
func NewATR(length int) (ATR, error) {
	atr := ATR{
		length: length,
	}

	if err := atr.validate(); err != nil {
		return ATR{}, err
	}

	return atr, nil
}

// validate checks whether the indicator has valid configuration properties.
func (atr *ATR) validate() error {
	if atr.length < 1 {
		return ErrInvalidLength
	}

	atr.valid = true

	return nil
}

// Calc calculates ATR from the provided candles ordered from the oldest
// to the newest. The first length true ranges seed a simple average
// which is then smoothed with the remaining true ranges using the
// Wilder smoothing method.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/a/atr.asp.
// All credits are due to J. Welles Wilder Jr. who developed ATR indicator.
func (atr ATR) Calc(cc []Candle) (decimal.Decimal, error) {
	if !atr.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(cc) != atr.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	trs := trueRanges(cc)
	length := decimal.NewFromInt(int64(atr.length))

	res := Average(trs[:atr.length])

	for _, tr := range trs[atr.length:] {
		res = res.Mul(length.Sub(_one)).Add(tr).Div(length)
	}

	return res, nil
}

// Count determines the total amount of data points needed for ATR
// calculation.
func (atr ATR) Count() int {
	return atr.length*2 + 1
}

// trueRanges calculates the true range of every candle that has a
// preceding candle. The candles must be ordered from the oldest to the
// newest.
func trueRanges(cc []Candle) []decimal.Decimal {
	trs := make([]decimal.Decimal, 0, len(cc)-1)

	for i := 1; i < len(cc); i++ {
		trs = append(trs, trueRange(cc[i-1].Close, cc[i]))
	}

	return trs
}

// trueRange calculates the true range of the candle by using the
// closing price of the preceding candle.
func trueRange(prevClose decimal.Decimal, c Candle) decimal.Decimal {
	res := c.High.Sub(c.Low)

	if v := c.High.Sub(prevClose).Abs(); v.GreaterThan(res) {
		res = v
	}

	if v := c.Low.Sub(prevClose).Abs(); v.GreaterThan(res) {
		res = v
	}

	return res
}

// Donchian holds all the necessary information needed to calculate
// Donchian channels.
// The zero value is not usable.
type Donchian struct {
	// valid specifies whether Donchian parameters were validated.
	valid bool

	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewDonchian validates provided configuration options and
// creates new Donchian indicator instance.
func NewDonchian(length int) (Donchian, error) {
	don := Donchian{
		length: length,
	}

	if err := don.validate(); err != nil {
		return Donchian{}, err
	}

	return don, nil
}

// validate checks whether the indicator has valid configuration properties.
func (don *Donchian) validate() error {
	if don.length < 1 {
		return ErrInvalidLength
	}

	don.valid = true

	return nil
}

// Calc calculates all Donchian channel values from the provided candles
// ordered from the oldest to the newest. The upper channel is the
// highest high, the lower channel is the lowest low and the middle
// channel is the average of the two.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/d/donchianchannels.asp.
// All credits are due to Richard Donchian who developed the indicator.
func (don Donchian) Calc(cc []Candle) (
	upper decimal.Decimal,
	middle decimal.Decimal,
	lower decimal.Decimal,
	err error,
) {

	if !don.valid {
		return decimal.Zero, decimal.Zero, decimal.Zero, ErrInvalidIndicator
	}

	if len(cc) != don.Count() {
		return decimal.Zero, decimal.Zero, decimal.Zero, ErrInvalidDataSize
	}

	upper = cc[0].High
	lower = cc[0].Low

	for _, c := range cc[1:] {
		if c.High.GreaterThan(upper) {
			upper = c.High
		}

		if c.Low.LessThan(lower) {
			lower = c.Low
		}
	}

	return upper, upper.Add(lower).Div(decimal.NewFromInt(2)), lower, nil
}

// Count determines the total amount of data points needed for Donchian
// calculation.
func (don Donchian) Count() int {
	return don.length
}

// Keltner holds all the necessary information needed to calculate
// Keltner channels.
// The zero value is not usable.
type Keltner struct {
	// valid specifies whether Keltner parameters were validated.
	valid bool

	// multiplier specifies by how much the average true range is
	// scaled before it is applied to the middle channel.
	multiplier decimal.Decimal

	// ma specifies the moving average used for the middle channel.
	ma MA

	// atr specifies the average true range used for the channel
	// offsets.
	atr ATR
}

// NewKeltner validates provided configuration options and
// creates new Keltner indicator instance.
func NewKeltner(mat MAType, multiplier decimal.Decimal, length int) (Keltner, error) {
	ma, err := NewMA(mat, length)
	if err != nil {
		return Keltner{}, err
	}

	atr, err := NewATR(length)
	if err != nil {
		return Keltner{}, err
	}

	kelt := Keltner{
		multiplier: multiplier,
		ma:         ma,
		atr:        atr,
	}

	if err := kelt.validate(); err != nil {
		return Keltner{}, err
	}

	return kelt, nil
}

// validate checks whether the indicator has valid configuration properties.
func (kelt *Keltner) validate() error {
	if kelt.multiplier.Cmp(decimal.Zero) <= 0 {
		return ErrInvalidMultiplier
	}

	kelt.valid = true

	return nil
}

// Calc calculates all Keltner channel values from the provided candles
// ordered from the oldest to the newest. The middle channel is a moving
// average of the closing prices and the upper and lower channels are
// offset from it by the scaled average true range.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/k/keltnerchannel.asp.
// All credits are due to Chester Keltner who developed the indicator.
func (kelt Keltner) Calc(cc []Candle) (
	upper decimal.Decimal,
	middle decimal.Decimal,
	lower decimal.Decimal,
	err error,
) {

	if !kelt.valid {
		return decimal.Zero, decimal.Zero, decimal.Zero, ErrInvalidIndicator
	}

	if len(cc) != kelt.Count() {
		return decimal.Zero, decimal.Zero, decimal.Zero, ErrInvalidDataSize
	}

	dd := make([]decimal.Decimal, 0, kelt.ma.Count())
	for _, c := range cc[len(cc)-kelt.ma.Count():] {
		dd = append(dd, c.Close)
	}

	middle, err = kelt.ma.Calc(dd)
	if err != nil {
		// unlikely to happen
		return decimal.Zero, decimal.Zero, decimal.Zero, err
	}

	atr, err := kelt.atr.Calc(cc[len(cc)-kelt.atr.Count():])
	if err != nil {
		// unlikely to happen
		return decimal.Zero, decimal.Zero, decimal.Zero, err
	}

	offset := atr.Mul(kelt.multiplier)

	return middle.Add(offset), middle, middle.Sub(offset), nil
}

// Count determines the total amount of data points needed for Keltner
// calculation.
func (kelt Keltner) Count() int {
	return max(kelt.ma.Count(), kelt.atr.Count())
}
