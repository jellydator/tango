package indc

import (
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

	// ErrSrcIndicatorNotSet is returned when src indicator field is nil.
	ErrSrcNotSet = errors.New("source indicator is not set")

	// ErrInvalidSrcName is returned when provided indicator name
	// isn't recognized.
	ErrInvalidSrcName = errors.New("unrecognized source indicator name")

	// ErrSrcIndicatorNotSet is returned when indicator field is nil.
	ErrMAIndicatorNotSet = errors.New("ma indicator value not set")

	// ErrInvalidType is returned when indicator type doesn't match any
	// of the available types.
	ErrInvalidType = errors.New("invalid indicator type")
)

// resize cuts given array based on length to use for
// calculations.
func resize(dd []decimal.Decimal, l int) ([]decimal.Decimal, error) {
	if l < 1 {
		return nil, ErrInvalidLength
	}

	if l > len(dd) {
		return nil, ErrInvalidDataPointCount
	}

	return dd[len(dd)-l:], nil
}

// resizeCandles cuts given array based on length to use for
// calculations.
func resizeCandles(cc []chartype.Candle, l int) ([]chartype.Candle, error) {
	if l < 1 {
		return nil, ErrInvalidLength
	}

	if l > len(cc) || l < 1 {
		return nil, ErrInvalidDataPointCount
	}

	return cc[len(cc)-l:], nil
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

// newIndicator finds and returns an empty indicator of the specified name.
func newIndicator(n string) (Indicator, error) {
	switch n {
	case "aroon":
		return Aroon{}, nil
	case "cci":
		return CCI{}, nil
	case "dema":
		return DEMA{}, nil
	case "ema":
		return EMA{}, nil
	case "hma":
		return HMA{}, nil
	case "macd":
		return MACD{}, nil
	case "roc":
		return ROC{}, nil
	case "rsi":
		return RSI{}, nil
	case "sma":
		return SMA{}, nil
	case "stoch":
		return Stoch{}, nil
	case "wma":
		return WMA{}, nil
	}

	return nil, ErrInvalidType
}

// indicatorname determines the name of the specified indicator.
func indicatorName(ind Indicator) (string, error) {
	switch ind.(type) {
	case Aroon:
		return "aroon", nil
	case CCI:
		return "cci", nil
	case DEMA:
		return "dema", nil
	case EMA:
		return "ema", nil
	case HMA:
		return "hma", nil
	case MACD:
		return "macd", nil
	case ROC:
		return "roc", nil
	case RSI:
		return "rsi", nil
	case SMA:
		return "sma", nil
	case Stoch:
		return "stoch", nil
	case WMA:
		return "wma", nil
	}

	return "", ErrInvalidType
}
