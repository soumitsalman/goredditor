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
	if configFile, err := os.Open(filepath); err != nil {
		log.Printf("Failed loading configuration file %v. Error: %v\n", filepath, err)
		var data T
		return data, err
	} else {
		defer configFile.Close()
		return DeserialzeJson[T](configFile)
	}
}

func SaveData(userId string, topic string, data any) {
	var content = map[string]any{
		"topic": topic,
		topic:   data,
	}
	var filename = os.Getenv("DATASTORE_LOCATION") + userId + "_" + topic + ".json"
	if WriteDataToJsonFile(&content, filename) == nil {
		log.Printf("Saved %s in %s\n", topic, filename)
	} else {
		log.Printf("Failed to save %s\n", topic)
	}
}

func ReadData(userId string, topic string) ([]map[string]any, error) {
	var filename = os.Getenv("DATASTORE_LOCATION") + userId + "_" + topic + ".json"
	if content, err := ReadDataFromJsonFile[[]map[string]any](filename); err != nil {
		return content, nil
	} else {
		log.Println("Failed loading file ", filename)
		return nil, err
	}
}
