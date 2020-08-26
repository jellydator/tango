package indc

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

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
		"Invalid Aroon": {
			ByteArray: []byte(`{"name":"aroon","trend":"up","length":-1,"offset":2}`),
			Error:     assert.AnError,
		},
		"Invalid BB": {
			ByteArray: []byte(`{"name":"bb","band":"upper","std_dev":"2","length":-1,"offset":2}`),
			Error:     assert.AnError,
		},
		"Invalid CCI": {
			ByteArray: []byte(`{"name":"cci","source":{"name":"sma","length":-1,"offset":3}}`),
			Error:     assert.AnError,
		},
		"Invalid DEMA": {
			ByteArray: []byte(`{"name":"dema","ema":{"length":-1,"offset":1}}`),
			Error:     assert.AnError,
		},
		"Invalid EMA": {
			ByteArray: []byte(`{"name":"ema","length":-1,"offset":3}`),
			Error:     assert.AnError,
		},
		"Invalid HMA": {
			ByteArray: []byte(`{"name":"hma", "wma":{"name":"wma","length":-2, "offset":3}}`),
			Error:     assert.AnError,
		},
		"Invalid CD": {
			ByteArray: []byte(`{"name":"cd",
			"source1":{"name":"sma","length":-2,"offset":2},
			"source2":{"name":"sma","length":3,"offset":4},
			"offset":3}`),
			Error: assert.AnError,
		},
		"Invalid ROC": {
			ByteArray: []byte(`{"name":"roc","length":-1,"offset":3}`),
			Error:     assert.AnError,
		},
		"Invalid RSI": {
			ByteArray: []byte(`{"name":"rsi","length":-1,"offset":2}`),
			Error:     assert.AnError,
		},
		"Invalid SMA": {
			ByteArray: []byte(`{"name":"sma","length":-1,"offset":3}`),
			Error:     assert.AnError,
		},
		"Invalid SRSI": {
			ByteArray: []byte(`{"name":"srsi", "rsi":{"name":"rsi","length":-1,"offset":1}}`),
			Error:     assert.AnError,
		},
		"Invalid Stoch": {
			ByteArray: []byte(`{"name":"stoch","length":-1,"offset":4}`),
			Error:     assert.AnError,
		},
		"Invalid WMA": {
			ByteArray: []byte(`{"name":"wma","length":-1,"offset":5}`),
			Error:     assert.AnError,
		},
		"Successful Aroon unmarshal": {
			ByteArray: []byte(`{"name":"aroon","trend":"up","length":1,"offset":2}`),
			Result:    Aroon{trend: TrendUp, length: 1, offset: 2, valid: true},
		},
		"Successful BB unmarshal": {
			ByteArray: []byte(`{"name":"bb","band":"upper","std_dev":"2","length":1,"offset":2}`),
			Result:    BB{band: BandUpper, stdDev: decimal.RequireFromString("2"), length: 1, offset: 2, valid: true},
		},
		"Successful CCI unmarshal": {
			ByteArray: []byte(`{"name":"cci","source":{"name":"sma","length":1,"offset":3}}`),
			Result:    CCI{source: SMA{length: 1, offset: 3, valid: true}, factor: decimal.RequireFromString("0.015"), valid: true},
		},
		"Successful DEMA unmarshal": {
			ByteArray: []byte(`{"name":"dema","ema":{"length":1,"offset":1}}`),
			Result:    DEMA{ema: EMA{SMA{length: 1, offset: 1, valid: true}}, valid: true},
		},
		"Successful EMA unmarshal": {
			ByteArray: []byte(`{"name":"ema","length":1,"offset":3}`),
			Result:    EMA{SMA{length: 1, offset: 3, valid: true}},
		},
		"Successful HMA unmarshal": {
			ByteArray: []byte(`{"name":"hma", "wma":{"name":"wma","length":2, "offset":3}}`),
			Result:    HMA{wma: WMA{length: 2, offset: 3, valid: true}, valid: true},
		},
		"Successful CD unmarshal": {
			ByteArray: []byte(`{"name":"cd",
			"source1":{"name":"sma","length":2,"offset":2},
			"source2":{"name":"sma","length":3,"offset":4},
			"offset":3}`),
			Result: CD{percent: false, source1: SMA{length: 2, offset: 2, valid: true},
				source2: SMA{length: 3, offset: 4, valid: true}, offset: 3, valid: true},
		},
		"Successful ROC unmarshal": {
			ByteArray: []byte(`{"name":"roc","length":1,"offset":3}`),
			Result:    ROC{length: 1, offset: 3, valid: true},
		},
		"Successful RSI unmarshal": {
			ByteArray: []byte(`{"name":"rsi","length":1,"offset":2}`),
			Result:    RSI{length: 1, offset: 2, valid: true},
		},
		"Successful SMA unmarshal": {
			ByteArray: []byte(`{"name":"sma","length":1,"offset":3}`),
			Result:    SMA{length: 1, offset: 3, valid: true},
		},
		"Successful SRSI unmarshal": {
			ByteArray: []byte(`{"name":"srsi", "rsi":{"name":"rsi","length":1,"offset":1}}`),
			Result:    SRSI{rsi: RSI{length: 1, offset: 1, valid: true}, valid: true},
		},
		"Successful Stoch unmarshal": {
			ByteArray: []byte(`{"name":"stoch","length":1,"offset":4}`),
			Result:    Stoch{length: 1, offset: 4, valid: true},
		},
		"Successful WMA unmarshal": {
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
