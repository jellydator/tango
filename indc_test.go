package indc

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type IndicatorMock struct{}

func (im IndicatorMock) validate() error { return assert.AnError }

func (im IndicatorMock) Calc(dd []decimal.Decimal) (decimal.Decimal, error) {
	return decimal.Zero, assert.AnError
}

func (im IndicatorMock) Count() int { return 1 }

func TestAroonNew(t *testing.T) {
	cc := map[string]struct {
		Trend  string
		Length int
		Result Aroon
		Error  error
	}{
		"Aroon throws an error": {
			Error: assert.AnError,
		},
		"Successful Aroon creation": {
			Trend:  "down",
			Length: 5,
			Result: Aroon{trend: "down", length: 5},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			a, err := NewAroon(c.Trend, c.Length)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, a)
			}
		})
	}
}

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

			a := Aroon{trend: c.Trend, length: c.Length}
			err := a.validate()
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
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

			a := Aroon{trend: c.Trend, length: c.Length}
			res, err := a.Calc(c.Data)
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result.String(), res.String())
			}
		})
	}
}

func TestAroonCount(t *testing.T) {
	a := Aroon{trend: "down", length: 5}
	assert.Equal(t, 5, a.Count())
}

func TestAroonUnmarshal(t *testing.T) {
	cc := map[string]struct {
		ByteArray []byte
		Result    Aroon
		Error     error
	}{
		"Unmarshal throws an error": {
			ByteArray: []byte(`{"trend": "down"}`),
			Error:     assert.AnError,
		},
		"Successful Aroon unmarshal": {
			ByteArray: []byte(`{"name":"aroon","trend":"up","length":1}`),
			Result:    Aroon{trend: "up", length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			a := Aroon{}
			err := json.Unmarshal(c.ByteArray, &a)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, a)
			}
		})
	}
}

func TestAroonMarshal(t *testing.T) {
	c := struct {
		Aroon  Aroon
		Result []byte
	}{
		Aroon:  Aroon{trend: "down", length: 1},
		Result: []byte(`{"name":"aroon","trend":"down","length":1}`),
	}

	d, _ := json.Marshal(c.Aroon)

	assert.Equal(t, c.Result, d)
}

// func TestCCINew(t *testing.T) {
// 	cc := map[string]struct {
// 		Source Indicator
// 		Result CCI
// 		Error  error
// 	}{
// 		"CCI throws an error": {
// 			Error: assert.AnError,
// 		},
// 		"Successful CCI creation": {
// 			Source: Aroon{trend: "down", length: 1},
// 			Result: CCI{source: Aroon{trend: "down", length: 1}},
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			cci, err := NewCCI(c.Source)
// 			if c.Error != nil {
// 				assert.NotNil(t, err)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result, cci)
// 			}
// 		})
// 	}
// }

// func TestCCIValidation(t *testing.T) {
// 	cc := map[string]struct {
// 		Source Indicator
// 		Error  error
// 	}{
// 		"Source returns an error": {
// 			Error: assert.AnError,
// 		},
// 		"Successful validation": {
// 			Source: Aroon{trend: "down", length: 1},
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			cci := CCI{source: c.Source}
// 			err := cci.validate()
// 			if c.Error != nil {
// 				if c.Error == assert.AnError {
// 					assert.NotNil(t, err)
// 				} else {
// 					assert.Equal(t, c.Error, err)
// 				}
// 			} else {
// 				assert.Nil(t, err)
// 			}
// 		})
// 	}
// }

// func TestCCICalc(t *testing.T) {
// 	cc := map[string]struct {
// 		Source Source
// 		Data   []decimal.Decimal
// 		Result decimal.Decimal
// 		Error  error
// 	}{
// 		"Insufficient amount of data points": {
// 			Source: Source{EMA{Length: 10}},
// 			Data: []decimal.Decimal{
// 				decimal.NewFromInt(30),
// 			},
// 			Error: ErrInvalidDataPointCount,
// 		},
// 		"Source returns an error": {
// 			Source: Source{IndicatorMock{}},
// 			Data: []decimal.Decimal{
// 				decimal.NewFromInt(30),
// 			},
// 			Error: assert.AnError,
// 		},
// 		"Successful calculation": {
// 			Source: Source{SMA{Length: 20}},
// 			Data: []decimal.Decimal{
// 				decimal.NewFromFloat(23.98),
// 				decimal.NewFromFloat(23.92),
// 				decimal.NewFromFloat(23.79),
// 				decimal.NewFromFloat(23.67),
// 				decimal.NewFromFloat(23.54),
// 				decimal.NewFromFloat(23.36),
// 				decimal.NewFromFloat(23.65),
// 				decimal.NewFromFloat(23.72),
// 				decimal.NewFromFloat(24.16),
// 				decimal.NewFromFloat(23.91),
// 				decimal.NewFromFloat(23.81),
// 				decimal.NewFromFloat(23.92),
// 				decimal.NewFromFloat(23.74),
// 				decimal.NewFromFloat(24.68),
// 				decimal.NewFromFloat(24.94),
// 				decimal.NewFromFloat(24.93),
// 				decimal.NewFromFloat(25.10),
// 				decimal.NewFromFloat(25.12),
// 				decimal.NewFromFloat(25.20),
// 				decimal.NewFromFloat(25.06),
// 			},
// 			Result: decimal.NewFromFloat(101.91846522781775),
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			cci := CCI{source: c.Source}
// 			res, err := cci.Calc(c.Data)
// 			if c.Error != nil {
// 				if c.Error == assert.AnError {
// 					assert.NotNil(t, err)
// 				} else {
// 					assert.Equal(t, c.Error, err)
// 				}
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result.String(), res.Round(14).String())
// 			}
// 		})
// 	}
// }

// func TestCCICount(t *testing.T) {
// 	c := CCI{source: Aroon{trend: "down", length: 10}}
// 	assert.Equal(t, c.source.Count(), c.Count())
// }

// func TestCCIUnmarshal(t *testing.T) {
// 	cc := map[string]struct {
// 		ByteArray []byte
// 		Result    CCI
// 		Error     error
// 	}{
// 		"Unmarshal throws an error": {
// 			ByteArray: []byte(`{}`),
// 			Error:     assert.AnError,
// 		},
// 		"Successful CCI unmarshal": {
// 			ByteArray: []byte(`{"name":"cci",
// 			"source":{"name":"aroon","trend":"up","length":1}}`),
// 			Result: CCI{Aroon{trend: "up", length: 1}},
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			cci := CCI{}
// 			err := json.Unmarshal(c.ByteArray, &cci)
// 			if c.Error != nil {
// 				assert.NotNil(t, err)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result, cci)
// 			}
// 		})
// 	}
// }

// func TestCCIMarshal(t *testing.T) {
// 	c := struct {
// 		CCI    CCI
// 		Result []byte
// 	}{
// 		CCI:    CCI{source: Aroon{trend: "down", length: 1}},
// 		Result: []byte(`{"name":"aroon","trend":"down","length":1}`),
// 	}

// 	d, _ := json.Marshal(c.CCI)

// 	assert.Equal(t, c.Result, d)
// }

// func TestDEMANew(t *testing.T) {
// 	cc := map[string]struct {
// 		Length int
// 		Result DEMA
// 		Error  error
// 	}{
// 		"DEMA throws an error": {
// 			Error: assert.AnError,
// 		},
// 		"Successful DEMA creation": {
// 			Length: 1,
// 			Result: DEMA{Length: 1},
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			d, err := NewDEMA(c.Length)
// 			if c.Error != nil {
// 				if c.Error == assert.AnError {
// 					assert.NotNil(t, err)
// 				} else {
// 					assert.Equal(t, c.Result, d)
// 				}
// 			} else {
// 				assert.Nil(t, err)
// 			}
// 		})
// 	}
// }

// func TestDEMAValidation(t *testing.T) {
// 	cc := map[string]struct {
// 		Length int
// 		Error  error
// 	}{
// 		"Length cannot be less than 1": {
// 			Length: 0,
// 			Error:  ErrInvalidLength,
// 		},
// 		"Successful validation": {
// 			Length: 1,
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			d := DEMA{Length: c.Length}
// 			err := d.Validate()
// 			if c.Error != nil {
// 				if c.Error == assert.AnError {
// 					assert.NotNil(t, err)
// 				} else {
// 					assert.Equal(t, c.Error, err)
// 				}
// 			} else {
// 				assert.Nil(t, err)
// 			}
// 		})
// 	}
// }

// func TestDEMACalc(t *testing.T) {
// 	cc := map[string]struct {
// 		Length int
// 		Data   []decimal.Decimal
// 		Result decimal.Decimal
// 		Error  error
// 	}{
// 		"Insufficient amount of data points": {
// 			Length: 3,
// 			Data: []decimal.Decimal{
// 				decimal.NewFromInt(30),
// 			},
// 			Error: ErrInvalidDataPointCount,
// 		},
// 		"Successful calculation": {
// 			Length: 2,
// 			Data: []decimal.Decimal{
// 				decimal.NewFromInt(30),
// 				decimal.NewFromInt(31),
// 				decimal.NewFromInt(32),
// 				decimal.NewFromInt(30),
// 				decimal.NewFromInt(31),
// 				decimal.NewFromInt(31),
// 			},
// 			Result: decimal.NewFromFloat(30.72222222222222),
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			d := DEMA{Length: c.Length}
// 			res, err := d.Calc(c.Data)
// 			if c.Error != nil {
// 				if c.Error == assert.AnError {
// 					assert.NotNil(t, err)
// 				} else {
// 					assert.Equal(t, c.Error, err)
// 				}
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result.String(), res.Round(14).String())
// 			}
// 		})
// 	}
// }

// func TestDEMACount(t *testing.T) {
// 	d := DEMA{Length: 15}
// 	assert.Equal(t, 29, d.Count())
// }

func TestEMANew(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result EMA
		Error  error
	}{
		"EMA throws an error": {
			Error: assert.AnError,
		},
		"Successful EMA creation": {
			Length: 1,
			Result: EMA{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			e, err := NewEMA(c.Length)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, e)
			}
		})
	}
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

			e := EMA{length: c.Length}
			err := e.validate()
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
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

			e := EMA{length: c.Length}
			res, err := e.Calc(c.Data)
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result.String(), res.Round(14).String())
			}
		})
	}
}

func TestEMACount(t *testing.T) {
	e := EMA{length: 15}
	assert.Equal(t, 29, e.Count())
}

func TestEMAMultiplier(t *testing.T) {
	e := EMA{length: 3}
	assert.Equal(t, decimal.NewFromFloat(0.5), e.multiplier())
}

func TestEMAUnmarshal(t *testing.T) {
	cc := map[string]struct {
		ByteArray []byte
		Result    EMA
		Error     error
	}{
		"Unmarshal throws an error": {
			ByteArray: []byte(`{"length": "down"}`),
			Error:     assert.AnError,
		},
		"Successful EMA unmarshal": {
			ByteArray: []byte(`{"name":"ema","length":1}`),
			Result:    EMA{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			e := EMA{}
			err := json.Unmarshal(c.ByteArray, &e)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, e)
			}
		})
	}
}

func TestEMAMarshal(t *testing.T) {
	c := struct {
		EMA    EMA
		Result []byte
	}{
		EMA:    EMA{length: 1},
		Result: []byte(`{"name":"ema","length":1}`),
	}

	d, _ := json.Marshal(c.EMA)

	assert.Equal(t, c.Result, d)
}

// func TestHMANew(t *testing.T) {
// 	cc := map[string]struct {
// 		WMA    WMA
// 		Result HMA
// 		Error  error
// 	}{
// 		"HMA throws an error": {
// 			Error: assert.AnError,
// 		},
// 		"Successful HMA creation": {
// 			WMA:    WMA{Length: 1},
// 			Result: HMA{WMA: WMA{Length: 1}},
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			h, err := NewHMA(c.WMA)
// 			if c.Error != nil {
// 				if c.Error == assert.AnError {
// 					assert.NotNil(t, err)
// 				} else {
// 					assert.Equal(t, c.Result, h)
// 				}
// 			} else {
// 				assert.Nil(t, err)
// 			}
// 		})
// 	}
// }

// func TestHMAValidation(t *testing.T) {
// 	cc := map[string]struct {
// 		WMA   WMA
// 		Error error
// 	}{
// 		"WMA returns an error": {
// 			WMA:   WMA{Length: -1},
// 			Error: assert.AnError,
// 		},
// 		"WMA not set": {
// 			Error: ErrMANotSet,
// 		},
// 		"Successful validation": {
// 			WMA: WMA{Length: 1},
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			h := HMA{WMA: c.WMA}
// 			err := h.Validate()
// 			if c.Error != nil {
// 				if c.Error == assert.AnError {
// 					assert.NotNil(t, err)
// 				} else {
// 					assert.Equal(t, c.Error, err)
// 				}
// 			} else {
// 				assert.Nil(t, err)
// 			}
// 		})
// 	}
// }

// func TestHMACalc(t *testing.T) {
// 	cc := map[string]struct {
// 		WMA    WMA
// 		Data   []decimal.Decimal
// 		Result decimal.Decimal
// 		Error  error
// 	}{
// 		"Insufficient amount of data points": {
// 			WMA: WMA{Length: 5},
// 			Data: []decimal.Decimal{
// 				decimal.NewFromInt(30),
// 			},
// 			Error: ErrInvalidDataPointCount,
// 		},
// 		"Successful calculation": {
// 			WMA: WMA{Length: 3},
// 			Data: []decimal.Decimal{
// 				decimal.NewFromInt(30),
// 				decimal.NewFromInt(31),
// 				decimal.NewFromInt(32),
// 				decimal.NewFromInt(30),
// 				decimal.NewFromInt(30),
// 				decimal.NewFromInt(31),
// 			},
// 			Result: decimal.NewFromFloat(31.5),
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			h := HMA{WMA: c.WMA}
// 			res, err := h.Calc(c.Data)
// 			if c.Error != nil {
// 				if c.Error == assert.AnError {
// 					assert.NotNil(t, err)
// 				} else {
// 					assert.Equal(t, c.Error, err)
// 				}
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result.String(), res.String())
// 			}
// 		})
// 	}
// }

// func TestHMACount(t *testing.T) {
// 	h := HMA{WMA: WMA{Length: 15}}
// 	assert.Equal(t, 29, h.Count())
// }

// func TestMACDNew(t *testing.T) {
// 	cc := map[string]struct {
// 		Source1 Source
// 		Source2 Source
// 		Result  MACD
// 		Error   error
// 	}{
// 		"MACD throws an error": {
// 			Error: assert.AnError,
// 		},
// 		"Successful MACD creation": {
// 			Source1: Source{WMA{Length: 1}},
// 			Source2: Source{WMA{Length: 1}},
// 			Result:  MACD{Source1: Source{WMA{Length: 1}}, Source2: Source{WMA{Length: 1}}},
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			m, err := NewMACD(c.Source1, c.Source2)
// 			if c.Error != nil {
// 				if c.Error == assert.AnError {
// 					assert.NotNil(t, err)
// 				} else {
// 					assert.Equal(t, c.Result, m)
// 				}
// 			} else {
// 				assert.Nil(t, err)
// 			}
// 		})
// 	}
// }

// func TestMACDValidation(t *testing.T) {
// 	cc := map[string]struct {
// 		Source1 Source
// 		Source2 Source
// 		Error   error
// 	}{
// 		"Source1 returns an error": {
// 			Source1: Source{EMA{Length: -1}},
// 			Source2: Source{EMA{Length: 1}},
// 			Error:   assert.AnError,
// 		},
// 		"Source2 returns an error": {
// 			Source1: Source{EMA{Length: 1}},
// 			Source2: Source{EMA{Length: -1}},
// 			Error:   assert.AnError,
// 		},
// 		"Successful validation": {
// 			Source1: Source{EMA{Length: 1}},
// 			Source2: Source{EMA{Length: 1}},
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			m := MACD{Source1: c.Source1, Source2: c.Source2}
// 			err := m.Validate()
// 			if c.Error != nil {
// 				if c.Error == assert.AnError {
// 					assert.NotNil(t, err)
// 				} else {
// 					assert.Equal(t, c.Error, err)
// 				}
// 			} else {
// 				assert.Nil(t, err)
// 			}
// 		})
// 	}
// }

// func TestMACDCalc(t *testing.T) {
// 	cc := map[string]struct {
// 		Source1 Source
// 		Source2 Source
// 		Data    []decimal.Decimal
// 		Result  decimal.Decimal
// 		Error   error
// 	}{
// 		"Source1 insufficient amount of data points": {
// 			Source1: Source{EMA{Length: 4}},
// 			Source2: Source{EMA{Length: 1}},
// 			Data: []decimal.Decimal{
// 				decimal.NewFromInt(30),
// 			},
// 			Error: ErrInvalidDataPointCount,
// 		},
// 		"Source2 insufficient amount of data points": {
// 			Source1: Source{EMA{Length: 1}},
// 			Source2: Source{EMA{Length: 4}},
// 			Data: []decimal.Decimal{
// 				decimal.NewFromInt(30),
// 			},
// 			Error: ErrInvalidDataPointCount,
// 		},
// 		"Source1 returns an error": {
// 			Source1: Source{IndicatorMock{}},
// 			Source2: Source{SMA{Length: 3}},
// 			Data: []decimal.Decimal{
// 				decimal.NewFromInt(30),
// 				decimal.NewFromInt(31),
// 				decimal.NewFromInt(32),
// 				decimal.NewFromInt(30),
// 				decimal.NewFromInt(31),
// 				decimal.NewFromInt(32),
// 			},
// 			Error: assert.AnError,
// 		},
// 		"Source2 returns an error": {
// 			Source1: Source{SMA{Length: 3}},
// 			Source2: Source{IndicatorMock{}},
// 			Data: []decimal.Decimal{
// 				decimal.NewFromInt(30),
// 				decimal.NewFromInt(31),
// 				decimal.NewFromInt(32),
// 				decimal.NewFromInt(30),
// 				decimal.NewFromInt(31),
// 				decimal.NewFromInt(32),
// 			},
// 			Error: assert.AnError,
// 		},
// 		"Successful calculation": {
// 			Source1: Source{SMA{Length: 2}},
// 			Source2: Source{SMA{Length: 3}},
// 			Data: []decimal.Decimal{
// 				decimal.NewFromInt(30),
// 				decimal.NewFromInt(31),
// 				decimal.NewFromInt(32),
// 				decimal.NewFromInt(30),
// 				decimal.NewFromInt(31),
// 				decimal.NewFromInt(32),
// 			},
// 			Result: decimal.NewFromFloat(0.5),
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			m := MACD{Source1: c.Source1, Source2: c.Source2}
// 			res, err := m.Calc(c.Data)
// 			if c.Error != nil {
// 				if c.Error == assert.AnError {
// 					assert.NotNil(t, err)
// 				} else {
// 					assert.Equal(t, c.Error, err)
// 				}
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result.String(), res.String())
// 			}
// 		})
// 	}
// }

// func TestMACDCount(t *testing.T) {
// 	m := MACD{Source1: Source{EMA{Length: 10}}, Source2: Source{EMA{Length: 1}}}
// 	assert.Equal(t, m.Count(), m.Source1.Count())

// 	m = MACD{Source1: Source{EMA{Length: 2}}, Source2: Source{EMA{Length: 9}}}
// 	assert.Equal(t, m.Count(), m.Source2.Count())
// }

func TestROCNew(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result ROC
		Error  error
	}{
		"ROC throws an error": {
			Error: assert.AnError,
		},
		"Successful ROC creation": {
			Length: 1,
			Result: ROC{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			r, err := NewROC(c.Length)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, r)
			}
		})
	}
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

			r := ROC{length: c.Length}
			err := r.validate()
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
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

			r := ROC{length: c.Length}
			res, err := r.Calc(c.Data)
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result.String(), res.String())
			}
		})
	}
}

func TestROCCount(t *testing.T) {
	r := ROC{length: 15}
	assert.Equal(t, 15, r.Count())
}

func TestROCUnmarshal(t *testing.T) {
	cc := map[string]struct {
		ByteArray []byte
		Result    ROC
		Error     error
	}{
		"Unmarshal throws an error": {
			ByteArray: []byte(`{"length": "down"}`),
			Error:     assert.AnError,
		},
		"Successful ROC unmarshal": {
			ByteArray: []byte(`{"name":"roc","length":1}`),
			Result:    ROC{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			r := ROC{}
			err := json.Unmarshal(c.ByteArray, &r)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, r)
			}
		})
	}
}

func TestROCMarshal(t *testing.T) {
	c := struct {
		ROC    ROC
		Result []byte
	}{
		ROC:    ROC{length: 1},
		Result: []byte(`{"name":"roc","length":1}`),
	}

	d, _ := json.Marshal(c.ROC)

	assert.Equal(t, c.Result, d)
}

func TestRSINew(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result RSI
		Error  error
	}{
		"RSI throws an error": {
			Error: assert.AnError,
		},
		"Successful RSI creation": {
			Length: 1,
			Result: RSI{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			r, err := NewRSI(c.Length)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, r)
			}
		})
	}
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

			r := RSI{length: c.Length}
			err := r.validate()
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
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

			r := RSI{length: c.Length}
			res, err := r.Calc(c.Data)
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result.String(), res.Round(14).String())
			}
		})
	}
}

func TestRSICount(t *testing.T) {
	r := RSI{length: 15}
	assert.Equal(t, 15, r.Count())
}

func TestRSIUnmarshal(t *testing.T) {
	cc := map[string]struct {
		ByteArray []byte
		Result    RSI
		Error     error
	}{
		"Unmarshal throws an error": {
			ByteArray: []byte(`{"length": "down"}`),
			Error:     assert.AnError,
		},
		"Successful RSI unmarshal": {
			ByteArray: []byte(`{"name":"rsi","length":1}`),
			Result:    RSI{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			r := RSI{}
			err := json.Unmarshal(c.ByteArray, &r)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, r)
			}
		})
	}
}

func TestRSIMarshal(t *testing.T) {
	c := struct {
		RSI    RSI
		Result []byte
	}{
		RSI:    RSI{length: 1},
		Result: []byte(`{"name":"rsi","length":1}`),
	}

	d, _ := json.Marshal(c.RSI)

	assert.Equal(t, c.Result, d)
}

func TestSMANew(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result SMA
		Error  error
	}{
		"SMA throws an error": {
			Error: assert.AnError,
		},
		"Successful SMA creation": {
			Length: 1,
			Result: SMA{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			s, err := NewSMA(c.Length)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, s)
			}
		})
	}
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

			s := SMA{length: c.Length}
			err := s.validate()
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
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

			s := SMA{length: c.Length}
			res, err := s.Calc(c.Data)
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result.String(), res.String())
			}
		})
	}
}

func TestSMACount(t *testing.T) {
	s := SMA{length: 15}
	assert.Equal(t, 15, s.Count())
}

func TestSMAUnmarshal(t *testing.T) {
	cc := map[string]struct {
		ByteArray []byte
		Result    SMA
		Error     error
	}{
		"Unmarshal throws an error": {
			ByteArray: []byte(`{"length": "down"}`),
			Error:     assert.AnError,
		},
		"Successful SMA unmarshal": {
			ByteArray: []byte(`{"name":"sma","length":1}`),
			Result:    SMA{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			s := SMA{}
			err := json.Unmarshal(c.ByteArray, &s)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, s)
			}
		})
	}
}

func TestSMAMarshal(t *testing.T) {
	c := struct {
		SMA    SMA
		Result []byte
	}{
		SMA:    SMA{length: 1},
		Result: []byte(`{"name":"sma","length":1}`),
	}

	d, _ := json.Marshal(c.SMA)

	assert.Equal(t, c.Result, d)
}

func TestStochNew(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result Stoch
		Error  error
	}{
		"Stoch throws an error": {
			Error: assert.AnError,
		},
		"Successful Stoch creation": {
			Length: 1,
			Result: Stoch{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			s, err := NewStoch(c.Length)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, s)
			}
		})
	}
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

			s := Stoch{length: c.Length}
			err := s.validate()

			if c.Error != nil {
				assert.Equal(t, c.Error, err)
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
		"Successful calculation v1": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(150),
				decimal.NewFromInt(125),
				decimal.NewFromInt(145),
			},
			Result: decimal.NewFromInt(80),
		},
		"Successful calculation v2": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(120),
				decimal.NewFromInt(145),
				decimal.NewFromInt(135),
			},
			Result: decimal.NewFromInt(60),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			s := Stoch{length: c.Length}
			res, err := s.Calc(c.Data)

			if c.Error != nil {
				assert.Equal(t, c.Error, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result.String(), res.String())
			}
		})
	}
}

func TestStochCount(t *testing.T) {
	s := Stoch{length: 15}
	assert.Equal(t, 15, s.Count())
}

func TestStochUnmarshal(t *testing.T) {
	cc := map[string]struct {
		ByteArray []byte
		Result    Stoch
		Error     error
	}{
		"Unmarshal throws an error": {
			ByteArray: []byte(`{"length": "down"}`),
			Error:     assert.AnError,
		},
		"Successful Stoch unmarshal": {
			ByteArray: []byte(`{"name":"stoch","length":1}`),
			Result:    Stoch{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			s := Stoch{}
			err := json.Unmarshal(c.ByteArray, &s)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, s)
			}
		})
	}
}

func TestStochMarshal(t *testing.T) {
	c := struct {
		Stoch  Stoch
		Result []byte
	}{
		Stoch:  Stoch{length: 1},
		Result: []byte(`{"name":"stoch","length":1}`),
	}

	d, _ := json.Marshal(c.Stoch)

	assert.Equal(t, c.Result, d)
}

// func TestWMANew(t *testing.T) {
// 	cc := map[string]struct {
// 		Length int
// 		Result WMA
// 		Error  error
// 	}{
// 		"WMA throws an error": {
// 			Error: assert.AnError,
// 		},
// 		"Successful WMA creation": {
// 			Length: 1,
// 			Result: WMA{Length: 1},
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			w, err := NewWMA(c.Length)
// 			if c.Error != nil {
// 				if c.Error == assert.AnError {
// 					assert.NotNil(t, err)
// 				} else {
// 					assert.Equal(t, c.Result, w)
// 				}
// 			} else {
// 				assert.Nil(t, err)
// 			}
// 		})
// 	}
// }

// func TestWMAValidation(t *testing.T) {
// 	cc := map[string]struct {
// 		Length int
// 		Error  error
// 	}{
// 		"Length cannot be less than 1": {
// 			Length: 0,
// 			Error:  ErrInvalidLength,
// 		},
// 		"Successful validation": {
// 			Length: 1,
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			w := WMA{Length: c.Length}
// 			err := w.Validate()
// 			if c.Error != nil {
// 				if c.Error == assert.AnError {
// 					assert.NotNil(t, err)
// 				} else {
// 					assert.Equal(t, c.Error, err)
// 				}
// 			} else {
// 				assert.Nil(t, err)
// 			}
// 		})
// 	}
// }

// func TestWMACalc(t *testing.T) {
// 	cc := map[string]struct {
// 		Length int
// 		Data   []decimal.Decimal
// 		Result decimal.Decimal
// 		Error  error
// 	}{
// 		"Insufficient amount of data points": {
// 			Length: 3,
// 			Data: []decimal.Decimal{
// 				decimal.NewFromInt(30),
// 			},
// 			Error: ErrInvalidDataPointCount,
// 		},
// 		"Successful calculation": {
// 			Length: 3,
// 			Data: []decimal.Decimal{
// 				decimal.NewFromInt(420),
// 				decimal.NewFromInt(420),
// 				decimal.NewFromInt(420),
// 				decimal.NewFromInt(30),
// 				decimal.NewFromInt(30),
// 				decimal.NewFromInt(32),
// 			},
// 			Result: decimal.NewFromFloat(31),
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			w := WMA{Length: c.Length}
// 			res, err := w.Calc(c.Data)
// 			if c.Error != nil {
// 				if c.Error == assert.AnError {
// 					assert.NotNil(t, err)
// 				} else {
// 					assert.Equal(t, c.Error, err)
// 				}
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result.String(), res.String())
// 			}
// 		})
// 	}
// }

// func TestWMACount(t *testing.T) {
// 	w := WMA{Length: 15}
// 	assert.Equal(t, 15, w.Count())
// }
