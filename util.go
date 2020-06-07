package indc

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/swithek/chartype"
)

var (
	// ErrInvalidLength is returned when incorrect length is provided.
	ErrInvalidLength = errors.New("invalid length")

	// ErrInvalidDataSize is returned when incorrect data size is provided.
	ErrInvalidDataSize = errors.New("invalid data size")

	// ErrInvalidSource is returned when source doesn't match any
	// of the available sources.
	ErrInvalidSource = errors.New("invalid source")
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
func resize(dd []decimal.Decimal, length int) ([]decimal.Decimal, error) {
	if length < 1 {
		return dd, nil
	}

	if length > len(dd) {
		return nil, ErrInvalidDataSize
	}

	return dd[len(dd)-length:], nil
}

// resizeCandles cuts given array based on length to use for
// calculations.
func resizeCandles(cc []chartype.Candle, length int) ([]chartype.Candle, error) {
	if length < 1 {
		return cc, nil
	}

	if length > len(cc) {
		return nil, ErrInvalidDataSize
	}

	return cc[len(cc)-length:], nil
}

// typicalPrice recalculates array of candles into an array of typical prices.
func typicalPrice(cc []chartype.Candle) []decimal.Decimal {
	tp := make([]decimal.Decimal, len(cc))

	for i := 0; i < len(cc); i++ {
		tp[i] = cc[i].High.Add(cc[i].Low.Add(cc[i].Close)).Div(decimal.NewFromInt(3))
	}

	return tp
}

// meanDeviation calculates mean deviation of given array.
func meanDeviation(dd []decimal.Decimal) decimal.Decimal {
	s := decimal.Zero
	rez := decimal.Zero
	length := decimal.NewFromInt(int64(len(dd)))

	if length.Equal(decimal.Zero) {
		return decimal.Zero
	}

	for i := 0; i < len(dd); i++ {
		s = s.Add(dd[i])
	}

	s = s.Div(length)

	for i := 0; i < len(dd); i++ {
		rez = rez.Add(dd[i].Sub(s).Abs())
	}

	return rez.Div(length)
}

// calcMultiple calculates specified amount of indicator within given list.
func calcMultiple(dd []decimal.Decimal, a int, s Indicator) ([]decimal.Decimal, error) {
	if a < 1 {
		return []decimal.Decimal{}, nil
	}

	dd, err := resize(dd, s.Count()+a-1)
	if err != nil {
		return nil, ErrInvalidDataSize
	}

	v := make([]decimal.Decimal, a)

	for i := 0; i < a; i++ {
		v[i], err = s.Calc(dd[:len(dd)-i])
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
		a := Aroon{}
		err := json.Unmarshal(d, &a)

		return a, err
	case NameCCI:
		c := CCI{}
		err := json.Unmarshal(d, &c)

		return c, err
	case NameDEMA:
		dm := DEMA{}
		err := json.Unmarshal(d, &dm)

		return dm, err
	case NameEMA:
		e := EMA{}
		err := json.Unmarshal(d, &e)

		return e, err
	case NameHMA:
		h := HMA{}
		err := json.Unmarshal(d, &h)

		return h, err
	case NameMACD:
		m := MACD{}
		err := json.Unmarshal(d, &m)

		return m, err
	case NameROC:
		r := ROC{}
		err := json.Unmarshal(d, &r)

		return r, err
	case NameRSI:
		r := RSI{}
		err := json.Unmarshal(d, &r)

		return r, err
	case NameSMA:
		s := SMA{}
		err := json.Unmarshal(d, &s)

		return s, err
	case NameSRSI:
		s := SRSI{}
		err := json.Unmarshal(d, &s)

		return s, err
	case NameStoch:
		s := Stoch{}
		err := json.Unmarshal(d, &s)

		return s, err
	case NameWMA:
		w := WMA{}
		err := json.Unmarshal(d, &w)

		return w, err
	// remove me in production
	case NameIndicatorMock:
		return &IndicatorMock{}, nil
	}

	return nil, ErrInvalidSource
}
