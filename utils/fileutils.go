package utils

import (
	"encoding/json"
	"log"
	"os"
)

func WriteDataToJsonFile[T any](data *T, filepath string) error {
	if j_bytes, err := json.Marshal(data); err != nil {
		return err
	} else {
		return os.WriteFile(filepath, j_bytes, os.FileMode(0644))
	}
}

// loading application configuration. In future making this retrieve from a DB
func ReadDataFromJsonFile[T any](filepath string) (T, error) {
	if file, err := os.Open(filepath); err != nil {
		log.Printf("Failed loading file %v. Error: %v\n", filepath, err)
		var data T
		return data, err
	} else {
		defer file.Close()
		return DeserialzeJson[T](file)
	}
}
