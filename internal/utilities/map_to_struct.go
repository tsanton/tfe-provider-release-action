package utilities

import "encoding/json"

func MapToStruct[T any](m map[string]interface{}) (T, error) {
	var ret T

	// Convert map to JSON
	jsonStr, err := json.Marshal(m)
	if err != nil {
		return ret, err
	}

	// Unmarshal JSON into struct
	err = json.Unmarshal(jsonStr, &ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}
