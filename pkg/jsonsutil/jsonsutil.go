package jsonsutil

import "encoding/json"

// ConvertData converts the generic data field to the specific type
func ConvertData[T any](data any) (T, error) {
	var result T

	// Marshal the data back to JSON and unmarshal into the specific type
	bytes, err := json.Marshal(data)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(bytes, &result)
	return result, err
}
