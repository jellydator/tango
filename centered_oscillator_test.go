package indc

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestMACDValidation(t *testing.T) {
	cc := map[string]struct {
		MA1   MA
		MA2   MA
		Error error
	}{
		"MA1 returns an error": {
			MA1:   EMA{Length: -1},
			MA2:   EMA{Length: 1},
			Error: assert.AnError,
		},
		"MA2 returns an error": {
			MA1:   EMA{Length: 1},
			MA2:   EMA{Length: -1},
			Error: assert.AnError,
		},
		"MA1 is nil": {
			MA2:   EMA{Length: 1},
			Error: ErrMANotSet,
		},
		"MA2 is nil": {
			MA1:   EMA{Length: 1},
			Error: ErrMANotSet,
		},
		"Successful validation": {
			MA1: EMA{Length: 1},
			MA2: EMA{Length: 1},
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
		MA1    MA
		MA2    MA
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"MA1 insufficient amount of data points": {
			MA1: EMA{Length: 4},
			MA2: EMA{Length: 1},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataPointCount,
		},
		"MA2 insufficient amount of data points": {
			MA1: EMA{Length: 1},
			MA2: EMA{Length: 4},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataPointCount,
		},
		"Successful calculation": {
			MA1: SMA{Length: 2},
			MA2: SMA{Length: 3},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
				decimal.NewFromInt(31),
				decimal.NewFromInt(32),
				decimal.NewFromInt(30),
				decimal.NewFromInt(31),
				decimal.NewFromInt(32),
			},
			Result: decimal.NewFromFloat(0.5),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			macd := MACD{MA1: c.MA1, MA2: c.MA2}
			res, err := macd.Calc(c.Data)
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

			res, err = CalcMACD(c.Data, c.MA1, c.MA2)
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

func TestMACDCount(t *testing.T) {
	macd := MACD{MA1: EMA{Length: 10}, MA2: EMA{Length: 1}}
	assert.Equal(t, macd.MA1.Count(), macd.Count())
	assert.Equal(t, macd.MA1.Count(), CountMACD(macd.MA1, macd.MA2))

	macd = MACD{MA1: EMA{Length: 2}, MA2: EMA{Length: 9}}
	assert.Equal(t, macd.MA2.Count(), macd.Count())
	assert.Equal(t, macd.MA2.Count(), CountMACD(macd.MA1, macd.MA2))
}

func TestCCIValidation(t *testing.T) {
	cc := map[string]struct {
		MA    MA
		Error error
	}{
		"MA returns an error": {
			MA:    EMA{Length: -1},
			Error: assert.AnError,
		},
		"MA is nil": {
			Error: ErrMANotSet,
		},
		"Successful validation": {
			MA: EMA{Length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			cci := CCI{MA: c.MA}
			err := cci.Validate()
			if c.Error != nil {
				if c.Error == assert.AnError {
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, c.Error, err)
				}
			} else {
				assert.Nil(t, err)
			}

			err = ValidateCCI(cci.MA)
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

func TestCCICalc(t *testing.T) {
	cc := map[string]struct {
		MA     MA
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Insufficient amount of data points": {
			MA: EMA{Length: 10},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataPointCount,
		},
		"Successful calculation": {
			MA: SMA{Length: 20},
			Data: []decimal.Decimal{
				decimal.NewFromFloat(23.98),
				decimal.NewFromFloat(23.92),
				decimal.NewFromFloat(23.79),
				decimal.NewFromFloat(23.67),
				decimal.NewFromFloat(23.54),
				decimal.NewFromFloat(23.36),
				decimal.NewFromFloat(23.65),
				decimal.NewFromFloat(23.72),
				decimal.NewFromFloat(24.16),
				decimal.NewFromFloat(23.91),
				decimal.NewFromFloat(23.81),
				decimal.NewFromFloat(23.92),
				decimal.NewFromFloat(23.74),
				decimal.NewFromFloat(24.68),
				decimal.NewFromFloat(24.94),
				decimal.NewFromFloat(24.93),
				decimal.NewFromFloat(25.10),
				decimal.NewFromFloat(25.12),
				decimal.NewFromFloat(25.20),
				decimal.NewFromFloat(25.06),
			},
			Result: decimal.NewFromFloat(101.91846523),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			cci := CCI{MA: c.MA}
			res, err := cci.Calc(c.Data)
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

			res, err = CalcCCI(c.Data, c.MA)
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

func TestCCICount(t *testing.T) {
	c := CCI{MA: EMA{Length: 10}}
	assert.Equal(t, c.MA.Count(), c.Count())
	assert.Equal(t, c.MA.Count(), CountCCI(c.MA))
}
