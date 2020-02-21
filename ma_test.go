package indc

import (
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/swithek/chartype"
)

func TestSMAValidation(t *testing.T) {
	cc := map[string]struct {
		Length  int
		Offset  int
		Src     chartype.CandleField
		Candles []chartype.Candle
		Result  decimal.Decimal
		Error   error
	}{
		"Length cannot be less than 1": {
			Length: 0,
			Error:  ErrInvalidLength,
		},
		"Offset cannot be less than 0": {
			Length: 1,
			Offset: -1,
			Error:  ErrInvalidOffset,
		},
		"Invalid CandleField value": {
			Length: 1,
			Offset: 0,
			Src:    -69,
			Error:  errors.New(""),
		},
		"Unexpected internal error has occured": {
			Length: 1,
			Offset: 0,
			Src:    1,
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			s := SMA{Length: c.Length, Offset: c.Offset, Src: c.Src}
			err := s.Validate()
			if c.Error != nil {
				if c.Error.Error() == "" {
					assert.Error(t, err)
				} else {
					assert.Equal(t, c.Error, err)
					return
				}
			} else {
				assert.Nil(t, err)
			}

			err = ValidateSMA(c.Length, c.Offset, c.Src)
			if c.Error != nil {
				if c.Error.Error() == "" {
					assert.Error(t, err)
				} else {
					assert.Equal(t, c.Error, err)
					return
				}
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestSMACalc(t *testing.T) {
	cc := map[string]struct {
		Length  int
		Offset  int
		Src     chartype.CandleField
		Candles []chartype.Candle
		Result  decimal.Decimal
		Error   error
	}{
		"Insufficient amount of candles": {
			Length: 3,
			Src:    chartype.CandleClose,
			Candles: []chartype.Candle{
				{Close: decimal.NewFromInt(30)},
			},
			Error: ErrInvalidCandleCount,
		},
		"Successful calculation with offset": {
			Length: 3,
			Offset: 1,
			Src:    chartype.CandleClose,
			Candles: []chartype.Candle{
				{Close: decimal.NewFromInt(30)},
				{Close: decimal.NewFromInt(31)},
				{Close: decimal.NewFromInt(32)},
				{Close: decimal.NewFromInt(39)},
			},
			Result: decimal.NewFromInt(31),
		},
		"Successful calculation without offset": {
			Length: 3,
			Src:    chartype.CandleClose,
			Candles: []chartype.Candle{
				{Close: decimal.NewFromInt(30)},
				{Close: decimal.NewFromInt(31)},
				{Close: decimal.NewFromInt(32)},
			},
			Result: decimal.NewFromInt(31),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			s := SMA{Length: c.Length, Offset: c.Offset, Src: c.Src}
			res, err := s.Calc(c.Candles)
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
				return
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result.String(), res.String())
			}

			res, err = CalcSMA(c.Candles, c.Length, c.Offset, c.Src)
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
				return
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result.String(), res.String())
			}
		})
	}
}

func TestSMACandleCount(t *testing.T) {
	s := SMA{Length: 15, Offset: 10}
	assert.Equal(t, 25, s.CandleCount())
	assert.Equal(t, 25, CandleCountSMA(15, 10))
}

func TestEMAValidation(t *testing.T) {
	cc := map[string]struct {
		Length  int
		Offset  int
		Src     chartype.CandleField
		Candles []chartype.Candle
		Result  decimal.Decimal
		Error   error
	}{
		"Length cannot be less than 1": {
			Length: 0,
			Error:  ErrInvalidLength,
		},
		"Offset cannot be less than 0": {
			Length: 1,
			Offset: -1,
			Error:  ErrInvalidOffset,
		},
		"Invalid CandleField value": {
			Length: 1,
			Offset: 0,
			Src:    -69,
			Error:  errors.New(""),
		},
		"Unexpected internal error has occured": {
			Length: 1,
			Offset: 0,
			Src:    1,
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			e := EMA{Length: c.Length, Offset: c.Offset, Src: c.Src}
			err := e.Validate()
			if c.Error != nil {
				if c.Error.Error() == "" {
					assert.Error(t, err)
				} else {
					assert.Equal(t, c.Error, err)
					return
				}
			} else {
				assert.Nil(t, err)
			}

			err = ValidateEMA(c.Length, c.Offset, c.Src)
			if c.Error != nil {
				if c.Error.Error() == "" {
					assert.Error(t, err)
				} else {
					assert.Equal(t, c.Error, err)
					return
				}
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestEMACalc(t *testing.T) {
	cc := map[string]struct {
		Length  int
		Offset  int
		Src     chartype.CandleField
		Candles []chartype.Candle
		Result  decimal.Decimal
		Error   error
	}{
		"Insufficient amount of candles": {
			Length: 3,
			Src:    chartype.CandleClose,
			Candles: []chartype.Candle{
				{Close: decimal.NewFromInt(30)},
			},
			Error: ErrInvalidCandleCount,
		},
		"Successful calculation with offset": {
			Length: 3,
			Offset: 1,
			Src:    chartype.CandleClose,
			Candles: []chartype.Candle{
				{Close: decimal.NewFromInt(30)},
				{Close: decimal.NewFromInt(31)},
				{Close: decimal.NewFromInt(32)},
				{Close: decimal.NewFromInt(30)},
				{Close: decimal.NewFromInt(31)},
				{Close: decimal.NewFromInt(32)},
				{Close: decimal.NewFromInt(39)},
			},
			Result: decimal.NewFromFloat(31.375),
		},
		"Successful calculation without offset": {
			Length: 3,
			Src:    chartype.CandleClose,
			Candles: []chartype.Candle{
				{Close: decimal.NewFromInt(30)},
				{Close: decimal.NewFromInt(31)},
				{Close: decimal.NewFromInt(32)},
				{Close: decimal.NewFromInt(30)},
				{Close: decimal.NewFromInt(31)},
				{Close: decimal.NewFromInt(32)},
			},
			Result: decimal.NewFromFloat(31.375),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			e := EMA{Length: c.Length, Offset: c.Offset, Src: c.Src}
			res, err := e.Calc(c.Candles)
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
				return
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result.String(), res.String())
			}

			res, err = CalcEMA(c.Candles, c.Length, c.Offset, c.Src)
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
				return
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result.String(), res.String())
			}
		})
	}
}

func TestEMACandleCount(t *testing.T) {
	e := EMA{Length: 15, Offset: 10}
	assert.Equal(t, 40, e.CandleCount())
	assert.Equal(t, 40, CandleCountEMA(15, 10))
}

func TestEMAMultiplier(t *testing.T) {
	e := EMA{Length: 3}
	assert.Equal(t, decimal.NewFromFloat(0.5), e.multiplier())
}
