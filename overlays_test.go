package tango

import (
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_NewBB(t *testing.T) {
	cc := map[string]struct {
		MAType MAType
		StdDev decimal.Decimal
		Length int
		Result BB
		Error  error
	}{
		"Invalid moving average": {
			Error: ErrInvalidMA,
		},
		"Validate returns an error": {
			MAType: MATypeSimple,
			Length: 1,
			Error:  errors.New("invalid standard deviation"),
		},
		"Successfully created new BB": {
			MAType: MATypeSimple,
			StdDev: decimal.RequireFromString("2.5"),
			Length: 5,
			Result: BB{
				valid:  true,
				stdDev: decimal.RequireFromString("2.5"),
				ma: SMA{
					length: 5,
					valid:  true,
				},
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewBB(c.MAType, c.StdDev, c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_BB_validate(t *testing.T) {
	cc := map[string]struct {
		BB    BB
		Error error
	}{
		"Invalid standard deviation": {
			BB: BB{
				stdDev: decimal.Decimal{},
				ma: SMA{
					length: 5,
				},
			},
			Error: errors.New("invalid standard deviation"),
		},
		"Successfully validated": {
			BB: BB{
				stdDev: decimal.NewFromInt(5),
				ma: SMA{
					length: 1,
				},
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.BB.validate())
			if c.Error == nil {
				assert.True(t, c.BB.valid)
			}
		})
	}
}

func Test_BB_Calc(t *testing.T) {
	cc := map[string]struct {
		BB          BB
		Data        []decimal.Decimal
		UpperResult decimal.Decimal
		LowerResult decimal.Decimal
		WidthResult decimal.Decimal
		Error       error
	}{
		"Invalid indicator": {
			BB: BB{
				valid: false,
			},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			BB: BB{
				valid: true,
				ma: SMA{
					valid:  true,
					length: 5,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
			BB: BB{
				valid:  true,
				stdDev: decimal.RequireFromString("2"),
				ma: SMA{
					length: 20,
					valid:  true,
				},
			},
			Data: []decimal.Decimal{
				decimal.RequireFromString("63.98"),
				decimal.RequireFromString("64.17"),
				decimal.RequireFromString("64.71"),
				decimal.RequireFromString("64.75"),
				decimal.RequireFromString("63.94"),
				decimal.RequireFromString("63.82"),
				decimal.RequireFromString("63.19"),
				decimal.RequireFromString("62.84"),
				decimal.RequireFromString("62.25"),
				decimal.RequireFromString("63.20"),
				decimal.RequireFromString("63.02"),
				decimal.RequireFromString("63.35"),
				decimal.RequireFromString("64.21"),
				decimal.RequireFromString("64.91"),
				decimal.RequireFromString("64.05"),
				decimal.RequireFromString("63.28"),
				decimal.RequireFromString("62.78"),
				decimal.RequireFromString("62.36"),
				decimal.RequireFromString("63.19"),
				decimal.RequireFromString("64.69"),
			},
			UpperResult: decimal.RequireFromString("65.19977921"),
			LowerResult: decimal.RequireFromString("62.06922079"),
			WidthResult: decimal.RequireFromString("4.91959301"),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			upperBand, lowerBand, widthBand, err := c.BB.Calc(c.Data)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.UpperResult.Round(8).String(), upperBand.Round(8).String())
			assert.Equal(t, c.LowerResult.Round(8).String(), lowerBand.Round(8).String())
			assert.Equal(t, c.WidthResult.Round(8).String(), widthBand.Round(8).String())
		})
	}
}

func Test_BB_CalcBand(t *testing.T) {
	cc := map[string]struct {
		BB     BB
		Band   Band
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid band": {
			BB: BB{
				valid: false,
			},
			Error: ErrInvalidBand,
		},
		"Invalid indicator": {
			BB: BB{
				valid: false,
			},
			Band:  BandUpper,
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			BB: BB{
				valid: true,
				ma: SMA{
					valid:  true,
					length: 5,
				},
			},
			Band: BandUpper,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation with BandUpper": {
			BB: BB{
				valid:  true,
				stdDev: decimal.RequireFromString("1"),
				ma: SMA{
					length: 5,
					valid:  true,
				},
			},
			Band: BandUpper,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
				decimal.NewFromInt(35),
				decimal.NewFromInt(40),
				decimal.NewFromInt(38),
				decimal.NewFromInt(32),
			},
			Result: decimal.RequireFromString("38.68781778"),
		},
		"Successful calculation with BandLower": {
			BB: BB{
				valid:  true,
				stdDev: decimal.RequireFromString("1"),
				ma: SMA{
					length: 5,
					valid:  true,
				},
			},
			Band: BandLower,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
				decimal.NewFromInt(35),
				decimal.NewFromInt(40),
				decimal.NewFromInt(38),
				decimal.NewFromInt(32),
			},
			Result: decimal.RequireFromString("31.31218222"),
		},
		"Successful calculation with BandWidth": {
			BB: BB{
				valid:  true,
				stdDev: decimal.RequireFromString("2"),
				ma: SMA{
					length: 20,
					valid:  true,
				},
			},
			Band: BandWidth,
			Data: []decimal.Decimal{
				decimal.RequireFromString("63.98"),
				decimal.RequireFromString("64.17"),
				decimal.RequireFromString("64.71"),
				decimal.RequireFromString("64.75"),
				decimal.RequireFromString("63.94"),
				decimal.RequireFromString("63.82"),
				decimal.RequireFromString("63.19"),
				decimal.RequireFromString("62.84"),
				decimal.RequireFromString("62.25"),
				decimal.RequireFromString("63.20"),
				decimal.RequireFromString("63.02"),
				decimal.RequireFromString("63.35"),
				decimal.RequireFromString("64.21"),
				decimal.RequireFromString("64.91"),
				decimal.RequireFromString("64.05"),
				decimal.RequireFromString("63.28"),
				decimal.RequireFromString("62.78"),
				decimal.RequireFromString("62.36"),
				decimal.RequireFromString("63.19"),
				decimal.RequireFromString("64.69"),
			},
			Result: decimal.RequireFromString("4.91959301"),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.BB.CalcBand(c.Data, c.Band)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.Round(8).String(), res.Round(8).String())
		})
	}
}

func Test_BB_Count(t *testing.T) {
	assert.Equal(t, 1, BB{ma: SMA{length: 1}}.Count())
}

func Test_NewDEMA(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result DEMA
		Error  error
	}{
		"NewEMA returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new DEMA": {
			Length: 1,
			Result: DEMA{
				valid: true,
				ema: EMA{
					sma: SMA{
						length: 1,
						valid:  true,
					},
					valid: true,
				},
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewDEMA(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_DEMA_Calc(t *testing.T) {
	cc := map[string]struct {
		DEMA   DEMA
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			DEMA:  DEMA{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			DEMA: DEMA{
				valid: true,
				ema: EMA{
					sma: SMA{
						length: 3,
						valid:  true,
					},
					valid: true,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
			DEMA: DEMA{
				valid: true,
				ema: EMA{
					sma: SMA{
						length: 3,
						valid:  true,
					},
					valid: true,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(31),
				decimal.NewFromInt(1),
				decimal.NewFromInt(1),
				decimal.NewFromInt(2),
				decimal.NewFromInt(3),
			},
			Result: decimal.RequireFromString("6.75"),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.DEMA.Calc(c.Data)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_DEMA_Count(t *testing.T) {
	assert.Equal(t, 29, DEMA{
		valid: false,
		ema: EMA{
			sma: SMA{
				length: 15,
			},
		},
	}.Count())
}

func Test_NewEMA(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result EMA
		Error  error
	}{
		"Invalid parameters": {
			Error: assert.AnError,
		},
		"Successfully created new EMA": {
			Length: 1,
			Result: EMA{
				valid: true,
				sma: SMA{
					length: 1,
					valid:  true,
				},
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewEMA(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_EMA_Calc(t *testing.T) {
	cc := map[string]struct {
		EMA    EMA
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			EMA:   EMA{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			EMA: EMA{
				valid: true,
				sma: SMA{
					length: 3,
					valid:  true,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
			EMA: EMA{
				valid: true,
				sma: SMA{
					length: 3,
					valid:  true,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(31),
				decimal.NewFromInt(1),
				decimal.NewFromInt(1),
				decimal.NewFromInt(2),
				decimal.NewFromInt(3),
			},
			Result: decimal.RequireFromString("4.75"),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.EMA.Calc(c.Data)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_EMA_CalcNext(t *testing.T) {
	cc := map[string]struct {
		EMA    EMA
		Last   decimal.Decimal
		Next   decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			EMA:   EMA{},
			Error: ErrInvalidIndicator,
		},
		"Successful calculation": {
			EMA: EMA{
				valid: true,
				sma: SMA{
					length: 3,
					valid:  true,
				},
			},
			Last:   decimal.NewFromInt(5),
			Next:   decimal.NewFromInt(5),
			Result: decimal.NewFromInt(5),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.EMA.CalcNext(c.Last, c.Next)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_EMA_Count(t *testing.T) {
	assert.Equal(t, 29, EMA{
		sma: SMA{
			length: 15,
		},
	}.Count())
}

func Test_EMA_multiplier(t *testing.T) {
	assert.Equal(t, decimal.RequireFromString("0.5").String(), EMA{
		sma: SMA{
			length: 3,
		},
	}.multiplier().String())
}

func Test_NewHMA(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result HMA
		Error  error
	}{
		"NewWMA returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new HMA": {
			Length: 2,
			Result: HMA{
				valid: true,
				wma: WMA{
					length: 2,
					valid:  true,
				},
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewHMA(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_HMA_Calc(t *testing.T) {
	cc := map[string]struct {
		HMA    HMA
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			HMA:   HMA{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			HMA: HMA{
				valid: true,
				wma: WMA{
					length: 5,
					valid:  true,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
			HMA: HMA{
				valid: true,
				wma: WMA{
					length: 4,
					valid:  true,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(32),
				decimal.NewFromInt(29),
				decimal.NewFromInt(38),
				decimal.NewFromInt(34),
				decimal.NewFromInt(29),
			},
			Result: decimal.RequireFromString("33.8"),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.HMA.Calc(c.Data)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.Round(8).String(), res.Round(8).String())
		})
	}
}

func Test_HMA_Count(t *testing.T) {
	assert.Equal(t, 17, HMA{
		wma: WMA{
			length: 15,
		},
	}.Count())
}

func Test_NewSMA(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result SMA
		Error  error
	}{
		"Validate returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new SMA": {
			Length: 1,
			Result: SMA{
				valid:  true,
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewSMA(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_SMA_validate(t *testing.T) {
	cc := map[string]struct {
		SMA   SMA
		Error error
	}{
		"Invalid length": {
			SMA: SMA{
				length: 0,
			},
			Error: ErrInvalidLength,
		},
		"Successfully validated": {
			SMA: SMA{
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.SMA.validate())
			if c.Error == nil {
				assert.True(t, c.SMA.valid)
			}
		})
	}
}

func Test_SMA_Calc(t *testing.T) {
	cc := map[string]struct {
		SMA    SMA
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			SMA:   SMA{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			SMA: SMA{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
			SMA: SMA{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
				decimal.NewFromInt(31),
				decimal.NewFromInt(32),
			},
			Result: decimal.NewFromInt(31),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.SMA.Calc(c.Data)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_SMA_Count(t *testing.T) {
	assert.Equal(t, 15, SMA{
		length: 15,
	}.Count())
}

func Test_NewWMA(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result WMA
		Error  error
	}{
		"Validate returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new WMA": {
			Length: 1,
			Result: WMA{
				valid:  true,
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewWMA(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_WMA_validate(t *testing.T) {
	cc := map[string]struct {
		WMA   WMA
		Error error
	}{
		"Invalid length": {
			WMA: WMA{
				length: 0,
			},
			Error: ErrInvalidLength,
		},
		"Successfully validated": {
			WMA: WMA{
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.WMA.validate())
			if c.Error == nil {
				assert.True(t, c.WMA.valid)
			}
		})
	}
}

func Test_WMA_Calc(t *testing.T) {
	cc := map[string]struct {
		WMA    WMA
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			WMA:   WMA{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			WMA: WMA{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
			WMA: WMA{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
				decimal.NewFromInt(30),
				decimal.NewFromInt(32),
			},
			Result: decimal.NewFromInt(31),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.WMA.Calc(c.Data)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_WMA_Count(t *testing.T) {
	assert.Equal(t, 15, WMA{
		length: 15,
	}.Count())
}
