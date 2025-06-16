package tango

import (
	"errors"

	"github.com/shopspring/decimal"
)

// CandlestickPattern represents a candlestick pattern in technical analysis.
type CandlestickPattern string

// A list of supported candlestick patterns.
const (
	CandlestickPatternHammer         CandlestickPattern = "hammer"
	CandlestickPatternHangingMan     CandlestickPattern = "hanging-man"
	CandlestickPatternInvertedHammer CandlestickPattern = "inverted-hammer"
	CandlestickPatternShootingStar   CandlestickPattern = "shooting-star"
	CandlestickPatternLongLeggedDoji CandlestickPattern = "long-legged-doji"
	CandlestickPatternDragonflyDoji  CandlestickPattern = "dragonfly-doji"
	CandlestickPatternGravestoneDoji CandlestickPattern = "gravestone-doji"
)

// ErrInvalidCandlestickPattern indicates that the provided candlestick pattern is not valid.
var ErrInvalidCandlestickPattern = errors.New("invalid candlestick pattern")

// Validate checks if the candlestick pattern is valid.
func (cp CandlestickPattern) Validate() error {
	switch cp {
	case CandlestickPatternHammer,
		CandlestickPatternHangingMan,
		CandlestickPatternInvertedHammer,
		CandlestickPatternShootingStar,
		CandlestickPatternLongLeggedDoji,
		CandlestickPatternDragonflyDoji,
		CandlestickPatternGravestoneDoji:

		return nil
	default:
		return ErrInvalidCandlestickPattern
	}
}

// Eval evaluates whether the given data matches the candlestick pattern.
func (cp CandlestickPattern) Eval(cc []Candle) bool {
	if len(cc) != cp.Count() {
		return false
	}

	switch cp {
	case CandlestickPatternHammer:
		return evalHammer(cc[0])
	case CandlestickPatternHangingMan:
		return evalHangingMan(cc[0])
	case CandlestickPatternInvertedHammer:
		return evalInvertedHammer(cc[0])
	case CandlestickPatternShootingStar:
		return evalShootingStar(cc[0])
	case CandlestickPatternLongLeggedDoji:
		return evalLongLeggedDoji(cc[0])
	case CandlestickPatternDragonflyDoji:
		return evalDragonflyDoji(cc[0])
	case CandlestickPatternGravestoneDoji:
		return evalGravestoneDoji(cc[0])
	default:
		return false
	}
}

// Count returns the number of occurrences of the candlestick pattern.
func (cp CandlestickPattern) Count() int {
	switch cp {
	case CandlestickPatternHammer,
		CandlestickPatternHangingMan,
		CandlestickPatternInvertedHammer,
		CandlestickPatternShootingStar,
		CandlestickPatternLongLeggedDoji,
		CandlestickPatternDragonflyDoji,
		CandlestickPatternGravestoneDoji:

		return 1
	default:
		return 0
	}
}

// evalHammer evaluates whether the given candle matches the Hammer candlestick pattern.
// The candle must be positive, the body must be less than 20% of the total candle size,
// and greater than 5% of the total candle size. The close price must be
// close to the high of the candle within a certain leeway.
// It is considered a bullish pattern.
func evalHammer(c Candle) bool {
	return isWithinCandleLeewayRange(
		c.High,
		c.Low,
		c.High,
		c.Close,
		decimal.NewFromFloat(0.10),
	) && c.Open.LessThan(c.High) &&
		isWithinCandleBodySize(c, decimal.NewFromFloat(0.2), decimal.NewFromFloat(0.05))
}

// evalHangingMan evaluates whether the given candle matches the Hanging Man candlestick pattern.
// The candle must be negative, the body must be less than 20% of the total candle size,
// and greater than 5% of the total candle size. The open price must be
// close to the high of the candle within a certain leeway.
// It is considered a bearish pattern.
func evalHangingMan(c Candle) bool {
	return isWithinCandleLeewayRange(
		c.High,
		c.Low,
		c.High,
		c.Open,
		decimal.NewFromFloat(0.10),
	) && c.Close.LessThan(c.Open) &&
		isWithinCandleBodySize(c, decimal.NewFromFloat(0.2), decimal.NewFromFloat(0.05))
}

// evalInvertedHammer evaluates whether the given candle matches the Inverted Hammer candlestick pattern.
// The candle must be positive, the body must be less than 20% of the total candle size,
// and greater than 5% of the total candle size. The open price must be
// close to the low of the candle within a certain leeway.
// It is considered a bullish pattern.
func evalInvertedHammer(c Candle) bool {
	return isWithinCandleLeewayRange(
		c.High,
		c.Low,
		c.Low,
		c.Open,
		decimal.NewFromFloat(0.10),
	) && c.Close.GreaterThan(c.Low) &&
		isWithinCandleBodySize(c, decimal.NewFromFloat(0.2), decimal.NewFromFloat(0.05))
}

// evalShootingStar evaluates whether the given candle matches the Shooting Star candlestick pattern.
// The candle must be negative, the body must be less than 20% of the total candle size,
// and greater than 5% of the total candle size. The close price must be
// close to the low of the candle within a certain leeway.
// It is considered a bearish pattern.
func evalShootingStar(c Candle) bool {
	return isWithinCandleLeewayRange(
		c.High,
		c.Low,
		c.Low,
		c.Close,
		decimal.NewFromFloat(0.10),
	) && c.Open.GreaterThan(c.Low) &&
		isWithinCandleBodySize(c, decimal.NewFromFloat(0.2), decimal.NewFromFloat(0.05))
}

// evalLongLeggedDoji evaluates whether the given candle matches the Long-Legged Doji candlestick pattern.
// The candle must have a close price that is in the middle of the high and low prices,
// and the body size must be less than 5% of the total candle size.
// It is considered a neutral pattern.
func evalLongLeggedDoji(c Candle) bool {
	return isWithinCandleLeewayRange(
		c.High,
		c.Low,
		c.High.Add(c.Low).Div(decimal.NewFromInt(2)),
		c.Close,
		decimal.NewFromFloat(0.05),
	) && isWithinCandleBodySize(c, decimal.NewFromFloat(0.05), decimal.NewFromFloat(0))
}

// evalDragonflyDoji evaluates whether the given candle matches the Dragonfly Doji candlestick pattern.
// The candle must have a close price that is near the high of the candle,
// and the body size must be less than 5% of the total candle size.
// It is considered a neutral pattern.
func evalDragonflyDoji(c Candle) bool {
	return isWithinCandleLeewayRange(
		c.High,
		c.Low,
		c.High,
		c.Close,
		decimal.NewFromFloat(0.05),
	) && isWithinCandleBodySize(c, decimal.NewFromFloat(0.05), decimal.NewFromFloat(0))
}

// evalGravestoneDoji evaluates whether the given candle matches the Gravestone Doji candlestick pattern.
// The candle must have a close price that is near the low of the candle,
// and the body size must be less than 5% of the total candle size.
// It is considered a neutral pattern.
func evalGravestoneDoji(c Candle) bool {
	return isWithinCandleLeewayRange(
		c.High,
		c.Low,
		c.Low,
		c.Close,
		decimal.NewFromFloat(0.05),
	) && isWithinCandleBodySize(c, decimal.NewFromFloat(0.05), decimal.NewFromFloat(0))
}

// Candle represents a single candlestick in a financial chart.
type Candle struct {
	// Open  is the opening price of the candle.
	Open decimal.Decimal

	// High  is the highest price of the candle.
	High decimal.Decimal

	// Low   is the lowest price of the candle.
	Low decimal.Decimal

	// Close is the closing price of the candle.
	Close decimal.Decimal
}

// isWithinCandleLeewayRange checks whether the actual value is within the
// range of high and low values with the given leeway multiplier which is
// derived from the high and low of the values.
func isWithinCandleLeewayRange(high, low, expected, actual, leewayMultiplier decimal.Decimal) bool {
	leeway := high.Sub(low).Mul(leewayMultiplier)

	upperBound := expected.Add(leeway)
	lowerBound := expected.Add(leeway.Neg())

	return actual.GreaterThanOrEqual(lowerBound) &&
		actual.LessThanOrEqual(upperBound)
}

// isWithinCandleBodySize calculates the size of the value based on the
// provided high and low values.
func isWithinCandleBodySize(c Candle, upperSize, lowerSize decimal.Decimal) bool {
	var size decimal.Decimal

	if !c.High.Equal(c.Low) {
		size = c.Close.Sub(c.Open).Abs().Div(c.High.Sub(c.Low))
	}

	return size.LessThanOrEqual(upperSize) && size.GreaterThanOrEqual(lowerSize)
}
