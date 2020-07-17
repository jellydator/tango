package indc

import (
	"testing"

	"github.com/jellydator/chartype"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func equalError(t *testing.T, exp, err error) {
	t.Helper()

	if exp != nil {
		if exp == assert.AnError { //nolint:goerr113 // direct check is needed
			assert.Error(t, err)
			return
		}

		assert.Equal(t, exp, err)

		return
	}

	assert.NoError(t, err)
}

func Test_CleanString(t *testing.T) {
	var e String = "aroon"

	r := CleanString(" aRooN ")

	assert.Equal(t, e, r)
}

func Test_String_UnmarshalText(t *testing.T) {
	var s String

	assert.NoError(t, s.UnmarshalText([]byte("   TEST       ")))
	assert.Equal(t, "test", string(s))
}

func Test_String_MarshalText(t *testing.T) {
	var s String = "test"
	v, err := s.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte("test"), v)
}

func Test_resize(t *testing.T) {
	cc := map[string]struct {
		Length int
		Offset int
		Data   []decimal.Decimal
		Result []decimal.Decimal
		Error  error
	}{
		"Invalid data size": {
			Length: 3,
			Offset: 0,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Unmodified slice returned when length is 1": {
			Length: 0,
			Offset: 0,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Result: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
		},
		"Successful computation": {
			Length: 3,
			Offset: 0,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
				decimal.NewFromInt(31),
				decimal.NewFromInt(32),
				decimal.NewFromInt(32),
				decimal.NewFromInt(32),
				decimal.NewFromInt(32),
			},
			Result: []decimal.Decimal{
				decimal.NewFromInt(32),
				decimal.NewFromInt(32),
				decimal.NewFromInt(32),
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := resize(c.Data, c.Length, c.Offset)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			for i := 0; i < len(c.Result); i++ {
				assert.Equal(t, c.Result[i].Round(8), res[i].Round(8))
			}
		})
	}
}

func Test_resizeCandles(t *testing.T) {
	cc := map[string]struct {
		Length int
		Offset int
		Data   []chartype.Candle
		Result []chartype.Candle
		Error  error
	}{
		"Invalid data size": {
			Length: 3,
			Offset: 0,
			Data: []chartype.Candle{
				{Close: decimal.NewFromInt(30)},
			},
			Error: ErrInvalidDataSize,
		},
		"Unmodified slice returned when length is 1": {
			Length: 0,
			Offset: 0,
			Data: []chartype.Candle{
				{Close: decimal.NewFromInt(30)},
			},
			Result: []chartype.Candle{
				{Close: decimal.NewFromInt(30)},
			},
		},
		"Successful computation": {
			Length: 3,
			Offset: 0,
			Data: []chartype.Candle{
				{Close: decimal.NewFromInt(30)},
				{Close: decimal.NewFromInt(31)},
				{Close: decimal.NewFromInt(32)},
				{Close: decimal.NewFromInt(32)},
				{Close: decimal.NewFromInt(32)},
				{Close: decimal.NewFromInt(32)},
			},
			Result: []chartype.Candle{
				{Close: decimal.NewFromInt(32)},
				{Close: decimal.NewFromInt(32)},
				{Close: decimal.NewFromInt(32)},
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := resizeCandles(c.Data, c.Length, c.Offset)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			for i := 0; i < len(c.Result); i++ {
				assert.Equal(t, c.Result[i].Close.Round(8), res[i].Close.Round(8))
			}
		})
	}
}

func Test_typicalPrice(t *testing.T) {
	cc := map[string]struct {
		Data   []chartype.Candle
		Result []decimal.Decimal
	}{
		"Successful calculation": {
			Data: []chartype.Candle{
				{High: decimal.RequireFromString("24.2"), Low: decimal.RequireFromString("23.85"), Close: decimal.RequireFromString("23.89")},
				{High: decimal.RequireFromString("24.07"), Low: decimal.RequireFromString("23.72"), Close: decimal.RequireFromString("23.95")},
				{High: decimal.RequireFromString("24.04"), Low: decimal.RequireFromString("23.64"), Close: decimal.RequireFromString("23.67")},
			},
			Result: []decimal.Decimal{
				decimal.RequireFromString("23.98"),
				decimal.RequireFromString("23.91333333"),
				decimal.RequireFromString("23.78333333"),
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res := typicalPrice(c.Data)

			for i := 0; i < len(c.Result); i++ {
				assert.Equal(t, c.Result[i].Round(8), res[i].Round(8))
			}
		})
	}
}

func Test_meanDeviation(t *testing.T) {
	cc := map[string]struct {
		Data   []decimal.Decimal
		Result decimal.Decimal
	}{
		"Successful calculation with no values": {
			Data:   []decimal.Decimal{},
			Result: decimal.NewFromInt(0),
		},
		"Successful calculation with one value": {
			Data: []decimal.Decimal{
				decimal.NewFromInt(2),
			},
			Result: decimal.NewFromInt(0),
		},
		"Successful calculation": {
			Data: []decimal.Decimal{
				decimal.NewFromInt(2),
				decimal.NewFromInt(5),
				decimal.NewFromInt(8),
			},
			Result: decimal.NewFromInt(2),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res := meanDeviation(c.Data)

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_calcMultiple(t *testing.T) {
	stubIndicator := func(ddv decimal.Decimal, count int, e error) *IndicatorMock {
		return &IndicatorMock{
			CalcFunc: func(dd []decimal.Decimal) (decimal.Decimal, error) {
				return ddv, e
			},
			CountFunc: func() int {
				return count
			},
		}
	}

	cc := map[string]struct {
		Data      []decimal.Decimal
		Amount    int
		Indicator Indicator
		Result    []decimal.Decimal
		Error     error
	}{
		"Invalid data size": {
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Indicator: stubIndicator(decimal.Zero, 2, nil),
			Amount:    1,
			Error:     ErrInvalidDataSize,
		},
		"Invalid Indicator": {
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Indicator: &IndicatorMock{
				CalcFunc: func(dd []decimal.Decimal) (decimal.Decimal, error) {
					return decimal.Zero, assert.AnError
				},
				CountFunc: func() int {
					return 1
				},
			},
			Amount: 1,
			Error:  assert.AnError,
		},
		"Successful calculation with amount less than 1": {
			Data: []decimal.Decimal{
				decimal.NewFromInt(2),
				decimal.NewFromInt(3),
				decimal.NewFromInt(4),
				decimal.NewFromInt(5),
				decimal.NewFromInt(6),
				decimal.NewFromInt(7),
			},
			Amount:    0,
			Indicator: stubIndicator(decimal.Zero, 2, nil),
			Result:    []decimal.Decimal{},
		},
		"Successful calculation with amount more than 1": {
			Data: []decimal.Decimal{
				decimal.NewFromInt(2),
				decimal.NewFromInt(3),
				decimal.NewFromInt(4),
				decimal.NewFromInt(5),
				decimal.NewFromInt(6),
				decimal.NewFromInt(7),
			},
			Amount:    3,
			Indicator: stubIndicator(decimal.NewFromInt(2), 2, nil),
			Result: []decimal.Decimal{
				decimal.NewFromInt(2),
				decimal.NewFromInt(2),
				decimal.NewFromInt(2),
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := calcMultiple(c.Indicator, c.Data, c.Amount)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			for i := 0; i < len(c.Result); i++ {
				assert.Equal(t, c.Result[i].Round(8), res[i].Round(8))
			}
		})
	}
}

func Test_fromJSON(t *testing.T) {
	cc := map[string]struct {
		ByteArray []byte
		Result    Indicator
		Error     error
	}{
		"Invalid JSON": {
			ByteArray: []byte(`{\"_"/`),
			Error:     assert.AnError,
		},
		"Invalid source name": {
			ByteArray: []byte(`{"name":"aa"}`),
			Error:     ErrInvalidSource,
		},
		"Successful creation of Aroon": {
			ByteArray: []byte(`{"name":"aroon","trend":"up","length":1,"offset":2}`),
			Result:    Aroon{trend: "up", length: 1, offset: 2, valid: true},
		},
		"Successful creation of CCI": {
			ByteArray: []byte(`{"name":"cci","source":{"name":"sma","length":1,"offset":3}}`),
			Result:    CCI{source: SMA{length: 1, offset: 3, valid: true}, factor: decimal.RequireFromString("0.015"), valid: true},
		},
		"Successful creation of DEMA": {
			ByteArray: []byte(`{"name":"dema","ema":{"length":1,"offset":1}}`),
			Result:    DEMA{ema: EMA{SMA{length: 1, offset: 1, valid: true}}, valid: true},
		},
		"Successful creation of EMA": {
			ByteArray: []byte(`{"name":"ema","length":1,"offset":3}`),
			Result:    EMA{SMA{length: 1, offset: 3, valid: true}},
		},
		"Successful creation of HMA": {
			ByteArray: []byte(`{"name":"hma", "wma":{"name":"wma","length":2, "offset":3}}`),
			Result:    HMA{wma: WMA{length: 2, offset: 3, valid: true}, valid: true},
		},
		"Successful creation of MACD": {
			ByteArray: []byte(`{"name":"macd",
			"source1":{"name":"sma","length":2,"offset":2},
			"source2":{"name":"sma","length":3,"offset":4},
			"offset":3}`),
			Result: MACD{source1: SMA{length: 2, offset: 2, valid: true},
				source2: SMA{length: 3, offset: 4, valid: true}, offset: 3, valid: true},
		},
		"Successful creation of ROC": {
			ByteArray: []byte(`{"name":"roc","length":1,"offset":3}`),
			Result:    ROC{length: 1, offset: 3, valid: true},
		},
		"Successful creation of RSI": {
			ByteArray: []byte(`{"name":"rsi","length":1,"offset":2}`),
			Result:    RSI{length: 1, offset: 2, valid: true},
		},
		"Successful creation of SMA": {
			ByteArray: []byte(`{"name":"sma","length":1,"offset":3}`),
			Result:    SMA{length: 1, offset: 3, valid: true},
		},
		"Successful creation of SRSI": {
			ByteArray: []byte(`{"name":"srsi", "rsi":{"name":"rsi","length":1,"offset":1}}`),
			Result:    SRSI{rsi: RSI{length: 1, offset: 1, valid: true}, valid: true},
		},
		"Successful creation of Stoch": {
			ByteArray: []byte(`{"name":"stoch","length":1,"offset":4}`),
			Result:    Stoch{length: 1, offset: 4, valid: true},
		},
		"Successful creation of WMA": {
			ByteArray: []byte(`{"name":"wma","length":1,"offset":5}`),
			Result:    WMA{length: 1, offset: 5, valid: true},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := fromJSON(c.ByteArray)
			equalError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result, res)
		})
	}
}
