package indc

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestSMAValidation(t *testing.T) {
	cc := map[string]struct {
		Length int
		Error  error
	}{
		"Length cannot be less than 1": {
			Length: 0,
			Error:  ErrInvalidLength,
		},
		"Successful validation": {
			Length: 1,
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			s := SMA{Length: c.Length}
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

			err = ValidateSMA(c.Length)
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
		Length int
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Insufficient amount of data points": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataPointCount,
		},
		"Successful calculation": {
			Length: 3,
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

			s := SMA{Length: c.Length}
			res, err := s.Calc(c.Data)
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

			res, err = CalcSMA(c.Data, c.Length)
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

func TestSMACount(t *testing.T) {
	s := SMA{Length: 15}
	assert.Equal(t, 15, s.Count())
	assert.Equal(t, 15, CountSMA(15))
}

func TestEMAValidation(t *testing.T) {
	cc := map[string]struct {
		Length int
		Error  error
	}{
		"Length cannot be less than 1": {
			Length: 0,
			Error:  ErrInvalidLength,
		},
		"Successful validation": {
			Length: 1,
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			e := EMA{Length: c.Length}
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

			err = ValidateEMA(c.Length)
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
		Length int
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Insufficient amount of data points": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataPointCount,
		},
		"Successful calculation": {
			Length: 2,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
				decimal.NewFromInt(31),
				decimal.NewFromInt(32),
				decimal.NewFromInt(30),
				decimal.NewFromInt(31),
				decimal.NewFromInt(31),
			},
			Result: decimal.NewFromFloat(31),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			e := EMA{Length: c.Length}
			res, err := e.Calc(c.Data)
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

			res, err = CalcEMA(c.Data, c.Length)
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

func TestEMACount(t *testing.T) {
	e := EMA{Length: 15}
	assert.Equal(t, 30, e.Count())
	assert.Equal(t, 30, CountEMA(15))
}

func TestEMAMultiplier(t *testing.T) {
	e := EMA{Length: 3}
	assert.Equal(t, decimal.NewFromFloat(0.5), e.multiplier())
}

func TestWMAValidation(t *testing.T) {
	cc := map[string]struct {
		Length int
		Error  error
	}{
		"Length cannot be less than 1": {
			Length: 0,
			Error:  ErrInvalidLength,
		},
		"Successful validation": {
			Length: 1,
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			w := WMA{Length: c.Length}
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

			err = ValidateWMA(c.Length)
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
		Length int
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Insufficient amount of data points": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataPointCount,
		},
		"Successful calculation": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(420),
				decimal.NewFromInt(420),
				decimal.NewFromInt(420),
				decimal.NewFromInt(30),
				decimal.NewFromInt(30),
				decimal.NewFromInt(32),
			},
			Result: decimal.NewFromFloat(31),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			w := WMA{Length: c.Length}
			res, err := w.Calc(c.Data)
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

			res, err = CalcWMA(c.Data, c.Length)
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

func TestWMACount(t *testing.T) {
	w := WMA{Length: 15}
	assert.Equal(t, 15, w.Count())
	assert.Equal(t, 15, CountWMA(15))
}
