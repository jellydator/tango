package tango

import (
	"errors"

	"github.com/shopspring/decimal"
)

// ErrInvalidPivotPointsType is returned when the provided pivot points
// type is not valid.
var ErrInvalidPivotPointsType = errors.New("invalid pivot points type")

// PivotPointsType is a custom type that validates it to be only of
// existing pivot points calculation variants.
type PivotPointsType int

// Available pivot points calculation variants.
const (
	PivotPointsTypeClassic PivotPointsType = iota + 1
	PivotPointsTypeFibonacci
	PivotPointsTypeCamarilla
)

// Validate checks whether the pivot points type is one of the
// supported calculation variants.
func (ppt PivotPointsType) Validate() error {
	switch ppt {
	case PivotPointsTypeClassic, PivotPointsTypeFibonacci, PivotPointsTypeCamarilla:
		return nil
	default:
		return ErrInvalidPivotPointsType
	}
}

// MarshalText turns the pivot points type into appropriate string
// representation.
func (ppt PivotPointsType) MarshalText() ([]byte, error) {
	var v string

	switch ppt {
	case PivotPointsTypeClassic:
		v = "classic"
	case PivotPointsTypeFibonacci:
		v = "fibonacci"
	case PivotPointsTypeCamarilla:
		v = "camarilla"
	default:
		return nil, ErrInvalidPivotPointsType
	}

	return []byte(v), nil
}

// UnmarshalText turns string to appropriate pivot points type value.
func (ppt *PivotPointsType) UnmarshalText(d []byte) error {
	switch string(d) {
	case "classic":
		*ppt = PivotPointsTypeClassic
	case "fibonacci":
		*ppt = PivotPointsTypeFibonacci
	case "camarilla":
		*ppt = PivotPointsTypeCamarilla
	default:
		return ErrInvalidPivotPointsType
	}

	return nil
}

// PivotLevels holds all the price levels produced by a pivot points
// calculation.
type PivotLevels struct {
	// Pivot specifies the central pivot point.
	Pivot decimal.Decimal

	// R1 specifies the first resistance level.
	R1 decimal.Decimal

	// R2 specifies the second resistance level.
	R2 decimal.Decimal

	// R3 specifies the third resistance level.
	R3 decimal.Decimal

	// S1 specifies the first support level.
	S1 decimal.Decimal

	// S2 specifies the second support level.
	S2 decimal.Decimal

	// S3 specifies the third support level.
	S3 decimal.Decimal
}

// PivotPoints holds all the necessary information needed to calculate
// pivot points.
// The zero value is not usable.
type PivotPoints struct {
	// valid specifies whether PivotPoints parameters were validated.
	valid bool

	// variant specifies the calculation variant.
	variant PivotPointsType
}

// NewPivotPoints validates provided configuration options and
// creates new PivotPoints indicator instance.
func NewPivotPoints(ppt PivotPointsType) (PivotPoints, error) {
	if err := ppt.Validate(); err != nil {
		return PivotPoints{}, err
	}

	return PivotPoints{
		valid:   true,
		variant: ppt,
	}, nil
}

// Calc calculates all pivot point levels from the provided candle of
// the preceding period.
// Calculation is based on formulas provided by investopedia.
// https://www.investopedia.com/terms/p/pivotpoint.asp.
func (pp PivotPoints) Calc(c Candle) (PivotLevels, error) {
	if !pp.valid {
		return PivotLevels{}, ErrInvalidIndicator
	}

	switch pp.variant {
	case PivotPointsTypeFibonacci:
		return calcFibonacciPivotPoints(c), nil
	case PivotPointsTypeCamarilla:
		return calcCamarillaPivotPoints(c), nil
	default: // PivotPoints is validated, only the classic variant is left.
		return calcClassicPivotPoints(c), nil
	}
}

// Count determines the total amount of data points needed for pivot
// points calculation.
func (pp PivotPoints) Count() int {
	return 1
}

// calcClassicPivotPoints calculates the classic pivot point levels
// from the provided candle.
func calcClassicPivotPoints(c Candle) PivotLevels {
	pivot := typicalPrice(c)
	spread := c.High.Sub(c.Low)
	two := decimal.NewFromInt(2)

	return PivotLevels{
		Pivot: pivot,
		R1:    pivot.Mul(two).Sub(c.Low),
		R2:    pivot.Add(spread),
		R3:    c.High.Add(pivot.Sub(c.Low).Mul(two)),
		S1:    pivot.Mul(two).Sub(c.High),
		S2:    pivot.Sub(spread),
		S3:    c.Low.Sub(c.High.Sub(pivot).Mul(two)),
	}
}

// calcFibonacciPivotPoints calculates the Fibonacci pivot point levels
// from the provided candle.
func calcFibonacciPivotPoints(c Candle) PivotLevels {
	pivot := typicalPrice(c)
	spread := c.High.Sub(c.Low)

	shallow := spread.Mul(decimal.NewFromFloat(0.382))
	golden := spread.Mul(decimal.NewFromFloat(0.618))

	return PivotLevels{
		Pivot: pivot,
		R1:    pivot.Add(shallow),
		R2:    pivot.Add(golden),
		R3:    pivot.Add(spread),
		S1:    pivot.Sub(shallow),
		S2:    pivot.Sub(golden),
		S3:    pivot.Sub(spread),
	}
}

// calcCamarillaPivotPoints calculates the Camarilla pivot point levels
// from the provided candle.
func calcCamarillaPivotPoints(c Candle) PivotLevels {
	scaled := c.High.Sub(c.Low).Mul(decimal.NewFromFloat(1.1))

	tight := scaled.Div(decimal.NewFromInt(12))
	medium := scaled.Div(decimal.NewFromInt(6))
	wide := scaled.Div(decimal.NewFromInt(4))

	return PivotLevels{
		Pivot: typicalPrice(c),
		R1:    c.Close.Add(tight),
		R2:    c.Close.Add(medium),
		R3:    c.Close.Add(wide),
		S1:    c.Close.Sub(tight),
		S2:    c.Close.Sub(medium),
		S3:    c.Close.Sub(wide),
	}
}
