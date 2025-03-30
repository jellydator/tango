package tango

import (
	"time"

	"github.com/shopspring/decimal"
)

// Point represents a value at a specific time.
type Point struct {
	// Value is the value at the specific time.
	Value decimal.Decimal

	// Timestamp is the time the value was recorded.
	Timestamp time.Time
}

// InverseHeadAndShoulders is a pattern that is used to identify
// a potential reversal in the market. The pattern consists of
// three peaks, with the middle peak being the highest. The
// left and right peaks should be approximately the same height.
// The pattern is considered to be a bullish reversal pattern.
type InverseHeadAndShoulders struct {
	// PeaksDeltaMultiplier is the multiplier used to determine
	// whether a peak is a local maximum or minimum.
	// The default value should be 0.95.
	PeaksDeltaMultiplier decimal.Decimal

	// ShoulderDifferenceMultiplier is the multiplier used to
	// determine whether the left and right shoulders are
	// approximately the same height.
	// The default value should be 0.05.
	ShoulderDifferenceMultiplier decimal.Decimal

	// MinNeckHeightMultiplier is the multiplier used to determine
	// whether the neck is less than the left and right shoulders.
	// The default value should be 0.95.
	MinNeckHeightMultiplier decimal.Decimal
}

// Calc returns the potential inverse head and shoulders patterns
// in the given slice of points. The pattern is identified by
// finding the peaks in the slice. The peaks are then used to
// determine whether the pattern is present.
func (ihas InverseHeadAndShoulders) Calc(pp []Point) [][]Point {
	peaks := findPeaks(pp, ihas.PeaksDeltaMultiplier)
	if len(peaks) < 5 {
		return [][]Point{}
	}

	var res [][]Point

	for i := 0; i < len(peaks)-4; i++ {
		lsStart := peaks[i]
		lsEnd := peaks[i+1]
		neck := peaks[i+2]
		rsEnd := peaks[i+3]
		rsStart := peaks[i+4]

		// NOTE: The left and right shoulders should be above
		// the neck. The starting positions of left and right
		// shoulders should be below the ending positions.
		if !lsStart.Value.LessThan(lsEnd.Value) ||
			!rsStart.Value.LessThan(rsEnd.Value) ||
			!lsEnd.Value.GreaterThan(neck.Value) ||
			!rsEnd.Value.GreaterThan(neck.Value) {

			continue
		}

		averageShoulderHeight := lsEnd.Value.Add(rsEnd.Value).Div(decimal.NewFromInt(2))

		// NOTE: The neck should be less than the left and right shoulders.
		if neck.Value.GreaterThan(
			averageShoulderHeight.Mul(ihas.MinNeckHeightMultiplier),
		) {
			continue
		}

		// NOTE: The shoulders should be approximately the same height.
		// We check that by seeing whether the difference between the left
		// shoulder and the right shoulder is less than 5%.
		if lsEnd.Value.Sub(rsEnd.Value).Abs().GreaterThan(
			averageShoulderHeight.Mul(ihas.ShoulderDifferenceMultiplier),
		) {
			continue
		}

		res = append(res, []Point{lsStart, lsEnd, neck, rsEnd, rsStart})
	}

	return res
}

// findPeaks returns the minimum and maximum values in the slice.
// The peaks are determined by following a trend and smoothing using
// the delta multiplier.
func findPeaks(values []Point, deltaMultiplier decimal.Decimal) []Point {
	var (
		searchMin bool

		result []Point

		minValue = values[0]
		maxValue = values[0]
	)

	for _, val := range values {
		if val.Value.GreaterThan(maxValue.Value) {
			maxValue = val
		}

		if val.Value.LessThan(minValue.Value) {
			minValue = val
		}

		if !searchMin {
			if val.Value.LessThan(maxValue.Value.Mul(deltaMultiplier)) {
				result = append(result, maxValue)

				minValue = val
				searchMin = true
			}

			continue
		}

		if val.Value.GreaterThan(
			minValue.Value.Add(
				minValue.Value.Mul(decimal.NewFromInt(1).Sub(deltaMultiplier)),
			),
		) {
			result = append(result, minValue)

			maxValue = val
			searchMin = false
		}
	}

	return result
}
