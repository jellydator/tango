package indc

import (
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/swithek/chartype"
)

func equalError(t *testing.T, exp, err error) {
	t.Helper()

	if exp != nil {
		if errors.Is(exp, assert.AnError) {
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
		Data   []decimal.Decimal
		Result []decimal.Decimal
		Error  error
	}{
		"Invalid data size": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Unmodified slice returned when length is 1": {
			Length: 0,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Result: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
		},
		"Successful computation": {
			Length: 3,
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

			res, err := resize(c.Data, c.Length)
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
		Data   []chartype.Candle
		Result []chartype.Candle
		Error  error
	}{
		"Invalid data size": {
			Length: 3,
			Data: []chartype.Candle{
				{Close: decimal.NewFromInt(30)},
			},
			Error: ErrInvalidDataSize,
		},
		"Unmodified slice returned when length is 1": {
			Length: 0,
			Data: []chartype.Candle{
				{Close: decimal.NewFromInt(30)},
			},
			Result: []chartype.Candle{
				{Close: decimal.NewFromInt(30)},
			},
		},
		"Successful computation": {
			Length: 3,
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

			res, err := resizeCandles(c.Data, c.Length)
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
				{High: decimal.NewFromFloat(24.2), Low: decimal.NewFromFloat(23.85), Close: decimal.NewFromFloat(23.89)},
				{High: decimal.NewFromFloat(24.07), Low: decimal.NewFromFloat(23.72), Close: decimal.NewFromFloat(23.95)},
				{High: decimal.NewFromFloat(24.04), Low: decimal.NewFromFloat(23.64), Close: decimal.NewFromFloat(23.67)},
			},
			Result: []decimal.Decimal{
				decimal.NewFromFloat(23.98),
				decimal.NewFromFloat(23.91333333),
				decimal.NewFromFloat(23.78333333),
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
			Result: decimal.NewFromFloat(0),
		},
		"Successful calculation with one value": {
			Data: []decimal.Decimal{
				decimal.NewFromInt(2),
			},
			Result: decimal.NewFromFloat(0),
		},
		"Successful calculation": {
			Data: []decimal.Decimal{
				decimal.NewFromInt(2),
				decimal.NewFromInt(5),
				decimal.NewFromInt(8),
			},
			Result: decimal.NewFromFloat(2),
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
			Indicator: stubIndicator(decimal.Zero, nil, 2),
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
			Indicator: stubIndicator(decimal.Zero, nil, 2),
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
			Indicator: stubIndicator(decimal.NewFromInt(2), nil, 2),
			Result: []decimal.Decimal{
				decimal.NewFromFloat(2),
				decimal.NewFromFloat(2),
				decimal.NewFromFloat(2),
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := calcMultiple(c.Data, c.Amount, c.Indicator)
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
			ByteArray: []byte(`{"name":"aroon","trend":"up","length":1}`),
			Result:    Aroon{trend: "up", length: 1},
		},
		"Successful creation of CCI": {
			ByteArray: []byte(`{"name":"cci",
			"source":{"name":"aroon","trend":"up","length":1}}`),
			Result: CCI{Aroon{trend: "up", length: 1}},
		},
		"Successful creation of DEMA": {
			ByteArray: []byte(`{"name":"dema","ema":{"sma":{"length":1}}}`),
			Result:    DEMA{ema: EMA{sma: SMA{length: 1}}},
		},
		"Successful creation of EMA": {
			ByteArray: []byte(`{"name":"ema","sma":{"length":1}}`),
			Result:    EMA{sma: SMA{length: 1}},
		},
		"Successful creation of HMA": {
			ByteArray: []byte(`{"name":"hma", "wma":{"name":"wma","length":2}}`),
			Result:    HMA{wma: WMA{length: 2}},
		},
		"Successful creation of MACD": {
			ByteArray: []byte(`{"name":"macd",
			"source1":{"name":"aroon","trend":"down","length":2},
			"source2":{"name":"cci","source":{"name":"ema",
			 "sma":{"length":2}}}}`),
			Result: MACD{Aroon{trend: "down", length: 2},
				CCI{EMA{SMA{length: 2}}}},
		},
		"Successful creation of ROC": {
			ByteArray: []byte(`{"name":"roc","length":1}`),
			Result:    ROC{length: 1},
		},
		"Successful creation of RSI": {
			ByteArray: []byte(`{"name":"rsi","length":1}`),
			Result:    RSI{length: 1},
		},
		"Successful creation of SMA": {
			ByteArray: []byte(`{"name":"sma","length":1}`),
			Result:    SMA{length: 1},
		},
		"Successful creation of SRSI": {
			ByteArray: []byte(`{"name":"srsi", "rsi":{"name":"rsi","length":1}}`),
			Result:    SRSI{rsi: RSI{length: 1}},
		},
		"Successful creation of Stoch": {
			ByteArray: []byte(`{"name":"stoch","length":1}`),
			Result:    Stoch{length: 1},
		},
		"Successful creation of WMA": {
			ByteArray: []byte(`{"name":"wma","length":1}`),
			Result:    WMA{length: 1},
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
