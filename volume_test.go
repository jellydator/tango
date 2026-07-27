package tango

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_NewOBV(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result OBV
		Error  error
	}{
		"Validate returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new OBV": {
			Length: 5,
			Result: OBV{
				valid:  true,
				length: 5,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewOBV(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_OBV_validate(t *testing.T) {
	cc := map[string]struct {
		OBV   OBV
		Error error
	}{
		"Invalid length": {
			OBV:   OBV{},
			Error: ErrInvalidLength,
		},
		"Successfully validated": {
			OBV: OBV{
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.OBV.validate())

			if c.Error == nil {
				assert.True(t, c.OBV.valid)
			}
		})
	}
}

func Test_OBV_Calc(t *testing.T) {
	cc := map[string]struct {
		OBV     OBV
		Candles []Candle
		Result  decimal.Decimal
		Error   error
	}{
		"Invalid indicator": {
			OBV: OBV{
				valid: false,
			},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			OBV: OBV{
				valid:  true,
				length: 3,
			},
			Candles: []Candle{
				{
					Close:  decimal.NewFromInt(10),
					Volume: decimal.NewFromInt(5),
				},
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
			OBV: OBV{
				valid:  true,
				length: 3,
			},
			Candles: []Candle{
				{
					Close:  decimal.NewFromInt(10),
					Volume: decimal.NewFromInt(9),
				},
				{
					Close:  decimal.NewFromInt(12),
					Volume: decimal.NewFromInt(5),
				},
				{
					Close:  decimal.NewFromInt(11),
					Volume: decimal.NewFromInt(3),
				},
				{
					Close:  decimal.NewFromInt(11),
					Volume: decimal.NewFromInt(7),
				},
			},
			Result: decimal.NewFromInt(2),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.OBV.Calc(c.Candles)
			assertEqualError(t, c.Error, err)

			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_OBV_Count(t *testing.T) {
	assert.Equal(t, 4, OBV{
		valid:  true,
		length: 3,
	}.Count())
}

func Test_NewMFI(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result MFI
		Error  error
	}{
		"Validate returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new MFI": {
			Length: 5,
			Result: MFI{
				valid:  true,
				length: 5,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewMFI(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_MFI_validate(t *testing.T) {
	cc := map[string]struct {
		MFI   MFI
		Error error
	}{
		"Invalid length": {
			MFI:   MFI{},
			Error: ErrInvalidLength,
		},
		"Successfully validated": {
			MFI: MFI{
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.MFI.validate())

			if c.Error == nil {
				assert.True(t, c.MFI.valid)
			}
		})
	}
}

func Test_MFI_Calc(t *testing.T) {
	cc := map[string]struct {
		MFI     MFI
		Candles []Candle
		Result  decimal.Decimal
		Error   error
	}{
		"Invalid indicator": {
			MFI: MFI{
				valid: false,
			},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			MFI: MFI{
				valid:  true,
				length: 3,
			},
			Candles: []Candle{
				{
					High:   decimal.NewFromInt(12),
					Low:    decimal.NewFromInt(6),
					Close:  decimal.NewFromInt(9),
					Volume: decimal.NewFromInt(10),
				},
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
			MFI: MFI{
				valid:  true,
				length: 3,
			},
			Candles: []Candle{
				{
					High:   decimal.NewFromInt(12),
					Low:    decimal.NewFromInt(6),
					Close:  decimal.NewFromInt(9),
					Volume: decimal.NewFromInt(10),
				},
				{
					High:   decimal.NewFromInt(15),
					Low:    decimal.NewFromInt(9),
					Close:  decimal.NewFromInt(12),
					Volume: decimal.NewFromInt(15),
				},
				{
					High:   decimal.NewFromInt(9),
					Low:    decimal.NewFromInt(3),
					Close:  decimal.NewFromInt(6),
					Volume: decimal.NewFromInt(10),
				},
				{
					High:   decimal.NewFromInt(9),
					Low:    decimal.NewFromInt(3),
					Close:  decimal.NewFromInt(6),
					Volume: decimal.NewFromInt(99),
				},
			},
			Result: decimal.NewFromInt(75),
		},
		"Successful calculation without negative money flow": {
			MFI: MFI{
				valid:  true,
				length: 1,
			},
			Candles: []Candle{
				{
					High:   decimal.NewFromInt(12),
					Low:    decimal.NewFromInt(6),
					Close:  decimal.NewFromInt(9),
					Volume: decimal.NewFromInt(10),
				},
				{
					High:   decimal.NewFromInt(15),
					Low:    decimal.NewFromInt(9),
					Close:  decimal.NewFromInt(12),
					Volume: decimal.NewFromInt(15),
				},
			},
			Result: decimal.NewFromInt(100),
		},
		"Successful calculation without positive money flow": {
			MFI: MFI{
				valid:  true,
				length: 1,
			},
			Candles: []Candle{
				{
					High:   decimal.NewFromInt(12),
					Low:    decimal.NewFromInt(6),
					Close:  decimal.NewFromInt(9),
					Volume: decimal.NewFromInt(10),
				},
				{
					High:   decimal.NewFromInt(9),
					Low:    decimal.NewFromInt(3),
					Close:  decimal.NewFromInt(6),
					Volume: decimal.NewFromInt(15),
				},
			},
			Result: decimal.NewFromInt(0),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.MFI.Calc(c.Candles)
			assertEqualError(t, c.Error, err)

			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_MFI_Count(t *testing.T) {
	assert.Equal(t, 4, MFI{
		valid:  true,
		length: 3,
	}.Count())
}

func Test_NewCMF(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result CMF
		Error  error
	}{
		"Validate returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new CMF": {
			Length: 5,
			Result: CMF{
				valid:  true,
				length: 5,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewCMF(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_CMF_validate(t *testing.T) {
	cc := map[string]struct {
		CMF   CMF
		Error error
	}{
		"Invalid length": {
			CMF:   CMF{},
			Error: ErrInvalidLength,
		},
		"Successfully validated": {
			CMF: CMF{
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.CMF.validate())

			if c.Error == nil {
				assert.True(t, c.CMF.valid)
			}
		})
	}
}

func Test_CMF_Calc(t *testing.T) {
	cc := map[string]struct {
		CMF     CMF
		Candles []Candle
		Result  decimal.Decimal
		Error   error
	}{
		"Invalid indicator": {
			CMF: CMF{
				valid: false,
			},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			CMF: CMF{
				valid:  true,
				length: 2,
			},
			Candles: []Candle{
				{
					High:   decimal.NewFromInt(10),
					Low:    decimal.NewFromInt(0),
					Close:  decimal.NewFromInt(10),
					Volume: decimal.NewFromInt(8),
				},
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
			CMF: CMF{
				valid:  true,
				length: 2,
			},
			Candles: []Candle{
				{
					High:   decimal.NewFromInt(10),
					Low:    decimal.NewFromInt(0),
					Close:  decimal.NewFromInt(10),
					Volume: decimal.NewFromInt(8),
				},
				{
					High:   decimal.NewFromInt(20),
					Low:    decimal.NewFromInt(10),
					Close:  decimal.NewFromInt(10),
					Volume: decimal.NewFromInt(2),
				},
			},
			Result: decimal.RequireFromString("0.6"),
		},
		"Successful calculation with a flat candle": {
			CMF: CMF{
				valid:  true,
				length: 2,
			},
			Candles: []Candle{
				{
					High:   decimal.NewFromInt(10),
					Low:    decimal.NewFromInt(0),
					Close:  decimal.NewFromInt(10),
					Volume: decimal.NewFromInt(8),
				},
				{
					High:   decimal.NewFromInt(5),
					Low:    decimal.NewFromInt(5),
					Close:  decimal.NewFromInt(5),
					Volume: decimal.NewFromInt(2),
				},
			},
			Result: decimal.RequireFromString("0.8"),
		},
		"Successful calculation with zero volume": {
			CMF: CMF{
				valid:  true,
				length: 1,
			},
			Candles: []Candle{
				{
					High:  decimal.NewFromInt(10),
					Low:   decimal.NewFromInt(0),
					Close: decimal.NewFromInt(10),
				},
			},
			Result: decimal.NewFromInt(0),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.CMF.Calc(c.Candles)
			assertEqualError(t, c.Error, err)

			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_CMF_Count(t *testing.T) {
	assert.Equal(t, 2, CMF{
		valid:  true,
		length: 2,
	}.Count())
}
