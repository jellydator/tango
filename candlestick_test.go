package tango

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_CandlestickPattern_Validate(t *testing.T) {
	patterns := []CandlestickPattern{
		CandlestickPatternHammer,
		CandlestickPatternHangingMan,
		CandlestickPatternInvertedHammer,
		CandlestickPatternShootingStar,
		CandlestickPatternLongLeggedDoji,
		CandlestickPatternDragonflyDoji,
		CandlestickPatternGravestoneDoji,
		CandlestickPatternBullishMarubozu,
		CandlestickPatternBearishMarubozu,
		CandlestickPatternSpinningTop,
		CandlestickPatternBullishEngulfing,
		CandlestickPatternBearishEngulfing,
		CandlestickPatternBullishHarami,
		CandlestickPatternBearishHarami,
		CandlestickPatternPiercingLine,
		CandlestickPatternDarkCloudCover,
		CandlestickPatternTweezerTop,
		CandlestickPatternTweezerBottom,
		CandlestickPatternMorningStar,
		CandlestickPatternEveningStar,
		CandlestickPatternThreeWhiteSoldiers,
		CandlestickPatternThreeBlackCrows,
	}

	for _, pattern := range patterns {
		assert.NoError(t, pattern.Validate())
	}

	assert.Error(t, CandlestickPattern("invalid").Validate(), ErrInvalidCandlestickPattern)
}

func Test_CandlestickPattern_Count(t *testing.T) {
	assert.Equal(t, 0, CandlestickPattern("invalid").Count())
	assert.Equal(t, 1, CandlestickPatternHammer.Count())
	assert.Equal(t, 2, CandlestickPatternBullishEngulfing.Count())
	assert.Equal(t, 3, CandlestickPatternMorningStar.Count())
}

func Test_CandlestickPattern_Eval(t *testing.T) {
	cc := map[string]struct {
		Pattern CandlestickPattern
		Candles []Candle
		Result  bool
	}{
		"Invalid pattern": {},
		"Invalid candle count": {
			Pattern: CandlestickPatternHammer,
		},
		"Invalid candle count for a multi candle pattern": {
			Pattern: CandlestickPatternBullishEngulfing,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(100),
					Close: decimal.NewFromFloat(95),
					Open:  decimal.NewFromFloat(80),
					Low:   decimal.NewFromFloat(20),
				},
			},
		},
		"Successfully evaluated hammer pattern with some leeway": {
			Pattern: CandlestickPatternHammer,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(100),
					Close: decimal.NewFromFloat(95),
					Open:  decimal.NewFromFloat(80),
					Low:   decimal.NewFromFloat(20),
				},
			},
			Result: true,
		},
		"Successfully evaluated hammer pattern": {
			Pattern: CandlestickPatternHammer,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(100),
					Close: decimal.NewFromFloat(100),
					Open:  decimal.NewFromFloat(90),
					Low:   decimal.NewFromFloat(40),
				},
			},
			Result: true,
		},
		"Successfully evaluated hanging man pattern with some leeway": {
			Pattern: CandlestickPatternHangingMan,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(100),
					Close: decimal.NewFromFloat(80),
					Open:  decimal.NewFromFloat(95),
					Low:   decimal.NewFromFloat(20),
				},
			},
			Result: true,
		},
		"Successfully evaluated hanging man pattern": {
			Pattern: CandlestickPatternHangingMan,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(100),
					Close: decimal.NewFromFloat(90),
					Open:  decimal.NewFromFloat(100),
					Low:   decimal.NewFromFloat(40),
				},
			},
			Result: true,
		},
		"Successfully evaluated inverted hammer pattern with some leeway": {
			Pattern: CandlestickPatternInvertedHammer,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(100),
					Close: decimal.NewFromFloat(40),
					Open:  decimal.NewFromFloat(25),
					Low:   decimal.NewFromFloat(20),
				},
			},
			Result: true,
		},
		"Successfully evaluated inverted hammer pattern": {
			Pattern: CandlestickPatternInvertedHammer,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(100),
					Close: decimal.NewFromFloat(50),
					Open:  decimal.NewFromFloat(40),
					Low:   decimal.NewFromFloat(40),
				},
			},
			Result: true,
		},
		"Successfully evaluated shooting star pattern with some leeway": {
			Pattern: CandlestickPatternShootingStar,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(100),
					Close: decimal.NewFromFloat(25),
					Open:  decimal.NewFromFloat(40),
					Low:   decimal.NewFromFloat(20),
				},
			},
			Result: true,
		},
		"Successfully evaluated shooting star pattern": {
			Pattern: CandlestickPatternShootingStar,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(100),
					Close: decimal.NewFromFloat(40),
					Open:  decimal.NewFromFloat(50),
					Low:   decimal.NewFromFloat(40),
				},
			},
			Result: true,
		},
		"Successfully evaluated long legged doji pattern with some leeway": {
			Pattern: CandlestickPatternLongLeggedDoji,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(100),
					Close: decimal.NewFromFloat(62),
					Open:  decimal.NewFromFloat(59),
					Low:   decimal.NewFromFloat(20),
				},
			},
			Result: true,
		},
		"Successfully evaluated long legged doji star pattern": {
			Pattern: CandlestickPatternLongLeggedDoji,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(100),
					Close: decimal.NewFromFloat(70),
					Open:  decimal.NewFromFloat(70),
					Low:   decimal.NewFromFloat(40),
				},
			},
			Result: true,
		},
		"Successfully evaluated dragonfly doji pattern with some leeway": {
			Pattern: CandlestickPatternDragonflyDoji,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(100),
					Close: decimal.NewFromFloat(98),
					Open:  decimal.NewFromFloat(99),
					Low:   decimal.NewFromFloat(20),
				},
			},
			Result: true,
		},
		"Successfully evaluated dragonfly doji star pattern": {
			Pattern: CandlestickPatternDragonflyDoji,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(100),
					Close: decimal.NewFromFloat(100),
					Open:  decimal.NewFromFloat(100),
					Low:   decimal.NewFromFloat(40),
				},
			},
			Result: true,
		},
		"Successfully evaluated gravestone doji pattern with some leeway": {
			Pattern: CandlestickPatternGravestoneDoji,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(100),
					Close: decimal.NewFromFloat(22),
					Open:  decimal.NewFromFloat(23),
					Low:   decimal.NewFromFloat(20),
				},
			},
			Result: true,
		},
		"Successfully evaluated gravestone doji star pattern": {
			Pattern: CandlestickPatternGravestoneDoji,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(100),
					Close: decimal.NewFromFloat(40),
					Open:  decimal.NewFromFloat(40),
					Low:   decimal.NewFromFloat(40),
				},
			},
			Result: true,
		},
		"Successfully evaluated bullish marubozu pattern": {
			Pattern: CandlestickPatternBullishMarubozu,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(120),
					Close: decimal.NewFromFloat(120),
					Open:  decimal.NewFromFloat(100),
					Low:   decimal.NewFromFloat(100),
				},
			},
			Result: true,
		},
		"Successfully evaluated bullish marubozu pattern with some leeway": {
			Pattern: CandlestickPatternBullishMarubozu,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(120.5),
					Close: decimal.NewFromFloat(120),
					Open:  decimal.NewFromFloat(100),
					Low:   decimal.NewFromFloat(99.8),
				},
			},
			Result: true,
		},
		"Failed to match bullish marubozu pattern with a long shadow": {
			Pattern: CandlestickPatternBullishMarubozu,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(100),
					Close: decimal.NewFromFloat(95),
					Open:  decimal.NewFromFloat(80),
					Low:   decimal.NewFromFloat(20),
				},
			},
		},
		"Successfully evaluated bearish marubozu pattern": {
			Pattern: CandlestickPatternBearishMarubozu,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(120),
					Close: decimal.NewFromFloat(100),
					Open:  decimal.NewFromFloat(120),
					Low:   decimal.NewFromFloat(100),
				},
			},
			Result: true,
		},
		"Failed to match bearish marubozu pattern with a positive candle": {
			Pattern: CandlestickPatternBearishMarubozu,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(120),
					Close: decimal.NewFromFloat(120),
					Open:  decimal.NewFromFloat(100),
					Low:   decimal.NewFromFloat(100),
				},
			},
		},
		"Successfully evaluated spinning top pattern": {
			Pattern: CandlestickPatternSpinningTop,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(110),
					Close: decimal.NewFromFloat(102),
					Open:  decimal.NewFromFloat(98),
					Low:   decimal.NewFromFloat(90),
				},
			},
			Result: true,
		},
		"Successfully evaluated spinning top pattern with some leeway": {
			Pattern: CandlestickPatternSpinningTop,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(110),
					Close: decimal.NewFromFloat(97),
					Open:  decimal.NewFromFloat(101),
					Low:   decimal.NewFromFloat(90),
				},
			},
			Result: true,
		},
		"Failed to match spinning top pattern with a large body": {
			Pattern: CandlestickPatternSpinningTop,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(110),
					Close: decimal.NewFromFloat(108),
					Open:  decimal.NewFromFloat(92),
					Low:   decimal.NewFromFloat(90),
				},
			},
		},
		"Successfully evaluated bullish engulfing pattern": {
			Pattern: CandlestickPatternBullishEngulfing,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(105),
					Close: decimal.NewFromFloat(95),
					Open:  decimal.NewFromFloat(100),
					Low:   decimal.NewFromFloat(90),
				},
				{
					High:  decimal.NewFromFloat(110),
					Close: decimal.NewFromFloat(103),
					Open:  decimal.NewFromFloat(94),
					Low:   decimal.NewFromFloat(92),
				},
			},
			Result: true,
		},
		"Failed to match bullish engulfing pattern without the prior body engulfed": {
			Pattern: CandlestickPatternBullishEngulfing,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(105),
					Close: decimal.NewFromFloat(95),
					Open:  decimal.NewFromFloat(100),
					Low:   decimal.NewFromFloat(90),
				},
				{
					High:  decimal.NewFromFloat(110),
					Close: decimal.NewFromFloat(103),
					Open:  decimal.NewFromFloat(96),
					Low:   decimal.NewFromFloat(92),
				},
			},
		},
		"Successfully evaluated bearish engulfing pattern": {
			Pattern: CandlestickPatternBearishEngulfing,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(105),
					Close: decimal.NewFromFloat(100),
					Open:  decimal.NewFromFloat(95),
					Low:   decimal.NewFromFloat(90),
				},
				{
					High:  decimal.NewFromFloat(103),
					Close: decimal.NewFromFloat(92),
					Open:  decimal.NewFromFloat(101),
					Low:   decimal.NewFromFloat(85),
				},
			},
			Result: true,
		},
		"Failed to match bearish engulfing pattern with a negative first candle": {
			Pattern: CandlestickPatternBearishEngulfing,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(105),
					Close: decimal.NewFromFloat(95),
					Open:  decimal.NewFromFloat(100),
					Low:   decimal.NewFromFloat(90),
				},
				{
					High:  decimal.NewFromFloat(103),
					Close: decimal.NewFromFloat(92),
					Open:  decimal.NewFromFloat(101),
					Low:   decimal.NewFromFloat(85),
				},
			},
		},
		"Successfully evaluated bullish harami pattern": {
			Pattern: CandlestickPatternBullishHarami,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(112),
					Close: decimal.NewFromFloat(90),
					Open:  decimal.NewFromFloat(110),
					Low:   decimal.NewFromFloat(88),
				},
				{
					High:  decimal.NewFromFloat(106),
					Close: decimal.NewFromFloat(105),
					Open:  decimal.NewFromFloat(95),
					Low:   decimal.NewFromFloat(94),
				},
			},
			Result: true,
		},
		"Failed to match bullish harami pattern with the body outside the prior body": {
			Pattern: CandlestickPatternBullishHarami,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(112),
					Close: decimal.NewFromFloat(90),
					Open:  decimal.NewFromFloat(110),
					Low:   decimal.NewFromFloat(88),
				},
				{
					High:  decimal.NewFromFloat(106),
					Close: decimal.NewFromFloat(105),
					Open:  decimal.NewFromFloat(85),
					Low:   decimal.NewFromFloat(84),
				},
			},
		},
		"Successfully evaluated bearish harami pattern": {
			Pattern: CandlestickPatternBearishHarami,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(112),
					Close: decimal.NewFromFloat(110),
					Open:  decimal.NewFromFloat(90),
					Low:   decimal.NewFromFloat(88),
				},
				{
					High:  decimal.NewFromFloat(106),
					Close: decimal.NewFromFloat(95),
					Open:  decimal.NewFromFloat(105),
					Low:   decimal.NewFromFloat(94),
				},
			},
			Result: true,
		},
		"Successfully evaluated piercing line pattern": {
			Pattern: CandlestickPatternPiercingLine,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(102),
					Close: decimal.NewFromFloat(90),
					Open:  decimal.NewFromFloat(100),
					Low:   decimal.NewFromFloat(88),
				},
				{
					High:  decimal.NewFromFloat(97),
					Close: decimal.NewFromFloat(96),
					Open:  decimal.NewFromFloat(89),
					Low:   decimal.NewFromFloat(87),
				},
			},
			Result: true,
		},
		"Failed to match piercing line pattern below the prior body midpoint": {
			Pattern: CandlestickPatternPiercingLine,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(102),
					Close: decimal.NewFromFloat(90),
					Open:  decimal.NewFromFloat(100),
					Low:   decimal.NewFromFloat(88),
				},
				{
					High:  decimal.NewFromFloat(95),
					Close: decimal.NewFromFloat(94),
					Open:  decimal.NewFromFloat(89),
					Low:   decimal.NewFromFloat(87),
				},
			},
		},
		"Successfully evaluated dark cloud cover pattern": {
			Pattern: CandlestickPatternDarkCloudCover,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(102),
					Close: decimal.NewFromFloat(100),
					Open:  decimal.NewFromFloat(90),
					Low:   decimal.NewFromFloat(88),
				},
				{
					High:  decimal.NewFromFloat(103),
					Close: decimal.NewFromFloat(94),
					Open:  decimal.NewFromFloat(101),
					Low:   decimal.NewFromFloat(93),
				},
			},
			Result: true,
		},
		"Successfully evaluated tweezer top pattern": {
			Pattern: CandlestickPatternTweezerTop,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(105),
					Close: decimal.NewFromFloat(100),
					Open:  decimal.NewFromFloat(90),
					Low:   decimal.NewFromFloat(88),
				},
				{
					High:  decimal.NewFromFloat(105.5),
					Close: decimal.NewFromFloat(92),
					Open:  decimal.NewFromFloat(99),
					Low:   decimal.NewFromFloat(91),
				},
			},
			Result: true,
		},
		"Failed to match tweezer top pattern with diverging highs": {
			Pattern: CandlestickPatternTweezerTop,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(105),
					Close: decimal.NewFromFloat(100),
					Open:  decimal.NewFromFloat(90),
					Low:   decimal.NewFromFloat(88),
				},
				{
					High:  decimal.NewFromFloat(110),
					Close: decimal.NewFromFloat(92),
					Open:  decimal.NewFromFloat(99),
					Low:   decimal.NewFromFloat(91),
				},
			},
		},
		"Successfully evaluated tweezer bottom pattern": {
			Pattern: CandlestickPatternTweezerBottom,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(102),
					Close: decimal.NewFromFloat(92),
					Open:  decimal.NewFromFloat(100),
					Low:   decimal.NewFromFloat(90),
				},
				{
					High:  decimal.NewFromFloat(100),
					Close: decimal.NewFromFloat(99),
					Open:  decimal.NewFromFloat(93),
					Low:   decimal.NewFromFloat(90.5),
				},
			},
			Result: true,
		},
		"Successfully evaluated morning star pattern": {
			Pattern: CandlestickPatternMorningStar,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(102),
					Close: decimal.NewFromFloat(80),
					Open:  decimal.NewFromFloat(100),
					Low:   decimal.NewFromFloat(78),
				},
				{
					High:  decimal.NewFromFloat(82),
					Close: decimal.NewFromFloat(81),
					Open:  decimal.NewFromFloat(79),
					Low:   decimal.NewFromFloat(77),
				},
				{
					High:  decimal.NewFromFloat(96),
					Close: decimal.NewFromFloat(95),
					Open:  decimal.NewFromFloat(81),
					Low:   decimal.NewFromFloat(80),
				},
			},
			Result: true,
		},
		"Failed to match morning star pattern below the prior body midpoint": {
			Pattern: CandlestickPatternMorningStar,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(102),
					Close: decimal.NewFromFloat(80),
					Open:  decimal.NewFromFloat(100),
					Low:   decimal.NewFromFloat(78),
				},
				{
					High:  decimal.NewFromFloat(82),
					Close: decimal.NewFromFloat(81),
					Open:  decimal.NewFromFloat(79),
					Low:   decimal.NewFromFloat(77),
				},
				{
					High:  decimal.NewFromFloat(89),
					Close: decimal.NewFromFloat(88),
					Open:  decimal.NewFromFloat(81),
					Low:   decimal.NewFromFloat(80),
				},
			},
		},
		"Successfully evaluated evening star pattern": {
			Pattern: CandlestickPatternEveningStar,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(102),
					Close: decimal.NewFromFloat(100),
					Open:  decimal.NewFromFloat(80),
					Low:   decimal.NewFromFloat(78),
				},
				{
					High:  decimal.NewFromFloat(103),
					Close: decimal.NewFromFloat(102),
					Open:  decimal.NewFromFloat(101),
					Low:   decimal.NewFromFloat(99),
				},
				{
					High:  decimal.NewFromFloat(100),
					Close: decimal.NewFromFloat(85),
					Open:  decimal.NewFromFloat(99),
					Low:   decimal.NewFromFloat(84),
				},
			},
			Result: true,
		},
		"Successfully evaluated three white soldiers pattern": {
			Pattern: CandlestickPatternThreeWhiteSoldiers,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(21),
					Close: decimal.NewFromFloat(20),
					Open:  decimal.NewFromFloat(10),
					Low:   decimal.NewFromFloat(9),
				},
				{
					High:  decimal.NewFromFloat(26),
					Close: decimal.NewFromFloat(25),
					Open:  decimal.NewFromFloat(15),
					Low:   decimal.NewFromFloat(14),
				},
				{
					High:  decimal.NewFromFloat(31),
					Close: decimal.NewFromFloat(30),
					Open:  decimal.NewFromFloat(22),
					Low:   decimal.NewFromFloat(21),
				},
			},
			Result: true,
		},
		"Failed to match three white soldiers pattern with an open above the prior body": {
			Pattern: CandlestickPatternThreeWhiteSoldiers,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(21),
					Close: decimal.NewFromFloat(20),
					Open:  decimal.NewFromFloat(10),
					Low:   decimal.NewFromFloat(9),
				},
				{
					High:  decimal.NewFromFloat(26),
					Close: decimal.NewFromFloat(25),
					Open:  decimal.NewFromFloat(15),
					Low:   decimal.NewFromFloat(14),
				},
				{
					High:  decimal.NewFromFloat(31),
					Close: decimal.NewFromFloat(30),
					Open:  decimal.NewFromFloat(26),
					Low:   decimal.NewFromFloat(25),
				},
			},
		},
		"Successfully evaluated three black crows pattern": {
			Pattern: CandlestickPatternThreeBlackCrows,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(31),
					Close: decimal.NewFromFloat(20),
					Open:  decimal.NewFromFloat(30),
					Low:   decimal.NewFromFloat(19),
				},
				{
					High:  decimal.NewFromFloat(26),
					Close: decimal.NewFromFloat(15),
					Open:  decimal.NewFromFloat(25),
					Low:   decimal.NewFromFloat(14),
				},
				{
					High:  decimal.NewFromFloat(19),
					Close: decimal.NewFromFloat(10),
					Open:  decimal.NewFromFloat(18),
					Low:   decimal.NewFromFloat(9),
				},
			},
			Result: true,
		},
		"Failed to match three black crows pattern with a positive last candle": {
			Pattern: CandlestickPatternThreeBlackCrows,
			Candles: []Candle{
				{
					High:  decimal.NewFromFloat(31),
					Close: decimal.NewFromFloat(20),
					Open:  decimal.NewFromFloat(30),
					Low:   decimal.NewFromFloat(19),
				},
				{
					High:  decimal.NewFromFloat(26),
					Close: decimal.NewFromFloat(15),
					Open:  decimal.NewFromFloat(25),
					Low:   decimal.NewFromFloat(14),
				},
				{
					High:  decimal.NewFromFloat(19),
					Close: decimal.NewFromFloat(18),
					Open:  decimal.NewFromFloat(16),
					Low:   decimal.NewFromFloat(15),
				},
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, c.Result, c.Pattern.Eval(c.Candles))
		})
	}
}

func Test_Candle_IsBullish(t *testing.T) {
	assert.True(t, Candle{
		Open:  decimal.NewFromFloat(10),
		Close: decimal.NewFromFloat(20),
	}.IsBullish())
	assert.False(t, Candle{
		Open:  decimal.NewFromFloat(20),
		Close: decimal.NewFromFloat(10),
	}.IsBullish())
	assert.False(t, Candle{
		Open:  decimal.NewFromFloat(10),
		Close: decimal.NewFromFloat(10),
	}.IsBullish())
}

func Test_Candle_IsBearish(t *testing.T) {
	assert.True(t, Candle{
		Open:  decimal.NewFromFloat(20),
		Close: decimal.NewFromFloat(10),
	}.IsBearish())
	assert.False(t, Candle{
		Open:  decimal.NewFromFloat(10),
		Close: decimal.NewFromFloat(20),
	}.IsBearish())
	assert.False(t, Candle{
		Open:  decimal.NewFromFloat(10),
		Close: decimal.NewFromFloat(10),
	}.IsBearish())
}

func Test_Candle_BodyTop(t *testing.T) {
	assert.Equal(t, "20", Candle{
		Open:  decimal.NewFromFloat(20),
		Close: decimal.NewFromFloat(10),
	}.BodyTop().String())
	assert.Equal(t, "20", Candle{
		Open:  decimal.NewFromFloat(10),
		Close: decimal.NewFromFloat(20),
	}.BodyTop().String())
}

func Test_Candle_BodyBottom(t *testing.T) {
	assert.Equal(t, "10", Candle{
		Open:  decimal.NewFromFloat(20),
		Close: decimal.NewFromFloat(10),
	}.BodyBottom().String())
	assert.Equal(t, "10", Candle{
		Open:  decimal.NewFromFloat(10),
		Close: decimal.NewFromFloat(20),
	}.BodyBottom().String())
}

func Test_Candle_BodySize(t *testing.T) {
	assert.Equal(t, "10", Candle{
		Open:  decimal.NewFromFloat(20),
		Close: decimal.NewFromFloat(10),
	}.BodySize().String())
	assert.Equal(t, "10", Candle{
		Open:  decimal.NewFromFloat(10),
		Close: decimal.NewFromFloat(20),
	}.BodySize().String())
}

func Test_Candle_BodyMidpoint(t *testing.T) {
	assert.Equal(t, "15", Candle{
		Open:  decimal.NewFromFloat(20),
		Close: decimal.NewFromFloat(10),
	}.BodyMidpoint().String())
}
