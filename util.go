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
	*s = String(CleanString(string(d)))
	return nil
}

// MarshalText converts String to a string ouput (works with JSON, etc).
func (s String) MarshalText() ([]byte, error) {
	return []byte(s), nil
}

// resize cuts given array based on length to use for
// calculations.
func resize(dd []decimal.Decimal, lh int) ([]decimal.Decimal, error) {
	if lh < 1 {
		return dd, nil
	}

	if lh > len(dd) {
		return nil, ErrInvalidDataSize
	}

	return dd[len(dd)-lh:], nil
}

// resizeCandles cuts given array based on length to use for
// calculations.
func resizeCandles(cc []chartype.Candle, lh int) ([]chartype.Candle, error) {
	if lh < 1 {
		return cc, nil
	}

	if lh > len(cc) {
		return nil, ErrInvalidDataSize
	}

	return cc[len(cc)-lh:], nil
}

// typicalPrice recalculates array of candles into an array of typical prices
func typicalPrice(cc []chartype.Candle) []decimal.Decimal {
	tp := make([]decimal.Decimal, len(cc))

	for i := 0; i < len(cc); i++ {
		tp[i] = cc[i].High.Add(cc[i].Low.Add(cc[i].Close)).Div(decimal.NewFromInt(3))
	}

	return tp
}

// meanDeviation calculates mean deviation of given array
func meanDeviation(dd []decimal.Decimal) decimal.Decimal {
	s := decimal.Zero
	rez := decimal.Zero

	for i := 0; i < len(dd); i++ {
		s = s.Add(dd[i])
	}

	s = s.Div(decimal.NewFromInt(int64(len(dd))))

	for i := 0; i < len(dd); i++ {
		rez = rez.Add(dd[i].Sub(s).Abs())
	}

	return rez.Div(decimal.NewFromInt(int64(len(dd)))).Round(8)
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
	}

	return nil, ErrInvalidSource
}
