package utils

import "encoding/json"

func ToMap(data interface{}) map[string]interface{} {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	var result map[string]interface{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil
	}
	return result
}
