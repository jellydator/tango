package indc

import (
	"testing"

	"github.com/swithek/chartype"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestResize(t *testing.T) {
	cc := map[string]struct {
		Length int
		Data   []decimal.Decimal
		Result []decimal.Decimal
		Error  error
	}{
		"Insufficient amount of data points": {
			Length: 3,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataPointCount,
		},
		"Invalid length": {
			Length: -3,
			Error:  ErrInvalidLength,
		},
		"Successful resize": {
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

func TestResizeCandles(t *testing.T) {
	cc := map[string]struct {
		Length int
		Data   []chartype.Candle
		Result []chartype.Candle
		Error  error
	}{
		"Insufficient amount of data points": {
			Length: 3,
			Data: []chartype.Candle{
				{Close: decimal.NewFromInt(30)},
			},
			Error: ErrInvalidDataPointCount,
		},
		"Invalid length": {
			Length: -3,
			Error:  ErrInvalidLength,
		},
		"Successful resize": {
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

func TestTypicalPrice(t *testing.T) {
	cc := map[string]struct {
		Data   []chartype.Candle
		Result []decimal.Decimal
	}{
		"Successful typical price": {
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

func TestMeanDeviation(t *testing.T) {
	cc := map[string]struct {
		Data   []decimal.Decimal
		Result decimal.Decimal
	}{
		"Successful mean deviation": {
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

func TestFromJSON(t *testing.T) {
	cc := map[string]struct {
		Name      string
		ByteArray []byte
		Result    Indicator
		Error     error
	}{
		"Successful creation of an aroon indicator": {
			Name:      "aroon",
			ByteArray: []byte(`{"name":"aroon","trend":"up","length":1}`),
			Result:    Aroon{Trend: "up", Length: 1},
		},
		"Successful creation of a cci indicator": {
			Name: "cci",
			ByteArray: []byte(`{"name":"cci",
			"source":{"name":"aroon","trend":"up","length":1}}`),
			Result: CCI{Source{Aroon{Trend: "up", Length: 1}}},
		},
		"Successful creation of a dema indicator": {
			Name:      "dema",
			ByteArray: []byte(`{"name":"dema","length":1}`),
			Result:    DEMA{Length: 1},
		},
		"Successful creation of an ema indicator": {
			Name:      "ema",
			ByteArray: []byte(`{"name":"ema","length":1}`),
			Result:    EMA{Length: 1},
		},
		"Successful creation of a hma indicator": {
			Name:      "hma",
			ByteArray: []byte(`{"name":"hma", "wma":{"length":1}}`),
			Result:    HMA{WMA: WMA{Length: 1}},
		},
		"Successful creation of a macd indicator": {
			Name: "macd",
			ByteArray: []byte(`{"name":"macd",
			"source1":{"name":"aroon","trend":"down","length":2},
			"source2":{"name":"cci","source":{"name":"ema", "length":2}}}`),
			Result: MACD{Source{Aroon{Trend: "down", Length: 2}},
				Source{CCI{Source{EMA{Length: 2}}}}},
		},
		"Successful creation of a roc indicator": {
			Name:      "roc",
			ByteArray: []byte(`{"name":"roc","length":1}`),
			Result:    ROC{Length: 1},
		},
		"Successful creation of a rsi indicator": {
			Name:      "rsi",
			ByteArray: []byte(`{"name":"rsi","length":1}`),
			Result:    RSI{Length: 1},
		},
		"Successful creation of a sma indicator": {
			Name:      "sma",
			ByteArray: []byte(`{"name":"sma","length":1}`),
			Result:    SMA{Length: 1},
		},
		"Successful creation of a stoch indicator": {
			Name:      "stoch",
			ByteArray: []byte(`{"name":"stoch","length":1}`),
			Result:    Stoch{Length: 1},
		},
		"Successful creation of an wma indicator": {
			Name:      "wma",
			ByteArray: []byte(`{"name":"wma","length":1}`),
			Result:    WMA{Length: 1},
		},
		"Invalid indicator name": {
			Name:  "tema",
			Error: ErrInvalidSourceName,
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := fromJSON(c.Name, c.ByteArray)
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

func TestToJSON(t *testing.T) {
	cc := map[string]struct {
		Indicator Indicator
		Result    []byte
		Error     error
	}{
		"Successful marshal of an aroon indicator": {
			Indicator: Aroon{Trend: "up", Length: 1},
			Result:    []byte(`{"name":"aroon","trend":"up","length":1}`),
		},
		"Successful marshal of a cci indicator": {
			Indicator: CCI{Source{Aroon{Trend: "up", Length: 1}}},
			Result:    []byte(`{"name":"cci","source":{"name":"aroon","trend":"up","length":1}}`),
		},
		"Successful marshal of a dema indicator": {
			Indicator: DEMA{Length: 1},
			Result:    []byte(`{"name":"dema","length":1}`),
		},
		"Successful marshal of an ema indicator": {
			Indicator: EMA{Length: 1},
			Result:    []byte(`{"name":"ema","length":1}`),
		},
		"Successful marshal of a hma indicator": {
			Indicator: HMA{WMA: WMA{Length: 1}},
			Result:    []byte(`{"name":"hma","wma":{"length":1}}`),
		},
		"Successful marshal of a macd indicator": {
			Indicator: MACD{Source{Aroon{Trend: "down", Length: 2}},
				Source{CCI{Source{EMA{Length: 2}}}}},
			Result: []byte(`{"name":"macd","source1":{"name":"aroon","trend":"down","length":2},"source2":{"name":"cci","source":{"name":"ema","length":2}}}`),
		},
		"Successful marshal of a roc indicator": {
			Indicator: ROC{Length: 1},
			Result:    []byte(`{"name":"roc","length":1}`),
		},
		"Successful marshal of a rsi indicator": {
			Indicator: RSI{Length: 1},
			Result:    []byte(`{"name":"rsi","length":1}`),
		},
		"Successful marshal of a sma indicator": {
			Indicator: SMA{Length: 1},
			Result:    []byte(`{"name":"sma","length":1}`),
		},
		"Successful marshal of a stoch indicator": {
			Indicator: Stoch{Length: 1},
			Result:    []byte(`{"name":"stoch","length":1}`),
		},
		"Successful marshal of a wma indicator": {
			Indicator: WMA{Length: 1},
			Result:    []byte(`{"name":"wma","length":1}`),
		},
		"Invalid indicator": {
			Indicator: IndicatorMock{},
			Error:     ErrInvalidType,
		},
	}

	for cn, c := range cc {
		c := c
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := toJSON(c.Indicator)
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
