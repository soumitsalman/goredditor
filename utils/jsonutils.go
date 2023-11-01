package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
)

func DeserialzeJson[T any](reader io.Reader) (T, error) {
	var decoder = json.NewDecoder(reader)
	var data T
	if err := decoder.Decode(&data); err != nil {
		log.Printf("Error deserializing to data of type %T: %v\n", data, err)
		return data, err
	} else {
		//log.Printf("Deserialized data: %v\n", data)
		return data, nil
	}
}

func SerializeJson(data map[string]any) (io.Reader, error) {
	if jsonData, err := json.Marshal(data); err != nil {
		return nil, err
	} else {
		return bytes.NewBuffer(jsonData), nil
	}
}
