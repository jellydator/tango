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
	}

	for _, pattern := range patterns {
		assert.NoError(t, pattern.Validate())
	}

	assert.Error(t, CandlestickPattern("invalid").Validate(), ErrInvalidCandlestickPattern)
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
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, c.Result, c.Pattern.Eval(c.Candles))
		})
	}
}
