package indc

import (
	"encoding/json"
	"errors"

	"github.com/shopspring/decimal"
	"github.com/swithek/chartype"
)

var (
	// ErrInvalidDataPointCount is returned when insufficient amount of
	// data points is provided.
	ErrInvalidDataPointCount = errors.New("insufficient amount of data points")

	// ErrInvalidLength is returned when provided length is less than 1.
	ErrInvalidLength = errors.New("length cannot be less than 1")

	// ErrSourceNotSet is returned when source indicator field is nil.
	ErrSourceNotSet = errors.New("source indicator is not set")

	// ErrInvalidSourceName is returned when provided indicator name
	// isn't recognized.
	ErrInvalidSourceName = errors.New("unrecognized source indicator name")

	// ErrMANotSet is returned when indicator field is nil.
	ErrMANotSet = errors.New("ma value not set")

	// ErrInvalidType is returned when indicator type doesn't match any
	// of the available types.
	ErrInvalidType = errors.New("invalid indicator type")
)

// resize cuts given array based on length to use for
// calculations.
func resize(dd []decimal.Decimal, lh int) ([]decimal.Decimal, error) {
	if lh < 1 {
		return nil, ErrInvalidLength
	}

	if lh > len(dd) {
		return nil, ErrInvalidDataPointCount
	}

	return dd[len(dd)-lh:], nil
}

// resizeCandles cuts given array based on length to use for
// calculations.
func resizeCandles(cc []chartype.Candle, lh int) ([]chartype.Candle, error) {
	if lh < 1 {
		return nil, ErrInvalidLength
	}

	if lh > len(cc) || lh < 1 {
		return nil, ErrInvalidDataPointCount
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

// fromJSON finds an indicator based on its name and returns it as Interface
// with its values.
func fromJSON(d []byte) (Indicator, error) {
	var i struct {
		N string `json:"name"`
	}

	if err := json.Unmarshal(d, &i); err != nil {
		return nil, err
	}

	switch i.N {
	case "aroon":
		a := Aroon{}
		err := json.Unmarshal(d, &a)
		return a, err
	case "cci":
		c := CCI{}
		err := json.Unmarshal(d, &c)
		return c, err
	case "dema":
		dm := DEMA{}
		err := json.Unmarshal(d, &dm)
		return dm, err
	case "ema":
		e := EMA{}
		err := json.Unmarshal(d, &e)
		return e, err
	case "hma":
		h := HMA{}
		err := json.Unmarshal(d, &h)
		return h, err
	// case "macd":
	// 	m := MACD{}
	// 	err := json.Unmarshal(d, &m)
	// 	return m, err
	case "roc":
		r := ROC{}
		err := json.Unmarshal(d, &r)
		return r, err
	case "rsi":
		r := RSI{}
		err := json.Unmarshal(d, &r)
		return r, err
	case "sma":
		s := SMA{}
		err := json.Unmarshal(d, &s)
		return s, err
	case "stoch":
		s := Stoch{}
		err := json.Unmarshal(d, &s)
		return s, err
	case "wma":
		w := WMA{}
		err := json.Unmarshal(d, &w)
		return w, err
	}

	return nil, ErrInvalidSourceName
}
