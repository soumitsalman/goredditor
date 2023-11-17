package utils

import (
	"encoding/json"
	"log"
	"os"
)

func WriteDataToJsonFile[T any](data *T, outFile string) error {
	if j_bytes, err := json.Marshal(data); err != nil {
		return err
	} else {
		return os.WriteFile(outFile, j_bytes, os.FileMode(0644))
	}
}

// loading application configuration. In future making this retrieve from a DB
func ReadDataFromJsonFile[T any](configFilePath string) (T, error) {
	if configFile, err := os.Open(configFilePath); err != nil {
		log.Printf("Failed loading configuration file %v. Error: %v\n", configFilePath, err)
		var data T
		return data, err
	} else {
		defer configFile.Close()
		return DeserialzeJson[T](configFile)
	}
}

func SaveToFile(userId string, topic string, data *[]map[string]any) {
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
