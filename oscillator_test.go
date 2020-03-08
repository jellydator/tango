package indc

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestRSIValidation(t *testing.T) {
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

			r := RSI{Length: c.Length}
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

			err = ValidateRSI(c.Length)
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
			Length: 14,
			Data: []decimal.Decimal{
				decimal.NewFromFloat32(44.34),
				decimal.NewFromFloat32(44.09),
				decimal.NewFromFloat32(44.15),
				decimal.NewFromFloat32(43.61),
				decimal.NewFromFloat32(44.33),
				decimal.NewFromFloat32(44.83),
				decimal.NewFromFloat32(45.10),
				decimal.NewFromFloat32(45.42),
				decimal.NewFromFloat32(45.84),
				decimal.NewFromFloat32(46.08),
				decimal.NewFromFloat32(45.89),
				decimal.NewFromFloat32(46.03),
				decimal.NewFromFloat32(45.61),
				decimal.NewFromFloat32(46.28),
			},
			Result: decimal.NewFromFloat(70.46413502),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			r := RSI{Length: c.Length}
			res, err := r.Calc(c.Data)
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

			res, err = CalcRSI(c.Data, c.Length)
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

func TestRSICount(t *testing.T) {
	r := RSI{Length: 15}
	assert.Equal(t, 15, r.Count())
	assert.Equal(t, 15, CountRSI(15))
}

func TestStochValidation(t *testing.T) {

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

			s := Stoch{Length: c.Length}
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

			err = ValidateStoch(c.Length)

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

func TestStochCalc(t *testing.T) {
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
				decimal.NewFromInt(150),
				decimal.NewFromInt(125),
				decimal.NewFromInt(145),
			},
			Result: decimal.NewFromInt(80),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			s := Stoch{Length: c.Length}
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

			res, err = CalcStoch(c.Data, c.Length)

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

func TestStochCount(t *testing.T) {
	s := Stoch{Length: 15}
	assert.Equal(t, 15, s.Count())
	assert.Equal(t, 15, CountStoch(15))
}

func TestROCValidation(t *testing.T) {
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

			r := ROC{Length: c.Length}
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

			err = ValidateROC(c.Length)
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

func TestROCCalc(t *testing.T) {
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
			Length: 5,
			Data: []decimal.Decimal{
				decimal.NewFromInt(7),
				decimal.NewFromInt(420),
				decimal.NewFromInt(420),
				decimal.NewFromInt(420),
				decimal.NewFromInt(10),
			},
			Result: decimal.NewFromFloat(42.85714286),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			r := ROC{Length: c.Length}
			res, err := r.Calc(c.Data)
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

			res, err = CalcROC(c.Data, c.Length)
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

func TestROCCount(t *testing.T) {
	r := ROC{Length: 15}
	assert.Equal(t, 15, r.Count())
	assert.Equal(t, 15, CountROC(15))
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
