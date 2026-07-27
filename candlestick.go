package tango

import (
	"errors"

	"github.com/shopspring/decimal"
)

// CandlestickPattern represents a candlestick pattern in technical analysis.
type CandlestickPattern string

// A list of supported candlestick patterns.
const (
	CandlestickPatternHammer             CandlestickPattern = "hammer"
	CandlestickPatternHangingMan         CandlestickPattern = "hanging-man"
	CandlestickPatternInvertedHammer     CandlestickPattern = "inverted-hammer"
	CandlestickPatternShootingStar       CandlestickPattern = "shooting-star"
	CandlestickPatternLongLeggedDoji     CandlestickPattern = "long-legged-doji"
	CandlestickPatternDragonflyDoji      CandlestickPattern = "dragonfly-doji"
	CandlestickPatternGravestoneDoji     CandlestickPattern = "gravestone-doji"
	CandlestickPatternBullishMarubozu    CandlestickPattern = "bullish-marubozu"
	CandlestickPatternBearishMarubozu    CandlestickPattern = "bearish-marubozu"
	CandlestickPatternSpinningTop        CandlestickPattern = "spinning-top"
	CandlestickPatternBullishEngulfing   CandlestickPattern = "bullish-engulfing"
	CandlestickPatternBearishEngulfing   CandlestickPattern = "bearish-engulfing"
	CandlestickPatternBullishHarami      CandlestickPattern = "bullish-harami"
	CandlestickPatternBearishHarami      CandlestickPattern = "bearish-harami"
	CandlestickPatternPiercingLine       CandlestickPattern = "piercing-line"
	CandlestickPatternDarkCloudCover     CandlestickPattern = "dark-cloud-cover"
	CandlestickPatternTweezerTop         CandlestickPattern = "tweezer-top"
	CandlestickPatternTweezerBottom      CandlestickPattern = "tweezer-bottom"
	CandlestickPatternMorningStar        CandlestickPattern = "morning-star"
	CandlestickPatternEveningStar        CandlestickPattern = "evening-star"
	CandlestickPatternThreeWhiteSoldiers CandlestickPattern = "three-white-soldiers"
	CandlestickPatternThreeBlackCrows    CandlestickPattern = "three-black-crows"
)

// ErrInvalidCandlestickPattern indicates that the provided candlestick pattern is not valid.
var ErrInvalidCandlestickPattern = errors.New("invalid candlestick pattern")

// patternDetails holds the evaluation properties of a single candlestick
// pattern.
type patternDetails struct {
	// count specifies how many candles the pattern is built from.
	count int

	// eval evaluates whether the given candles match the pattern. It is
	// always called with exactly count candles ordered from the oldest
	// to the newest.
	eval func(cc []Candle) bool
}

// _patterns maps all supported candlestick patterns to their evaluation
// details.
var _patterns = map[CandlestickPattern]patternDetails{
	CandlestickPatternHammer:             {count: 1, eval: single(evalHammer)},
	CandlestickPatternHangingMan:         {count: 1, eval: single(evalHangingMan)},
	CandlestickPatternInvertedHammer:     {count: 1, eval: single(evalInvertedHammer)},
	CandlestickPatternShootingStar:       {count: 1, eval: single(evalShootingStar)},
	CandlestickPatternLongLeggedDoji:     {count: 1, eval: single(evalLongLeggedDoji)},
	CandlestickPatternDragonflyDoji:      {count: 1, eval: single(evalDragonflyDoji)},
	CandlestickPatternGravestoneDoji:     {count: 1, eval: single(evalGravestoneDoji)},
	CandlestickPatternBullishMarubozu:    {count: 1, eval: single(evalBullishMarubozu)},
	CandlestickPatternBearishMarubozu:    {count: 1, eval: single(evalBearishMarubozu)},
	CandlestickPatternSpinningTop:        {count: 1, eval: single(evalSpinningTop)},
	CandlestickPatternBullishEngulfing:   {count: 2, eval: double(evalBullishEngulfing)},
	CandlestickPatternBearishEngulfing:   {count: 2, eval: double(evalBearishEngulfing)},
	CandlestickPatternBullishHarami:      {count: 2, eval: double(evalBullishHarami)},
	CandlestickPatternBearishHarami:      {count: 2, eval: double(evalBearishHarami)},
	CandlestickPatternPiercingLine:       {count: 2, eval: double(evalPiercingLine)},
	CandlestickPatternDarkCloudCover:     {count: 2, eval: double(evalDarkCloudCover)},
	CandlestickPatternTweezerTop:         {count: 2, eval: double(evalTweezerTop)},
	CandlestickPatternTweezerBottom:      {count: 2, eval: double(evalTweezerBottom)},
	CandlestickPatternMorningStar:        {count: 3, eval: triple(evalMorningStar)},
	CandlestickPatternEveningStar:        {count: 3, eval: triple(evalEveningStar)},
	CandlestickPatternThreeWhiteSoldiers: {count: 3, eval: triple(evalThreeWhiteSoldiers)},
	CandlestickPatternThreeBlackCrows:    {count: 3, eval: triple(evalThreeBlackCrows)},
}

// single adapts a one candle evaluation function to the common pattern
// evaluation signature.
func single(eval func(c Candle) bool) func(cc []Candle) bool {
	return func(cc []Candle) bool {
		return eval(cc[0])
	}
}

// double adapts a two candle evaluation function to the common pattern
// evaluation signature.
func double(eval func(first, second Candle) bool) func(cc []Candle) bool {
	return func(cc []Candle) bool {
		return eval(cc[0], cc[1])
	}
}

// triple adapts a three candle evaluation function to the common pattern
// evaluation signature.
func triple(eval func(first, second, third Candle) bool) func(cc []Candle) bool {
	return func(cc []Candle) bool {
		return eval(cc[0], cc[1], cc[2])
	}
}

// Validate checks if the candlestick pattern is valid.
func (cp CandlestickPattern) Validate() error {
	if _, ok := _patterns[cp]; !ok {
		return ErrInvalidCandlestickPattern
	}

	return nil
}

// Eval evaluates whether the given data matches the candlestick pattern.
// The candles must be ordered from the oldest to the newest.
func (cp CandlestickPattern) Eval(cc []Candle) bool {
	details, ok := _patterns[cp]
	if !ok || len(cc) != details.count {
		return false
	}

	return details.eval(cc)
}

// Count returns the number of candles required to evaluate the
// candlestick pattern.
func (cp CandlestickPattern) Count() int {
	return _patterns[cp].count
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

// evalBullishMarubozu evaluates whether the given candle matches the Bullish Marubozu
// candlestick pattern. The candle must be positive and the body must make up at
// least 95% of the total candle size, leaving no meaningful shadows.
// It is considered a bullish pattern.
func evalBullishMarubozu(c Candle) bool {
	return c.IsBullish() &&
		isWithinCandleBodySize(c, decimal.NewFromInt(1), decimal.NewFromFloat(0.95))
}

// evalBearishMarubozu evaluates whether the given candle matches the Bearish Marubozu
// candlestick pattern. The candle must be negative and the body must make up at
// least 95% of the total candle size, leaving no meaningful shadows.
// It is considered a bearish pattern.
func evalBearishMarubozu(c Candle) bool {
	return c.IsBearish() &&
		isWithinCandleBodySize(c, decimal.NewFromInt(1), decimal.NewFromFloat(0.95))
}

// evalSpinningTop evaluates whether the given candle matches the Spinning Top
// candlestick pattern. The body must be between 5% and 30% of the total candle
// size and its midpoint must be near the middle of the high and low prices,
// leaving a substantial shadow on both sides.
// It is considered a neutral pattern.
func evalSpinningTop(c Candle) bool {
	return isWithinCandleLeewayRange(
		c.High,
		c.Low,
		c.High.Add(c.Low).Div(decimal.NewFromInt(2)),
		c.BodyMidpoint(),
		decimal.NewFromFloat(0.15),
	) && isWithinCandleBodySize(c, decimal.NewFromFloat(0.3), decimal.NewFromFloat(0.05))
}

// evalBullishEngulfing evaluates whether the given candles match the Bullish
// Engulfing candlestick pattern. The first candle must be negative and the
// second candle must be positive with a strictly larger body that fully
// engulfs the body of the first candle.
// It is considered a bullish pattern.
func evalBullishEngulfing(first, second Candle) bool {
	return first.IsBearish() && second.IsBullish() &&
		second.Open.LessThanOrEqual(first.BodyBottom()) &&
		second.Close.GreaterThanOrEqual(first.BodyTop()) &&
		second.BodySize().GreaterThan(first.BodySize())
}

// evalBearishEngulfing evaluates whether the given candles match the Bearish
// Engulfing candlestick pattern. The first candle must be positive and the
// second candle must be negative with a strictly larger body that fully
// engulfs the body of the first candle.
// It is considered a bearish pattern.
func evalBearishEngulfing(first, second Candle) bool {
	return first.IsBullish() && second.IsBearish() &&
		second.Open.GreaterThanOrEqual(first.BodyTop()) &&
		second.Close.LessThanOrEqual(first.BodyBottom()) &&
		second.BodySize().GreaterThan(first.BodySize())
}

// evalBullishHarami evaluates whether the given candles match the Bullish
// Harami candlestick pattern. The first candle must be negative and the
// second candle must be positive with a strictly smaller body that is fully
// contained within the body of the first candle.
// It is considered a bullish pattern.
func evalBullishHarami(first, second Candle) bool {
	return first.IsBearish() && second.IsBullish() &&
		second.BodyTop().LessThanOrEqual(first.BodyTop()) &&
		second.BodyBottom().GreaterThanOrEqual(first.BodyBottom()) &&
		second.BodySize().LessThan(first.BodySize())
}

// evalBearishHarami evaluates whether the given candles match the Bearish
// Harami candlestick pattern. The first candle must be positive and the
// second candle must be negative with a strictly smaller body that is fully
// contained within the body of the first candle.
// It is considered a bearish pattern.
func evalBearishHarami(first, second Candle) bool {
	return first.IsBullish() && second.IsBearish() &&
		second.BodyTop().LessThanOrEqual(first.BodyTop()) &&
		second.BodyBottom().GreaterThanOrEqual(first.BodyBottom()) &&
		second.BodySize().LessThan(first.BodySize())
}

// evalPiercingLine evaluates whether the given candles match the Piercing
// Line candlestick pattern. The first candle must be negative and the second
// candle must be positive, opening at or below the close of the first candle
// and closing above the midpoint of the first candle body while staying
// below its open. The classic downward gap requirement is relaxed to an open
// at or below the prior close, since continuously traded markets rarely gap.
// It is considered a bullish pattern.
func evalPiercingLine(first, second Candle) bool {
	return first.IsBearish() && second.IsBullish() &&
		second.Open.LessThanOrEqual(first.Close) &&
		second.Close.GreaterThan(first.BodyMidpoint()) &&
		second.Close.LessThan(first.Open)
}

// evalDarkCloudCover evaluates whether the given candles match the Dark
// Cloud Cover candlestick pattern. The first candle must be positive and the
// second candle must be negative, opening at or above the close of the first
// candle and closing below the midpoint of the first candle body while
// staying above its open. The classic upward gap requirement is relaxed to
// an open at or above the prior close, since continuously traded markets
// rarely gap.
// It is considered a bearish pattern.
func evalDarkCloudCover(first, second Candle) bool {
	return first.IsBullish() && second.IsBearish() &&
		second.Open.GreaterThanOrEqual(first.Close) &&
		second.Close.LessThan(first.BodyMidpoint()) &&
		second.Close.GreaterThan(first.Open)
}

// evalTweezerTop evaluates whether the given candles match the Tweezer Top
// candlestick pattern. The first candle must be positive and the second
// candle must be negative, with both highs matching within a certain leeway.
// It is considered a bearish pattern.
func evalTweezerTop(first, second Candle) bool {
	return first.IsBullish() && second.IsBearish() &&
		isWithinCandleLeewayRange(
			first.High,
			first.Low,
			first.High,
			second.High,
			decimal.NewFromFloat(0.05),
		)
}

// evalTweezerBottom evaluates whether the given candles match the Tweezer
// Bottom candlestick pattern. The first candle must be negative and the
// second candle must be positive, with both lows matching within a certain
// leeway.
// It is considered a bullish pattern.
func evalTweezerBottom(first, second Candle) bool {
	return first.IsBearish() && second.IsBullish() &&
		isWithinCandleLeewayRange(
			first.High,
			first.Low,
			first.Low,
			second.Low,
			decimal.NewFromFloat(0.05),
		)
}

// evalMorningStar evaluates whether the given candles match the Morning Star
// candlestick pattern. The first candle must be negative, the second candle
// body must be less than 30% of the first candle body and the third candle
// must be positive, closing above the midpoint of the first candle body.
// It is considered a bullish pattern.
func evalMorningStar(first, second, third Candle) bool {
	return first.IsBearish() && third.IsBullish() &&
		second.BodySize().LessThan(first.BodySize().Mul(decimal.NewFromFloat(0.3))) &&
		third.Close.GreaterThan(first.BodyMidpoint())
}

// evalEveningStar evaluates whether the given candles match the Evening Star
// candlestick pattern. The first candle must be positive, the second candle
// body must be less than 30% of the first candle body and the third candle
// must be negative, closing below the midpoint of the first candle body.
// It is considered a bearish pattern.
func evalEveningStar(first, second, third Candle) bool {
	return first.IsBullish() && third.IsBearish() &&
		second.BodySize().LessThan(first.BodySize().Mul(decimal.NewFromFloat(0.3))) &&
		third.Close.LessThan(first.BodyMidpoint())
}

// evalThreeWhiteSoldiers evaluates whether the given candles match the Three
// White Soldiers candlestick pattern. Each candle must be positive and close
// near its high, and each subsequent candle must open within the body of the
// previous candle and close above the previous close.
// It is considered a bullish pattern.
func evalThreeWhiteSoldiers(first, second, third Candle) bool {
	return isSoldier(first) && isSoldier(second) && isSoldier(third) &&
		soldierAdvances(first, second) && soldierAdvances(second, third)
}

// isSoldier checks whether the candle is positive and closes near its high,
// making it a valid member of the Three White Soldiers pattern.
func isSoldier(c Candle) bool {
	return c.IsBullish() && isWithinCandleLeewayRange(
		c.High,
		c.Low,
		c.High,
		c.Close,
		decimal.NewFromFloat(0.2),
	)
}

// soldierAdvances checks whether the current candle opens within the body of
// the previous candle and closes above the previous close.
func soldierAdvances(prev, curr Candle) bool {
	return curr.Open.GreaterThanOrEqual(prev.BodyBottom()) &&
		curr.Open.LessThanOrEqual(prev.BodyTop()) &&
		curr.Close.GreaterThan(prev.Close)
}

// evalThreeBlackCrows evaluates whether the given candles match the Three
// Black Crows candlestick pattern. Each candle must be negative and close
// near its low, and each subsequent candle must open within the body of the
// previous candle and close below the previous close.
// It is considered a bearish pattern.
func evalThreeBlackCrows(first, second, third Candle) bool {
	return isCrow(first) && isCrow(second) && isCrow(third) &&
		crowAdvances(first, second) && crowAdvances(second, third)
}

// isCrow checks whether the candle is negative and closes near its low,
// making it a valid member of the Three Black Crows pattern.
func isCrow(c Candle) bool {
	return c.IsBearish() && isWithinCandleLeewayRange(
		c.High,
		c.Low,
		c.Low,
		c.Close,
		decimal.NewFromFloat(0.2),
	)
}

// crowAdvances checks whether the current candle opens within the body of
// the previous candle and closes below the previous close.
func crowAdvances(prev, curr Candle) bool {
	return curr.Open.LessThanOrEqual(prev.BodyTop()) &&
		curr.Open.GreaterThanOrEqual(prev.BodyBottom()) &&
		curr.Close.LessThan(prev.Close)
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

	// Volume is the traded volume of the candle. It is only used by
	// volume based indicators and may be left empty otherwise.
	Volume decimal.Decimal
}

// IsBullish reports whether the candle closed above its open.
func (c Candle) IsBullish() bool {
	return c.Close.GreaterThan(c.Open)
}

// IsBearish reports whether the candle closed below its open.
func (c Candle) IsBearish() bool {
	return c.Close.LessThan(c.Open)
}

// BodyTop returns the higher of the open and close prices.
func (c Candle) BodyTop() decimal.Decimal {
	if c.Open.GreaterThan(c.Close) {
		return c.Open
	}

	return c.Close
}

// BodyBottom returns the lower of the open and close prices.
func (c Candle) BodyBottom() decimal.Decimal {
	if c.Open.LessThan(c.Close) {
		return c.Open
	}

	return c.Close
}

// BodySize returns the absolute difference between the open and close
// prices.
func (c Candle) BodySize() decimal.Decimal {
	return c.Close.Sub(c.Open).Abs()
}

// BodyMidpoint returns the price halfway between the open and close
// prices.
func (c Candle) BodyMidpoint() decimal.Decimal {
	return c.Open.Add(c.Close).Div(decimal.NewFromInt(2))
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
