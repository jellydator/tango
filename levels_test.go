package tango

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_PivotPointsType_Validate(t *testing.T) {
	cc := map[string]struct {
		Type PivotPointsType
		Err  error
	}{
		"Invalid PivotPointsType": {
			Err: ErrInvalidPivotPointsType,
		},
		"Successful PivotPointsTypeClassic validation": {
			Type: PivotPointsTypeClassic,
		},
		"Successful PivotPointsTypeFibonacci validation": {
			Type: PivotPointsTypeFibonacci,
		},
		"Successful PivotPointsTypeCamarilla validation": {
			Type: PivotPointsTypeCamarilla,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			assertEqualError(t, c.Err, c.Type.Validate())
		})
	}
}

func Test_PivotPointsType_MarshalText(t *testing.T) {
	cc := map[string]struct {
		Type PivotPointsType
		Text string
		Err  error
	}{
		"Invalid PivotPointsType": {
			Err: ErrInvalidPivotPointsType,
		},
		"Successful PivotPointsTypeClassic marshal": {
			Type: PivotPointsTypeClassic,
			Text: "classic",
		},
		"Successful PivotPointsTypeFibonacci marshal": {
			Type: PivotPointsTypeFibonacci,
			Text: "fibonacci",
		},
		"Successful PivotPointsTypeCamarilla marshal": {
			Type: PivotPointsTypeCamarilla,
			Text: "camarilla",
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.Type.MarshalText()
			assertEqualError(t, c.Err, err)

			if err != nil {
				return
			}

			assert.Equal(t, c.Text, string(res))
		})
	}
}

func Test_PivotPointsType_UnmarshalText(t *testing.T) {
	cc := map[string]struct {
		Text   string
		Result PivotPointsType
		Err    error
	}{
		"Invalid PivotPointsType": {
			Err: ErrInvalidPivotPointsType,
		},
		"Successful PivotPointsTypeClassic unmarshal": {
			Text:   "classic",
			Result: PivotPointsTypeClassic,
		},
		"Successful PivotPointsTypeFibonacci unmarshal": {
			Text:   "fibonacci",
			Result: PivotPointsTypeFibonacci,
		},
		"Successful PivotPointsTypeCamarilla unmarshal": {
			Text:   "camarilla",
			Result: PivotPointsTypeCamarilla,
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			var ppt PivotPointsType
			err := ppt.UnmarshalText([]byte(c.Text))
			assertEqualError(t, c.Err, err)

			if err != nil {
				return
			}

			assert.Equal(t, c.Result, ppt)
		})
	}
}

func Test_NewPivotPoints(t *testing.T) {
	cc := map[string]struct {
		Type   PivotPointsType
		Result PivotPoints
		Error  error
	}{
		"Invalid pivot points type": {
			Error: ErrInvalidPivotPointsType,
		},
		"Successfully created new PivotPoints": {
			Type: PivotPointsTypeClassic,
			Result: PivotPoints{
				valid:   true,
				variant: PivotPointsTypeClassic,
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := NewPivotPoints(c.Type)
			assertEqualError(t, c.Error, err)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_PivotPoints_Calc(t *testing.T) {
	candle := Candle{
		High:  decimal.NewFromInt(15),
		Low:   decimal.NewFromInt(9),
		Close: decimal.NewFromInt(12),
	}

	cc := map[string]struct {
		PivotPoints PivotPoints
		Candle      Candle
		Result      PivotLevels
		Error       error
	}{
		"Invalid indicator": {
			PivotPoints: PivotPoints{
				valid: false,
			},
			Error: ErrInvalidIndicator,
		},
		"Successful classic calculation": {
			PivotPoints: PivotPoints{
				valid:   true,
				variant: PivotPointsTypeClassic,
			},
			Candle: candle,
			Result: PivotLevels{
				Pivot: decimal.NewFromInt(12),
				R1:    decimal.NewFromInt(15),
				R2:    decimal.NewFromInt(18),
				R3:    decimal.NewFromInt(21),
				S1:    decimal.NewFromInt(9),
				S2:    decimal.NewFromInt(6),
				S3:    decimal.NewFromInt(3),
			},
		},
		"Successful fibonacci calculation": {
			PivotPoints: PivotPoints{
				valid:   true,
				variant: PivotPointsTypeFibonacci,
			},
			Candle: candle,
			Result: PivotLevels{
				Pivot: decimal.NewFromInt(12),
				R1:    decimal.RequireFromString("14.292"),
				R2:    decimal.RequireFromString("15.708"),
				R3:    decimal.NewFromInt(18),
				S1:    decimal.RequireFromString("9.708"),
				S2:    decimal.RequireFromString("8.292"),
				S3:    decimal.NewFromInt(6),
			},
		},
		"Successful camarilla calculation": {
			PivotPoints: PivotPoints{
				valid:   true,
				variant: PivotPointsTypeCamarilla,
			},
			Candle: candle,
			Result: PivotLevels{
				Pivot: decimal.NewFromInt(12),
				R1:    decimal.RequireFromString("12.55"),
				R2:    decimal.RequireFromString("13.1"),
				R3:    decimal.RequireFromString("13.65"),
				S1:    decimal.RequireFromString("11.45"),
				S2:    decimal.RequireFromString("10.9"),
				S3:    decimal.RequireFromString("10.35"),
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res, err := c.PivotPoints.Calc(c.Candle)
			assertEqualError(t, c.Error, err)

			if err != nil {
				return
			}

			assert.Equal(t, c.Result.Pivot.String(), res.Pivot.String())
			assert.Equal(t, c.Result.R1.String(), res.R1.String())
			assert.Equal(t, c.Result.R2.String(), res.R2.String())
			assert.Equal(t, c.Result.R3.String(), res.R3.String())
			assert.Equal(t, c.Result.S1.String(), res.S1.String())
			assert.Equal(t, c.Result.S2.String(), res.S2.String())
			assert.Equal(t, c.Result.S3.String(), res.S3.String())
		})
	}
}

func Test_PivotPoints_Count(t *testing.T) {
	assert.Equal(t, 1, PivotPoints{
		valid:   true,
		variant: PivotPointsTypeClassic,
	}.Count())
}
