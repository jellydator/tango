package indc

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestAroonValidation(t *testing.T) {
	cc := map[string]struct {
		Trend  string
		Length int
		Error  error
	}{
		"Invalid Aroon trend": {
			Trend:  "downn",
			Length: 5,
			Error:  ErrInvalidType,
		},
		"Length cannot be less than 1": {
			Trend:  "down",
			Length: 0,
			Error:  ErrInvalidLength,
		},
		"Successful validation with up trend": {
			Trend:  "up",
			Length: 5,
		},
		"Successful validation with down trend": {
			Trend:  "down",
			Length: 5,
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			a := Aroon{Trend: c.Trend, Length: c.Length}
			err := a.Validate()
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

func TestAroonCalc(t *testing.T) {
	cc := map[string]struct {
		Trend  string
		Length int
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Insufficient amount of data points": {
			Trend:  "down",
			Length: 5,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataPointCount,
		},
		"Successful calculation with up trend": {
			Trend:  "up",
			Length: 5,
			Data: []decimal.Decimal{
				decimal.NewFromInt(25),
				decimal.NewFromInt(31),
				decimal.NewFromInt(38),
				decimal.NewFromInt(35),
				decimal.NewFromInt(29),
				decimal.NewFromInt(29),
			},
			Result: decimal.NewFromFloat(40),
		},
		"Successful calculation with down trend": {
			Trend:  "down",
			Length: 5,
			Data: []decimal.Decimal{
				decimal.NewFromInt(25),
				decimal.NewFromInt(31),
				decimal.NewFromInt(38),
				decimal.NewFromInt(35),
				decimal.NewFromInt(29),
				decimal.NewFromInt(29),
			},
			Result: decimal.NewFromFloat(100),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			a := Aroon{Trend: c.Trend, Length: c.Length}
			res, err := a.Calc(c.Data)
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

func TestAroonCount(t *testing.T) {
	a := Aroon{Trend: "down", Length: 5}
	assert.Equal(t, 5, a.Count())
}

func TestCCIValidation(t *testing.T) {
	cc := map[string]struct {
		Src    Source
		Error error
	}{
		"Source returns an error": {
			Src:    Source{Indicator: EMA{Length: -1}},
			Error: assert.AnError,
		},
		"Successful validation": {
			Src: Source{Indicator: EMA{Length: 1}},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			cci := CCI{Src: c.Src}
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
		})
	}
}

func TestCCICalc(t *testing.T) {
	cc := map[string]struct {
		Src     Source
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Insufficient amount of data points": {
			Src: Source{Indicator: EMA{Length: 10}},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataPointCount,
		},
		"Successful calculation": {
			Src: Source{Indicator: SMA{Length: 20}},
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
			Result: decimal.NewFromFloat(101.91846522781775),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			cci := CCI{Src: c.Src}
			res, err := cci.Calc(c.Data)
			if c.Error != nil {
				if c.Error == assert.AnError {
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, c.Error, err)
				}
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result.String(), res.Round(14).String())
			}
		})
	}
}

func TestCCICount(t *testing.T) {
	c := CCI{Src: Source{Indicator: EMA{Length: 10}}}
	assert.Equal(t, c.Src.Indicator.Count(), c.Count())
}

func TestDEMAValidation(t *testing.T) {
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

			d := DEMA{Length: c.Length}
			err := d.Validate()
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

func TestDEMACalc(t *testing.T) {
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
			Result: decimal.NewFromFloat(30.72222222222222),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			d := DEMA{Length: c.Length}
			res, err := d.Calc(c.Data)
			if c.Error != nil {
				if c.Error == assert.AnError {
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, c.Error, err)
				}
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result.String(), res.Round(14).String())
			}
		})
	}
}

func TestDEMACount(t *testing.T) {
	d := DEMA{Length: 15}
	assert.Equal(t, 29, d.Count())
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
			Result: decimal.NewFromFloat(30.83333333333333),
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
				assert.Equal(t, c.Result.String(), res.Round(14).String())
			}
		})
	}
}

func TestEMACount(t *testing.T) {
	e := EMA{Length: 15}
	assert.Equal(t, 29, e.Count())
}

func TestEMAMultiplier(t *testing.T) {
	e := EMA{Length: 3}
	assert.Equal(t, decimal.NewFromFloat(0.5), e.multiplier())
}

func TestHMAValidation(t *testing.T) {
	cc := map[string]struct {
		WMA WMA
		Error  error
	}{
		"WMA returns an error": {
			WMA:   WMA{Length: -1},
			Error: assert.AnError,
		},
		"WMA not set": {
			Error: ErrMAIndicatorNotSet,
		},
		"Successful validation": {
			WMA: WMA{Length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			h := HMA{WMA: c.WMA}
			err := h.Validate()
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

func TestHMACalc(t *testing.T) {
	cc := map[string]struct {
		WMA WMA
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Insufficient amount of data points": {
			WMA: WMA{Length: 5},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataPointCount,
		},
		"Successful calculation": {
			WMA: WMA{Length: 3},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
				decimal.NewFromInt(31),
				decimal.NewFromInt(32),
				decimal.NewFromInt(30),
				decimal.NewFromInt(30),
				decimal.NewFromInt(31),
			},
			Result: decimal.NewFromFloat(31.5),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			h := HMA{WMA: c.WMA}
			res, err := h.Calc(c.Data)
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

func TestHMACount(t *testing.T) {
	h := HMA{WMA: WMA{Length: 15}}
	assert.Equal(t, 29, h.Count())
}

func TestMACDValidation(t *testing.T) {
	cc := map[string]struct {
		Src1   Source
		Src2   Source
		Error error
	}{
		"Src1 returns an error": {
			Src1:   Source{Indicator: EMA{Length: -1}},
			Src2:   Source{Indicator: EMA{Length: 1}},
			Error: assert.AnError,
		},
		"Src2 returns an error": {
			Src1:   Source{Indicator: EMA{Length: 1}},
			Src2:   Source{Indicator: EMA{Length: -1}},
			Error: assert.AnError,
		},
		"Successful validation": {
			Src1: Source{Indicator: EMA{Length: 1}},
			Src2: Source{Indicator: EMA{Length: 1}},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			m := MACD{Src1: c.Src1, Src2: c.Src2}
			err := m.Validate()
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
		Src1    Source
		Src2    Source
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Src1 insufficient amount of data points": {
			Src1: Source{Indicator: EMA{Length: 4}},
			Src2: Source{Indicator: EMA{Length: 1}},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataPointCount,
		},
		"Src2 insufficient amount of data points": {
			Src1: Source{Indicator: EMA{Length: 1}},
			Src2: Source{Indicator: EMA{Length: 4}},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataPointCount,
		},
		"Successful calculation": {
			Src1: Source{Indicator: SMA{Length: 2}},
			Src2: Source{Indicator: SMA{Length: 3}},
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

			m := MACD{Src1: c.Src1, Src2: c.Src2}
			res, err := m.Calc(c.Data)
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
	m := MACD{Src1: Source{Indicator: EMA{Length: 10}}, Src2: Source{Indicator: EMA{Length: 1}}}
	assert.Equal(t, m.Count(), m.Src1.Indicator.Count())

	m = MACD{Src1: Source{Indicator: EMA{Length: 2}}, Src2: Source{Indicator: EMA{Length: 9}}}
	assert.Equal(t, m.Count(), m.Src2.Indicator.Count())
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
			Result: decimal.NewFromFloat(42.85714285714286),
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
		})
	}
}

func TestROCCount(t *testing.T) {
	r := ROC{Length: 15}
	assert.Equal(t, 15, r.Count())
}

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
			Result: decimal.NewFromFloat(70.46413502109705),
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
				assert.Equal(t, c.Result.String(), res.Round(14).String())
			}
		})
	}
}

func TestRSICount(t *testing.T) {
	r := RSI{Length: 15}
	assert.Equal(t, 15, r.Count())
}

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
		})
	}
}

func TestSMACount(t *testing.T) {
	s := SMA{Length: 15}
	assert.Equal(t, 15, s.Count())
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
		})
	}
}

func TestStochCount(t *testing.T) {
	s := Stoch{Length: 15}
	assert.Equal(t, 15, s.Count())
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
		})
	}
}

func TestWMACount(t *testing.T) {
	w := WMA{Length: 15}
	assert.Equal(t, 15, w.Count())
}
