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
func fromJSON(n string, d []byte) (Indicator, error) {
	switch n {
	case "aroon":
		a := Aroon{}
		err := json.Unmarshal(d, &a)
		return a, err
		// case "cci":
		// 	c := CCI{}
		// 	err := json.Unmarshal(d, &c)
		// 	return c, err
		// }
		// case "dema":
		// 	dm := DEMA{}
		// 	err := json.Unmarshal(d, &dm)
		// 	return dm, err
		// case "ema":
		// 	e := EMA{}
		// 	err := json.Unmarshal(d, &e)
		// 	return e, err
		// case "hma":
		// 	h := HMA{}
		// 	err := json.Unmarshal(d, &h)
		// 	return h, err
		// case "macd":
		// 	m := MACD{}
		// 	err := json.Unmarshal(d, &m)
		// 	return m, err
		// case "roc":
		// 	r := ROC{}
		// 	err := json.Unmarshal(d, &r)
		// 	return r, err
		// case "rsi":
		// 	r := RSI{}
		// 	err := json.Unmarshal(d, &r)
		// 	return r, err
		// case "sma":
		// 	s := SMA{}
		// 	err := json.Unmarshal(d, &s)
		// 	return s, err
		// case "stoch":
		// 	s := Stoch{}
		// 	err := json.Unmarshal(d, &s)
		// 	return s, err
		// case "wma":
		// 	w := WMA{}
		// 	err := json.Unmarshal(d, &w)
		// 	return w, err
	}

	return nil, ErrInvalidSourceName
}

// toJSON determines the name of the specified indicator and creates a byte
// slice array with its information.
func toJSON(ind Indicator) ([]byte, error) {
	switch i := ind.(type) {
	case Aroon:
		return json.Marshal(struct {
			Aroon
			Name string `json:"name"`
		}{Aroon: i, Name: "aroon"})
		// case CCI:
		// 	return json.Marshal(struct {
		// 		CCI
		// 		Name string `json:"name"`
		// 	}{CCI: i, Name: "cci"})
	}

	return nil, ErrInvalidType
}

// 	case DEMA:
// 		d, err := json.Marshal(struct {
// 			Name string `json:"name"`
// 			DEMA
// 		}{Name: "dema", DEMA: ind.(DEMA)})
// 		return d, err
// 	case EMA:
// 		d, err := json.Marshal(struct {
// 			Name string `json:"name"`
// 			EMA
// 		}{Name: "ema", EMA: ind.(EMA)})
// 		return d, err
// 	case HMA:
// 		d, err := json.Marshal(struct {
// 			Name string `json:"name"`
// 			HMA
// 		}{Name: "hma", HMA: ind.(HMA)})
// 		return d, err
// 	case MACD:
// 		d, err := json.Marshal(struct {
// 			Name string `json:"name"`
// 			MACD
// 		}{Name: "macd", MACD: ind.(MACD)})
// 		return d, err
// 	case ROC:
// 		d, err := json.Marshal(struct {
// 			Name string `json:"name"`
// 			ROC
// 		}{Name: "roc", ROC: ind.(ROC)})
// 		return d, err
// 	case RSI:
// 		d, err := json.Marshal(struct {
// 			Name string `json:"name"`
// 			RSI
// 		}{Name: "rsi", RSI: ind.(RSI)})
// 		return d, err
// 	case SMA:
// 		d, err := json.Marshal(struct {
// 			Name string `json:"name"`
// 			SMA
// 		}{Name: "sma", SMA: ind.(SMA)})
// 		return d, err
// 	case Stoch:
// 		d, err := json.Marshal(struct {
// 			Name string `json:"name"`
// 			Stoch
// 		}{Name: "stoch", Stoch: ind.(Stoch)})
// 		return d, err
// 	case WMA:
// 		d, err := json.Marshal(struct {
// 			Name string `json:"name"`
// 			WMA
// 		}{Name: "wma", WMA: ind.(WMA)})
// 		return d, err
// 	}

// 	return nil, ErrInvalidType
// }
