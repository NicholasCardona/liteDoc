package db

import (
	"encoding/gob"
	"fmt"
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

func Find(query map[string]interface{}) []map[string]interface{} {
	var results []map[string]interface{}
	for _, v := range Store {
		doc, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		if MatchDocument(doc, query) {
			results = append(results, doc)
		}
	}
	return results
}

func MatchDocument(doc map[string]interface{}, query map[string]interface{}) bool {
	for key, val := range query {
		docVal, ok := doc[key]
		if !ok {
			return false
		}
		switch t := val.(type) {
		case map[string]interface{}: // operators like $lte, $gte
			for op, v := range t {
				switch op {
				case "$eq":
					if !compare(docVal, v, "==") {
						return false
					}
				case "$lte":
					if !compare(docVal, v, "<=") {
						return false
					}
				case "$gte":
					if !compare(docVal, v, ">=") {
						return false
					}
				}
			}
		default:
			if docVal != val {
				return false
			}
		}
	}
	return true
}

func compare(a, b interface{}, op string) bool {
	af, aok := toFloat(a)
	bf, bok := toFloat(b)
	if !(aok && bok) {
		return false
	}

	switch op {
	case "==":
		return a == b
	case "<=":
		return af <= bf
	case ">=":
		return af >= bf
	}

	return false
}

func toFloat(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case float64:
		return v, true
	case string:
		var f float64
		n, err := fmt.Sscanf(v, "%f", &f)
		if n == 1 && err == nil {
			return f, true
		}
	}
	return 0, false
}
