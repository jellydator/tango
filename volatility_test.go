package tango

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_NewATR(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result ATR
		Error  error
	}{
		"Validate returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new ATR": {
			Length: 5,
			Result: ATR{
				valid:  true,
				length: 5,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewATR(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_ATR_validate(t *testing.T) {
	cc := map[string]struct {
		ATR   ATR
		Error error
	}{
		"Invalid length": {
			ATR:   ATR{},
			Error: ErrInvalidLength,
		},
		"Successfully validated": {
			ATR: ATR{
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.ATR.validate())

			if c.Error == nil {
				assert.True(t, c.ATR.valid)
			}
		})
	}
}

func Test_ATR_Calc(t *testing.T) {
	cc := map[string]struct {
		ATR     ATR
		Candles []Candle
		Result  decimal.Decimal
		Error   error
	}{
		"Invalid indicator": {
			ATR: ATR{
				valid: false,
			},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			ATR: ATR{
				valid:  true,
				length: 2,
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
		"Successful calculation": {
			ATR: ATR{
				valid:  true,
				length: 2,
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
					High:  decimal.NewFromInt(14),
					Low:   decimal.NewFromInt(6),
					Close: decimal.NewFromInt(10),
				},
				{
					High:  decimal.NewFromInt(13),
					Low:   decimal.NewFromInt(7),
					Close: decimal.NewFromInt(10),
				},
				{
					High:  decimal.NewFromInt(15),
					Low:   decimal.NewFromInt(5),
					Close: decimal.NewFromInt(10),
				},
			},
			Result: decimal.NewFromInt(8),
		},
		"Successful calculation with gapped candles": {
			ATR: ATR{
				valid:  true,
				length: 1,
			},
			Candles: []Candle{
				{
					High:  decimal.NewFromInt(21),
					Low:   decimal.NewFromInt(19),
					Close: decimal.NewFromInt(20),
				},
				{
					High:  decimal.NewFromInt(15),
					Low:   decimal.NewFromInt(12),
					Close: decimal.NewFromInt(13),
				},
				{
					High:  decimal.NewFromInt(30),
					Low:   decimal.NewFromInt(14),
					Close: decimal.NewFromInt(25),
				},
			},
			Result: decimal.NewFromInt(17),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.ATR.Calc(c.Candles)
			assertEqualError(t, c.Error, err)

			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_ATR_Count(t *testing.T) {
	assert.Equal(t, 5, ATR{
		valid:  true,
		length: 2,
	}.Count())
}

func Test_NewDonchian(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result Donchian
		Error  error
	}{
		"Validate returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new Donchian": {
			Length: 5,
			Result: Donchian{
				valid:  true,
				length: 5,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewDonchian(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_Donchian_validate(t *testing.T) {
	cc := map[string]struct {
		Donchian Donchian
		Error    error
	}{
		"Invalid length": {
			Donchian: Donchian{},
			Error:    ErrInvalidLength,
		},
		"Successfully validated": {
			Donchian: Donchian{
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.Donchian.validate())

			if c.Error == nil {
				assert.True(t, c.Donchian.valid)
			}
		})
	}
}

func Test_Donchian_Calc(t *testing.T) {
	cc := map[string]struct {
		Donchian     Donchian
		Candles      []Candle
		UpperResult  decimal.Decimal
		MiddleResult decimal.Decimal
		LowerResult  decimal.Decimal
		Error        error
	}{
		"Invalid indicator": {
			Donchian: Donchian{
				valid: false,
			},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			Donchian: Donchian{
				valid:  true,
				length: 3,
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
			Donchian: Donchian{
				valid:  true,
				length: 3,
			},
			Candles: []Candle{
				{
					High: decimal.NewFromInt(10),
					Low:  decimal.NewFromInt(5),
				},
				{
					High: decimal.NewFromInt(12),
					Low:  decimal.NewFromInt(7),
				},
				{
					High: decimal.NewFromInt(11),
					Low:  decimal.NewFromInt(4),
				},
			},
			UpperResult:  decimal.NewFromInt(12),
			MiddleResult: decimal.NewFromInt(8),
			LowerResult:  decimal.NewFromInt(4),
		},
		"Successful calculation with the first candle extremes": {
			Donchian: Donchian{
				valid:  true,
				length: 2,
			},
			Candles: []Candle{
				{
					High: decimal.NewFromInt(15),
					Low:  decimal.NewFromInt(3),
				},
				{
					High: decimal.NewFromInt(12),
					Low:  decimal.NewFromInt(7),
				},
			},
			UpperResult:  decimal.NewFromInt(15),
			MiddleResult: decimal.NewFromInt(9),
			LowerResult:  decimal.NewFromInt(3),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			upper, middle, lower, err := c.Donchian.Calc(c.Candles)
			assertEqualError(t, c.Error, err)

			if err != nil {
				return
			}

			assert.Equal(t, c.UpperResult.String(), upper.String())
			assert.Equal(t, c.MiddleResult.String(), middle.String())
			assert.Equal(t, c.LowerResult.String(), lower.String())
		})
	}
}

func Test_Donchian_Count(t *testing.T) {
	assert.Equal(t, 3, Donchian{
		valid:  true,
		length: 3,
	}.Count())
}

func Test_NewKeltner(t *testing.T) {
	cc := map[string]struct {
		Type       MAType
		Multiplier decimal.Decimal
		Length     int
		Result     Keltner
		Error      error
	}{
		"Invalid moving average type": {
			Multiplier: decimal.NewFromInt(2),
			Length:     1,
			Error:      ErrInvalidMA,
		},
		"Invalid length": {
			Type:       MATypeSimple,
			Multiplier: decimal.NewFromInt(2),
			Error:      ErrInvalidLength,
		},
		"Validate returns an error": {
			Type:   MATypeSimple,
			Length: 1,
			Error:  assert.AnError,
		},
		"Successfully created new Keltner": {
			Type:       MATypeSimple,
			Multiplier: decimal.NewFromInt(2),
			Length:     1,
			Result: Keltner{
				valid:      true,
				multiplier: decimal.NewFromInt(2),
				ma: SMA{
					valid:  true,
					length: 1,
				},
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

			res, err := NewKeltner(c.Type, c.Multiplier, c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_Keltner_validate(t *testing.T) {
	cc := map[string]struct {
		Keltner Keltner
		Error   error
	}{
		"Invalid multiplier": {
			Keltner: Keltner{},
			Error:   ErrInvalidMultiplier,
		},
		"Successfully validated": {
			Keltner: Keltner{
				multiplier: decimal.NewFromInt(2),
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.Keltner.validate())

			if c.Error == nil {
				assert.True(t, c.Keltner.valid)
			}
		})
	}
}

func Test_Keltner_Calc(t *testing.T) {
	cc := map[string]struct {
		Keltner      Keltner
		Candles      []Candle
		UpperResult  decimal.Decimal
		MiddleResult decimal.Decimal
		LowerResult  decimal.Decimal
		Error        error
	}{
		"Invalid indicator": {
			Keltner: Keltner{
				valid: false,
			},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			Keltner: Keltner{
				valid:      true,
				multiplier: decimal.NewFromInt(2),
				ma: SMA{
					valid:  true,
					length: 1,
				},
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
		"Successful calculation": {
			Keltner: Keltner{
				valid:      true,
				multiplier: decimal.NewFromInt(2),
				ma: SMA{
					valid:  true,
					length: 1,
				},
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
			UpperResult:  decimal.NewFromInt(21),
			MiddleResult: decimal.NewFromInt(9),
			LowerResult:  decimal.NewFromInt(-3),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			upper, middle, lower, err := c.Keltner.Calc(c.Candles)
			assertEqualError(t, c.Error, err)

			if err != nil {
				return
			}

			assert.Equal(t, c.UpperResult.String(), upper.String())
			assert.Equal(t, c.MiddleResult.String(), middle.String())
			assert.Equal(t, c.LowerResult.String(), lower.String())
		})
	}
}

func Test_Keltner_Count(t *testing.T) {
	assert.Equal(t, 21, Keltner{
		valid: true,
		ma: SMA{
			valid:  true,
			length: 10,
		},
		atr: ATR{
			valid:  true,
			length: 10,
		},
	}.Count())
}
