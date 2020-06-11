package indc

import (
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_NewAroon(t *testing.T) {
	cc := map[string]struct {
		Trend  String
		Length int
		Result Aroon
		Error  error
	}{
		"Invalid parameters": {
			Trend:  "",
			Length: 1,
			Error:  assert.AnError,
		},
		"Successful creation": {
			Trend:  "down",
			Length: 5,
			Result: Aroon{trend: "down", length: 5, valid: true},
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

func Test_Aroon_Length(t *testing.T) {
	assert.Equal(t, 1, Aroon{length: 1}.Length())
}

func Test_Aroon_Trend(t *testing.T) {
	assert.Equal(t, CleanString("up"), Aroon{trend: CleanString("up")}.Trend())
}

func Test_Aroon_validate(t *testing.T) {
	cc := map[string]struct {
		Aroon Aroon
		Error error
		Valid bool
	}{
		"Invalid trend": {
			Aroon: Aroon{trend: "downn", length: 5},
			Error: errors.New("invalid trend"),
			Valid: false,
		},
		"Invalid length": {
			Aroon: Aroon{trend: "down", length: 0},
			Error: ErrInvalidLength,
			Valid: false,
		},
		"Successful validation with trend being up": {
			Aroon: Aroon{trend: "up", length: 1},
			Valid: true,
		},
		"Successful validation with trend being down": {
			Aroon: Aroon{trend: "down", length: 1},
			Valid: true,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			equalError(t, c.Error, c.Aroon.validate())
			assert.Equal(t, c.Valid, c.Aroon.valid)
		})
	}
}

func Test_Aroon_Calc(t *testing.T) {
	cc := map[string]struct {
		Aroon  Aroon
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			Aroon: Aroon{valid: false},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			Aroon: Aroon{trend: "down", length: 5, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation with trend being up": {
			Aroon: Aroon{trend: "up", length: 5, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(25),
				decimal.NewFromInt(31),
				decimal.NewFromInt(38),
				decimal.NewFromInt(35),
				decimal.NewFromInt(29),
				decimal.NewFromInt(29),
			},
			Result: decimal.NewFromInt(40),
		},
		"Successful calculation with trend being down": {
			Aroon: Aroon{trend: "down", length: 5, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(25),
				decimal.NewFromInt(31),
				decimal.NewFromInt(38),
				decimal.NewFromInt(35),
				decimal.NewFromInt(29),
				decimal.NewFromInt(29),
			},
			Result: decimal.NewFromInt(100),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.Aroon.Calc(c.Data)
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
			Result: Aroon{trend: "up", length: 1, valid: true},
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
		Factor decimal.Decimal
		Result CCI
		Error  error
	}{
		"Invalid parameters": {
			Source: &IndicatorMock{},
			Factor: decimal.NewFromInt(-1),
			Error:  assert.AnError,
		},
		"Successful creation (default factor)": {
			Source: &IndicatorMock{},
			Factor: decimal.Zero,
			Result: CCI{source: &IndicatorMock{}, factor: decimal.RequireFromString("0.015"), valid: true},
		},
		"Successful creation": {
			Source: &IndicatorMock{},
			Factor: Hundred,
			Result: CCI{source: &IndicatorMock{}, factor: Hundred, valid: true},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			cci, err := NewCCI(c.Source, c.Factor)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, cci)
		})
	}
}

func Test_CCI_Sub(t *testing.T) {
	assert.Equal(t, &IndicatorMock{}, CCI{source: &IndicatorMock{}}.Sub())
}

func Test_CCI_Factor(t *testing.T) {
	assert.Equal(t, Hundred, CCI{factor: Hundred}.Factor())
}

func Test_CCI_validate(t *testing.T) {
	cc := map[string]struct {
		CCI   CCI
		Error error
		Valid bool
	}{
		"Invalid source": {
			CCI:   CCI{source: nil},
			Error: ErrInvalidSource,
			Valid: false,
		},
		"Invalid factor": {
			CCI:   CCI{source: &IndicatorMock{}, factor: decimal.NewFromInt(-1)},
			Error: errors.New("invalid factor"),
			Valid: false,
		},
		"Successful validation": {
			CCI:   CCI{source: &IndicatorMock{}, factor: decimal.RequireFromString("1")},
			Valid: true,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			equalError(t, c.Error, c.CCI.validate())
			assert.Equal(t, c.Valid, c.CCI.valid)
		})
	}
}

func Test_CCI_Calc(t *testing.T) {
	stubIndicator := func(v decimal.Decimal, e error, a int) *IndicatorMock {
		return &IndicatorMock{
			CalcFunc: func(dd []decimal.Decimal) (decimal.Decimal, error) {
				return v, e
			},
			CountFunc: func() int {
				return a
			},
		}
	}

	cc := map[string]struct {
		CCI    CCI
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			CCI:   CCI{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			CCI: CCI{source: stubIndicator(decimal.Zero, nil, 10), factor: decimal.RequireFromString("0.015"), valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Invalid source calc": {
			CCI: CCI{source: stubIndicator(decimal.Zero, assert.AnError, 1), factor: decimal.RequireFromString("0.015"), valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: assert.AnError,
		},
		"Successful handled division by 0": {
			CCI: CCI{source: stubIndicator(decimal.NewFromInt(3), nil, 1), factor: decimal.RequireFromString("0.015"), valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(3),
				decimal.NewFromInt(6),
				decimal.NewFromInt(9),
			},
			Result: decimal.Zero,
		},
		"Successful calculation": {
			CCI: CCI{source: stubIndicator(decimal.NewFromInt(3), nil, 3), factor: decimal.RequireFromString("0.015"), valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(3),
				decimal.NewFromInt(6),
				decimal.NewFromInt(9),
			},
			Result: decimal.NewFromInt(200),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.CCI.Calc(c.Data)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_CCI_Count(t *testing.T) {
	indicator := &IndicatorMock{
		CountFunc: func() int {
			return 10
		},
	}

	c := CCI{source: indicator}
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
		"Invalid source": {
			JSON:  `{}`,
			Error: assert.AnError,
		},
		"Successful unmarshal": {
			JSON:   `{"source":{"name":"sma","length":1}}`,
			Result: CCI{source: SMA{length: 1, valid: true}, factor: decimal.RequireFromString("0.015"), valid: true},
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
	stubIndicator := func(d []byte, e error) *IndicatorMock {
		return &IndicatorMock{
			namedMarshalJSONFunc: func() ([]byte, error) {
				return d, e
			},
		}
	}

	cc := map[string]struct {
		CCI    CCI
		Result string
		Error  error
	}{
		"Error returned during source marshalling": {
			CCI: CCI{
				source: stubIndicator(nil, assert.AnError),
			},
			Error: assert.AnError,
		},
		"Successful marshal": {
			CCI: CCI{
				source: stubIndicator([]byte(`{"name":"indicatormock"}`), nil),
				factor: Hundred,
			},
			Result: `{"source":{"name":"indicatormock"},"factor":"100"}`,
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
	stubIndicator := func(d []byte, e error) *IndicatorMock {
		return &IndicatorMock{
			namedMarshalJSONFunc: func() ([]byte, error) {
				return d, e
			},
		}
	}

	cc := map[string]struct {
		CCI    CCI
		Result string
		Error  error
	}{
		"Error returned during source marshalling": {
			CCI: CCI{
				source: stubIndicator(nil, assert.AnError),
			},
			Error: assert.AnError,
		},
		"Successful marshal": {
			CCI: CCI{
				source: stubIndicator([]byte(`{"name":"indicatormock"}`), nil),
				factor: Hundred,
			},
			Result: `{"name":"cci","source":{"name":"indicatormock"},"factor":"100"}`,
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
		EMA    EMA
		Result DEMA
		Error  error
	}{
		"Invalid parameters": {
			Error: assert.AnError,
		},
		"Successful creation": {
			EMA:    EMA{sma: SMA{length: 1, valid: true}, valid: true},
			Result: DEMA{ema: EMA{sma: SMA{length: 1, valid: true}, valid: true}, valid: true},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			dm, err := NewDEMA(c.EMA)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, dm)
		})
	}
}

func Test_DEMA_Length(t *testing.T) {
	dm := DEMA{ema: EMA{sma: SMA{length: 1}}}
	assert.Equal(t, 1, dm.Length())
}

func Test_DEMA_validate(t *testing.T) {
	cc := map[string]struct {
		DEMA  DEMA
		Error error
		Valid bool
	}{
		"Invalid EMA": {
			DEMA:  DEMA{ema: EMA{sma: SMA{length: -1}}},
			Error: assert.AnError,
			Valid: false,
		},
		"Successful validation": {
			DEMA:  DEMA{ema: EMA{sma: SMA{length: 1}}},
			Valid: true,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			equalError(t, c.Error, c.DEMA.validate())
			assert.Equal(t, c.Valid, c.DEMA.valid)
		})
	}
}

func Test_DEMA_Calc(t *testing.T) {
	cc := map[string]struct {
		DEMA   DEMA
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			DEMA:  DEMA{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			DEMA: DEMA{ema: EMA{sma: SMA{length: 3, valid: true}, valid: true}, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
			DEMA: DEMA{ema: EMA{sma: SMA{length: 3, valid: true}, valid: true}, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
				decimal.NewFromInt(31),
				decimal.NewFromInt(1),
				decimal.NewFromInt(1),
				decimal.NewFromInt(2),
				decimal.NewFromInt(3),
			},
			Result: decimal.RequireFromString("6.75"),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.DEMA.Calc(c.Data)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_DEMA_Count(t *testing.T) {
	d := DEMA{ema: EMA{sma: SMA{length: 15}}}
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
			JSON:   `{"ema":{"length":1}}`,
			Result: DEMA{ema: EMA{sma: SMA{length: 1, valid: true}, valid: true}, valid: true},
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
	dm := DEMA{ema: EMA{sma: SMA{length: 1}}}
	d, err := dm.MarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"ema":{"length":1}}`, string(d))
}

func Test_DEMA_namedMarshalJSON(t *testing.T) {
	dm := DEMA{ema: EMA{sma: SMA{length: 1}}}
	d, err := dm.namedMarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"name":"dema","ema":{"length":1}}`, string(d))
}

func Test_NewEMA(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result EMA
		Error  error
	}{
		"Invalid parameters": {
			Length: -1,
			Error:  assert.AnError,
		},
		"Successful creation": {
			Length: 1,
			Result: EMA{sma: SMA{length: 1, valid: true}, valid: true},
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

func Test_EMA_Length(t *testing.T) {
	e := EMA{sma: SMA{length: 1}}
	assert.Equal(t, 1, e.Length())
}

func Test_EMA_validate(t *testing.T) {
	cc := map[string]struct {
		EMA   EMA
		Error error
		Valid bool
	}{
		"Invalid SMA": {
			EMA:   EMA{sma: SMA{length: -1}},
			Error: assert.AnError,
			Valid: false,
		},
		"Successful validation": {
			EMA:   EMA{sma: SMA{length: 1}},
			Valid: true,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			equalError(t, c.Error, c.EMA.validate())
			assert.Equal(t, c.Valid, c.EMA.valid)
		})
	}
}

func Test_EMA_Calc(t *testing.T) {
	cc := map[string]struct {
		EMA    EMA
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			EMA:   EMA{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			EMA: EMA{sma: SMA{length: 3, valid: true}, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
			EMA: EMA{sma: SMA{length: 3, valid: true}, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(31),
				decimal.NewFromInt(1),
				decimal.NewFromInt(1),
				decimal.NewFromInt(2),
				decimal.NewFromInt(3),
			},
			Result: decimal.RequireFromString("4.75"),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.EMA.Calc(c.Data)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_EMA_CalcNext(t *testing.T) {
	cc := map[string]struct {
		EMA    EMA
		Last   decimal.Decimal
		Next   decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			EMA:   EMA{},
			Error: ErrInvalidIndicator,
		},
		"Successful calculation": {
			EMA:    EMA{sma: SMA{length: 3, valid: true}, valid: true},
			Last:   decimal.NewFromInt(5),
			Next:   decimal.NewFromInt(5),
			Result: decimal.NewFromInt(5),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.EMA.CalcNext(c.Last, c.Next)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_EMA_Count(t *testing.T) {
	e := EMA{sma: SMA{length: 15}}
	assert.Equal(t, 29, e.Count())
}

func Test_EMA_multiplier(t *testing.T) {
	e := EMA{sma: SMA{length: 3}}
	assert.Equal(t, decimal.RequireFromString("0.5").String(), e.multiplier().String())
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
			Result: EMA{sma: SMA{length: 1, valid: true}, valid: true},
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
	e := EMA{sma: SMA{length: 1}}
	d, err := e.MarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"length":1}`, string(d))
}

func Test_EMA_namedMarshalJSON(t *testing.T) {
	e := EMA{sma: SMA{length: 1}}
	d, err := e.namedMarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"name":"ema", "length":1}`, string(d))
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
			WMA:    WMA{length: 2, valid: true},
			Result: HMA{wma: WMA{length: 2, valid: true}, valid: true},
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

func Test_HMA_WMA(t *testing.T) {
	h := HMA{wma: WMA{length: 1}}
	assert.Equal(t, WMA{length: 1}, h.WMA())
}

func Test_HMA_validate(t *testing.T) {
	cc := map[string]struct {
		HMA   HMA
		Error error
		Valid bool
	}{
		"Invalid WMA": {
			HMA:   HMA{wma: WMA{length: -1}},
			Error: assert.AnError,
			Valid: false,
		},
		"Successful validation": {
			HMA:   HMA{wma: WMA{length: 2}},
			Valid: true,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			equalError(t, c.Error, c.HMA.validate())
			assert.Equal(t, c.Valid, c.HMA.valid)
		})
	}
}

func Test_HMA_Calc(t *testing.T) {
	cc := map[string]struct {
		HMA    HMA
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			HMA:   HMA{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			HMA: HMA{wma: WMA{length: 5, valid: true}, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
			HMA: HMA{wma: WMA{length: 3, valid: true}, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
				decimal.NewFromInt(31),
				decimal.NewFromInt(32),
				decimal.NewFromInt(30),
				decimal.NewFromInt(30),
				decimal.NewFromInt(31),
			},
			Result: decimal.RequireFromString("31.5"),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.HMA.Calc(c.Data)
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
		"Invalid HMA length": {
			JSON:  `{"wma":{"length":1}}`,
			Error: assert.AnError,
		},
		"Invalid WMA length": {
			JSON:  `{"wma":{"length":0}}`,
			Error: assert.AnError,
		},
		"Successful unmarshal": {
			JSON:   `{"wma":{"length":3}}`,
			Result: HMA{wma: WMA{length: 3, valid: true}, valid: true},
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
	h := HMA{wma: WMA{length: 3}}
	d, err := h.MarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"wma":{"length":3}}`, string(d))
}

func Test_HMA_namedMarshalJSON(t *testing.T) {
	h := HMA{wma: WMA{length: 3}}
	d, err := h.namedMarshalJSON()

	assert.NoError(t, err)
	assert.JSONEq(t, `{"name":"hma","wma":{"length":3}}`, string(d))
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
			Source1: &IndicatorMock{},
			Source2: &IndicatorMock{},
			Result:  MACD{source1: &IndicatorMock{}, source2: &IndicatorMock{}, valid: true},
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

func Test_MACD_Sub1(t *testing.T) {
	m := MACD{source1: &IndicatorMock{}, source2: nil}
	assert.Equal(t, &IndicatorMock{}, m.Sub1())
}

func Test_MACD_Sub2(t *testing.T) {
	m := MACD{source1: nil, source2: &IndicatorMock{}}
	assert.Equal(t, &IndicatorMock{}, m.Sub2())
}

func Test_MACD_validate(t *testing.T) {
	cc := map[string]struct {
		MACD  MACD
		Error error
		Valid bool
	}{
		"Invalid source1": {
			MACD:  MACD{source1: nil, source2: &IndicatorMock{}},
			Error: ErrInvalidSource,
			Valid: false,
		},
		"Invalid source2": {
			MACD:  MACD{source1: &IndicatorMock{}, source2: nil},
			Error: ErrInvalidSource,
			Valid: false,
		},
		"Successful MACD": {
			MACD:  MACD{source1: &IndicatorMock{}, source2: &IndicatorMock{}},
			Valid: true,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			equalError(t, c.Error, c.MACD.validate())
			assert.Equal(t, c.Valid, c.MACD.valid)
		})
	}
}

func Test_MACD_Calc(t *testing.T) {
	stubIndicator := func(v decimal.Decimal, e error, a int) *IndicatorMock {
		return &IndicatorMock{
			CalcFunc: func(dd []decimal.Decimal) (decimal.Decimal, error) {
				return v, e
			},
			CountFunc: func() int {
				return a
			},
		}
	}

	cc := map[string]struct {
		MACD   MACD
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			MACD:  MACD{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size for source1": {
			MACD: MACD{source1: stubIndicator(decimal.Zero, nil, 10), source2: stubIndicator(decimal.Zero, nil, 1), valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Invalid data size for source2": {
			MACD: MACD{source1: stubIndicator(decimal.Zero, nil, 1), source2: stubIndicator(decimal.Zero, nil, 10), valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Invalid source1": {
			MACD: MACD{source1: stubIndicator(decimal.Zero, assert.AnError, 1), source2: stubIndicator(decimal.Zero, nil, 1), valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: assert.AnError,
		},
		"Invalid source2": {
			MACD: MACD{source1: stubIndicator(decimal.Zero, nil, 1), source2: stubIndicator(decimal.Zero, assert.AnError, 1), valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: assert.AnError,
		},
		"Successful calculation": {
			MACD: MACD{source1: stubIndicator(decimal.NewFromInt(5), nil, 1), source2: stubIndicator(decimal.NewFromInt(10), nil, 1), valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
				decimal.NewFromInt(31),
				decimal.NewFromInt(32),
			},
			Result: decimal.NewFromInt(-5),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.MACD.Calc(c.Data)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_MACD_Count(t *testing.T) {
	stubIndicator := func(a int) *IndicatorMock {
		return &IndicatorMock{
			CountFunc: func() int {
				return a
			},
		}
	}

	m := MACD{source1: stubIndicator(5), source2: stubIndicator(10)}
	assert.Equal(t, m.Count(), 10)

	m = MACD{source1: stubIndicator(15), source2: stubIndicator(10)}
	assert.Equal(t, m.Count(), 15)
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
			JSON:  `{"source2":{"name":"sma","length":1}}`,
			Error: assert.AnError,
		},
		"Invalid source2": {
			JSON:  `{"source1":{"name":"sma","length":1}}`,
			Error: assert.AnError,
		},
		"Successful unmarshal": {
			JSON:   `{"source1":{"name":"sma","length":1},"source2":{"name":"sma","length":2}}`,
			Result: MACD{source1: SMA{length: 1, valid: true}, source2: SMA{length: 2, valid: true}, valid: true},
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
	stubIndicator := func(d []byte, e error) *IndicatorMock {
		return &IndicatorMock{
			namedMarshalJSONFunc: func() ([]byte, error) {
				return d, e
			},
		}
	}

	cc := map[string]struct {
		MACD   MACD
		Result string
		Error  error
	}{
		"Error returned during source1 marshalling": {
			MACD: MACD{
				source1: stubIndicator(nil, assert.AnError),
				source2: stubIndicator([]byte(`{"name":"indicatormock"}`), nil),
			},
			Error: assert.AnError,
		},
		"Error returned during source2 marshalling": {
			MACD: MACD{
				source1: stubIndicator([]byte(`{"name":"indicatormock"}`), nil),
				source2: stubIndicator(nil, assert.AnError),
			},
			Error: assert.AnError,
		},
		"Successful marshal": {
			MACD: MACD{
				source1: stubIndicator([]byte(`{"name":"indicatormock"}`), nil),
				source2: stubIndicator([]byte(`{"name":"indicatormock"}`), nil),
			},
			Result: `{"source1":{"name":"indicatormock"},"source2":{"name":"indicatormock"}}`,
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
	stubIndicator := func(d []byte, e error) *IndicatorMock {
		return &IndicatorMock{
			namedMarshalJSONFunc: func() ([]byte, error) {
				return d, e
			},
		}
	}

	cc := map[string]struct {
		MACD   MACD
		Result string
		Error  error
	}{
		"Error returned during source1 marshalling": {
			MACD: MACD{
				source1: stubIndicator(nil, assert.AnError),
				source2: stubIndicator([]byte(`{"name":"indicatormock"}`), nil),
			},
			Error: assert.AnError,
		},
		"Error returned during source2 marshalling": {
			MACD: MACD{
				source1: stubIndicator([]byte(`{"name":"indicatormock"}`), nil),
				source2: stubIndicator(nil, assert.AnError),
			},
			Error: assert.AnError,
		},
		"Successful marshal": {
			MACD: MACD{
				source1: stubIndicator([]byte(`{"name":"indicatormock"}`), nil),
				source2: stubIndicator([]byte(`{"name":"indicatormock"}`), nil),
			},
			Result: `{"name":"macd","source1":{"name":"indicatormock"},"source2":{"name":"indicatormock"}}`,
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
			Result: ROC{length: 1, valid: true},
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

func Test_ROC_Length(t *testing.T) {
	r := ROC{length: 1}
	assert.Equal(t, 1, r.Length())
}

func Test_ROC_validate(t *testing.T) {
	cc := map[string]struct {
		ROC   ROC
		Error error
		Valid bool
	}{
		"Invalid length": {
			ROC:   ROC{length: -1},
			Error: ErrInvalidLength,
			Valid: false,
		},
		"Successful validation": {
			ROC:   ROC{length: 1},
			Valid: true,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			equalError(t, c.Error, c.ROC.validate())
			assert.Equal(t, c.Valid, c.ROC.valid)
		})
	}
}

func Test_ROC_Calc(t *testing.T) {
	cc := map[string]struct {
		ROC    ROC
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			ROC:   ROC{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			ROC: ROC{length: 3, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful handled division by 0": {
			ROC: ROC{length: 5, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(420),
				decimal.NewFromInt(0),
				decimal.NewFromInt(420),
				decimal.NewFromInt(420),
				decimal.NewFromInt(420),
				decimal.NewFromInt(420),
			},
			Result: decimal.Zero,
		},
		"Successful calculation": {
			ROC: ROC{length: 5, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(7),
				decimal.NewFromInt(420),
				decimal.NewFromInt(420),
				decimal.NewFromInt(420),
				decimal.NewFromInt(10),
			},
			Result: decimal.RequireFromString("42.85714285714286"),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.ROC.Calc(c.Data)
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
			Result: ROC{length: 1, valid: true},
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
			Result: RSI{length: 1, valid: true},
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

func Test_RSI_Length(t *testing.T) {
	r := RSI{length: 1}
	assert.Equal(t, 1, r.Length())
}

func Test_RSI_validate(t *testing.T) {
	cc := map[string]struct {
		RSI   RSI
		Error error
		Valid bool
	}{
		"Invalid length": {
			RSI:   RSI{length: 0},
			Error: ErrInvalidLength,
			Valid: false,
		},
		"Successful validation": {
			RSI:   RSI{length: 1},
			Valid: true,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			equalError(t, c.Error, c.RSI.validate())
			assert.Equal(t, c.Valid, c.RSI.valid)
		})
	}
}

func Test_RSI_Calc(t *testing.T) {
	cc := map[string]struct {
		RSI    RSI
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			RSI:   RSI{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			RSI: RSI{length: 3, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation when average gain 0": {
			RSI: RSI{length: 3, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(16),
				decimal.NewFromInt(12),
				decimal.NewFromInt(8),
			},
			Result: decimal.NewFromInt(0),
		},
		"Successful calculation when average loss 0": {
			RSI: RSI{length: 3, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(2),
				decimal.NewFromInt(4),
				decimal.NewFromInt(8),
			},
			Result: decimal.NewFromInt(100),
		},
		"Successful calculation": {
			RSI: RSI{length: 3, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(8),
				decimal.NewFromInt(12),
				decimal.NewFromInt(8),
			},
			Result: decimal.NewFromInt(50),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.RSI.Calc(c.Data)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
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
			Result: RSI{length: 1, valid: true},
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
			Result: SMA{length: 1, valid: true},
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

func Test_SMA_Length(t *testing.T) {
	s := SMA{length: 1}
	assert.Equal(t, 1, s.Length())
}

func Test_SMA_validate(t *testing.T) {
	cc := map[string]struct {
		SMA   SMA
		Error error
		Valid bool
	}{
		"Invalid length": {
			SMA:   SMA{length: 0},
			Error: ErrInvalidLength,
			Valid: false,
		},
		"Successful validation": {
			SMA:   SMA{length: 1},
			Valid: true,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			equalError(t, c.Error, c.SMA.validate())
			assert.Equal(t, c.Valid, c.SMA.valid)
		})
	}
}

func Test_SMA_Calc(t *testing.T) {
	cc := map[string]struct {
		SMA    SMA
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			SMA:   SMA{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			SMA: SMA{length: 3, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
			SMA: SMA{length: 3, valid: true},
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

			res, err := c.SMA.Calc(c.Data)
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
			Result: SMA{length: 1, valid: true},
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
			RSI:    RSI{length: 1, valid: true},
			Result: SRSI{rsi: RSI{length: 1, valid: true}, valid: true},
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

func Test_SRSI_RSI(t *testing.T) {
	s := SRSI{rsi: RSI{length: 1}}
	assert.Equal(t, RSI{length: 1}, s.RSI())
}

func Test_SRSI_validate(t *testing.T) {
	cc := map[string]struct {
		SRSI  SRSI
		Error error
		Valid bool
	}{
		"Invalid RSI": {
			SRSI:  SRSI{},
			Error: assert.AnError,
			Valid: false,
		},
		"Invalid RSI length": {
			SRSI:  SRSI{rsi: RSI{length: -1}},
			Error: assert.AnError,
			Valid: false,
		},
		"Successful validation": {
			SRSI:  SRSI{rsi: RSI{length: 1}},
			Valid: true,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			equalError(t, c.Error, c.SRSI.validate())
			assert.Equal(t, c.Valid, c.SRSI.valid)
		})
	}
}

func Test_SRSI_Calc(t *testing.T) {
	cc := map[string]struct {
		SRSI   SRSI
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			SRSI:  SRSI{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			SRSI: SRSI{rsi: RSI{length: 5, valid: true}, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: assert.AnError,
		},
		"Successfully handled division by 0": {
			SRSI: SRSI{rsi: RSI{length: 3, valid: true}, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(8),
				decimal.NewFromInt(12),
				decimal.NewFromInt(8),
				decimal.NewFromInt(12),
				decimal.NewFromInt(8),
			},
			Result: decimal.Zero,
		},
		"Successful calculation": {
			SRSI: SRSI{rsi: RSI{length: 3, valid: true}, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(4),
				decimal.NewFromInt(10),
				decimal.NewFromInt(6),
				decimal.NewFromInt(8),
				decimal.NewFromInt(6),
			},
			Result: decimal.RequireFromString("0.625"),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.SRSI.Calc(c.Data)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
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
			Result: SRSI{rsi: RSI{length: 1, valid: true}, valid: true},
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
			Result: Stoch{length: 1, valid: true},
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

func Test_Stoch_Length(t *testing.T) {
	s := Stoch{length: 1}
	assert.Equal(t, 1, s.Length())
}

func Test_Stoch_validate(t *testing.T) {
	cc := map[string]struct {
		Stoch Stoch
		Error error
		Valid bool
	}{
		"Invalid length": {
			Stoch: Stoch{length: 0},
			Error: ErrInvalidLength,
			Valid: false,
		},
		"Successful validation": {
			Stoch: Stoch{length: 1},
			Valid: true,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			equalError(t, c.Error, c.Stoch.validate())
			assert.Equal(t, c.Valid, c.Stoch.valid)
		})
	}
}

func Test_Stoch_Calc(t *testing.T) {
	cc := map[string]struct {
		Stoch  Stoch
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			Stoch: Stoch{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			Stoch: Stoch{length: 3, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation when new lows are reached": {
			Stoch: Stoch{length: 3, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(150),
				decimal.NewFromInt(125),
				decimal.NewFromInt(145),
			},
			Result: decimal.NewFromInt(80),
		},
		"Successfully handled division by 0": {
			Stoch: Stoch{length: 3, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(150),
				decimal.NewFromInt(150),
				decimal.NewFromInt(150),
			},
			Result: decimal.Zero,
		},
		"Successful calculation when new highs are reached": {
			Stoch: Stoch{length: 3, valid: true},
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

			res, err := c.Stoch.Calc(c.Data)
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
			Result: Stoch{length: 1, valid: true},
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
			Result: WMA{length: 1, valid: true},
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

func Test_WMA_Length(t *testing.T) {
	w := WMA{length: 1}
	assert.Equal(t, 1, w.Length())
}

func Test_WMA_validate(t *testing.T) {
	cc := map[string]struct {
		WMA   WMA
		Error error
		Valid bool
	}{
		"Invalid length": {
			WMA:   WMA{length: 0},
			Error: ErrInvalidLength,
			Valid: false,
		},
		"Successful validation": {
			WMA:   WMA{length: 1},
			Valid: true,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			equalError(t, c.Error, c.WMA.validate())
			assert.Equal(t, c.Valid, c.WMA.valid)
		})
	}
}

func Test_WMA_Calc(t *testing.T) {
	cc := map[string]struct {
		WMA    WMA
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			WMA:   WMA{},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			WMA: WMA{length: 3, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
			WMA: WMA{length: 3, valid: true},
			Data: []decimal.Decimal{
				decimal.NewFromInt(420),
				decimal.NewFromInt(420),
				decimal.NewFromInt(420),
				decimal.NewFromInt(30),
				decimal.NewFromInt(30),
				decimal.NewFromInt(32),
			},
			Result: decimal.NewFromInt(31),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.WMA.Calc(c.Data)
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
			Result: WMA{length: 1, valid: true},
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
