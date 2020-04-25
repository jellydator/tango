package indc

import (
	"testing"

	"github.com/swithek/chartype"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_CleanString(t *testing.T) {
	var e String
	e = "aroon"
	r := CleanString(" aRooN ")
	assert.Equal(t, e, r)
}
func Test_resize(t *testing.T) {
	cc := map[string]struct {
		Length int
		Data   []decimal.Decimal
		Result []decimal.Decimal
		Error  error
	}{
		"Successfully resize returned an ErrInvalidDataSize with insufficient amount of data points": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successfully resize returned unchanged list with length less than 1": {
			Length: 0,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Result: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
		},
		"Successful resize computation": {
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
			if c.Error != nil {
				if c.Error == assert.AnError {
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, c.Error, err)
				}
			} else {
				assert.Nil(t, err)
				for i := 0; i < len(c.Result); i++ {
					assert.Equal(t, c.Result[i].Round(8), res[i].Round(8))
				}
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
		"Successfully resizeCandles returned an ErrInvalidDataSize with insufficient amount of data points": {
			Length: 3,
			Data: []chartype.Candle{
				{Close: decimal.NewFromInt(30)},
			},
			Error: ErrInvalidDataSize,
		},
		"Successfully resizeCandles returned unchanged list with length less than 1": {
			Length: 0,
			Data: []chartype.Candle{
				{Close: decimal.NewFromInt(30)},
			},
			Result: []chartype.Candle{
				{Close: decimal.NewFromInt(30)},
			},
		},
		"Successful resizeCandles computation": {
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
			if c.Error != nil {
				if c.Error == assert.AnError {
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, c.Error, err)
				}
			} else {
				assert.Nil(t, err)
				for i := 0; i < len(c.Result); i++ {
					assert.Equal(t, c.Result[i].Close.Round(8), res[i].Close.Round(8))
				}
			}
		})
	}
}

func Test_typicalPrice(t *testing.T) {
	cc := map[string]struct {
		Data   []chartype.Candle
		Result []decimal.Decimal
	}{
		"Successful typical price calculation": {
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
		"Successful mean deviation calculation": {
			Data: []decimal.Decimal{
				decimal.NewFromInt(2),
				decimal.NewFromInt(5),
				decimal.NewFromInt(7),
				decimal.NewFromInt(10),
				decimal.NewFromInt(12),
				decimal.NewFromInt(14),
			},
			Result: decimal.NewFromFloat(3.66666667),
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res := meanDeviation(c.Data)

			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_calcMultiple(t *testing.T) {
	cc := map[string]struct {
		Data      []decimal.Decimal
		Amount    int
		Indicator Indicator
		Result    []decimal.Decimal
		Error     error
	}{
		"Successfully calcMultiple returned an ErrInvalidDataSize with insufficient amount of data points": {
			Indicator: SMA{length: 2},
			Amount:    1,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successfully calcMultiple indicator returned an error": {
			Indicator: IndicatorMock{},
			Amount:    1,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: assert.AnError,
		},
		"Successful calcMultiple calculation with amount less than 1": {
			Data: []decimal.Decimal{
				decimal.NewFromInt(2),
				decimal.NewFromInt(3),
				decimal.NewFromInt(4),
				decimal.NewFromInt(5),
				decimal.NewFromInt(6),
				decimal.NewFromInt(7),
			},
			Amount:    0,
			Indicator: SMA{length: 2},
			Result:    []decimal.Decimal{},
		},
		"Successful calcMultiple calculation with amount more than 1": {
			Data: []decimal.Decimal{
				decimal.NewFromInt(2),
				decimal.NewFromInt(3),
				decimal.NewFromInt(4),
				decimal.NewFromInt(5),
				decimal.NewFromInt(6),
				decimal.NewFromInt(7),
			},
			Amount:    3,
			Indicator: SMA{length: 2},
			Result: []decimal.Decimal{
				decimal.NewFromFloat(6.5),
				decimal.NewFromFloat(5.5),
				decimal.NewFromFloat(4.5),
			},
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := calcMultiple(c.Data, c.Amount, c.Indicator)

			if c.Error != nil {
				if c.Error == assert.AnError {
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, c.Error, err)
				}
			} else {
				assert.Nil(t, err)

				for i := 0; i < len(c.Result); i++ {
					assert.Equal(t, c.Result[i].Round(8), res[i].Round(8))
				}
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
		"Successful creation of an Aroon indicator": {
			ByteArray: []byte(`{"name":"aroon","trend":"up","length":1}`),
			Result:    Aroon{trend: "up", length: 1},
		},
		"Successful creation of a CCI indicator": {
			ByteArray: []byte(`{"name":"cci",
			"source":{"name":"aroon","trend":"up","length":1}}`),
			Result: CCI{Aroon{trend: "up", length: 1}},
		},
		"Successful creation of a DEMA indicator": {
			ByteArray: []byte(`{"name":"dema","length":1}`),
			Result:    DEMA{length: 1},
		},
		"Successful creation of an EMA indicator": {
			ByteArray: []byte(`{"name":"ema","length":1}`),
			Result:    EMA{length: 1},
		},
		"Successful creation of a HMA indicator": {
			ByteArray: []byte(`{"name":"hma", "wma":{"name":"wma","length":1}}`),
			Result:    HMA{wma: WMA{length: 1}},
		},
		"Successful creation of a MACD indicator": {
			ByteArray: []byte(`{"name":"macd",
			"source1":{"name":"aroon","trend":"down","length":2},
			"source2":{"name":"cci","source":{"name":"ema", "length":2}}}`),
			Result: MACD{Aroon{trend: "down", length: 2},
				CCI{EMA{length: 2}}},
		},
		"Successful creation of a ROC indicator": {
			ByteArray: []byte(`{"name":"roc","length":1}`),
			Result:    ROC{length: 1},
		},
		"Successful creation of a RSI indicator": {
			ByteArray: []byte(`{"name":"rsi","length":1}`),
			Result:    RSI{length: 1},
		},
		"Successful creation of a SMA indicator": {
			ByteArray: []byte(`{"name":"sma","length":1}`),
			Result:    SMA{length: 1},
		},
		"Successful creation of a SRSI indicator": {
			ByteArray: []byte(`{"name":"srsi", "rsi":{"name":"rsi","length":1}}`),
			Result:    SRSI{rsi: RSI{length: 1}},
		},
		"Successful creation of a Stoch indicator": {
			ByteArray: []byte(`{"name":"stoch","length":1}`),
			Result:    Stoch{length: 1},
		},
		"Successful creation of an WMA indicator": {
			ByteArray: []byte(`{"name":"wma","length":1}`),
			Result:    WMA{length: 1},
		},
		"Successfully fromJSON JSON unmarshal returned an error": {
			ByteArray: []byte(`{\"_"/`),
			Error:     assert.AnError,
		},
		"Successfully fromJSON returned an ErrInvalidDataPointCount with invalid source name": {
			ByteArray: []byte(`{"name":"aa"}`),
			Error:     ErrInvalidSource,
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := fromJSON(c.ByteArray)
			if c.Error != nil {
				if c.Error == assert.AnError {
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, c.Error, err)
				}
			} else {
				assert.Nil(t, err)
				assert.Equal(t, c.Result, res)
			}
		})
	}
}
