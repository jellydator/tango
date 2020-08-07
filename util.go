package indc

import (
	"encoding/json"
	"errors"
	"math"
	"strings"

	"github.com/shopspring/decimal"
)

var (
	// Hundred is just plain 100 in decimal format.
	Hundred = decimal.NewFromInt(100)

	// One is just plain 1 in decimal format.
	One = decimal.NewFromInt(1)
)

var (
	// ErrInvalidIndicator is returned when indicator is invalid.
	ErrInvalidIndicator = errors.New("invalid indicator")

	// ErrInvalidLength is returned when incorrect length is provided.
	ErrInvalidLength = errors.New("invalid length")

	// ErrInvalidOffset is returned when incorrect offset is provided.
	ErrInvalidOffset = errors.New("invalid offset")

	// ErrInvalidDataSize is returned when incorrect data size is provided.
	ErrInvalidDataSize = errors.New("invalid data size")

	// ErrInvalidSource is returned when source doesn't match any of the
	// available sources.
	ErrInvalidSource = errors.New("invalid source")

	// ErrInvalidTrend is returned when trend doesn't match any of the
	// available trends.
	ErrInvalidTrend = errors.New("invalid trend")

	// ErrInvalidBand is returned when band doesn't match any of the
	// available bands.
	ErrInvalidBand = errors.New("invalid band")
)

// String is a custom string that helps prevent capitalization issues by
// lowercasing provided string.
type String string

// CleanString returns a properly formatted string.
func CleanString(s string) String {
	return String(strings.ToLower(strings.TrimSpace(s)))
}

// UnmarshalText parses String from a string form input (works with JSON, etc).
func (s *String) UnmarshalText(d []byte) error {
	*s = CleanString(string(d))
	return nil
}

// MarshalText converts String to a string output (works with JSON, etc).
func (s String) MarshalText() ([]byte, error) {
	return []byte(s), nil
}

// resize cuts given array based on length to use for
// calculations.
func resize(dd []decimal.Decimal, length, offset int) ([]decimal.Decimal, error) {
	if length < 1 || offset < 0 {
		return dd, nil
	}

	if length+offset > len(dd) {
		return nil, ErrInvalidDataSize
	}

	return dd[len(dd)-length-offset : len(dd)-offset], nil
}

// average calculates average decimal number of given array.
func average(dd []decimal.Decimal) decimal.Decimal {
	var sum decimal.Decimal

	for i := range dd {
		sum = sum.Add(dd[i])
	}

	return sum.Div(decimal.NewFromInt(int64(len(dd))))
}

// sqrt is used to get a square root of decimal number.
func sqrt(d decimal.Decimal) decimal.Decimal {
	f, _ := d.Float64()

	return decimal.NewFromFloat(math.Sqrt(f))
}

// meanDeviation calculates mean deviation of given array.
func meanDeviation(dd []decimal.Decimal) decimal.Decimal {
	length := decimal.NewFromInt(int64(len(dd)))

	if length.Equal(decimal.Zero) {
		return decimal.Zero
	}

	res := decimal.Zero
	mean := average(dd)

	for i := range dd {
		res = res.Add(dd[i].Sub(mean).Abs().Div(length))
	}

	return res
}

// standardDeviation calculates standart deviation of given array.
func standardDeviation(dd []decimal.Decimal) decimal.Decimal {
	length := decimal.NewFromInt(int64(len(dd)))

	if length.Equal(decimal.Zero) {
		return decimal.Zero
	}

	res := decimal.Zero
	mean := average(dd)

	for i := range dd {
		res = res.Add(dd[i].Sub(mean).Pow(decimal.NewFromInt(2)).Div(length))
	}

	return sqrt(res)
}

// calcMultiple calculates specified amount of values by using specified
// indicator.
func calcMultiple(src Indicator, amount int, dd []decimal.Decimal) ([]decimal.Decimal, error) {
	if amount < 1 {
		return []decimal.Decimal{}, nil
	}

	dd, err := resize(dd, src.Count()+amount-1, 0)
	if err != nil {
		return nil, ErrInvalidDataSize
	}

	v := make([]decimal.Decimal, amount)

	for i := 0; i < amount; i++ {
		v[i], err = src.Calc(dd[:len(dd)-i])
		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

// fromJSON finds a source indicator based on its name and returns it.
func fromJSON(d []byte) (Indicator, error) {
	var i struct {
		N String `json:"name"`
	}

	if err := json.Unmarshal(d, &i); err != nil {
		return nil, err
	}

	switch i.N {
	case NameAroon:
		v := Aroon{}
		err := json.Unmarshal(d, &v)

		return v, err
	case NameBB:
		v := BB{}
		err := json.Unmarshal(d, &v)

		return v, err
	case NameCCI:
		v := CCI{}
		err := json.Unmarshal(d, &v)

		return v, err
	case NameDEMA:
		v := DEMA{}
		err := json.Unmarshal(d, &v)

		return v, err
	case NameEMA:
		v := EMA{}
		err := json.Unmarshal(d, &v)

		return v, err
	case NameHMA:
		v := HMA{}
		err := json.Unmarshal(d, &v)

		return v, err
	case NameCD:
		v := CD{}
		err := json.Unmarshal(d, &v)

		return v, err
	case NameROC:
		v := ROC{}
		err := json.Unmarshal(d, &v)

		return v, err
	case NameRSI:
		v := RSI{}
		err := json.Unmarshal(d, &v)

		return v, err
	case NameSMA:
		v := SMA{}
		err := json.Unmarshal(d, &v)

		return v, err
	case NameSRSI:
		v := SRSI{}
		err := json.Unmarshal(d, &v)

		return v, err
	case NameStoch:
		v := Stoch{}
		err := json.Unmarshal(d, &v)

		return v, err
	case NameWMA:
		v := WMA{}
		err := json.Unmarshal(d, &v)

		return v, err
	}

	return nil, ErrInvalidSource
}

const (
	// TrendUp specifies increasing value trend.
	TrendUp Trend = iota + 1

	// TrendDown specifies decreasing value value.
	TrendDown
)

// Trend specifies which trend should be used.
type Trend int

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

// Available Bollinger Band indicator types.
const (
	BandUpper Band = iota + 1
	BandMiddle
	BandLower
	BandWidth
)

// Band specifies which band should be used.
type Band int

// Validate checks whether the band is one of
// supported band types or not.
func (b Band) Validate() error {
	switch b {
	case BandUpper, BandMiddle, BandLower, BandWidth:
		return nil
	default:
		return ErrInvalidBand
	}
}

// MarshalText turns band into appropriate string
// representation in JSON.
func (b Band) MarshalText() ([]byte, error) {
	var v string

	switch b {
	case BandUpper:
		v = "upper"
	case BandMiddle:
		v = "middle"
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
	case "upper", "u":
		*b = BandUpper
	case "middle", "m":
		*b = BandMiddle
	case "lower", "l":
		*b = BandLower
	case "width", "w":
		*b = BandWidth
	default:
		return ErrInvalidBand
	}

	return nil
}
