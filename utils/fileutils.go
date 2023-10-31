package utils

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func DeserialzeJsonBlob[T any](reader io.Reader) (T, error) {
	var decoder = json.NewDecoder(reader)
	var data T
	if err := decoder.Decode(&data); err != nil {
		log.Printf("Error deserializing to data of type %T: %v\n", data, err)
		return data, err
	} else {
		return data, nil
	}
}

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
		return DeserialzeJsonBlob[T](configFile)
	}
}
