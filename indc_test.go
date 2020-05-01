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
		"Invalid parameters": {
			Error: assert.AnError,
		},
		"Successful creation": {
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
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, a)
		})
	}
}

func Test_Aroon_validate(t *testing.T) {
	cc := map[string]struct {
		Trend  String
		Length int
		Error  error
	}{
		"Invalid trend": {
			Trend:  "downn",
			Length: 5,
			Error:  assert.AnError,
		},
		"Invalid length": {
			Trend:  "down",
			Length: 0,
			Error:  ErrInvalidLength,
		},
		"Successful validation with 'up' trend": {
			Trend:  "up",
			Length: 5,
		},
		"Successful validation with 'down' trend": {
			Trend:  "down",
			Length: 5,
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			a := Aroon{trend: c.Trend, length: c.Length}
			equalError(t, c.Error, a.validate())
		})
	}
}

func Test_Aroon_Calc(t *testing.T) {
	cc := map[string]struct {
		Trend  String
		Length int
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid data size": {
			Trend:  "down",
			Length: 5,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation with 'up' trend": {
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
		"Successful calculation with 'down' trend": {
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
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_Aroon_Count(t *testing.T) {
	a := Aroon{trend: "down", length: 5}
	assert.Equal(t, 5, a.Count())
}

func Test_Aroon_UnmarshalJSON(t *testing.T) {
	cc := map[string]struct {
		JSON   string
		Result Aroon
		Error  error
	}{
		"Invalid JSON": {
			JSON:  `{\"_"/`,
			Error: assert.AnError,
		},
		"Invalid trend": {
			JSON:  `{"trend":"upp","length":1}`,
			Error: assert.AnError,
		},
		"Successful unmarshal": {
			JSON:   `{"trend":"up","length":1}`,
			Result: Aroon{trend: "up", length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			a := Aroon{}
			err := a.UnmarshalJSON([]byte(c.JSON))
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, a)
		})
	}
}

func Test_Aroon_MarshalJSON(t *testing.T) {
	a := Aroon{trend: "down", length: 1}
	d, err := a.MarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"trend":"down","length":1}`, string(d))
}

func Test_Aroon_namedMarshalJSON(t *testing.T) {
	a := Aroon{trend: "down", length: 1}
	d, err := a.namedMarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"name":"aroon","trend":"down","length":1}`, string(d))
}

func Test_NewCCI(t *testing.T) {
	cc := map[string]struct {
		Source Indicator
		Result CCI
		Error  error
	}{
		"Invalid parameters": {
			Error: assert.AnError,
		},
		"Successful creation": {
			Source: Aroon{trend: "down", length: 1},
			Result: CCI{source: Aroon{trend: "down", length: 1}},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			cci, err := NewCCI(c.Source)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, cci)
		})
	}
}

func Test_CCI_validate(t *testing.T) {
	cc := map[string]struct {
		Source Indicator
		Error  error
	}{
		"Invalid source": {
			Error: ErrInvalidSource,
		},
		"Successful validation": {
			Source: Aroon{trend: "down", length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			cci := CCI{source: c.Source}
			equalError(t, c.Error, cci.validate())
		})
	}
}

func Test_CCI_Calc(t *testing.T) {
	cc := map[string]struct {
		Source Indicator
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid data size": {
			Source: EMA{length: 10},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Error returned by the source": {
			Source: IndicatorMock{},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: assert.AnError,
		},
		"Successful calculation": {
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
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.Round(14).String())
		})
	}
}

func Test_CCI_Count(t *testing.T) {
	c := CCI{source: Aroon{trend: "down", length: 10}}
	assert.Equal(t, c.source.Count(), c.Count())
}

func Test_CCI_UnmarshalJSON(t *testing.T) {
	cc := map[string]struct {
		JSON   string
		Result CCI
		Error  error
	}{
		"Invalid JSON": {
			JSON:  `{\-_-/}`,
			Error: assert.AnError,
		},
		"Error returned during source parsing": {
			JSON:  `{"trend":"up"}`,
			Error: assert.AnError,
		},
		"Successful CCI": {
			JSON:   `{"source":{"name":"aroon","trend":"up","length":1}}`,
			Result: CCI{Aroon{trend: "up", length: 1}},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			cci := CCI{}
			err := cci.UnmarshalJSON([]byte(c.JSON))
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, cci)
		})
	}
}

func Test_CCI_MarshalJSON(t *testing.T) {
	cc := map[string]struct {
		CCI    CCI
		Result string
		Error  error
	}{
		"Error returned during source marshalling": {
			CCI:   CCI{source: IndicatorMock{}},
			Error: assert.AnError,
		},
		"Successful marshal": {
			CCI:    CCI{source: Aroon{trend: "down", length: 1}},
			Result: `{"source":{"name":"aroon","trend":"down","length":1}}`,
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			d, err := c.CCI.MarshalJSON()
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.JSONEq(t, c.Result, string(d))
		})
	}
}

func Test_CCI_namedMarshalJSON(t *testing.T) {
	cc := map[string]struct {
		CCI    CCI
		Result string
		Error  error
	}{
		"Error returned during source marshalling": {
			CCI:   CCI{source: IndicatorMock{}},
			Error: assert.AnError,
		},
		"Successful marshal": {
			CCI:    CCI{source: Aroon{trend: "down", length: 1}},
			Result: `{"name":"cci","source":{"name":"aroon","trend":"down","length":1}}`,
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			d, err := c.CCI.namedMarshalJSON()
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.JSONEq(t, c.Result, string(d))
		})
	}
}

func Test_NewDEMA(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result DEMA
		Error  error
	}{
		"Invalid parameters": {
			Error: assert.AnError,
		},
		"Successful creation": {
			Length: 1,
			Result: DEMA{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			dm, err := NewDEMA(c.Length)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, dm)
		})
	}
}

func Test_DEMA_validate(t *testing.T) {
	cc := map[string]struct {
		Length int
		Error  error
	}{
		"Invalid length": {
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

			d := DEMA{length: c.Length}
			equalError(t, c.Error, d.validate())
		})
	}
}

func Test_DEMA_Calc(t *testing.T) {
	cc := map[string]struct {
		Length int
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid data size": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
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

			d := DEMA{length: c.Length}
			res, err := d.Calc(c.Data)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.Round(14).String())
		})
	}
}

func Test_DEMA_Count(t *testing.T) {
	d := DEMA{length: 15}
	assert.Equal(t, 29, d.Count())
}

func Test_DEMA_UnmarshalJSON(t *testing.T) {
	cc := map[string]struct {
		JSON   string
		Result DEMA
		Error  error
	}{
		"Invalid JSON": {
			JSON:  `{\"_"/`,
			Error: assert.AnError,
		},
		"Invalid length": {
			JSON:  `{"length":0}`,
			Error: assert.AnError,
		},
		"Successful unmarshal": {
			JSON:   `{"length":1}`,
			Result: DEMA{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			dm := DEMA{}
			err := dm.UnmarshalJSON([]byte(c.JSON))
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, dm)
		})
	}
}

func Test_DEMA_MarshalJSON(t *testing.T) {
	dm := DEMA{length: 1}
	d, err := dm.MarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"length":1}`, string(d))
}

func Test_DEMA_namedMarshalJSON(t *testing.T) {
	dm := DEMA{length: 1}
	d, err := dm.namedMarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"name":"dema","length":1}`, string(d))
}

func Test_NewEMA(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result EMA
		Error  error
	}{
		"Invalid parameters": {
			Error: assert.AnError,
		},
		"Successful creation": {
			Length: 1,
			Result: EMA{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			e, err := NewEMA(c.Length)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, e)
		})
	}
}

func Test_EMA_validate(t *testing.T) {
	cc := map[string]struct {
		Length int
		Error  error
	}{
		"Invalid length": {
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
			equalError(t, c.Error, e.validate())
		})
	}
}

func Test_EMA_Calc(t *testing.T) {
	cc := map[string]struct {
		Length int
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid data size": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
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
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.Round(14).String())
		})
	}
}

func Test_EMA_Count(t *testing.T) {
	e := EMA{length: 15}
	assert.Equal(t, 29, e.Count())
}

func Test_EMA_multiplier(t *testing.T) {
	e := EMA{length: 3}
	assert.Equal(t, decimal.NewFromFloat(0.5), e.multiplier())
}

func Test_EMA_UnmarshalJSON(t *testing.T) {
	cc := map[string]struct {
		JSON   string
		Result EMA
		Error  error
	}{
		"Invalid JSON": {
			JSON:  `{\"_"/`,
			Error: assert.AnError,
		},
		"Invalid length": {
			JSON:  `{"length":0}`,
			Error: assert.AnError,
		},
		"Successful unmarshal": {
			JSON:   `{"length":1}`,
			Result: EMA{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			e := EMA{}
			err := e.UnmarshalJSON([]byte(c.JSON))
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, e)
		})
	}
}

func Test_EMA_MarshalJSON(t *testing.T) {
	e := EMA{length: 1}
	d, err := e.MarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"length":1}`, string(d))
}

func Test_EMA_namedMarshalJSON(t *testing.T) {
	e := EMA{length: 1}
	d, err := e.namedMarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"name":"ema","length":1}`, string(d))
}

func Test_NewHMA(t *testing.T) {
	cc := map[string]struct {
		WMA    WMA
		Result HMA
		Error  error
	}{
		"Invalid parameters": {
			Error: assert.AnError,
		},
		"Successful creation": {
			WMA:    WMA{length: 1},
			Result: HMA{wma: WMA{length: 1}},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			h, err := NewHMA(c.WMA)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, h)
		})
	}
}

func Test_HMA_validate(t *testing.T) {
	cc := map[string]struct {
		WMA   WMA
		Error error
	}{
		"Invalid WMA": {
			Error: assert.AnError,
		},
		"Invalid WMA length": {
			WMA:   WMA{length: -1},
			Error: assert.AnError,
		},
		"Successful validation": {
			WMA: WMA{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			h := HMA{wma: c.WMA}
			equalError(t, c.Error, h.validate())
		})
	}
}

func Test_HMA_Calc(t *testing.T) {
	cc := map[string]struct {
		WMA    WMA
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid data size": {
			WMA: WMA{length: 5},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
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
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_HMA_Count(t *testing.T) {
	h := HMA{wma: WMA{length: 15}}
	assert.Equal(t, 29, h.Count())
}

func Test_HMA_UnmarshalJSON(t *testing.T) {
	cc := map[string]struct {
		JSON   string
		Result HMA
		Error  error
	}{
		"Invalid JSON": {
			JSON:  `{\"_"/}`,
			Error: assert.AnError,
		},
		"Invalid length": {
			JSON:  `{"length":0}`,
			Error: assert.AnError,
		},
		"Successful unmarshal": {
			JSON:   `{"wma":{"length":1}}`,
			Result: HMA{wma: WMA{length: 1}},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			h := HMA{}
			err := h.UnmarshalJSON([]byte(c.JSON))
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, h)
		})
	}
}

func Test_HMA_MarshalJSON(t *testing.T) {
	h := HMA{wma: WMA{length: 1}}
	d, err := h.MarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"wma":{"length":1}}`, string(d))
}

func Test_HMA_namedMarshalJSON(t *testing.T) {
	h := HMA{wma: WMA{length: 1}}
	d, err := h.namedMarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"name":"hma","wma":{"length":1}}`, string(d))
}

func Test_NewMACD(t *testing.T) {
	cc := map[string]struct {
		Source1 Indicator
		Source2 Indicator
		Result  MACD
		Error   error
	}{
		"Invalid parameters": {
			Error: assert.AnError,
		},
		"Successful creation": {
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
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, m)
		})
	}
}

func Test_MACD_validate(t *testing.T) {
	cc := map[string]struct {
		Source1 Indicator
		Source2 Indicator
		Error   error
	}{
		"Invalid source1": {
			Source1: EMA{length: 1},
			Error:   ErrInvalidSource,
		},
		"Invalid source2": {
			Source2: EMA{length: 1},
			Error:   ErrInvalidSource,
		},
		"Successful MACD": {
			Source1: EMA{length: 1},
			Source2: EMA{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			m := MACD{source1: c.Source1, source2: c.Source2}
			equalError(t, c.Error, m.validate())
		})
	}
}

func Test_MACD_Calc(t *testing.T) {
	cc := map[string]struct {
		Source1 Indicator
		Source2 Indicator
		Data    []decimal.Decimal
		Result  decimal.Decimal
		Error   error
	}{
		"Invalid data size for source1": {
			Source1: EMA{length: 4},
			Source2: EMA{length: 1},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Invalid data size for source2": {
			Source1: EMA{length: 1},
			Source2: EMA{length: 4},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Error returned by source1": {
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
		"Error returned by source2": {
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
		"Successful calculation": {
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
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_MACD_Count(t *testing.T) {
	m := MACD{source1: EMA{length: 10}, source2: EMA{length: 1}}
	assert.Equal(t, m.Count(), m.source1.Count())

	m = MACD{source1: EMA{length: 2}, source2: EMA{length: 9}}
	assert.Equal(t, m.Count(), m.source2.Count())
}

func Test_MACD_UnmarshalJSON(t *testing.T) {
	cc := map[string]struct {
		JSON   string
		Result MACD
		Error  error
	}{
		"Invalid JSON": {
			JSON:  `{\-_-/}`,
			Error: assert.AnError,
		},
		"Invalid source1": {
			JSON: `{"source1":{"name":"aroon","trend":"dsown","length":2},
			"source2":{"name":"cci","source":{"name":"ema", "length":2}}}`,
			Error: assert.AnError,
		},
		"Invalid source2": {
			JSON: `{"source1":{"name":"aroon","trend":"down","length":2},
			"source2":{"name":"ccis","source":{"name":"ema", "length":2}}}`,
			Error: assert.AnError,
		},
		"Successful unmarshal": {
			JSON: `{"source1":{"name":"aroon","trend":"down","length":2},
			"source2":{"name":"cci","source":{"name":"ema", "length":2}}}`,
			Result: MACD{source1: Aroon{trend: "down", length: 2},
				source2: CCI{source: EMA{length: 2}}},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			m := MACD{}
			err := m.UnmarshalJSON([]byte(c.JSON))
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, m)
		})
	}
}

func Test_MACD_MarshalJSON(t *testing.T) {
	cc := map[string]struct {
		MACD   MACD
		Result string
		Error  error
	}{
		"Error returned during source1 marshalling": {
			MACD: MACD{source1: IndicatorMock{},
				source2: CCI{source: EMA{length: 2}}},
			Error: assert.AnError,
		},
		"Error returned during source2 marshalling": {
			MACD: MACD{source1: Aroon{trend: "down", length: 2},
				source2: IndicatorMock{}},
			Error: assert.AnError,
		},
		"Successful marshal": {
			MACD: MACD{source1: Aroon{trend: "down", length: 2},
				source2: CCI{source: EMA{length: 2}}},
			Result: `{"source1":{"name":"aroon","trend":"down","length":2},"source2":{"name":"cci","source":{"name":"ema","length":2}}}`,
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			d, err := c.MACD.MarshalJSON()
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.JSONEq(t, c.Result, string(d))
		})
	}
}

func Test_MACD_namedMarshalJSON(t *testing.T) {
	cc := map[string]struct {
		MACD   MACD
		Result string
		Error  error
	}{
		"Error returned during source1 marshalling": {
			MACD: MACD{source1: IndicatorMock{},
				source2: CCI{source: EMA{length: 2}}},
			Error: assert.AnError,
		},
		"Error returned during source2 marshalling": {
			MACD: MACD{source1: Aroon{trend: "down", length: 2},
				source2: IndicatorMock{}},
			Error: assert.AnError,
		},
		"Successful marshal": {
			MACD: MACD{source1: Aroon{trend: "down", length: 2},
				source2: CCI{source: EMA{length: 2}}},
			Result: `{"name":"macd","source1":{"name":"aroon","trend":"down","length":2},"source2":{"name":"cci","source":{"name":"ema","length":2}}}`,
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			d, err := c.MACD.namedMarshalJSON()
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.JSONEq(t, c.Result, string(d))
		})
	}
}

func Test_NewROC(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result ROC
		Error  error
	}{
		"Invalid parameters": {
			Error: assert.AnError,
		},
		"Successful creation": {
			Length: 1,
			Result: ROC{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			r, err := NewROC(c.Length)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, r)
		})
	}
}

func Test_ROC_validate(t *testing.T) {
	cc := map[string]struct {
		Length int
		Error  error
	}{
		"Invalid length": {
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
			equalError(t, c.Error, r.validate())
		})
	}
}

func Test_ROC_Calc(t *testing.T) {
	cc := map[string]struct {
		Length int
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid data size": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
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
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_ROC_Count(t *testing.T) {
	r := ROC{length: 15}
	assert.Equal(t, 15, r.Count())
}

func Test_ROC_UnmarshalJSON(t *testing.T) {
	cc := map[string]struct {
		JSON   string
		Result ROC
		Error  error
	}{
		"Invalid JSON": {
			JSON:  `{\"_"/`,
			Error: assert.AnError,
		},
		"Invalid length": {
			JSON:  `{"length":0}`,
			Error: assert.AnError,
		},
		"Successful unmarshal": {
			JSON:   `{"length":1}`,
			Result: ROC{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			r := ROC{}
			err := r.UnmarshalJSON([]byte(c.JSON))
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, r)
		})
	}
}

func Test_ROC_MarshalJSON(t *testing.T) {
	rc := ROC{length: 1}
	d, err := rc.MarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"length":1}`, string(d))
}

func Test_ROC_namedMarshalJSON(t *testing.T) {
	rc := ROC{length: 1}
	d, err := rc.namedMarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"name":"roc","length":1}`, string(d))
}

func Test_NewRSI(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result RSI
		Error  error
	}{
		"Invalid parameters": {
			Error: assert.AnError,
		},
		"Successful creation": {
			Length: 1,
			Result: RSI{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			r, err := NewRSI(c.Length)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, r)
		})
	}
}

func Test_RSI_validate(t *testing.T) {
	cc := map[string]struct {
		Length int
		Error  error
	}{
		"Invalid length": {
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
			equalError(t, c.Error, r.validate())
		})
	}
}

func Test_RSI_Calc(t *testing.T) {
	cc := map[string]struct {
		Length int
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid data size": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
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
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.Round(14).String())
		})
	}
}

func Test_RSI_Count(t *testing.T) {
	r := RSI{length: 15}
	assert.Equal(t, 15, r.Count())
}

func Test_RSI_UnmarshalJSON(t *testing.T) {
	cc := map[string]struct {
		JSON   string
		Result RSI
		Error  error
	}{
		"Invalid JSON": {
			JSON:  `{\"_"/`,
			Error: assert.AnError,
		},
		"Invalid length": {
			JSON:  `{"length":0}`,
			Error: assert.AnError,
		},
		"Successful unmarshal": {
			JSON:   `{"length":1}`,
			Result: RSI{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			r := RSI{}
			err := r.UnmarshalJSON([]byte(c.JSON))
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, r)
		})
	}
}

func Test_RSI_MarshalJSON(t *testing.T) {
	rs := RSI{length: 1}
	d, err := rs.MarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"length":1}`, string(d))
}

func Test_RSI_namedMarshalJSON(t *testing.T) {
	rs := RSI{length: 1}
	d, err := rs.namedMarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"name":"rsi","length":1}`, string(d))
}

func Test_NewSMA(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result SMA
		Error  error
	}{
		"Invalid parameters": {
			Error: assert.AnError,
		},
		"Successful creation": {
			Length: 1,
			Result: SMA{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			s, err := NewSMA(c.Length)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, s)
		})
	}
}

func Test_SMA_validate(t *testing.T) {
	cc := map[string]struct {
		Length int
		Error  error
	}{
		"Invalid length": {
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
			equalError(t, c.Error, s.validate())
		})
	}
}

func Test_SMA_Calc(t *testing.T) {
	cc := map[string]struct {
		Length int
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid data size": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
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
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_SMA_Count(t *testing.T) {
	s := SMA{length: 15}
	assert.Equal(t, 15, s.Count())
}

func Test_SMA_UnmarshalJSON(t *testing.T) {
	cc := map[string]struct {
		JSON   string
		Result SMA
		Error  error
	}{
		"Invalid JSON": {
			JSON:  `{\"_"/`,
			Error: assert.AnError,
		},
		"Invalid length": {
			JSON:  `{"length":0}`,
			Error: assert.AnError,
		},
		"Successful unmarshal": {
			JSON:   `{"length":1}`,
			Result: SMA{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			s := SMA{}
			err := s.UnmarshalJSON([]byte(c.JSON))
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, s)
		})
	}
}

func Test_SMA_MarshalJSON(t *testing.T) {
	s := SMA{length: 1}
	d, err := s.MarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"length":1}`, string(d))
}

func Test_SMA_namedMarshalJSON(t *testing.T) {
	s := SMA{length: 1}
	d, err := s.namedMarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"name":"sma","length":1}`, string(d))
}

func Test_NewSRSI(t *testing.T) {
	cc := map[string]struct {
		RSI    RSI
		Result SRSI
		Error  error
	}{
		"Invalid parameters": {
			Error: assert.AnError,
		},
		"Successful creation": {
			RSI:    RSI{length: 1},
			Result: SRSI{rsi: RSI{length: 1}},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			s, err := NewSRSI(c.RSI)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, s)
		})
	}
}

func Test_SRSI_validate(t *testing.T) {
	cc := map[string]struct {
		RSI   RSI
		Error error
	}{
		"Invalid RSI": {
			Error: assert.AnError,
		},
		"Invalid RSI length": {
			RSI:   RSI{length: -1},
			Error: assert.AnError,
		},
		"Successful validation": {
			RSI: RSI{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			s := SRSI{rsi: c.RSI}
			equalError(t, c.Error, s.validate())
		})
	}
}

func Test_SRSI_Calc(t *testing.T) {
	cc := map[string]struct {
		RSI    RSI
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid data size": {
			RSI: RSI{length: 5},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: assert.AnError,
		},
		"Successful calculation": {
			RSI: RSI{length: 5},
			Data: []decimal.Decimal{
				decimal.NewFromFloat(54.09),
				decimal.NewFromFloat(59.9),
				decimal.NewFromFloat(58.2),
				decimal.NewFromFloat(59.76),
				decimal.NewFromFloat(52.35),
				decimal.NewFromFloat(52.82),
				decimal.NewFromFloat(56.94),
				decimal.NewFromFloat(57.47),
				decimal.NewFromFloat(55.26),
				decimal.NewFromFloat(57.51),
				decimal.NewFromFloat(54.8),
				decimal.NewFromFloat(51.47),
				decimal.NewFromFloat(56.16),
				decimal.NewFromFloat(58.34),
			},
			Result: decimal.NewFromFloat(0.58524401),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			s := SRSI{rsi: c.RSI}
			res, err := s.Calc(c.Data)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.Round(8).String())
		})
	}
}

func Test_SRSI_Count(t *testing.T) {
	s := SRSI{rsi: RSI{length: 15}}
	assert.Equal(t, 29, s.Count())
}

func Test_SRSI_UnmarshalJSON(t *testing.T) {
	cc := map[string]struct {
		JSON   string
		Result SRSI
		Error  error
	}{
		"Invalid JSON": {
			JSON:  `{\"_"/`,
			Error: assert.AnError,
		},
		"Invalid length": {
			JSON:  `{"length":0}`,
			Error: assert.AnError,
		},
		"Successful unmarshal": {
			JSON:   `{"rsi":{"length":1}}`,
			Result: SRSI{rsi: RSI{length: 1}},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			s := SRSI{}
			err := s.UnmarshalJSON([]byte(c.JSON))
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, s)
		})
	}
}

func Test_SRSI_MarshalJSON(t *testing.T) {
	s := SRSI{rsi: RSI{length: 1}}
	d, err := s.MarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"rsi":{"length":1}}`, string(d))
}

func Test_SRSI_namedMarshalJSON(t *testing.T) {
	s := SRSI{rsi: RSI{length: 1}}
	d, err := s.namedMarshalJSON()

	assert.NoError(t, err)
	assert.Equal(t, `{"name":"srsi","rsi":{"length":1}}`, string(d))
}

func Test_NewStoch(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result Stoch
		Error  error
	}{
		"Invalid parameters": {
			Error: assert.AnError,
		},
		"Successful creation": {
			Length: 1,
			Result: Stoch{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			s, err := NewStoch(c.Length)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, s)
		})
	}
}

func Test_Stoch_validate(t *testing.T) {
	cc := map[string]struct {
		Length int
		Error  error
	}{
		"Invalid length": {
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
			equalError(t, c.Error, s.validate())
		})
	}
}

func Test_Stoch_Calc(t *testing.T) {
	cc := map[string]struct {
		Length int
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid data size": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation when new lows are reached": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(150),
				decimal.NewFromInt(125),
				decimal.NewFromInt(145),
			},
			Result: decimal.NewFromInt(80),
		},
		"Successful calculation when new highs are reached": {
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
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_Stoch_Count(t *testing.T) {
	s := Stoch{length: 15}
	assert.Equal(t, 15, s.Count())
}

func Test_Stoch_UnmarshalJSON(t *testing.T) {
	cc := map[string]struct {
		JSON   string
		Result Stoch
		Error  error
	}{
		"Invalid JSON": {
			JSON:  `{"length": "down"}`,
			Error: assert.AnError,
		},
		"Invalid length": {
			JSON:  `{"length":0}`,
			Error: assert.AnError,
		},
		"Successful unmarshal": {
			JSON:   `{"length":1}`,
			Result: Stoch{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			s := Stoch{}
			err := s.UnmarshalJSON([]byte(c.JSON))
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, s)
		})
	}
}

func Test_Stoch_MarshalJSON(t *testing.T) {
	s := Stoch{length: 1}
	d, err := s.MarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"length":1}`, string(d))
}

func Test_Stoch_namedMarshalJSON(t *testing.T) {
	s := Stoch{length: 1}
	d, err := s.namedMarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"name":"stoch","length":1}`, string(d))
}

func Test_NewWMA(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result WMA
		Error  error
	}{
		"Invalid parameters": {
			Error: assert.AnError,
		},
		"Successful creation": {
			Length: 1,
			Result: WMA{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			w, err := NewWMA(c.Length)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, w)
		})
	}
}

func Test_WMA_validate(t *testing.T) {
	cc := map[string]struct {
		Length int
		Error  error
	}{
		"Invalid length": {
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

			w := WMA{length: c.Length}
			equalError(t, c.Error, w.validate())
		})
	}
}

func Test_WMA_Calc(t *testing.T) {
	cc := map[string]struct {
		Length int
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid data size": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
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

			w := WMA{length: c.Length}
			res, err := w.Calc(c.Data)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_WMA_Count(t *testing.T) {
	w := WMA{length: 15}
	assert.Equal(t, 15, w.Count())
}

func Test_WMA_UnmarshalJSON(t *testing.T) {
	cc := map[string]struct {
		JSON   string
		Result WMA
		Error  error
	}{
		"Invalid JSON": {
			JSON:  `{\"_"/`,
			Error: assert.AnError,
		},
		"Invalid length": {
			JSON:  `{"length":0}`,
			Error: assert.AnError,
		},
		"Successful unmarshal": {
			JSON:   `{"length":1}`,
			Result: WMA{length: 1},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			w := WMA{}
			err := w.UnmarshalJSON([]byte(c.JSON))
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, w)
		})
	}
}

func Test_WMA_MarshalJSON(t *testing.T) {
	w := WMA{length: 1}
	d, err := w.MarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"length":1}`, string(d))
}

func Test_WMA_namedMarshalJSON(t *testing.T) {
	w := WMA{length: 1}
	d, err := w.namedMarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"name":"wma","length":1}`, string(d))
}
