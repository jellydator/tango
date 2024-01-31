## Tango - (Trading Analysis & Go)

[![Go Reference](https://pkg.go.dev/badge/github.com/jellydator/tango.svg)](https://pkg.go.dev/github.com/jellydator/tango)
[![Build Status](https://github.com/jellydator/tango/actions/workflows/go.yml/badge.svg)](https://github.com/jellydator/tango/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jellydator/tango)](https://goreportcard.com/report/github.com/jellydator/tango)

## Features
- Simple API
- Built-in parameters validation
- Includes thorough documentation
- A wide variety of [Oscillators](#Oscillators) and [Overlays](#Overlays).

## Installation
```
go get github.com/jellydator/tango
```

## Usage
All of the tools must be created using `New*` function. It performs
parameters validation and returns an object that is capable of working
with data slices.

The main calculations are done using `Calc` method. The return types varies
based on the tool.

A simple use case could look like this:
```go
func main() {
  sma, err := tango.NewSMA(3)
  if err != nil {
    // handle the error.
  }

  dataPoints := []decimal.Decimal{
    decimal.NewFromInt(2),
    decimal.NewFromInt(3),
    decimal.NewFromInt(4),
  }

  // the value is 3
  value, err := sma.Calc(dataPoints)
  if err != nil {
    // handle the error.
  }
}
```

For the calculation to be successful, the `Calc` method should receive only the
information that it requires. In some scenarios, it might not be known how many
data points is needed, for this, a `Count` method may be used.

```go
func CalculateMA(ma tango.MA, values []decimal.Decimal) (decimal.Decimal, error) {
  requiredValues := ma.Count()

  if len(values) < requiredValues {
    return decimal.Zero, errors.New("invalid count of values")
  }

  return ma.Calc(values[:requiredValues])
}
```

## Oscillators
- [Aroon](https://www.investopedia.com/terms/a/aroon.asp)
- [CCI (Commodity Channel Index)](https://www.investopedia.com/terms/c/commoditychannelindex.asp)
- [ROC (Rate of Change)](https://www.investopedia.com/terms/p/pricerateofchange.asp)
- [RSI (Relative Strength Index)](https://www.investopedia.com/terms/r/rsi.asp)
- [StochRSI (Stochastic Relative Strength Index)](https://www.investopedia.com/terms/s/stochrsi.asp)
- [Stoch (Stochastic)](https://www.investopedia.com/terms/s/stochasticoscillator.asp)

## Overlays
- [BB (Bollinger Bands)](https://www.investopedia.com/terms/b/bollingerbands.asp)
- [DEMA (Double Exponential Moving Average)](https://www.investopedia.com/terms/d/double-exponential-moving-average.asp)
- [EMA (Exponential Moving Average)](https://www.investopedia.com/terms/e/ema.asp)
- [HMA (Hull Moving Average)](https://www.fidelity.com/learning-center/trading-investing/technical-analysis/technical-indicator-guide/hull-moving-average)
- [SMA (Simple Moving Average)](https://www.investopedia.com/terms/s/sma.asp)
- [WMA (Weighted Moving Average)](https://www.investopedia.com/articles/technical/060401.asp)
