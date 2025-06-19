package tango

import (
	"errors"
	"math"

	"github.com/shopspring/decimal"
)

var (
	// _hundred is 100 in decimal format.
	_hundred = decimal.NewFromInt(100)

	// _one is 1 in decimal format.
	_one = decimal.NewFromInt(1)
)

var (
	// ErrInvalidIndicator is returned when indicator is invalid.
	ErrInvalidIndicator = errors.New("invalid indicator")

	// ErrInvalidLength is returned when incorrect length is provided.
	ErrInvalidLength = errors.New("invalid length")

	// ErrInvalidDataSize is returned when incorrect data size is provided.
	ErrInvalidDataSize = errors.New("invalid data size")

	// ErrInvalidLevel is returned when level doesn't match any of the
	// available levels.
	ErrInvalidLevel = errors.New("invalid level")

	// ErrInvalidTrend is returned when trend doesn't match any of the
	// available trends.
	ErrInvalidTrend = errors.New("invalid trend")

	// ErrInvalidBand is returned when band doesn't match any of the
	// available bands.
	ErrInvalidBand = errors.New("invalid band")

	// ErrInvalidMA is returned when ma doesn't match any of the
	// availabble ma types.
	ErrInvalidMA = errors.New("invalid moving average")

	// ErrInvalidStandardDeviation is returned when standard deviation
	// is invalid.
	ErrInvalidStandardDeviation = errors.New("invalid standard deviation")
)

// Average is a helper function that calculates average decimal number of
// given slice.
func Average(dd []decimal.Decimal) decimal.Decimal {
	var sum decimal.Decimal

	for i := range dd {
		sum = sum.Add(dd[i])
	}

	return sum.Div(decimal.NewFromInt(int64(len(dd))))
}

// SquareRoot is a helper function that calculated the square root of decimal number.
func SquareRoot(d decimal.Decimal) decimal.Decimal {
	f, _ := d.Float64()
	return decimal.NewFromFloat(math.Sqrt(f))
}

// MeanDeviation calculates mean deviation of given slice.
func MeanDeviation(dd []decimal.Decimal) decimal.Decimal {
	length := decimal.NewFromInt(int64(len(dd)))

	if length.Equal(decimal.Zero) {
		return decimal.Zero
	}

	res := decimal.Zero
	mean := Average(dd)

	for i := range dd {
		res = res.Add(dd[i].Sub(mean).Abs().Div(length))
	}

	return res
}

// StandardDeviation calculates standard deviation of given slice.
func StandardDeviation(dd []decimal.Decimal) decimal.Decimal {
	length := decimal.NewFromInt(int64(len(dd)))

	if length.Equal(decimal.Zero) {
		return decimal.Zero
	}

	res := decimal.Zero
	mean := Average(dd)

	for i := range dd {
		res = res.Add(dd[i].Sub(mean).Pow(decimal.NewFromInt(2)).Div(length))
	}

	return SquareRoot(res)
}

// Trend specifies which trend should be used.
type Trend int

const (
	// TrendUp specifies increasing value trend.
	TrendUp Trend = iota + 1

	// TrendDown specifies decreasing value value.
	TrendDown
)

// Validate checks whether the trend is one of
// supported trend types or not.
func (t Trend) Validate() error {
	switch t {
	case TrendUp, TrendDown:
		return nil
	default:
		return ErrInvalidTrend
	}
}

// MarshalText turns trend into appropriate string
// representation.
func (t Trend) MarshalText() ([]byte, error) {
	var v string

	switch t {
	case TrendUp:
		v = "up"
	case TrendDown:
		v = "down"
	default:
		return nil, ErrInvalidTrend
	}

	return []byte(v), nil
}

// UnmarshalText turns string to appropriate trend value.
func (t *Trend) UnmarshalText(d []byte) error {
	switch string(d) {
	case "up", "u":
		*t = TrendUp
	case "down", "d":
		*t = TrendDown
	default:
		return ErrInvalidTrend
	}

	return nil
}

// Band specifies which band should be used.
type Band int

// Available Bollinger Band indicator types.
const (
	BandUpper Band = iota + 1
	BandLower
	BandWidth
)

// Validate checks whether band is one of supported band types.
func (b Band) Validate() error {
	switch b {
	case BandUpper, BandLower, BandWidth:
		return nil
	default:
		return ErrInvalidBand
	}
}

// MarshalText turns band into appropriate string representation in JSON.
func (b Band) MarshalText() ([]byte, error) {
	var v string

	switch b {
	case BandUpper:
		v = "upper"
	case BandLower:
		v = "lower"
	case BandWidth:
		v = "width"
	default:
		return nil, ErrInvalidBand
	}

	return []byte(v), nil
}

// UnmarshalText turns JSON string to appropriate band value.
func (b *Band) UnmarshalText(d []byte) error {
	switch string(d) {
	case "upper":
		*b = BandUpper
	case "lower":
		*b = BandLower
	case "width":
		*b = BandWidth
	default:
		return ErrInvalidBand
	}

	return nil
}

// MAType is a custom type that validates it to be only of existing
// moving average types.
type MAType int

// Available moving average indicator types.
const (
	MATypeDoubleExponential MAType = iota + 1
	MATypeExponential
	MATypeHull
	MATypeSimple
	MATypeWeighted
)

// NewMA constructs new moving average based on the provided type.
func NewMA(mat MAType, length int) (MA, error) {
	switch mat {
	case MATypeDoubleExponential:
		return NewDEMA(length)
	case MATypeExponential:
		return NewEMA(length)
	case MATypeHull:
		return NewHMA(length)
	case MATypeSimple:
		return NewSMA(length)
	case MATypeWeighted:
		return NewWMA(length)
	default:
		return nil, ErrInvalidMA
	}
}

// MarshalText turns MAType into appropriate string representation in JSON.
func (mat MAType) MarshalText() ([]byte, error) {
	var v string

	switch mat {
	case MATypeDoubleExponential:
		v = "double-exponential"
	case MATypeExponential:
		v = "exponential"
	case MATypeHull:
		v = "hull"
	case MATypeSimple:
		v = "simple"
	case MATypeWeighted:
		v = "weighted"
	default:
		return nil, ErrInvalidMA
	}

	return []byte(v), nil
}

// UnmarshalText turns JSON string to appropriate moving average type value.
func (mat *MAType) UnmarshalText(d []byte) error {
	switch string(d) {
	case "double-exponential":
		*mat = MATypeDoubleExponential
	case "exponential":
		*mat = MATypeExponential
	case "hull":
		*mat = MATypeHull
	case "simple":
		*mat = MATypeSimple
	case "weighted":
		*mat = MATypeWeighted
	default:
		return ErrInvalidMA
	}

	return nil
}

// MA is an interface that all moving averages implement.
type MA interface {
	// Calc should return calculation results based on provided data
	// points slice.
	Calc([]decimal.Decimal) (decimal.Decimal, error)

	// Count should determine the total amount data points required for
	// the calculation.
	Count() int
}
