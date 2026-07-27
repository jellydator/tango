package tango

import "github.com/shopspring/decimal"

// OBV holds all the necessary information needed to calculate on-balance
// volume.
// The zero value is not usable.
type OBV struct {
	// valid specifies whether OBV parameters were validated.
	valid bool

	// length specifies how many closing price changes should be used
	// during the calculations.
	length int
}

// NewOBV validates provided configuration options and
// creates new OBV indicator instance.
func NewOBV(length int) (OBV, error) {
	obv := OBV{
		length: length,
	}

	if err := obv.validate(); err != nil {
		return OBV{}, err
	}

	return obv, nil
}

// validate checks whether the indicator has valid configuration properties.
func (obv *OBV) validate() error {
	if obv.length < 1 {
		return ErrInvalidLength
	}

	obv.valid = true

	return nil
}

// Calc calculates OBV from the provided candles ordered from the oldest
// to the newest. The running total starts at zero, every candle that
// closes above the preceding close adds its volume and every candle
// that closes below the preceding close subtracts it.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/o/onbalancevolume.asp.
// All credits are due to Joseph Granville who developed OBV indicator.
func (obv OBV) Calc(cc []Candle) (decimal.Decimal, error) {
	if !obv.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(cc) != obv.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	res := decimal.Zero

	for i := 1; i < len(cc); i++ {
		if cc[i].Close.GreaterThan(cc[i-1].Close) {
			res = res.Add(cc[i].Volume)
		}

		if cc[i].Close.LessThan(cc[i-1].Close) {
			res = res.Sub(cc[i].Volume)
		}
	}

	return res, nil
}

// Count determines the total amount of data points needed for OBV
// calculation.
func (obv OBV) Count() int {
	return obv.length + 1
}

// MFI holds all the necessary information needed to calculate money
// flow index.
// The zero value is not usable.
type MFI struct {
	// valid specifies whether MFI parameters were validated.
	valid bool

	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewMFI validates provided configuration options and
// creates new MFI indicator instance.
func NewMFI(length int) (MFI, error) {
	mfi := MFI{
		length: length,
	}

	if err := mfi.validate(); err != nil {
		return MFI{}, err
	}

	return mfi, nil
}

// validate checks whether the indicator has valid configuration properties.
func (mfi *MFI) validate() error {
	if mfi.length < 1 {
		return ErrInvalidLength
	}

	mfi.valid = true

	return nil
}

// Calc calculates MFI from the provided candles ordered from the oldest
// to the newest. The money flow of every candle is its typical price
// multiplied by its volume, counted as positive when the typical price
// rises above the preceding one and as negative when it falls below it.
// Calculation is based on formula provided by investopedia.
// https://www.investopedia.com/terms/m/mfi.asp.
// All credits are due to Gene Quong and Avrum Soudack who developed
// MFI indicator.
func (mfi MFI) Calc(cc []Candle) (decimal.Decimal, error) {
	if !mfi.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(cc) != mfi.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	positive := decimal.Zero
	negative := decimal.Zero

	prev := typicalPrice(cc[0])

	for _, c := range cc[1:] {
		curr := typicalPrice(c)
		flow := curr.Mul(c.Volume)

		if curr.GreaterThan(prev) {
			positive = positive.Add(flow)
		}

		if curr.LessThan(prev) {
			negative = negative.Add(flow)
		}

		prev = curr
	}

	if positive.IsZero() {
		return decimal.Zero, nil
	}

	if negative.IsZero() {
		return _hundred, nil
	}

	return _hundred.Sub(_hundred.Div(_one.Add(positive.Div(negative)))), nil
}

// Count determines the total amount of data points needed for MFI
// calculation.
func (mfi MFI) Count() int {
	return mfi.length + 1
}

// typicalPrice calculates the typical price of the candle.
func typicalPrice(c Candle) decimal.Decimal {
	return c.High.Add(c.Low).Add(c.Close).Div(decimal.NewFromInt(3))
}

// CMF holds all the necessary information needed to calculate Chaikin
// money flow.
// The zero value is not usable.
type CMF struct {
	// valid specifies whether CMF parameters were validated.
	valid bool

	// length specifies how many data points should be used
	// during the calculations.
	length int
}

// NewCMF validates provided configuration options and
// creates new CMF indicator instance.
func NewCMF(length int) (CMF, error) {
	cmf := CMF{
		length: length,
	}

	if err := cmf.validate(); err != nil {
		return CMF{}, err
	}

	return cmf, nil
}

// validate checks whether the indicator has valid configuration properties.
func (cmf *CMF) validate() error {
	if cmf.length < 1 {
		return ErrInvalidLength
	}

	cmf.valid = true

	return nil
}

// Calc calculates CMF from the provided candles ordered from the oldest
// to the newest. The money flow volume of every candle is its volume
// scaled by the position of the close between the low and the high, and
// the result is the total money flow volume divided by the total
// volume. The result is zero when the total volume is zero and candles
// with an empty price range contribute no money flow volume.
// Calculation is based on formula provided by fidelity.
// https://www.fidelity.com/learning-center/trading-investing/technical-analysis/technical-indicator-guide/cmf.
// All credits are due to Marc Chaikin who developed CMF indicator.
func (cmf CMF) Calc(cc []Candle) (decimal.Decimal, error) {
	if !cmf.valid {
		return decimal.Zero, ErrInvalidIndicator
	}

	if len(cc) != cmf.Count() {
		return decimal.Zero, ErrInvalidDataSize
	}

	flow := decimal.Zero
	volume := decimal.Zero

	for _, c := range cc {
		volume = volume.Add(c.Volume)

		spread := c.High.Sub(c.Low)
		if spread.IsZero() {
			continue
		}

		multiplier := c.Close.Sub(c.Low).Sub(c.High.Sub(c.Close)).Div(spread)
		flow = flow.Add(multiplier.Mul(c.Volume))
	}

	if volume.IsZero() {
		return decimal.Zero, nil
	}

	return flow.Div(volume), nil
}

// Count determines the total amount of data points needed for CMF
// calculation.
func (cmf CMF) Count() int {
	return cmf.length
}
