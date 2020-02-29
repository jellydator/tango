package indc

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/swithek/chartype"
)

func TestSMAValidation(t *testing.T) {
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

			s := SMA{Length: c.Length, Offset: c.Offset, Src: c.Src}
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

			err = ValidateSMA(c.Length, c.Offset, c.Src)
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
				{Close: decimal.NewFromInt(420)},
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
			}
			assert.Nil(t, err)
			assert.Equal(t, c.Result.String(), res.String())

			res, err = CalcSMA(c.Candles, c.Length, c.Offset, c.Src)
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, c.Result.String(), res.String())
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

			e := EMA{Length: c.Length, Offset: c.Offset, Src: c.Src}
			err := e.Validate()
			if c.Error != nil {
				if c.Error == assert.AnError {
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, c.Error, err)
				}
			} else {
				assert.Nil(t, err)
			}

			err = ValidateEMA(c.Length, c.Offset, c.Src)
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
				{Close: decimal.NewFromInt(420)},
			},
			Result: decimal.NewFromFloat(31.375),
		},
		"Successful calculation without offset": {
			Length: 2,
			Src:    chartype.CandleClose,
			Candles: []chartype.Candle{
				{Close: decimal.NewFromInt(30)},
				{Close: decimal.NewFromInt(31)},
				{Close: decimal.NewFromInt(32)},
				{Close: decimal.NewFromInt(30)},
				{Close: decimal.NewFromInt(31)},
				{Close: decimal.NewFromInt(31)},
			},
			Result: decimal.NewFromFloat(31),
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
			}
			assert.Nil(t, err)
			assert.Equal(t, c.Result.String(), res.String())

			res, err = CalcEMA(c.Candles, c.Length, c.Offset, c.Src)
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, c.Result.String(), res.String())
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

func TestWMAValidation(t *testing.T) {
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

			w := WMA{Length: c.Length, Offset: c.Offset, Src: c.Src}
			err := w.Validate()
			if c.Error != nil {
				if c.Error == assert.AnError {
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, c.Error, err)
				}
			} else {
				assert.Nil(t, err)
			}

			err = ValidateWMA(c.Length, c.Offset, c.Src)
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

func TestWMACalc(t *testing.T) {
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
				{Close: decimal.NewFromInt(420)},
				{Close: decimal.NewFromInt(420)},
				{Close: decimal.NewFromInt(420)},
				{Close: decimal.NewFromInt(30)},
				{Close: decimal.NewFromInt(30)},
				{Close: decimal.NewFromInt(32)},
				{Close: decimal.NewFromInt(420)},
			},
			Result: decimal.NewFromFloat(31),
		},
		"Successful calculation without offset": {
			Length: 3,
			Src:    chartype.CandleClose,
			Candles: []chartype.Candle{
				{Close: decimal.NewFromInt(420)},
				{Close: decimal.NewFromInt(420)},
				{Close: decimal.NewFromInt(420)},
				{Close: decimal.NewFromInt(30)},
				{Close: decimal.NewFromInt(30)},
				{Close: decimal.NewFromInt(32)},
			},
			Result: decimal.NewFromFloat(31),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			w := WMA{Length: c.Length, Offset: c.Offset, Src: c.Src}
			res, err := w.Calc(c.Candles)
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, c.Result.String(), res.String())

			res, err = CalcWMA(c.Candles, c.Length, c.Offset, c.Src)
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func TestWMACandleCount(t *testing.T) {
	w := WMA{Length: 15, Offset: 10}
	assert.Equal(t, 25, w.CandleCount())
	assert.Equal(t, 25, CandleCountWMA(15, 10))
}

func TestMACDValidation(t *testing.T) {
	cc := map[string]struct {
		MA1   MA
		MA2   MA
		Error error
	}{
		"MA1 returns an error": {
			MA1:   EMA{Length: -1, Offset: 0, Src: chartype.CandleClose},
			MA2:   EMA{Length: 1, Offset: 0, Src: chartype.CandleClose},
			Error: assert.AnError,
		},
		"MA2 returns an error": {
			MA1:   EMA{Length: 1, Offset: 0, Src: chartype.CandleClose},
			MA2:   EMA{Length: -1, Offset: 0, Src: chartype.CandleClose},
			Error: assert.AnError,
		},
		"MA1 is nil": {
			MA2:   EMA{Length: 1, Offset: 0, Src: chartype.CandleClose},
			Error: ErrMANotSet,
		},
		"MA2 is nil": {
			MA1:   EMA{Length: 1, Offset: 0, Src: chartype.CandleClose},
			Error: ErrMANotSet,
		},
		"Successful validation": {
			MA1: EMA{Length: 1, Offset: 0, Src: chartype.CandleClose},
			MA2: EMA{Length: 1, Offset: 0, Src: chartype.CandleClose},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			macd := MACD{MA1: c.MA1, MA2: c.MA2}
			err := macd.Validate()
			if c.Error != nil {
				if c.Error == assert.AnError {
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, c.Error, err)
				}
			} else {
				assert.Nil(t, err)
			}

			err = ValidateMACD(c.MA1, c.MA2)
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

func TestMACDCalc(t *testing.T) {
	cc := map[string]struct {
		MA1     MA
		MA2     MA
		Candles []chartype.Candle
		Result  decimal.Decimal
		Error   error
	}{
		"MA1 insufficient amount of candles": {
			MA1: EMA{Length: 4, Offset: 0, Src: chartype.CandleClose},
			MA2: EMA{Length: 1, Offset: 0, Src: chartype.CandleClose},
			Candles: []chartype.Candle{
				{Close: decimal.NewFromInt(30)},
			},
			Error: ErrInvalidCandleCount,
		},
		"MA2 insufficient amount of candles": {
			MA1: EMA{Length: 1, Offset: 0, Src: chartype.CandleClose},
			MA2: EMA{Length: 4, Offset: 0, Src: chartype.CandleClose},
			Candles: []chartype.Candle{
				{Close: decimal.NewFromInt(30)},
			},
			Error: ErrInvalidCandleCount,
		},
		"Successful calculation with offset": {
			MA1: EMA{Length: 3, Offset: 1, Src: chartype.CandleOpen},
			MA2: EMA{Length: 2, Offset: 1, Src: chartype.CandleClose},
			Candles: []chartype.Candle{
				{Open: decimal.NewFromInt(30), Close: decimal.NewFromInt(30)},
				{Open: decimal.NewFromInt(31), Close: decimal.NewFromInt(31)},
				{Open: decimal.NewFromInt(32), Close: decimal.NewFromInt(32)},
				{Open: decimal.NewFromInt(30), Close: decimal.NewFromInt(30)},
				{Open: decimal.NewFromInt(31), Close: decimal.NewFromInt(31)},
				{Open: decimal.NewFromInt(32), Close: decimal.NewFromInt(31)},
				{Open: decimal.NewFromInt(420), Close: decimal.NewFromInt(420)},
			},
			Result: decimal.NewFromFloat(0.375),
		},

		"Successful calculation without offset": {
			MA1: EMA{Length: 3, Offset: 0, Src: chartype.CandleOpen},
			MA2: EMA{Length: 2, Offset: 0, Src: chartype.CandleClose},
			Candles: []chartype.Candle{
				{Open: decimal.NewFromInt(30), Close: decimal.NewFromInt(30)},
				{Open: decimal.NewFromInt(31), Close: decimal.NewFromInt(31)},
				{Open: decimal.NewFromInt(32), Close: decimal.NewFromInt(32)},
				{Open: decimal.NewFromInt(30), Close: decimal.NewFromInt(30)},
				{Open: decimal.NewFromInt(31), Close: decimal.NewFromInt(31)},
				{Open: decimal.NewFromInt(32), Close: decimal.NewFromInt(31)},
			},
			Result: decimal.NewFromFloat(0.375),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			macd := MACD{MA1: c.MA1, MA2: c.MA2}
			res, err := macd.Calc(c.Candles)
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, c.Result.String(), res.String())

			res, err = CalcMACD(c.Candles, c.MA1, c.MA2)
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func TestMACDCandleCount(t *testing.T) {
	macd := MACD{MA1: EMA{Length: 10, Offset: 5}, MA2: EMA{Length: 1, Offset: 0}}
	assert.Equal(t, macd.MA1.CandleCount(), macd.CandleCount())
	assert.Equal(t, macd.MA1.CandleCount(), CandleCountMACD(macd.MA1, macd.MA2))

	macd = MACD{MA1: EMA{Length: 2, Offset: 3}, MA2: EMA{Length: 9, Offset: 15}}
	assert.Equal(t, macd.MA2.CandleCount(), macd.CandleCount())
	assert.Equal(t, macd.MA2.CandleCount(), CandleCountMACD(macd.MA1, macd.MA2))
}
