package indc

import "encoding/json"

// fromJSON finds a source indicator by name and parses its data from json.
// Should be used in places where wrapped unknown indicators are parsed.
//nolint:gocognit,gocyclo // many switch cases are needed to cover all of
// the indicators.
func fromJSON(data []byte) (Indicator, error) {
	var id struct {
		Name String `json:"name"`
	}

	if err := json.Unmarshal(data, &id); err != nil {
		return nil, err
	}

	switch id.Name {
	case NameAroon:
		var v Aroon
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}

		return v, nil
	case NameBB:
		var v BB
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}

		return v, nil
	case NameCCI:
		var v CCI
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}

		return v, nil
	case NameDEMA:
		var v DEMA
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}

		return v, nil
	case NameEMA:
		var v EMA
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}

		return v, nil
	case NameHMA:
		var v HMA
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}

		return v, nil
	case NameCD:
		var v CD
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}

		return v, nil
	case NameROC:
		var v ROC
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}

		return v, nil
	case NameRSI:
		var v RSI
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}

		return v, nil
	case NameSMA:
		var v SMA
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}

		return v, nil
	case NameSRSI:
		var v SRSI
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}

		return v, nil
	case NameStoch:
		var v Stoch
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}

		return v, nil
	case NameWMA:
		var v WMA
		if err := json.Unmarshal(data, &v); err != nil {
			return nil, err
		}

		return v, nil
	}

	return nil, ErrInvalidSource
}
