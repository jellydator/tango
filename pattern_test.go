package tango

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_InverseHeadAndShoulders_Calc(t *testing.T) {
	cc := map[string]struct {
		InverseHeadAndShoulders InverseHeadAndShoulders
		Values                  []Point
		Result                  [][]Point
	}{
		"Successfully found inverse head & shoulders": {
			Values: []Point{
				{
					Value: decimal.NewFromFloat(30),
				},
				{
					Value: decimal.NewFromFloat(40),
				},
				{
					Value: decimal.NewFromFloat(50),
				},
				{
					Value: decimal.NewFromFloat(40),
				},
				{
					Value: decimal.NewFromFloat(30),
				},
				{
					Value: decimal.NewFromFloat(31), // False positive, delta skips this.
				},
				{
					Value: decimal.NewFromFloat(30),
				},
				{
					Value: decimal.NewFromFloat(29), // False positive, delta skips this.
				},
				{
					Value: decimal.NewFromFloat(60),
				},
				{
					Value: decimal.NewFromFloat(30),
				},
				{
					Value: decimal.NewFromFloat(10),
				},
				{
					Value: decimal.NewFromFloat(30),
				},
				{
					Value: decimal.NewFromFloat(57),
				},
				{
					Value: decimal.NewFromFloat(49),
				},
				{
					Value: decimal.NewFromFloat(28),
				},
				{
					Value: decimal.NewFromFloat(100),
				},
			},
			InverseHeadAndShoulders: InverseHeadAndShoulders{
				PeaksDeltaMultiplier:         decimal.NewFromFloat(0.95),
				ShoulderDifferenceMultiplier: decimal.NewFromFloat(0.1),
				MinNeckHeightMultiplier:      decimal.NewFromFloat(0.95),
			},
			Result: [][]Point{
				{
					{
						Value: decimal.NewFromFloat(29),
					},
					{
						Value: decimal.NewFromFloat(60),
					},
					{
						Value: decimal.NewFromFloat(10),
					},
					{
						Value: decimal.NewFromFloat(57),
					},
					{
						Value: decimal.NewFromFloat(28),
					},
				},
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			res := c.InverseHeadAndShoulders.Calc(c.Values)
			assert.Equal(t, c.Result, res)
		})
	}
}

func Test_findPeaks(t *testing.T) {
	cc := map[string]struct {
		Values          []Point
		DeltaMultiplier decimal.Decimal
		Result          []Point
	}{
		"Successfully found min and max values": {
			Values: []Point{
				{
					Value: decimal.NewFromFloat(300),
				},
				{
					Value: decimal.NewFromFloat(400),
				},
				{
					Value: decimal.NewFromFloat(500),
				},
				{
					Value: decimal.NewFromFloat(400),
				},
				{
					Value: decimal.NewFromFloat(300),
				},
				{
					Value: decimal.NewFromFloat(315), // False positive, delta skips this.
				},
				{
					Value: decimal.NewFromFloat(200),
				},
				{
					Value: decimal.NewFromFloat(300),
				},
				{
					Value: decimal.NewFromFloat(290), // False positive, delta skips this.
				},
				{
					Value: decimal.NewFromFloat(600),
				},
				{
					Value: decimal.NewFromFloat(300),
				},
			},
			DeltaMultiplier: decimal.NewFromFloat(0.95),
			Result: []Point{
				{
					Value: decimal.NewFromFloat(500),
				},
				{
					Value: decimal.NewFromFloat(200),
				},
				{
					Value: decimal.NewFromFloat(600),
				},
			},
		},
	}

	for cn, c := range cc {
		c := c

		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			result := findPeaks(c.Values, c.DeltaMultiplier)
			assert.Equal(t, c.Result, result)
		})
	}
}
