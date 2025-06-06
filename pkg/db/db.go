package db

import (
	"encoding/gob"
	"os"
)

var Store = make(map[string]interface{})

func Init() {
	// register base types for json
	gob.Register(map[string]interface{}{})
	gob.Register(map[string]map[interface{}]interface{}{})
	gob.Register([]any{})
}

func SaveData(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(Store)
}

func LoadData(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	return decoder.Decode(&Store)
}
