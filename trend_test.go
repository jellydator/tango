package tango

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_NewADX(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result ADX
		Error  error
	}{
		"Validate returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new ADX": {
			Length: 5,
			Result: ADX{
				valid:  true,
				length: 5,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewADX(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_ADX_validate(t *testing.T) {
	cc := map[string]struct {
		ADX   ADX
		Error error
	}{
		"Invalid length": {
			ADX:   ADX{},
			Error: ErrInvalidLength,
		},
		"Successfully validated": {
			ADX: ADX{
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.ADX.validate())

			if c.Error == nil {
				assert.True(t, c.ADX.valid)
			}
		})
	}
}

func Test_ADX_Calc(t *testing.T) {
	trendingCandles := []Candle{
		{High: decimal.NewFromInt(10), Low: decimal.NewFromInt(5), Close: decimal.NewFromFloat(7.5)},
		{High: decimal.NewFromInt(12), Low: decimal.NewFromInt(7), Close: decimal.NewFromFloat(9.5)},
		{High: decimal.NewFromInt(14), Low: decimal.NewFromInt(9), Close: decimal.NewFromFloat(11.5)},
		{High: decimal.NewFromInt(16), Low: decimal.NewFromInt(11), Close: decimal.NewFromFloat(13.5)},
		{High: decimal.NewFromInt(18), Low: decimal.NewFromInt(13), Close: decimal.NewFromFloat(15.5)},
		{High: decimal.NewFromInt(20), Low: decimal.NewFromInt(15), Close: decimal.NewFromFloat(17.5)},
		{High: decimal.NewFromInt(22), Low: decimal.NewFromInt(17), Close: decimal.NewFromFloat(19.5)},
	}

	risingCandle := Candle{
		High:  decimal.NewFromInt(22),
		Low:   decimal.NewFromInt(17),
		Close: decimal.NewFromFloat(19.5),
	}

	fallingCandle := Candle{
		High:  decimal.NewFromInt(20),
		Low:   decimal.NewFromInt(15),
		Close: decimal.NewFromFloat(17.5),
	}

	flatCandle := Candle{
		High:  decimal.NewFromInt(10),
		Low:   decimal.NewFromInt(10),
		Close: decimal.NewFromInt(10),
	}

	cc := map[string]struct {
		ADX           ADX
		Candles       []Candle
		Result        decimal.Decimal
		PlusDIResult  decimal.Decimal
		MinusDIResult decimal.Decimal
		Error         error
	}{
		"Invalid indicator": {
			ADX: ADX{
				valid: false,
			},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			ADX: ADX{
				valid:  true,
				length: 2,
			},
			Candles: trendingCandles[:2],
			Error:   ErrInvalidDataSize,
		},
		"Successful calculation with a steady uptrend": {
			ADX: ADX{
				valid:  true,
				length: 2,
			},
			Candles:       trendingCandles,
			Result:        decimal.NewFromInt(100),
			PlusDIResult:  decimal.NewFromInt(40),
			MinusDIResult: decimal.NewFromInt(0),
		},
		"Successful calculation with an alternating market": {
			ADX: ADX{
				valid:  true,
				length: 2,
			},
			Candles: []Candle{
				fallingCandle,
				risingCandle,
				fallingCandle,
				risingCandle,
				fallingCandle,
				risingCandle,
				fallingCandle,
			},
			Result:        decimal.RequireFromString("31.25"),
			PlusDIResult:  decimal.RequireFromString("13.75"),
			MinusDIResult: decimal.RequireFromString("26.25"),
		},
		"Successful calculation with flat candles": {
			ADX: ADX{
				valid:  true,
				length: 1,
			},
			Candles: []Candle{
				flatCandle,
				flatCandle,
				flatCandle,
				flatCandle,
			},
			Result:        decimal.NewFromInt(0),
			PlusDIResult:  decimal.NewFromInt(0),
			MinusDIResult: decimal.NewFromInt(0),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			value, plusDI, minusDI, err := c.ADX.Calc(c.Candles)
			assertEqualError(t, c.Error, err)

			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), value.String())
			assert.Equal(t, c.PlusDIResult.String(), plusDI.String())
			assert.Equal(t, c.MinusDIResult.String(), minusDI.String())
		})
	}
}

func Test_ADX_Count(t *testing.T) {
	assert.Equal(t, 7, ADX{
		valid:  true,
		length: 2,
	}.Count())
}

func Test_NewSuperTrend(t *testing.T) {
	cc := map[string]struct {
		Multiplier decimal.Decimal
		Length     int
		Result     SuperTrend
		Error      error
	}{
		"Invalid length": {
			Multiplier: decimal.NewFromInt(1),
			Error:      ErrInvalidLength,
		},
		"Validate returns an error": {
			Length: 1,
			Error:  assert.AnError,
		},
		"Successfully created new SuperTrend": {
			Multiplier: decimal.NewFromInt(1),
			Length:     1,
			Result: SuperTrend{
				valid:      true,
				multiplier: decimal.NewFromInt(1),
				atr: ATR{
					valid:  true,
					length: 1,
				},
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewSuperTrend(c.Multiplier, c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_SuperTrend_validate(t *testing.T) {
	cc := map[string]struct {
		SuperTrend SuperTrend
		Error      error
	}{
		"Invalid multiplier": {
			SuperTrend: SuperTrend{},
			Error:      ErrInvalidMultiplier,
		},
		"Successfully validated": {
			SuperTrend: SuperTrend{
				multiplier: decimal.NewFromInt(1),
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.SuperTrend.validate())

			if c.Error == nil {
				assert.True(t, c.SuperTrend.valid)
			}
		})
	}
}

func Test_SuperTrend_Calc(t *testing.T) {
	cc := map[string]struct {
		SuperTrend SuperTrend
		Candles    []Candle
		Result     decimal.Decimal
		Trend      Trend
		Error      error
	}{
		"Invalid indicator": {
			SuperTrend: SuperTrend{
				valid: false,
			},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			SuperTrend: SuperTrend{
				valid:      true,
				multiplier: decimal.NewFromInt(1),
				atr: ATR{
					valid:  true,
					length: 1,
				},
			},
			Candles: []Candle{
				{
					High:  decimal.NewFromInt(11),
					Low:   decimal.NewFromInt(9),
					Close: decimal.NewFromInt(10),
				},
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation with an uptrend": {
			SuperTrend: SuperTrend{
				valid:      true,
				multiplier: decimal.NewFromInt(1),
				atr: ATR{
					valid:  true,
					length: 1,
				},
			},
			Candles: []Candle{
				{
					High:  decimal.NewFromInt(11),
					Low:   decimal.NewFromInt(9),
					Close: decimal.NewFromInt(10),
				},
				{
					High:  decimal.NewFromInt(12),
					Low:   decimal.NewFromInt(8),
					Close: decimal.NewFromInt(10),
				},
				{
					High:  decimal.NewFromInt(13),
					Low:   decimal.NewFromInt(7),
					Close: decimal.NewFromInt(11),
				},
			},
			Result: decimal.NewFromInt(4),
			Trend:  TrendUp,
		},
		"Successful calculation with a downtrend": {
			SuperTrend: SuperTrend{
				valid:      true,
				multiplier: decimal.NewFromInt(1),
				atr: ATR{
					valid:  true,
					length: 1,
				},
			},
			Candles: []Candle{
				{
					High:  decimal.NewFromInt(11),
					Low:   decimal.NewFromInt(9),
					Close: decimal.NewFromInt(10),
				},
				{
					High:  decimal.NewFromInt(12),
					Low:   decimal.NewFromInt(8),
					Close: decimal.NewFromInt(10),
				},
				{
					High:  decimal.NewFromInt(13),
					Low:   decimal.NewFromInt(7),
					Close: decimal.NewFromInt(9),
				},
			},
			Result: decimal.NewFromInt(16),
			Trend:  TrendDown,
		},
		"Successful calculation with a steady uptrend without band changes": {
			SuperTrend: SuperTrend{
				valid:      true,
				multiplier: decimal.NewFromInt(1),
				atr: ATR{
					valid:  true,
					length: 2,
				},
			},
			Candles: []Candle{
				{
					High:  decimal.NewFromInt(12),
					Low:   decimal.NewFromInt(8),
					Close: decimal.NewFromInt(11),
				},
				{
					High:  decimal.NewFromInt(12),
					Low:   decimal.NewFromInt(8),
					Close: decimal.NewFromInt(11),
				},
				{
					High:  decimal.NewFromInt(12),
					Low:   decimal.NewFromInt(8),
					Close: decimal.NewFromInt(11),
				},
				{
					High:  decimal.NewFromInt(12),
					Low:   decimal.NewFromInt(8),
					Close: decimal.NewFromInt(11),
				},
				{
					High:  decimal.NewFromInt(12),
					Low:   decimal.NewFromInt(8),
					Close: decimal.NewFromInt(11),
				},
				{
					High:  decimal.NewFromInt(12),
					Low:   decimal.NewFromInt(8),
					Close: decimal.NewFromInt(11),
				},
			},
			Result: decimal.NewFromInt(6),
			Trend:  TrendUp,
		},
		"Successful calculation with a reversal to a downtrend": {
			SuperTrend: SuperTrend{
				valid:      true,
				multiplier: decimal.NewFromInt(1),
				atr: ATR{
					valid:  true,
					length: 2,
				},
			},
			Candles: []Candle{
				{
					High:  decimal.NewFromInt(12),
					Low:   decimal.NewFromInt(8),
					Close: decimal.NewFromInt(10),
				},
				{
					High:  decimal.NewFromInt(12),
					Low:   decimal.NewFromInt(8),
					Close: decimal.NewFromInt(10),
				},
				{
					High:  decimal.NewFromInt(12),
					Low:   decimal.NewFromInt(8),
					Close: decimal.NewFromInt(10),
				},
				{
					High:  decimal.NewFromInt(12),
					Low:   decimal.NewFromInt(8),
					Close: decimal.NewFromInt(10),
				},
				{
					High:  decimal.NewFromInt(12),
					Low:   decimal.NewFromInt(8),
					Close: decimal.NewFromInt(11),
				},
				{
					High:  decimal.NewFromInt(4),
					Low:   decimal.NewFromInt(0),
					Close: decimal.NewFromInt(1),
				},
			},
			Result: decimal.RequireFromString("9.5"),
			Trend:  TrendDown,
		},
		"Successful calculation with a trend reversal and band carryover": {
			SuperTrend: SuperTrend{
				valid:      true,
				multiplier: decimal.NewFromInt(1),
				atr: ATR{
					valid:  true,
					length: 2,
				},
			},
			Candles: []Candle{
				{
					High:  decimal.NewFromInt(12),
					Low:   decimal.NewFromInt(8),
					Close: decimal.NewFromInt(10),
				},
				{
					High:  decimal.NewFromInt(12),
					Low:   decimal.NewFromInt(8),
					Close: decimal.NewFromInt(10),
				},
				{
					High:  decimal.NewFromInt(12),
					Low:   decimal.NewFromInt(8),
					Close: decimal.NewFromInt(10),
				},
				{
					High:  decimal.NewFromInt(12),
					Low:   decimal.NewFromInt(8),
					Close: decimal.NewFromInt(10),
				},
				{
					High:  decimal.NewFromInt(12),
					Low:   decimal.NewFromInt(8),
					Close: decimal.NewFromInt(10),
				},
				{
					High:  decimal.NewFromInt(20),
					Low:   decimal.NewFromInt(16),
					Close: decimal.NewFromInt(19),
				},
			},
			Result: decimal.NewFromInt(11),
			Trend:  TrendUp,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, trend, err := c.SuperTrend.Calc(c.Candles)
			assertEqualError(t, c.Error, err)

			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
			assert.Equal(t, c.Trend, trend)
		})
	}
}

func Test_SuperTrend_Count(t *testing.T) {
	assert.Equal(t, 6, SuperTrend{
		valid: true,
		atr: ATR{
			valid:  true,
			length: 2,
		},
	}.Count())
}

func Test_NewPSAR(t *testing.T) {
	cc := map[string]struct {
		Step   decimal.Decimal
		Limit  decimal.Decimal
		Length int
		Result PSAR
		Error  error
	}{
		"Validate returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new PSAR": {
			Step:   decimal.NewFromFloat(0.5),
			Limit:  decimal.NewFromInt(1),
			Length: 5,
			Result: PSAR{
				valid:  true,
				step:   decimal.NewFromFloat(0.5),
				limit:  decimal.NewFromInt(1),
				length: 5,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewPSAR(c.Step, c.Limit, c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_PSAR_validate(t *testing.T) {
	cc := map[string]struct {
		PSAR  PSAR
		Error error
	}{
		"Invalid length": {
			PSAR:  PSAR{},
			Error: ErrInvalidLength,
		},
		"Invalid zero step": {
			PSAR: PSAR{
				length: 2,
			},
			Error: ErrInvalidStep,
		},
		"Invalid step above the limit": {
			PSAR: PSAR{
				step:   decimal.NewFromInt(2),
				limit:  decimal.NewFromInt(1),
				length: 2,
			},
			Error: ErrInvalidStep,
		},
		"Successfully validated": {
			PSAR: PSAR{
				step:   decimal.NewFromFloat(0.5),
				limit:  decimal.NewFromInt(1),
				length: 2,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.PSAR.validate())

			if c.Error == nil {
				assert.True(t, c.PSAR.valid)
			}
		})
	}
}

func Test_PSAR_Calc(t *testing.T) {
	cc := map[string]struct {
		PSAR    PSAR
		Candles []Candle
		Result  decimal.Decimal
		Trend   Trend
		Error   error
	}{
		"Invalid indicator": {
			PSAR: PSAR{
				valid: false,
			},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			PSAR: PSAR{
				valid:  true,
				step:   decimal.NewFromFloat(0.5),
				limit:  decimal.NewFromInt(1),
				length: 4,
			},
			Candles: []Candle{
				{
					High:  decimal.NewFromInt(10),
					Low:   decimal.NewFromInt(5),
					Close: decimal.NewFromInt(7),
				},
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation with an uptrend": {
			PSAR: PSAR{
				valid:  true,
				step:   decimal.NewFromFloat(0.5),
				limit:  decimal.NewFromInt(1),
				length: 4,
			},
			Candles: []Candle{
				{
					High:  decimal.NewFromInt(10),
					Low:   decimal.NewFromInt(5),
					Close: decimal.NewFromInt(7),
				},
				{
					High:  decimal.NewFromInt(12),
					Low:   decimal.NewFromInt(7),
					Close: decimal.NewFromInt(11),
				},
				{
					High:  decimal.NewFromInt(14),
					Low:   decimal.NewFromInt(9),
					Close: decimal.NewFromInt(13),
				},
				{
					High:  decimal.NewFromInt(13),
					Low:   decimal.NewFromInt(8),
					Close: decimal.NewFromInt(12),
				},
			},
			Result: decimal.NewFromInt(7),
			Trend:  TrendUp,
		},
		"Successful calculation with a trend reversal": {
			PSAR: PSAR{
				valid:  true,
				step:   decimal.NewFromFloat(0.5),
				limit:  decimal.NewFromInt(1),
				length: 4,
			},
			Candles: []Candle{
				{
					High:  decimal.NewFromInt(10),
					Low:   decimal.NewFromInt(5),
					Close: decimal.NewFromInt(7),
				},
				{
					High:  decimal.NewFromInt(12),
					Low:   decimal.NewFromInt(7),
					Close: decimal.NewFromInt(11),
				},
				{
					High:  decimal.NewFromInt(14),
					Low:   decimal.NewFromInt(9),
					Close: decimal.NewFromInt(13),
				},
				{
					High:  decimal.NewFromInt(13),
					Low:   decimal.NewFromInt(6),
					Close: decimal.NewFromFloat(6.5),
				},
			},
			Result: decimal.NewFromInt(14),
			Trend:  TrendDown,
		},
		"Successful calculation with a downtrend": {
			PSAR: PSAR{
				valid:  true,
				step:   decimal.NewFromFloat(0.5),
				limit:  decimal.NewFromInt(1),
				length: 2,
			},
			Candles: []Candle{
				{
					High:  decimal.NewFromInt(10),
					Low:   decimal.NewFromInt(5),
					Close: decimal.NewFromInt(7),
				},
				{
					High:  decimal.NewFromInt(8),
					Low:   decimal.NewFromInt(4),
					Close: decimal.NewFromInt(5),
				},
			},
			Result: decimal.NewFromInt(10),
			Trend:  TrendDown,
		},
		"Successful calculation with a downtrend continuation": {
			PSAR: PSAR{
				valid:  true,
				step:   decimal.NewFromFloat(0.5),
				limit:  decimal.NewFromInt(1),
				length: 4,
			},
			Candles: []Candle{
				{
					High:  decimal.NewFromInt(20),
					Low:   decimal.NewFromInt(15),
					Close: decimal.NewFromInt(17),
				},
				{
					High:  decimal.NewFromInt(18),
					Low:   decimal.NewFromInt(13),
					Close: decimal.NewFromInt(14),
				},
				{
					High:  decimal.NewFromInt(16),
					Low:   decimal.NewFromInt(11),
					Close: decimal.NewFromInt(12),
				},
				{
					High:  decimal.NewFromInt(15),
					Low:   decimal.NewFromInt(10),
					Close: decimal.NewFromInt(11),
				},
			},
			Result: decimal.NewFromInt(18),
			Trend:  TrendDown,
		},
		"Successful calculation with a downtrend reversal": {
			PSAR: PSAR{
				valid:  true,
				step:   decimal.NewFromFloat(0.5),
				limit:  decimal.NewFromInt(1),
				length: 4,
			},
			Candles: []Candle{
				{
					High:  decimal.NewFromInt(20),
					Low:   decimal.NewFromInt(15),
					Close: decimal.NewFromInt(17),
				},
				{
					High:  decimal.NewFromInt(18),
					Low:   decimal.NewFromInt(13),
					Close: decimal.NewFromInt(14),
				},
				{
					High:  decimal.NewFromInt(16),
					Low:   decimal.NewFromInt(11),
					Close: decimal.NewFromInt(12),
				},
				{
					High:  decimal.NewFromInt(19),
					Low:   decimal.NewFromInt(14),
					Close: decimal.NewFromInt(18),
				},
			},
			Result: decimal.NewFromInt(11),
			Trend:  TrendUp,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, trend, err := c.PSAR.Calc(c.Candles)
			assertEqualError(t, c.Error, err)

			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
			assert.Equal(t, c.Trend, trend)
		})
	}
}

func Test_PSAR_Count(t *testing.T) {
	assert.Equal(t, 4, PSAR{
		valid:  true,
		length: 4,
	}.Count())
}

func Test_NewIchimoku(t *testing.T) {
	cc := map[string]struct {
		Conversion int
		Base       int
		Span       int
		Result     Ichimoku
		Error      error
	}{
		"Validate returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new Ichimoku": {
			Conversion: 9,
			Base:       26,
			Span:       52,
			Result: Ichimoku{
				valid:      true,
				conversion: 9,
				base:       26,
				span:       52,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewIchimoku(c.Conversion, c.Base, c.Span)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_Ichimoku_validate(t *testing.T) {
	cc := map[string]struct {
		Ichimoku Ichimoku
		Error    error
	}{
		"Invalid conversion length": {
			Ichimoku: Ichimoku{
				base: 1,
				span: 1,
			},
			Error: ErrInvalidLength,
		},
		"Invalid base length": {
			Ichimoku: Ichimoku{
				conversion: 1,
				span:       1,
			},
			Error: ErrInvalidLength,
		},
		"Invalid span length": {
			Ichimoku: Ichimoku{
				conversion: 1,
				base:       1,
			},
			Error: ErrInvalidLength,
		},
		"Successfully validated": {
			Ichimoku: Ichimoku{
				conversion: 1,
				base:       1,
				span:       1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.Ichimoku.validate())

			if c.Error == nil {
				assert.True(t, c.Ichimoku.valid)
			}
		})
	}
}

func Test_Ichimoku_Calc(t *testing.T) {
	cc := map[string]struct {
		Ichimoku         Ichimoku
		Candles          []Candle
		ConversionResult decimal.Decimal
		BaseResult       decimal.Decimal
		SpanAResult      decimal.Decimal
		SpanBResult      decimal.Decimal
		Error            error
	}{
		"Invalid indicator": {
			Ichimoku: Ichimoku{
				valid: false,
			},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			Ichimoku: Ichimoku{
				valid:      true,
				conversion: 2,
				base:       3,
				span:       4,
			},
			Candles: []Candle{
				{
					High: decimal.NewFromInt(10),
					Low:  decimal.NewFromInt(5),
				},
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
			Ichimoku: Ichimoku{
				valid:      true,
				conversion: 2,
				base:       3,
				span:       4,
			},
			Candles: []Candle{
				{
					High: decimal.NewFromInt(20),
					Low:  decimal.NewFromInt(2),
				},
				{
					High: decimal.NewFromInt(8),
					Low:  decimal.NewFromInt(4),
				},
				{
					High: decimal.NewFromInt(10),
					Low:  decimal.NewFromInt(6),
				},
				{
					High: decimal.NewFromInt(9),
					Low:  decimal.NewFromInt(5),
				},
			},
			ConversionResult: decimal.RequireFromString("7.5"),
			BaseResult:       decimal.NewFromInt(7),
			SpanAResult:      decimal.RequireFromString("7.25"),
			SpanBResult:      decimal.NewFromInt(11),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			conversion, base, spanA, spanB, err := c.Ichimoku.Calc(c.Candles)
			assertEqualError(t, c.Error, err)

			if err != nil {
				return
			}

			assert.Equal(t, c.ConversionResult.String(), conversion.String())
			assert.Equal(t, c.BaseResult.String(), base.String())
			assert.Equal(t, c.SpanAResult.String(), spanA.String())
			assert.Equal(t, c.SpanBResult.String(), spanB.String())
		})
	}
}

func Test_Ichimoku_Count(t *testing.T) {
	assert.Equal(t, 52, Ichimoku{
		valid:      true,
		conversion: 9,
		base:       26,
		span:       52,
	}.Count())
}
