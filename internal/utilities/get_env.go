package utilities

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func GetEnv[T any](key string, fallback T) T {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	var ret T
	switch any(fallback).(type) {
	case string:
		hack := fmt.Sprintf(`{"my_string": "%s"}`, value)
		var result struct {
			MyString T `json:"my_string"`
		}
		err := json.Unmarshal([]byte(hack), &result)
		if err != nil {
			log.Panicf("unable to unmarshal value for key %s. Error: %s", key, err.Error())
		}
		return result.MyString
	default:
		if err := json.Unmarshal([]byte(value), &ret); err != nil {
			log.Panicf("unable to unmarshal value for key %s. Error: %s", key, err.Error())
		}
	}
	return ret
}
