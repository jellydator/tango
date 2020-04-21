package indc

import (
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
		"Successfully Aroon creation threw an error when no values were provided": {
			Error: assert.AnError,
		},
		"Successful creation of Aroon": {
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
		"Successfully Aroon threw an ErrInvalidType with incorrect trend": {
			Trend:  "downn",
			Length: 5,
			Error:  ErrInvalidType,
		},
		"Successfully Aroon threw an ErrInvalidLength with less than 1 length": {
			Trend:  "down",
			Length: 0,
			Error:  ErrInvalidLength,
		},
		"Successful Aroon validation of trend parameter with 'up' value": {
			Trend:  "up",
			Length: 5,
		},
		"Successful Aroon validation of trend parameter with 'down' value": {
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
		"Successfully Aroon threw an ErrInvalidDataPointCount with insufficient amount of data points": {
			Trend:  "down",
			Length: 5,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataPointCount,
		},
		"Successful Aroon calculation with trend parameter set to 'up'": {
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
		"Successful Aroon calculation with trend parameter set to 'down'": {
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
		"Successfully Aroon unmarshal threw an error": {
			ByteArray: []byte(`{\"_"/`),
			Error:     assert.AnError,
		},
		"Successfully Aroon validate threw an error": {
			ByteArray: []byte(`{"trend":"upp","length":1}`),
			Error:     assert.AnError,
		},
		"Successful unmarshal of an Aroon": {
			ByteArray: []byte(`{"trend":"up","length":1}`),
			Result:    Aroon{trend: "up", length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			a := Aroon{}
			err := a.UnmarshalJSON(c.ByteArray)
			if c.Error != nil {
				if c.Error == assert.AnError {
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, c.Error, err)
				}
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, a)
			}
		})
	}
}

func TestAroonMarshal(t *testing.T) {
	a := Aroon{trend: "down", length: 1}
	r := []byte(`{"trend":"down","length":1}`)

	d, _ := a.MarshalJSON()

	assert.Equal(t, r, d)
}

func TestCCINew(t *testing.T) {
	cc := map[string]struct {
		Source Indicator
		Result CCI
		Error  error
	}{
		"Successfully CCI creation threw an error when no values were provided": {
			Error: assert.AnError,
		},
		"Successful creation of CCI": {
			Source: Aroon{trend: "down", length: 1},
			Result: CCI{source: Aroon{trend: "down", length: 1}},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			cci, err := NewCCI(c.Source)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, cci)
			}
		})
	}
}

func TestCCIValidation(t *testing.T) {
	cc := map[string]struct {
		Source Indicator
		Error  error
	}{
		"Successfully CCI threw an error when source wasn't provided": {
			Error: ErrSourceNotSet,
		},
		"Successfully CCI source validation threw an error": {
			Source: Aroon{length: -1, trend: "s"},
			Error:  assert.AnError,
		},
		"Successful CCI validation": {
			Source: Aroon{trend: "down", length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			cci := CCI{source: c.Source}
			err := cci.validate()
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
		Source Indicator
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Successfully CCI threw an ErrInvalidDataPointCount with insufficient amount of data points": {
			Source: EMA{length: 10},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataPointCount,
		},
		"Successfully CCI source threw an error": {
			Source: IndicatorMock{},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: assert.AnError,
		},
		"Successful CCI calculation with given source": {
			Source: SMA{length: 20},
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

			cci := CCI{source: c.Source}
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
	c := CCI{source: Aroon{trend: "down", length: 10}}
	assert.Equal(t, c.source.Count(), c.Count())
}

func TestCCIUnmarshal(t *testing.T) {
	cc := map[string]struct {
		ByteArray []byte
		Result    CCI
		Error     error
	}{
		"Successfully CCI unmarshal threw an error": {
			ByteArray: []byte(`{\-_-/}`),
			Error:     assert.AnError,
		},
		"Successfully CCI validate threw an error": {
			ByteArray: []byte(`{"trend":"up"}`),
			Error:     assert.AnError,
		},
		"Successful CCI unmarshal": {
			ByteArray: []byte(`{"source":{"name":"aroon","trend":"up","length":1}}`),
			Result:    CCI{Aroon{trend: "up", length: 1}},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			cci := CCI{}
			err := cci.UnmarshalJSON(c.ByteArray)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, cci)
			}
		})
	}
}

func TestCCIMarshal(t *testing.T) {
	cc := map[string]struct {
		CCI    CCI
		Result []byte
		Error  error
	}{
		"Successsfully CCI toJSON threw an error": {
			CCI:   CCI{source: IndicatorMock{}},
			Error: assert.AnError,
		},
		"Successfully CCI unmarshal threw an error": {
			CCI:    CCI{source: Aroon{trend: "down", length: 1}},
			Result: []byte(`{"source":{"name":"aroon","trend":"down","length":1}}`),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			d, err := c.CCI.MarshalJSON()
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, d)
			}
		})
	}
}

func TestDEMANew(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result DEMA
		Error  error
	}{
		"Successfully DEMA threw an error when no values were provided": {
			Error: assert.AnError,
		},
		"Successful creation of DEMA": {
			Length: 1,
			Result: DEMA{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			dm, err := NewDEMA(c.Length)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, dm)
			}
		})
	}
}

func TestDEMAValidation(t *testing.T) {
	cc := map[string]struct {
		Length int
		Error  error
	}{
		"Successfully DEMA threw an ErrInvalidLength with less than 1 length": {
			Length: 0,
			Error:  ErrInvalidLength,
		},
		"Successful DEMA validation": {
			Length: 1,
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			d := DEMA{length: c.Length}
			err := d.validate()
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
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
		"Successfully DEMA threw an ErrInvalidDataPointCount with insufficient amount of data points": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataPointCount,
		},
		"Successful DEMA calculation": {
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

			d := DEMA{length: c.Length}
			res, err := d.Calc(c.Data)
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result.String(), res.Round(14).String())
			}
		})
	}
}

func TestDEMACount(t *testing.T) {
	d := DEMA{length: 15}
	assert.Equal(t, 29, d.Count())
}

func TestDEMAUnmarshal(t *testing.T) {
	cc := map[string]struct {
		ByteArray []byte
		Result    DEMA
		Error     error
	}{
		"Successfully DEMA unmarshal threw an error": {
			ByteArray: []byte(`{\"_"/`),
			Error:     assert.AnError,
		},
		"Successfully DEMA validate threw an error": {
			ByteArray: []byte(`{"length":0}`),
			Error:     assert.AnError,
		},
		"Successful unmarshal of a DEMA": {
			ByteArray: []byte(`{"length":1}`),
			Result:    DEMA{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			dm := DEMA{}
			err := dm.UnmarshalJSON(c.ByteArray)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, dm)
			}
		})
	}
}

func TestDEMAMarshal(t *testing.T) {
	dm := DEMA{length: 1}
	r := []byte(`{"length":1}`)

	d, _ := dm.MarshalJSON()

	assert.Equal(t, r, d)
}

func TestEMANew(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result EMA
		Error  error
	}{
		"Successfully EMA creation threw an error when no values were provided": {
			Error: assert.AnError,
		},
		"Successful creation of EMA": {
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
		"Successfully EMA threw an ErrInvalidLength with less than 1 length": {
			Length: 0,
			Error:  ErrInvalidLength,
		},
		"Successful EMA validation": {
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
		"Successfully EMA threw an ErrInvalidDataPointCount with insufficient amount of data points": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataPointCount,
		},
		"Successful EMA calculation": {
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
		"Successfully EMA unmarshal threw an error": {
			ByteArray: []byte(`{\"_"/`),
			Error:     assert.AnError,
		},
		"Successfully EMA validate threw an error": {
			ByteArray: []byte(`{"length":0}`),
			Error:     assert.AnError,
		},
		"Successful unmarshal of a EMA": {
			ByteArray: []byte(`{"length":1}`),
			Result:    EMA{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			e := EMA{}
			err := e.UnmarshalJSON(c.ByteArray)
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
	e := EMA{length: 1}
	r := []byte(`{"length":1}`)

	d, _ := e.MarshalJSON()

	assert.Equal(t, r, d)
}

func TestHMANew(t *testing.T) {
	cc := map[string]struct {
		WMA    WMA
		Result HMA
		Error  error
	}{
		"Successfully HMA threw an error when no values were provided": {
			Error: assert.AnError,
		},
		"Successful HMA creation": {
			WMA:    WMA{length: 1},
			Result: HMA{wma: WMA{length: 1}},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			h, err := NewHMA(c.WMA)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, h)
			}
		})
	}
}

func TestHMAValidation(t *testing.T) {
	cc := map[string]struct {
		WMA   WMA
		Error error
	}{
		"Successfully HMA wma threw an error": {
			WMA:   WMA{length: -1},
			Error: assert.AnError,
		},
		"Successfully HMA threw an ErrMANotSet when WMA wasn't set": {
			Error: ErrMANotSet,
		},
		"Successful HMA validation": {
			WMA: WMA{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			h := HMA{wma: c.WMA}
			err := h.validate()
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
		WMA    WMA
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Successfully HMA threw an ErrInvalidDataPointCount with insufficient amount of data points": {
			WMA: WMA{length: 5},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataPointCount,
		},
		"Successful HMA calculation": {
			WMA: WMA{length: 3},
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

			h := HMA{wma: c.WMA}
			res, err := h.Calc(c.Data)
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result.String(), res.String())
			}
		})
	}
}

func TestHMACount(t *testing.T) {
	h := HMA{wma: WMA{length: 15}}
	assert.Equal(t, 29, h.Count())
}

func TestHMAUnmarshal(t *testing.T) {
	cc := map[string]struct {
		ByteArray []byte
		Result    HMA
		Error     error
	}{
		"Successfully HMA unmarshal threw an error": {
			ByteArray: []byte(`{\"_"/`),
			Error:     assert.AnError,
		},
		"Successfully HMA validate threw an error": {
			ByteArray: []byte(`{"length":0}`),
			Error:     assert.AnError,
		},
		"Successful unmarshal of a HMA": {
			ByteArray: []byte(`{"wma":{"length":1}}`),
			Result:    HMA{wma: WMA{length: 1}},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			h := HMA{}
			err := h.UnmarshalJSON(c.ByteArray)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, h)
			}
		})
	}
}

func TestHMAMarshal(t *testing.T) {
	h := HMA{wma: WMA{length: 1}}
	r := []byte(`{"wma":{"length":1}}`)

	d, _ := h.MarshalJSON()

	assert.Equal(t, r, d)
}

// func TestMACDNew(t *testing.T) {
// 	cc := map[string]struct {
// 		Source1 Indicator
// 		Source2 Indicator
// 		Result  MACD
// 		Error   error
// 	}{
// 		"MACD throws an error": {
// 			Error: assert.AnError,
// 		},
// 		"Successful MACD creation": {
// 			Source1: WMA{length: 1},
// 			Source2: WMA{length: 1},
// 			Result:  MACD{source1: WMA{length: 1}, source2: WMA{length: 1}},
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			m, err := NewMACD(c.Source1, c.Source2)
// 			if c.Error != nil {
// 				assert.NotNil(t, err)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result, m)
// 			}
// 		})
// 	}
// }

// func TestMACDValidation(t *testing.T) {
// 	cc := map[string]struct {
// 		Source1 Indicator
// 		Source2 Indicator
// 		Error   error
// 	}{
// 		"Source1 returns an error": {
// 			Source1: EMA{length: -1},
// 			Source2: EMA{length: 1},
// 			Error:   assert.AnError,
// 		},
// 		"Source2 returns an error": {
// 			Source1: EMA{length: 1},
// 			Source2: EMA{length: -1},
// 			Error:   assert.AnError,
// 		},
// 		"Successful validation": {
// 			Source1: EMA{length: 1},
// 			Source2: EMA{length: 1},
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			m := MACD{source1: c.Source1, source2: c.Source2}
// 			err := m.validate()
// 			if c.Error != nil {
// 				assert.NotNil(t, err)
// 			} else {
// 				assert.Nil(t, err)
// 			}
// 		})
// 	}
// }

// func TestMACDCalc(t *testing.T) {
// 	cc := map[string]struct {
// 		Source1 Indicator
// 		Source2 Indicator
// 		Data    []decimal.Decimal
// 		Result  decimal.Decimal
// 		Error   error
// 	}{
// 		"Source1 insufficient amount of data points": {
// 			Source1: EMA{length: 4},
// 			Source2: EMA{length: 1},
// 			Data: []decimal.Decimal{
// 				decimal.NewFromInt(30),
// 			},
// 			Error: ErrInvalidDataPointCount,
// 		},
// 		"Source2 insufficient amount of data points": {
// 			Source1: EMA{length: 1},
// 			Source2: EMA{length: 4},
// 			Data: []decimal.Decimal{
// 				decimal.NewFromInt(30),
// 			},
// 			Error: ErrInvalidDataPointCount,
// 		},
// 		"Source1 returns an error": {
// 			Source1: IndicatorMock{},
// 			Source2: SMA{length: 3},
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
// 			Source1: SMA{length: 3},
// 			Source2: IndicatorMock{},
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
// 			Source1: SMA{length: 2},
// 			Source2: SMA{length: 3},
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

// 			m := MACD{source1: c.Source1, source2: c.Source2}
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
// 	m := MACD{source1: EMA{length: 10}, source2: EMA{length: 1}}
// 	assert.Equal(t, m.Count(), m.source1.Count())

// 	m = MACD{source1: EMA{length: 2}, source2: EMA{length: 9}}
// 	assert.Equal(t, m.Count(), m.source2.Count())
// }

// func TestMACDUnmarshal(t *testing.T) {
// 	cc := map[string]struct {
// 		ByteArray []byte
// 		Result    MACD
// 		Error     error
// 	}{
// 		"Unmarshal throws an error": {
// 			ByteArray: []byte(`{}`),
// 			Error:     assert.AnError,
// 		},
// 		"Successful MACD unmarshal": {
// 			ByteArray: []byte(`{"name":"macd",
// 			"source1":{"name":"aroon","trend":"down","length":2},
// 			"source2":{"name":"cci","source":{"name":"ema", "length":2}}}`),
// 			Result: MACD{source1: Aroon{trend: "down", length: 2},
// 				source2: CCI{source: EMA{length: 2}}},
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			m := MACD{}
// 			err := json.Unmarshal(c.ByteArray, &m)
// 			if c.Error != nil {
// 				assert.NotNil(t, err)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result, m)
// 			}
// 		})
// 	}
// }

// func TestMACDMarshal(t *testing.T) {
// 	c := struct {
// 		MACD   MACD
// 		Result []byte
// 	}{
// 		MACD: MACD{source1: Aroon{trend: "down", length: 2},
// 			source2: CCI{source: EMA{length: 2}}},
// 		Result: []byte(`{"source1":{"name":"aroon","trend":"down","length":2},"source2":{"name":"cci","source":{"name":"ema","length":2}}}`),
// 	}

// 	d, _ := json.Marshal(c.MACD)

// 	assert.Equal(t, c.Result, d)
// }

// func TestROCNew(t *testing.T) {
// 	cc := map[string]struct {
// 		Length int
// 		Result ROC
// 		Error  error
// 	}{
// 		"ROC throws an error": {
// 			Error: assert.AnError,
// 		},
// 		"Successful ROC creation": {
// 			Length: 1,
// 			Result: ROC{length: 1},
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			r, err := NewROC(c.Length)
// 			if c.Error != nil {
// 				assert.NotNil(t, err)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result, r)
// 			}
// 		})
// 	}
// }

// func TestROCValidation(t *testing.T) {
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

// 			r := ROC{length: c.Length}
// 			err := r.validate()
// 			if c.Error != nil {
// 				assert.Equal(t, c.Error, err)
// 			} else {
// 				assert.Nil(t, err)
// 			}
// 		})
// 	}
// }

// func TestROCCalc(t *testing.T) {
// 	cc := map[string]struct {
// 		Length int
// 		Data   []decimal.Decimal
// 		Result decimal.Decimal
// 		Error  error
// 	}{
// "Successfully ROC threw an ErrInvalidDataPointCount with insufficient amount of data points": {
// 			Length: 3,
// 			Data: []decimal.Decimal{
// 				decimal.NewFromInt(30),
// 			},
// 			Error: ErrInvalidDataPointCount,
// 		},
// 		"Successful calculation": {
// 			Length: 5,
// 			Data: []decimal.Decimal{
// 				decimal.NewFromInt(7),
// 				decimal.NewFromInt(420),
// 				decimal.NewFromInt(420),
// 				decimal.NewFromInt(420),
// 				decimal.NewFromInt(10),
// 			},
// 			Result: decimal.NewFromFloat(42.85714285714286),
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			r := ROC{length: c.Length}
// 			res, err := r.Calc(c.Data)
// 			if c.Error != nil {
// 				assert.Equal(t, c.Error, err)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result.String(), res.String())
// 			}
// 		})
// 	}
// }

// func TestROCCount(t *testing.T) {
// 	r := ROC{length: 15}
// 	assert.Equal(t, 15, r.Count())
// }

// func TestROCUnmarshal(t *testing.T) {
// 	cc := map[string]struct {
// 		ByteArray []byte
// 		Result    ROC
// 		Error     error
// 	}{
// 		"Unmarshal throws an error": {
// 			ByteArray: []byte(`{"length": "down"}`),
// 			Error:     assert.AnError,
// 		},
// 		"Successful ROC unmarshal": {
// 			ByteArray: []byte(`{"name":"roc","length":1}`),
// 			Result:    ROC{length: 1},
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			r := ROC{}
// 			err := json.Unmarshal(c.ByteArray, &r)
// 			if c.Error != nil {
// 				assert.NotNil(t, err)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result, r)
// 			}
// 		})
// 	}
// }

// func TestROCMarshal(t *testing.T) {
// 	c := struct {
// 		ROC    ROC
// 		Result []byte
// 	}{
// 		ROC:    ROC{length: 1},
// 		Result: []byte(`{"name":"roc","length":1}`),
// 	}

// 	d, _ := json.Marshal(c.ROC)

// 	assert.Equal(t, c.Result, d)
// }

// func TestRSINew(t *testing.T) {
// 	cc := map[string]struct {
// 		Length int
// 		Result RSI
// 		Error  error
// 	}{
// 		"RSI throws an error": {
// 			Error: assert.AnError,
// 		},
// 		"Successful RSI creation": {
// 			Length: 1,
// 			Result: RSI{length: 1},
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			r, err := NewRSI(c.Length)
// 			if c.Error != nil {
// 				assert.NotNil(t, err)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result, r)
// 			}
// 		})
// 	}
// }

// func TestRSIValidation(t *testing.T) {
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

// 			r := RSI{length: c.Length}
// 			err := r.validate()
// 			if c.Error != nil {
// 				assert.Equal(t, c.Error, err)
// 			} else {
// 				assert.Nil(t, err)
// 			}
// 		})
// 	}
// }

// func TestRSICalc(t *testing.T) {
// 	cc := map[string]struct {
// 		Length int
// 		Data   []decimal.Decimal
// 		Result decimal.Decimal
// 		Error  error
// 	}{
// "Successfully RSI threw an ErrInvalidDataPointCount with insufficient amount of data points": {
// 			Length: 3,
// 			Data: []decimal.Decimal{
// 				decimal.NewFromInt(30),
// 			},
// 			Error: ErrInvalidDataPointCount,
// 		},
// 		"Successful calculation": {
// 			Length: 14,
// 			Data: []decimal.Decimal{
// 				decimal.NewFromFloat32(44.34),
// 				decimal.NewFromFloat32(44.09),
// 				decimal.NewFromFloat32(44.15),
// 				decimal.NewFromFloat32(43.61),
// 				decimal.NewFromFloat32(44.33),
// 				decimal.NewFromFloat32(44.83),
// 				decimal.NewFromFloat32(45.10),
// 				decimal.NewFromFloat32(45.42),
// 				decimal.NewFromFloat32(45.84),
// 				decimal.NewFromFloat32(46.08),
// 				decimal.NewFromFloat32(45.89),
// 				decimal.NewFromFloat32(46.03),
// 				decimal.NewFromFloat32(45.61),
// 				decimal.NewFromFloat32(46.28),
// 			},
// 			Result: decimal.NewFromFloat(70.46413502109705),
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			r := RSI{length: c.Length}
// 			res, err := r.Calc(c.Data)
// 			if c.Error != nil {
// 				assert.Equal(t, c.Error, err)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result.String(), res.Round(14).String())
// 			}
// 		})
// 	}
// }

// func TestRSICount(t *testing.T) {
// 	r := RSI{length: 15}
// 	assert.Equal(t, 15, r.Count())
// }

// func TestRSIUnmarshal(t *testing.T) {
// 	cc := map[string]struct {
// 		ByteArray []byte
// 		Result    RSI
// 		Error     error
// 	}{
// 		"Unmarshal throws an error": {
// 			ByteArray: []byte(`{"length": "down"}`),
// 			Error:     assert.AnError,
// 		},
// 		"Successful RSI unmarshal": {
// 			ByteArray: []byte(`{"name":"rsi","length":1}`),
// 			Result:    RSI{length: 1},
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			r := RSI{}
// 			err := json.Unmarshal(c.ByteArray, &r)
// 			if c.Error != nil {
// 				assert.NotNil(t, err)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result, r)
// 			}
// 		})
// 	}
// }

// func TestRSIMarshal(t *testing.T) {
// 	c := struct {
// 		RSI    RSI
// 		Result []byte
// 	}{
// 		RSI:    RSI{length: 1},
// 		Result: []byte(`{"name":"rsi","length":1}`),
// 	}

// 	d, _ := json.Marshal(c.RSI)

// 	assert.Equal(t, c.Result, d)
// }

// func TestSMANew(t *testing.T) {
// 	cc := map[string]struct {
// 		Length int
// 		Result SMA
// 		Error  error
// 	}{
// 		"SMA throws an error": {
// 			Error: assert.AnError,
// 		},
// 		"Successful SMA creation": {
// 			Length: 1,
// 			Result: SMA{length: 1},
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			s, err := NewSMA(c.Length)
// 			if c.Error != nil {
// 				assert.NotNil(t, err)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result, s)
// 			}
// 		})
// 	}
// }

// func TestSMAValidation(t *testing.T) {
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

// 			s := SMA{length: c.Length}
// 			err := s.validate()
// 			if c.Error != nil {
// 				assert.Equal(t, c.Error, err)
// 			} else {
// 				assert.Nil(t, err)
// 			}
// 		})
// 	}
// }

// func TestSMACalc(t *testing.T) {
// 	cc := map[string]struct {
// 		Length int
// 		Data   []decimal.Decimal
// 		Result decimal.Decimal
// 		Error  error
// 	}{
// "Successfully SMA threw an ErrInvalidDataPointCount with insufficient amount of data points": {
// 			Length: 3,
// 			Data: []decimal.Decimal{
// 				decimal.NewFromInt(30),
// 			},
// 			Error: ErrInvalidDataPointCount,
// 		},
// 		"Successful calculation": {
// 			Length: 3,
// 			Data: []decimal.Decimal{
// 				decimal.NewFromInt(30),
// 				decimal.NewFromInt(31),
// 				decimal.NewFromInt(32),
// 			},
// 			Result: decimal.NewFromInt(31),
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			s := SMA{length: c.Length}
// 			res, err := s.Calc(c.Data)
// 			if c.Error != nil {
// 				assert.Equal(t, c.Error, err)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result.String(), res.String())
// 			}
// 		})
// 	}
// }

// func TestSMACount(t *testing.T) {
// 	s := SMA{length: 15}
// 	assert.Equal(t, 15, s.Count())
// }

// func TestSMAUnmarshal(t *testing.T) {
// 	cc := map[string]struct {
// 		ByteArray []byte
// 		Result    SMA
// 		Error     error
// 	}{
// 		"Unmarshal throws an error": {
// 			ByteArray: []byte(`{"length": "down"}`),
// 			Error:     assert.AnError,
// 		},
// 		"Successful SMA unmarshal": {
// 			ByteArray: []byte(`{"name":"sma","length":1}`),
// 			Result:    SMA{length: 1},
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			s := SMA{}
// 			err := json.Unmarshal(c.ByteArray, &s)
// 			if c.Error != nil {
// 				assert.NotNil(t, err)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result, s)
// 			}
// 		})
// 	}
// }

// func TestSMAMarshal(t *testing.T) {
// 	c := struct {
// 		SMA    SMA
// 		Result []byte
// 	}{
// 		SMA:    SMA{length: 1},
// 		Result: []byte(`{"name":"sma","length":1}`),
// 	}

// 	d, _ := json.Marshal(c.SMA)

// 	assert.Equal(t, c.Result, d)
// }

// func TestStochNew(t *testing.T) {
// 	cc := map[string]struct {
// 		Length int
// 		Result Stoch
// 		Error  error
// 	}{
// 		"Stoch throws an error": {
// 			Error: assert.AnError,
// 		},
// 		"Successful Stoch creation": {
// 			Length: 1,
// 			Result: Stoch{length: 1},
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			s, err := NewStoch(c.Length)
// 			if c.Error != nil {
// 				assert.NotNil(t, err)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result, s)
// 			}
// 		})
// 	}
// }

// func TestStochValidation(t *testing.T) {

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

// 			s := Stoch{length: c.Length}
// 			err := s.validate()

// 			if c.Error != nil {
// 				assert.Equal(t, c.Error, err)
// 			} else {
// 				assert.Nil(t, err)
// 			}
// 		})
// 	}
// }

// func TestStochCalc(t *testing.T) {
// 	cc := map[string]struct {
// 		Length int
// 		Data   []decimal.Decimal
// 		Result decimal.Decimal
// 		Error  error
// 	}{
// "Successfully Stoch threw an ErrInvalidDataPointCount with insufficient amount of data points": {
// 			Length: 3,
// 			Data: []decimal.Decimal{
// 				decimal.NewFromInt(30),
// 			},
// 			Error: ErrInvalidDataPointCount,
// 		},
// 		"Successful calculation v1": {
// 			Length: 3,
// 			Data: []decimal.Decimal{
// 				decimal.NewFromInt(150),
// 				decimal.NewFromInt(125),
// 				decimal.NewFromInt(145),
// 			},
// 			Result: decimal.NewFromInt(80),
// 		},
// 		"Successful calculation v2": {
// 			Length: 3,
// 			Data: []decimal.Decimal{
// 				decimal.NewFromInt(120),
// 				decimal.NewFromInt(145),
// 				decimal.NewFromInt(135),
// 			},
// 			Result: decimal.NewFromInt(60),
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			s := Stoch{length: c.Length}
// 			res, err := s.Calc(c.Data)

// 			if c.Error != nil {
// 				assert.Equal(t, c.Error, err)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result.String(), res.String())
// 			}
// 		})
// 	}
// }

// func TestStochCount(t *testing.T) {
// 	s := Stoch{length: 15}
// 	assert.Equal(t, 15, s.Count())
// }

// func TestStochUnmarshal(t *testing.T) {
// 	cc := map[string]struct {
// 		ByteArray []byte
// 		Result    Stoch
// 		Error     error
// 	}{
// 		"Unmarshal throws an error": {
// 			ByteArray: []byte(`{"length": "down"}`),
// 			Error:     assert.AnError,
// 		},
// 		"Successful Stoch unmarshal": {
// 			ByteArray: []byte(`{"name":"stoch","length":1}`),
// 			Result:    Stoch{length: 1},
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			s := Stoch{}
// 			err := json.Unmarshal(c.ByteArray, &s)
// 			if c.Error != nil {
// 				assert.NotNil(t, err)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result, s)
// 			}
// 		})
// 	}
// }

// func TestStochMarshal(t *testing.T) {
// 	c := struct {
// 		Stoch  Stoch
// 		Result []byte
// 	}{
// 		Stoch:  Stoch{length: 1},
// 		Result: []byte(`{"name":"stoch","length":1}`),
// 	}

// 	d, _ := json.Marshal(c.Stoch)

// 	assert.Equal(t, c.Result, d)
// }

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
// 			Result: WMA{length: 1},
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			w, err := NewWMA(c.Length)
// 			if c.Error != nil {
// 				assert.NotNil(t, err)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result, w)
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

// 			w := WMA{length: c.Length}
// 			err := w.validate()
// 			if c.Error != nil {
// 				assert.Equal(t, c.Error, err)
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
// "Successfully WMA threw an ErrInvalidDataPointCount with insufficient amount of data points": {
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

// 			w := WMA{length: c.Length}
// 			res, err := w.Calc(c.Data)
// 			if c.Error != nil {
// 				assert.Equal(t, c.Error, err)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result.String(), res.String())
// 			}
// 		})
// 	}
// }

// func TestWMACount(t *testing.T) {
// 	w := WMA{length: 15}
// 	assert.Equal(t, 15, w.Count())
// }

// func TestWMAUnmarshal(t *testing.T) {
// 	cc := map[string]struct {
// 		ByteArray []byte
// 		Result    WMA
// 		Error     error
// 	}{
// 		"Unmarshal throws an error": {
// 			ByteArray: []byte(`{"length": "down"}`),
// 			Error:     assert.AnError,
// 		},
// 		"Successful WMA unmarshal": {
// 			ByteArray: []byte(`{"name":"wma","length":1}`),
// 			Result:    WMA{length: 1},
// 		},
// 	}

// 	for cn, c := range cc {
// 		c := c
// 		t.Run(cn, func(t *testing.T) {
// 			t.Parallel()

// 			w := WMA{}
// 			err := json.Unmarshal(c.ByteArray, &w)
// 			if c.Error != nil {
// 				assert.NotNil(t, err)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, c.Result, w)
// 			}
// 		})
// 	}
// }

// func TestWMAMarshal(t *testing.T) {
// 	c := struct {
// 		WMA    WMA
// 		Result []byte
// 	}{
// 		WMA:    WMA{length: 1},
// 		Result: []byte(`{"name":"wma","length":1}`),
// 	}

// 	d, _ := json.Marshal(c.WMA)

// 	assert.Equal(t, c.Result, d)
// }
