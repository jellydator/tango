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

func (im IndicatorMock) namedMarshalJSON() ([]byte, error) { return nil, assert.AnError }

func Test_NewAroon(t *testing.T) {
	cc := map[string]struct {
		Trend  String
		Length int
		Result Aroon
		Error  error
	}{
		"Successfully Aroon creation returned an error when no parameters were provided": {
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

func TestAroon_validate(t *testing.T) {
	cc := map[string]struct {
		Trend  String
		Length int
		Error  error
	}{
		"Successfully Aroon returned an ErrInvalidType with incorrect trend": {
			Trend:  "downn",
			Length: 5,
			Error:  assert.AnError,
		},
		"Successfully Aroon returned an ErrInvalidLength with less than 1 length": {
			Trend:  "down",
			Length: 0,
			Error:  ErrInvalidLength,
		},
		"Successful Aroon validation of trend parameter with 'up' trend": {
			Trend:  "up",
			Length: 5,
		},
		"Successful Aroon validation of trend parameter with 'down' trend": {
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

func TestAroon_Calc(t *testing.T) {
	cc := map[string]struct {
		Trend  String
		Length int
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Successfully Aroon returned an ErrInvalidDataSize with insufficient amount of data points": {
			Trend:  "down",
			Length: 5,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
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

func TestAroon_Count(t *testing.T) {
	a := Aroon{trend: "down", length: 5}
	assert.Equal(t, 5, a.Count())
}

func TestAroon_UnmarshalJSON(t *testing.T) {
	cc := map[string]struct {
		ByteArray []byte
		Result    Aroon
		Error     error
	}{
		"Successfully Aroon unmarshal returned an error": {
			ByteArray: []byte(`{\"_"/`),
			Error:     assert.AnError,
		},
		"Successfully Aroon validate returned an error": {
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

func TestAroon_MarshalJSON(t *testing.T) {
	a := Aroon{trend: "down", length: 1}
	r := []byte(`{"trend":"down","length":1}`)

	d, _ := a.MarshalJSON()

	assert.Equal(t, r, d)
}

func TestAroon_namedMarshalJSON(t *testing.T) {
	a := Aroon{trend: "down", length: 1}
	r := []byte(`{"name":"aroon","trend":"down","length":1}`)

	d, _ := a.namedMarshalJSON()

	assert.Equal(t, r, d)
}

func Test_NewCCI(t *testing.T) {
	cc := map[string]struct {
		Source Indicator
		Result CCI
		Error  error
	}{
		"Successfully CCI creation returned an error when no parameters were provided": {
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

func TestCCI_validate(t *testing.T) {
	cc := map[string]struct {
		Source Indicator
		Error  error
	}{
		"Successfully CCI returned an error when invalid source was provided": {
			Error: ErrInvalidSource,
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

func TestCCI_Calc(t *testing.T) {
	cc := map[string]struct {
		Source Indicator
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Successfully CCI returned an ErrInvalidDataSize with insufficient amount of data points": {
			Source: EMA{length: 10},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successfully CCI source returned an error": {
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

func TestCCI_Count(t *testing.T) {
	c := CCI{source: Aroon{trend: "down", length: 10}}
	assert.Equal(t, c.source.Count(), c.Count())
}

func TestCCI_UnmarshalJSON(t *testing.T) {
	cc := map[string]struct {
		ByteArray []byte
		Result    CCI
		Error     error
	}{
		"Successfully CCI unmarshal returned an error": {
			ByteArray: []byte(`{\-_-/}`),
			Error:     assert.AnError,
		},
		"Successfully CCI fromJSON returned an error": {
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

func TestCCI_MarshalJSON(t *testing.T) {
	cc := map[string]struct {
		CCI    CCI
		Result []byte
		Error  error
	}{
		"Successfully CCI source marshal returned an error": {
			CCI:   CCI{source: IndicatorMock{}},
			Error: assert.AnError,
		},
		"Successful CCI unmarshal": {
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

func TestCCI_namedMarshalJSON(t *testing.T) {
	cc := map[string]struct {
		CCI    CCI
		Result []byte
		Error  error
	}{
		"Successfully CCI source marshal returned an error": {
			CCI:   CCI{source: IndicatorMock{}},
			Error: assert.AnError,
		},
		"Successful CCI unmarshal": {
			CCI:    CCI{source: Aroon{trend: "down", length: 1}},
			Result: []byte(`{"name":"cci","source":{"name":"aroon","trend":"down","length":1}}`),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			d, err := c.CCI.namedMarshalJSON()
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, d)
			}
		})
	}
}

func Test_NewDEMA(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result DEMA
		Error  error
	}{
		"Successfully DEMA returned an error when no parameters were provided": {
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

func TestDEMA_validate(t *testing.T) {
	cc := map[string]struct {
		Length int
		Error  error
	}{
		"Successfully DEMA returned an ErrInvalidLength with less than 1 length": {
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

func TestDEMA_Calc(t *testing.T) {
	cc := map[string]struct {
		Length int
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Successfully DEMA returned an ErrInvalidDataSize with insufficient amount of data points": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
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

func TestDEMA_Count(t *testing.T) {
	d := DEMA{length: 15}
	assert.Equal(t, 29, d.Count())
}

func TestDEMA_UnmarshalJSON(t *testing.T) {
	cc := map[string]struct {
		ByteArray []byte
		Result    DEMA
		Error     error
	}{
		"Successfully DEMA unmarshal returned an error": {
			ByteArray: []byte(`{\"_"/`),
			Error:     assert.AnError,
		},
		"Successfully DEMA validate returned an error": {
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

func TestDEMA_MarshalJSON(t *testing.T) {
	dm := DEMA{length: 1}
	r := []byte(`{"length":1}`)

	d, _ := dm.MarshalJSON()

	assert.Equal(t, r, d)
}

func TestDEMA_namedMarshalJSON(t *testing.T) {
	dm := DEMA{length: 1}
	r := []byte(`{"name":"dema","length":1}`)

	d, _ := dm.namedMarshalJSON()

	assert.Equal(t, r, d)
}

func Test_NewEMA(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result EMA
		Error  error
	}{
		"Successfully EMA creation returned an error when no parameters were provided": {
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

func TestEMA_validate(t *testing.T) {
	cc := map[string]struct {
		Length int
		Error  error
	}{
		"Successfully EMA returned an ErrInvalidLength with less than 1 length": {
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

func TestEMA_Calc(t *testing.T) {
	cc := map[string]struct {
		Length int
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Successfully EMA returned an ErrInvalidDataSize with insufficient amount of data points": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
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

func TestEMA_Count(t *testing.T) {
	e := EMA{length: 15}
	assert.Equal(t, 29, e.Count())
}

func TestEMA_multiplier(t *testing.T) {
	e := EMA{length: 3}
	assert.Equal(t, decimal.NewFromFloat(0.5), e.multiplier())
}

func TestEMA_UnmarshalJSON(t *testing.T) {
	cc := map[string]struct {
		ByteArray []byte
		Result    EMA
		Error     error
	}{
		"Successfully EMA unmarshal returned an error": {
			ByteArray: []byte(`{\"_"/`),
			Error:     assert.AnError,
		},
		"Successfully EMA validate returned an error": {
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

func TestEMA_MarshalJSON(t *testing.T) {
	e := EMA{length: 1}
	r := []byte(`{"length":1}`)

	d, _ := e.MarshalJSON()

	assert.Equal(t, r, d)
}

func TestEMA_namedMarshalJSON(t *testing.T) {
	e := EMA{length: 1}
	r := []byte(`{"name":"ema","length":1}`)

	d, _ := e.namedMarshalJSON()

	assert.Equal(t, r, d)
}

func Test_NewHMA(t *testing.T) {
	cc := map[string]struct {
		WMA    WMA
		Result HMA
		Error  error
	}{
		"Successfully HMA returned an error when no parameters were provided": {
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

func TestHMA_validate(t *testing.T) {
	cc := map[string]struct {
		WMA   WMA
		Error error
	}{
		"Successfully HMA wma returned an error": {
			WMA:   WMA{length: -1},
			Error: assert.AnError,
		},
		"Successfully HMA returned an error when WMA wasn't set": {
			Error: assert.AnError,
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
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestHMA_Calc(t *testing.T) {
	cc := map[string]struct {
		WMA    WMA
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Successfully HMA returned an ErrInvalidDataSize with insufficient amount of data points": {
			WMA: WMA{length: 5},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
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

func TestHMA_Count(t *testing.T) {
	h := HMA{wma: WMA{length: 15}}
	assert.Equal(t, 29, h.Count())
}

func TestHMA_UnmarshalJSON(t *testing.T) {
	cc := map[string]struct {
		ByteArray []byte
		Result    HMA
		Error     error
	}{
		"Successfully HMA unmarshal returned an error": {
			ByteArray: []byte(`{\"_"/`),
			Error:     assert.AnError,
		},
		"Successfully HMA validate returned an error": {
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

func TestHMA_MarshalJSON(t *testing.T) {
	h := HMA{wma: WMA{length: 1}}
	r := []byte(`{"wma":{"length":1}}`)

	d, _ := h.MarshalJSON()

	assert.Equal(t, r, d)
}

func TestHMA_namedMarshalJSON(t *testing.T) {
	h := HMA{wma: WMA{length: 1}}
	r := []byte(`{"name":"hma","wma":{"length":1}}`)

	d, _ := h.namedMarshalJSON()

	assert.Equal(t, r, d)
}

func Test_NewMACD(t *testing.T) {
	cc := map[string]struct {
		Source1 Indicator
		Source2 Indicator
		Result  MACD
		Error   error
	}{
		"Successfully MACD creation returned an error when no parameters were provided": {
			Error: assert.AnError,
		},
		"Successful creation of MACD": {
			Source1: WMA{length: 1},
			Source2: WMA{length: 1},
			Result:  MACD{source1: WMA{length: 1}, source2: WMA{length: 1}},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			m, err := NewMACD(c.Source1, c.Source2)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, m)
			}
		})
	}
}

func TestMACD_validate(t *testing.T) {
	cc := map[string]struct {
		Source1 Indicator
		Source2 Indicator
		Error   error
	}{
		"Successfully MACD returned an error when source1 wasn't provided": {
			Source1: EMA{length: 1},
			Error:   ErrInvalidSource,
		},
		"Successfully MACD returned an error when source2 wasn't provided": {
			Source2: EMA{length: 1},
			Error:   ErrInvalidSource,
		},
		"Successful MACD validation": {
			Source1: EMA{length: 1},
			Source2: EMA{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			m := MACD{source1: c.Source1, source2: c.Source2}
			err := m.validate()
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestMACD_Calc(t *testing.T) {
	cc := map[string]struct {
		Source1 Indicator
		Source2 Indicator
		Data    []decimal.Decimal
		Result  decimal.Decimal
		Error   error
	}{
		"Successfully MACD returned an ErrInvalidDataSize with insufficient amount of data points for source1": {
			Source1: EMA{length: 4},
			Source2: EMA{length: 1},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successfully MACD returned an ErrInvalidDataSize with insufficient amount of data points for source2": {
			Source1: EMA{length: 1},
			Source2: EMA{length: 4},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successfully MACD source1 returned an error": {
			Source1: IndicatorMock{},
			Source2: SMA{length: 3},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
				decimal.NewFromInt(31),
				decimal.NewFromInt(32),
				decimal.NewFromInt(30),
				decimal.NewFromInt(31),
				decimal.NewFromInt(32),
			},
			Error: assert.AnError,
		},
		"Successfully MACD source2 returned an error": {
			Source1: SMA{length: 3},
			Source2: IndicatorMock{},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
				decimal.NewFromInt(31),
				decimal.NewFromInt(32),
				decimal.NewFromInt(30),
				decimal.NewFromInt(31),
				decimal.NewFromInt(32),
			},
			Error: assert.AnError,
		},
		"Successful MACD calculation with given sources": {
			Source1: SMA{length: 2},
			Source2: SMA{length: 3},
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

			m := MACD{source1: c.Source1, source2: c.Source2}
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

func TestMACD_Count(t *testing.T) {
	m := MACD{source1: EMA{length: 10}, source2: EMA{length: 1}}
	assert.Equal(t, m.Count(), m.source1.Count())

	m = MACD{source1: EMA{length: 2}, source2: EMA{length: 9}}
	assert.Equal(t, m.Count(), m.source2.Count())
}

func TestMACD_UnmarshalJSON(t *testing.T) {
	cc := map[string]struct {
		ByteArray []byte
		Result    MACD
		Error     error
	}{
		"Successfully MACD unmarshal returned an error": {
			ByteArray: []byte(`{\-_-/}`),
			Error:     assert.AnError,
		},
		"Successfully MACD fromJSON returned an error with invalid source1 parameters": {
			ByteArray: []byte(`{"source1":{"name":"aroon","trend":"dsown","length":2},
			"source2":{"name":"cci","source":{"name":"ema", "length":2}}}`),
			Error: assert.AnError,
		},
		"Successfully MACD fromJSON returned an error with invalid source2 parameters": {
			ByteArray: []byte(`{"source1":{"name":"aroon","trend":"down","length":2},
			"source2":{"name":"ccis","source":{"name":"ema", "length":2}}}`),
			Error: assert.AnError,
		},
		"Successful MACD unmarshal": {
			ByteArray: []byte(`{"source1":{"name":"aroon","trend":"down","length":2},
			"source2":{"name":"cci","source":{"name":"ema", "length":2}}}`),
			Result: MACD{source1: Aroon{trend: "down", length: 2},
				source2: CCI{source: EMA{length: 2}}},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			m := MACD{}
			err := m.UnmarshalJSON(c.ByteArray)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, m)
			}
		})
	}
}

func TestMACD_MarshalJSON(t *testing.T) {
	cc := map[string]struct {
		MACD   MACD
		Result []byte
		Error  error
	}{
		"Successfully MACD source1 marshal returned an error": {
			MACD: MACD{source1: IndicatorMock{},
				source2: CCI{source: EMA{length: 2}}},
			Error: assert.AnError,
		},
		"Successfully MACD source2 marshal returned an error": {
			MACD: MACD{source1: Aroon{trend: "down", length: 2},
				source2: IndicatorMock{}},
			Error: assert.AnError,
		},
		"Successful MACD unmarshal": {
			MACD: MACD{source1: Aroon{trend: "down", length: 2},
				source2: CCI{source: EMA{length: 2}}},
			Result: []byte(`{"source1":{"name":"aroon","trend":"down","length":2},"source2":{"name":"cci","source":{"name":"ema","length":2}}}`),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			d, err := c.MACD.MarshalJSON()
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, d)
			}
		})
	}
}

func TestMACD_namedMarshalJSON(t *testing.T) {
	cc := map[string]struct {
		MACD   MACD
		Result []byte
		Error  error
	}{
		"Successfully MACD source1 marshal returned an error": {
			MACD: MACD{source1: IndicatorMock{},
				source2: CCI{source: EMA{length: 2}}},
			Error: assert.AnError,
		},
		"Successfully MACD source2 marshal returned an error": {
			MACD: MACD{source1: Aroon{trend: "down", length: 2},
				source2: IndicatorMock{}},
			Error: assert.AnError,
		},
		"Successful MACD unmarshal": {
			MACD: MACD{source1: Aroon{trend: "down", length: 2},
				source2: CCI{source: EMA{length: 2}}},
			Result: []byte(`{"name":"macd","source1":{"name":"aroon","trend":"down","length":2},"source2":{"name":"cci","source":{"name":"ema","length":2}}}`),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			d, err := c.MACD.namedMarshalJSON()
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, d)
			}
		})
	}
}

func Test_NewROC(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result ROC
		Error  error
	}{
		"Successfully ROC returned an error when no parameters were provided": {
			Error: assert.AnError,
		},
		"Successful creation of ROC": {
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

func TestROC_validate(t *testing.T) {
	cc := map[string]struct {
		Length int
		Error  error
	}{
		"Successfully ROC returned an ErrInvalidLength with less than 1 length": {
			Length: 0,
			Error:  ErrInvalidLength,
		},
		"Successful ROC validation": {
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

func TestROC_Calc(t *testing.T) {
	cc := map[string]struct {
		Length int
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Successfully ROC returned an ErrInvalidDataSize with insufficient amount of data points": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful ROC calculation": {
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

func TestROC_Count(t *testing.T) {
	r := ROC{length: 15}
	assert.Equal(t, 15, r.Count())
}

func TestROC_UnmarshalJSON(t *testing.T) {
	cc := map[string]struct {
		ByteArray []byte
		Result    ROC
		Error     error
	}{
		"Successfully ROC unmarshal returned an error": {
			ByteArray: []byte(`{\"_"/`),
			Error:     assert.AnError,
		},
		"Successfully ROC validate returned an error": {
			ByteArray: []byte(`{"length":0}`),
			Error:     assert.AnError,
		},
		"Successful unmarshal of a ROC": {
			ByteArray: []byte(`{"length":1}`),
			Result:    ROC{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			r := ROC{}
			err := r.UnmarshalJSON(c.ByteArray)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, r)
			}
		})
	}
}

func TestROC_MarshalJSON(t *testing.T) {
	rc := ROC{length: 1}
	r := []byte(`{"length":1}`)

	d, _ := rc.MarshalJSON()

	assert.Equal(t, r, d)
}

func TestROC_namedMarshalJSON(t *testing.T) {
	rc := ROC{length: 1}
	r := []byte(`{"name":"roc","length":1}`)

	d, _ := rc.namedMarshalJSON()

	assert.Equal(t, r, d)
}

func Test_NewRSI(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result RSI
		Error  error
	}{
		"Successfully RSI returned an error when no parameters were provided": {
			Error: assert.AnError,
		},
		"Successful creation of RSI": {
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

func TestRSI_validate(t *testing.T) {
	cc := map[string]struct {
		Length int
		Error  error
	}{
		"Successfully RSI returned an ErrInvalidLength with less than 1 length": {
			Length: 0,
			Error:  ErrInvalidLength,
		},
		"Successful RSI validation": {
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

func TestRSI_Calc(t *testing.T) {
	cc := map[string]struct {
		Length int
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Successfully RSI returned an ErrInvalidDataSize with insufficient amount of data points": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful RSI calculation": {
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

func TestRSI_Count(t *testing.T) {
	r := RSI{length: 15}
	assert.Equal(t, 15, r.Count())
}

func TestRSI_UnmarshalJSON(t *testing.T) {
	cc := map[string]struct {
		ByteArray []byte
		Result    RSI
		Error     error
	}{
		"Successfully RSI unmarshal returned an error": {
			ByteArray: []byte(`{\"_"/`),
			Error:     assert.AnError,
		},
		"Successfully RSI validate returned an error": {
			ByteArray: []byte(`{"length":0}`),
			Error:     assert.AnError,
		},
		"Successful unmarshal of a RSI": {
			ByteArray: []byte(`{"length":1}`),
			Result:    RSI{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			r := RSI{}
			err := r.UnmarshalJSON(c.ByteArray)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, r)
			}
		})
	}
}

func TestRSI_MarshalJSON(t *testing.T) {

	rs := RSI{length: 1}
	r := []byte(`{"length":1}`)

	d, _ := rs.MarshalJSON()

	assert.Equal(t, r, d)
}

func TestRSI_namedMarshalJSON(t *testing.T) {

	rs := RSI{length: 1}
	r := []byte(`{"name":"rsi","length":1}`)

	d, _ := rs.namedMarshalJSON()

	assert.Equal(t, r, d)
}

func Test_NewSMA(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result SMA
		Error  error
	}{
		"Successfully SMA returned an error when no parameters were provided": {
			Error: assert.AnError,
		},
		"Successful creation of SMA": {
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

func TestSMA_validate(t *testing.T) {
	cc := map[string]struct {
		Length int
		Error  error
	}{
		"Successfully SMA returned an ErrInvalidLength with less than 1 length": {
			Length: 0,
			Error:  ErrInvalidLength,
		},
		"Successful SMA validation": {
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

func TestSMA_Calc(t *testing.T) {
	cc := map[string]struct {
		Length int
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Successfully SMA returned an ErrInvalidDataSize with insufficient amount of data points": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful SMA calculation": {
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

func TestSMA_Count(t *testing.T) {
	s := SMA{length: 15}
	assert.Equal(t, 15, s.Count())
}

func TestSMA_UnmarshalJSON(t *testing.T) {
	cc := map[string]struct {
		ByteArray []byte
		Result    SMA
		Error     error
	}{
		"Successfully SMA unmarshal returned an error": {
			ByteArray: []byte(`{\"_"/`),
			Error:     assert.AnError,
		},
		"Successfully SMA validate returned an error": {
			ByteArray: []byte(`{"length":0}`),
			Error:     assert.AnError,
		},
		"Successful unmarshal of a SMA": {
			ByteArray: []byte(`{"length":1}`),
			Result:    SMA{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			s := SMA{}
			err := s.UnmarshalJSON(c.ByteArray)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, s)
			}
		})
	}
}

func TestSMA_MarshalJSON(t *testing.T) {
	s := SMA{length: 1}
	r := []byte(`{"length":1}`)

	d, _ := s.MarshalJSON()

	assert.Equal(t, r, d)
}

func TestSMA_namedMarshalJSON(t *testing.T) {
	s := SMA{length: 1}
	r := []byte(`{"name":"sma","length":1}`)

	d, _ := s.namedMarshalJSON()

	assert.Equal(t, r, d)
}

func Test_NewStoch(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result Stoch
		Error  error
	}{
		"Successfully Stoch returned an error when no parameters were provided": {
			Error: assert.AnError,
		},
		"Successful creation of Stoch": {
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

func TestStoch_validate(t *testing.T) {
	cc := map[string]struct {
		Length int
		Error  error
	}{
		"Successfully Stoch returned an ErrInvalidLength with less than 1 length": {
			Length: 0,
			Error:  ErrInvalidLength,
		},
		"Successful Stoch validation": {
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

func TestStoch_Calc(t *testing.T) {
	cc := map[string]struct {
		Length int
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Successfully Stoch returned an ErrInvalidDataSize with insufficient amount of data points": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful stoch calculation when lower lows are made": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(150),
				decimal.NewFromInt(125),
				decimal.NewFromInt(145),
			},
			Result: decimal.NewFromInt(80),
		},
		"Successful stoch calculation when higher highs are made": {
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

func TestStoch_Count(t *testing.T) {
	s := Stoch{length: 15}
	assert.Equal(t, 15, s.Count())
}

func TestStoch_UnmarshalMarshal(t *testing.T) {
	cc := map[string]struct {
		ByteArray []byte
		Result    Stoch
		Error     error
	}{
		"Successfully Stoch unmarshal returned an error": {
			ByteArray: []byte(`{"length": "down"}`),
			Error:     assert.AnError,
		},
		"Successfully Stoch validate returned an error": {
			ByteArray: []byte(`{"length":0}`),
			Error:     assert.AnError,
		},
		"Successful unmarshal of a Stoch": {
			ByteArray: []byte(`{"length":1}`),
			Result:    Stoch{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			s := Stoch{}
			err := s.UnmarshalJSON(c.ByteArray)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, s)
			}
		})
	}
}

func TestStoch_MarshalJSON(t *testing.T) {
	s := Stoch{length: 1}
	r := []byte(`{"length":1}`)

	d, _ := s.MarshalJSON()

	assert.Equal(t, r, d)
}

func TestStoch_namedMarshalJSON(t *testing.T) {
	s := Stoch{length: 1}
	r := []byte(`{"name":"stoch","length":1}`)

	d, _ := s.namedMarshalJSON()

	assert.Equal(t, r, d)
}

func Test_NewWMA(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result WMA
		Error  error
	}{
		"Successfully WMA returned an error when no parameters were provided": {
			Error: assert.AnError,
		},
		"Successful creation of WMA": {
			Length: 1,
			Result: WMA{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			w, err := NewWMA(c.Length)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, w)
			}
		})
	}
}

func TestWMA_validate(t *testing.T) {
	cc := map[string]struct {
		Length int
		Error  error
	}{
		"Successfully WMA returned an ErrInvalidLength with less than 1 length": {
			Length: 0,
			Error:  ErrInvalidLength,
		},
		"Successful WMA validation": {
			Length: 1,
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			w := WMA{length: c.Length}
			err := w.validate()
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestWMA_Calc(t *testing.T) {
	cc := map[string]struct {
		Length int
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Successfully WMA returned an ErrInvalidDataSize with insufficient amount of data points": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful WMA calculation": {
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

			w := WMA{length: c.Length}
			res, err := w.Calc(c.Data)
			if c.Error != nil {
				assert.Equal(t, c.Error, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result.String(), res.String())
			}
		})
	}
}

func TestWMA_Count(t *testing.T) {
	w := WMA{length: 15}
	assert.Equal(t, 15, w.Count())
}

func TestWMA_UnmarshalJSON(t *testing.T) {
	cc := map[string]struct {
		ByteArray []byte
		Result    WMA
		Error     error
	}{
		"Successfully DEMA unmarshal returned an error": {
			ByteArray: []byte(`{\"_"/`),
			Error:     assert.AnError,
		},
		"Successfully WMA validate returned an error": {
			ByteArray: []byte(`{"length":0}`),
			Error:     assert.AnError,
		},
		"Successful unmarshal of a WMA": {
			ByteArray: []byte(`{"length":1}`),
			Result:    WMA{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			w := WMA{}
			err := w.UnmarshalJSON(c.ByteArray)
			if c.Error != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, w)
			}
		})
	}
}

func TestWMA_MarshalJSON(t *testing.T) {
	w := WMA{length: 1}
	r := []byte(`{"length":1}`)

	d, _ := w.MarshalJSON()

	assert.Equal(t, r, d)
}

func TestWMA_namedMarshalJSON(t *testing.T) {
	w := WMA{length: 1}
	r := []byte(`{"name":"wma","length":1}`)

	d, _ := w.namedMarshalJSON()

	assert.Equal(t, r, d)
}
