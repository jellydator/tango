package indc

import (
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_NewAroon(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result Aroon
		Error  error
	}{
		"Validate returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new Aroon": {
			Length: 5,
			Result: Aroon{
				valid:  true,
				length: 5,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewAroon(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_Aroon_validate(t *testing.T) {
	cc := map[string]struct {
		Aroon Aroon
		Error error
	}{
		"Invalid length": {
			Aroon: Aroon{},
			Error: ErrInvalidLength,
		},
		"Successfully validated": {
			Aroon: Aroon{
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.Aroon.validate())
			if c.Error == nil {
				assert.True(t, c.Aroon.valid)
			}
		})
	}
}

func Test_Aroon_Calc(t *testing.T) {
	cc := map[string]struct {
		Aroon      Aroon
		Data       []decimal.Decimal
		UpResult   decimal.Decimal
		DownResult decimal.Decimal
		Error      error
	}{
		"Invalid indicator": {
			Aroon: Aroon{
				valid: false,
			},
			Error: ErrInvalidIndicator,
		},
		"Invalid data size": {
			Aroon: Aroon{
				valid:  true,
				length: 5,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
			Aroon: Aroon{
				valid:  true,
				length: 5,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(31),
				decimal.NewFromInt(38),
				decimal.NewFromInt(35),
				decimal.NewFromInt(29),
				decimal.NewFromInt(29),
			},
			UpResult:   decimal.NewFromInt(40),
			DownResult: _hundred,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			uptrend, downtrend, err := c.Aroon.Calc(c.Data)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.UpResult.String(), uptrend.String())
			assert.Equal(t, c.DownResult.String(), downtrend.String())
		})
	}
}

func Test_Aroon_CalcTrend(t *testing.T) {
	cc := map[string]struct {
		Aroon  Aroon
		Trend  Trend
		Data   []decimal.Decimal
		Result decimal.Decimal
		Error  error
	}{
		"Invalid indicator": {
			Aroon: Aroon{
				valid: false,
			},
			Trend: TrendDown,
			Error: ErrInvalidIndicator,
		},
		"Invalid trend": {
			Aroon: Aroon{
				valid:  true,
				length: 1,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidTrend,
		},
		"Invalid data size": {
			Aroon: Aroon{
				valid:  true,
				length: 5,
			},
			Trend: TrendDown,
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation with TrendUp": {
			Aroon: Aroon{
				valid:  true,
				length: 5,
			},
			Trend: TrendUp,
			Data: []decimal.Decimal{
				decimal.NewFromInt(31),
				decimal.NewFromInt(38),
				decimal.NewFromInt(35),
				decimal.NewFromInt(29),
				decimal.NewFromInt(29),
			},
			Result: decimal.NewFromInt(40),
		},
		"Successful calculation with TrendDown": {
			Aroon: Aroon{
				valid:  true,
				length: 5,
			},
			Trend: TrendDown,
			Data: []decimal.Decimal{
				decimal.NewFromInt(31),
				decimal.NewFromInt(38),
				decimal.NewFromInt(35),
				decimal.NewFromInt(29),
				decimal.NewFromInt(29),
			},
			Result: _hundred,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.Aroon.CalcTrend(c.Data, c.Trend)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_Aroon_Count(t *testing.T) {
	assert.Equal(t, 5, Aroon{
		length: 5,
	}.Count())
}

func Test_NewCCI(t *testing.T) {
	cc := map[string]struct {
		Type   MAType
		Length int
		Result CCI
		Error  error
	}{
		"NewSMA returns an error": {
			Error: assert.AnError,
		},
		"Invalid provided moving average type": {
			Length: 1,
			Error:  errors.New("invalid moving average"),
		},
		"Successfully created new CCI with default factor": {
			Type:   MATypeSimple,
			Length: 10,
			Result: CCI{
				valid: true,
				ma: SMA{
					length: 10,
					valid:  true,
				},
			},
		},
		"Successfully created new CCI": {
			Type:   MATypeSimple,
			Length: 10,
			Result: CCI{
				valid: true,
				ma: SMA{
					length: 10,
					valid:  true,
				},
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewCCI(c.Type, c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_CCI_Calc(t *testing.T) {
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
			CCI: CCI{
				valid: true,
				ma: SMA{
					length: 31,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Invalid SMA calc": {
			CCI: CCI{
				valid: true,
				ma:    SMA{},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: assert.AnError,
		},
		"Successful handled division by 0": {
			CCI: CCI{
				valid: true,
				ma: SMA{
					length: 1,
					valid:  true,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(3),
			},
			Result: decimal.Zero,
		},
		"Successful calculation": {
			CCI: CCI{
				valid: true,
				ma: SMA{
					length: 3,
					valid:  true,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(3),
				decimal.NewFromInt(6),
				decimal.NewFromInt(9),
			},
			Result: decimal.NewFromInt(100),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.CCI.Calc(c.Data)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_CCI_Count(t *testing.T) {
	assert.Equal(t, 10, CCI{
		ma: SMA{
			length: 10,
		},
	}.Count())
}

func Test_NewROC(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result ROC
		Error  error
	}{
		"Validate returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new ROC": {
			Length: 1,
			Result: ROC{
				valid:  true,
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewROC(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_ROC_validate(t *testing.T) {
	cc := map[string]struct {
		ROC   ROC
		Error error
	}{
		"Invalid length": {
			ROC: ROC{
				length: -1,
			},
			Error: ErrInvalidLength,
		},
		"Successfully validated": {
			ROC: ROC{
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.ROC.validate())
			if c.Error == nil {
				assert.True(t, c.ROC.valid)
			}
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
			ROC: ROC{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation": {
			ROC: ROC{
				valid:  true,
				length: 5,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(7),
				decimal.NewFromInt(16),
				decimal.NewFromInt(24),
				decimal.NewFromInt(16),
				decimal.NewFromInt(10),
			},
			Result: decimal.RequireFromString("-30"),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.ROC.Calc(c.Data)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_ROC_Count(t *testing.T) {
	assert.Equal(t, 15, ROC{
		length: 15,
	}.Count())
}

func Test_NewRSI(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result RSI
		Error  error
	}{
		"Validate returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new RSI": {
			Length: 1,
			Result: RSI{
				valid:  true,
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewRSI(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_RSI_validate(t *testing.T) {
	cc := map[string]struct {
		RSI   RSI
		Error error
	}{
		"Invalid length": {
			RSI: RSI{
				length: 0,
			},
			Error: ErrInvalidLength,
		},
		"Successfully validated": {
			RSI: RSI{
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.RSI.validate())
			if c.Error == nil {
				assert.True(t, c.RSI.valid)
			}
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
			RSI: RSI{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation when average gain 0": {
			RSI: RSI{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(16),
				decimal.NewFromInt(12),
				decimal.NewFromInt(8),
			},
			Result: decimal.NewFromInt(0),
		},
		"Successful calculation when average loss 0": {
			RSI: RSI{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(2),
				decimal.NewFromInt(4),
				decimal.NewFromInt(8),
			},
			Result: _hundred,
		},
		"Successful calculation": {
			RSI: RSI{
				valid:  true,
				length: 3,
			},
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
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_RSI_Count(t *testing.T) {
	assert.Equal(t, 15, RSI{
		length: 15,
	}.Count())
}

func Test_NewStochRSI(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result StochRSI
		Error  error
	}{
		"Validate returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new StochRSI": {
			Length: 1,
			Result: StochRSI{
				valid: true,
				rsi: RSI{
					length: 1,
					valid:  true,
				},
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewStochRSI(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_StochRSI_Calc(t *testing.T) {
	cc := map[string]struct {
		StochRSI StochRSI
		Data     []decimal.Decimal
		Result   decimal.Decimal
		Error    error
	}{
		"Invalid indicator": {
			StochRSI: StochRSI{},
			Error:    ErrInvalidIndicator,
		},
		"Invalid data size": {
			StochRSI: StochRSI{
				valid: true,
				rsi: RSI{
					length: 5,
					valid:  true,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: assert.AnError,
		},
		"Successfully handled division by 0": {
			StochRSI: StochRSI{
				valid: true,
				rsi: RSI{
					length: 3,
					valid:  true,
				},
			},
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
			StochRSI: StochRSI{
				valid: true,
				rsi: RSI{
					length: 3,
					valid:  true,
				},
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(11),
				decimal.NewFromInt(12),
				decimal.NewFromInt(11),
				decimal.NewFromInt(11),
				decimal.NewFromInt(11),
			},
			Result: decimal.RequireFromString("0.5"),
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.StochRSI.Calc(c.Data)
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}

func Test_StochRSI_Count(t *testing.T) {
	assert.Equal(t, 29, StochRSI{
		valid: false,
		rsi: RSI{
			length: 15,
		},
	}.Count())
}

func Test_NewStoch(t *testing.T) {
	cc := map[string]struct {
		Length int
		Result Stoch
		Error  error
	}{
		"Validate returns an error": {
			Error: assert.AnError,
		},
		"Successfully created new Stoch": {
			Length: 1,
			Result: Stoch{
				valid:  true,
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewStoch(c.Length)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_Stoch_validate(t *testing.T) {
	cc := map[string]struct {
		Stoch Stoch
		Error error
	}{
		"Invalid length": {
			Stoch: Stoch{
				length: 0,
			},
			Error: ErrInvalidLength,
		},
		"Successfully validated": {
			Stoch: Stoch{
				length: 1,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Error, c.Stoch.validate())
			if c.Error == nil {
				assert.True(t, c.Stoch.valid)
			}
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
			Stoch: Stoch{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(30),
			},
			Error: ErrInvalidDataSize,
		},
		"Successful calculation when new lows are reached": {
			Stoch: Stoch{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(150),
				decimal.NewFromInt(125),
				decimal.NewFromInt(145),
			},
			Result: decimal.NewFromInt(80),
		},
		"Successfully handled division by 0": {
			Stoch: Stoch{
				valid:  true,
				length: 3,
			},
			Data: []decimal.Decimal{
				decimal.NewFromInt(150),
				decimal.NewFromInt(150),
				decimal.NewFromInt(150),
			},
			Result: decimal.Zero,
		},
		"Successful calculation when new highs are reached": {
			Stoch: Stoch{
				valid:  true,
				length: 3,
			},
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
			assertEqualError(t, c.Error, err)
			if err != nil {
				return
			}

			assert.Equal(t, c.Result.String(), res.String())
		})
	}
}
