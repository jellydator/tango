package indc

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/swithek/chartype"
)

func TestRSIValidation(t *testing.T) {
	cc := map[string]struct {
		Length int
		Offset int
		Src    chartype.CandleField
		Error  error
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
			Error:  assert.AnError,
		},
		"Successful validation": {
			Length: 1,
			Offset: 0,
			Src:    chartype.CandleClose,
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			r := RSI{Length: c.Length, Offset: c.Offset, Src: c.Src}
			err := r.Validate()
			if c.Error != nil {
				if c.Error == assert.AnError {
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, c.Error, err)
				}
			} else {
				assert.Nil(t, err)
			}

			err = ValidateRSI(c.Length, c.Offset, c.Src)
			if c.Error != nil {
				if c.Error == assert.AnError {
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, c.Error, err)
				}
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestRSICalc(t *testing.T) {
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
			Length: 14,
			Offset: 1,
			Src:    chartype.CandleClose,
			Candles: []chartype.Candle{
				{Close: decimal.NewFromFloat32(44.34)},
				{Close: decimal.NewFromFloat32(44.09)},
				{Close: decimal.NewFromFloat32(44.15)},
				{Close: decimal.NewFromFloat32(43.61)},
				{Close: decimal.NewFromFloat32(44.33)},
				{Close: decimal.NewFromFloat32(44.83)},
				{Close: decimal.NewFromFloat32(45.10)},
				{Close: decimal.NewFromFloat32(45.42)},
				{Close: decimal.NewFromFloat32(45.84)},
				{Close: decimal.NewFromFloat32(46.08)},
				{Close: decimal.NewFromFloat32(45.89)},
				{Close: decimal.NewFromFloat32(46.03)},
				{Close: decimal.NewFromFloat32(45.61)},
				{Close: decimal.NewFromFloat32(46.28)},
				{Close: decimal.NewFromInt(420)},
			},
			Result: decimal.NewFromFloat(70.46413502),
		},
		"Successful calculation without offset": {
			Length: 14,
			Src:    chartype.CandleClose,
			Candles: []chartype.Candle{
				{Close: decimal.NewFromFloat32(44.34)},
				{Close: decimal.NewFromFloat32(44.09)},
				{Close: decimal.NewFromFloat32(44.15)},
				{Close: decimal.NewFromFloat32(43.61)},
				{Close: decimal.NewFromFloat32(44.33)},
				{Close: decimal.NewFromFloat32(44.83)},
				{Close: decimal.NewFromFloat32(45.10)},
				{Close: decimal.NewFromFloat32(45.42)},
				{Close: decimal.NewFromFloat32(45.84)},
				{Close: decimal.NewFromFloat32(46.08)},
				{Close: decimal.NewFromFloat32(45.89)},
				{Close: decimal.NewFromFloat32(46.03)},
				{Close: decimal.NewFromFloat32(45.61)},
				{Close: decimal.NewFromFloat32(46.28)},
			},
			Result: decimal.NewFromFloat(70.46413502),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			r := RSI{Length: c.Length, Offset: c.Offset, Src: c.Src}
			res, err := r.Calc(c.Candles)
			if c.Error != nil {
				if c.Error == assert.AnError {
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, c.Error, err)
				}
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result.String(), res.String())
			}

			res, err = CalcRSI(c.Candles, c.Length, c.Offset, c.Src)
			if c.Error != nil {
				if c.Error == assert.AnError {
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, c.Error, err)
				}
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result.String(), res.String())
			}
		})
	}
}

func TestRSICandleCount(t *testing.T) {
	r := RSI{Length: 15, Offset: 10}
	assert.Equal(t, 25, r.CandleCount())
	assert.Equal(t, 25, CandleCountRSI(15, 10))
}

func TestSTOCHValidation(t *testing.T) {
	cc := map[string]struct {
		Length int
		Offset int
		Src    chartype.CandleField
		Error  error
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
			Error:  assert.AnError,
		},
		"Successful validation": {
			Length: 1,
			Offset: 0,
			Src:    chartype.CandleClose,
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			s := STOCH{Length: c.Length, Offset: c.Offset, Src: c.Src}
			err := s.Validate()
			if c.Error != nil {
				if c.Error == assert.AnError {
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, c.Error, err)
				}
			} else {
				assert.Nil(t, err)
			}

			err = ValidateSTOCH(c.Length, c.Offset, c.Src)
			if c.Error != nil {
				if c.Error == assert.AnError {
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, c.Error, err)
				}
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestSTOCHCalc(t *testing.T) {
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
				{Close: decimal.NewFromInt(150)},
				{Close: decimal.NewFromInt(125)},
				{Close: decimal.NewFromInt(145)},
				{Close: decimal.NewFromInt(420)},
			},
			Result: decimal.NewFromInt(80),
		},
		"Successful calculation without offset": {
			Length: 3,
			Src:    chartype.CandleClose,
			Candles: []chartype.Candle{
				{Close: decimal.NewFromInt(150)},
				{Close: decimal.NewFromInt(125)},
				{Close: decimal.NewFromInt(145)},
			},
			Result: decimal.NewFromInt(80),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			s := STOCH{Length: c.Length, Offset: c.Offset, Src: c.Src}
			res, err := s.Calc(c.Candles)
			if c.Error != nil {
				if c.Error == assert.AnError {
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, c.Error, err)
				}
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result.String(), res.String())
			}

			res, err = CalcSTOCH(c.Candles, c.Length, c.Offset, c.Src)
			if c.Error != nil {
				if c.Error == assert.AnError {
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, c.Error, err)
				}
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result.String(), res.String())
			}
		})
	}
}

func TestSTOCHCandleCount(t *testing.T) {
	s := STOCH{Length: 15, Offset: 10}
	assert.Equal(t, 25, s.CandleCount())
	assert.Equal(t, 25, CandleCountSTOCH(15, 10))
}
