package test

import (
	"encoding/json"
	"os"
	"testing"

	"jsonCache/pkg/db"
)

var test_path string = "test_data.gob"

func init() {
	db.Init()
}

func TestStoreAndRetrieve(t *testing.T) {
	// Clear the store
	db.Store = make(map[string]interface{})

	doc := map[string]interface{}{"name": "Ben", "age": 44}
	id := "test-id"

	// insert and save
	db.Store[id] = doc
	err := db.SaveData(test_path)
	if err != nil {
		t.Fatalf("Error saving: %v", err)
	}

	// Clear the store
	db.Store = make(map[string]interface{})
	// Reload data
	err = db.LoadData(test_path)
	if err != nil {
		t.Fatalf("Error loading: %v", err)
	}

	// check if data is there
	loaded, ok := db.Store[id]
	if !ok {
		t.Fatal("Doc not found after load")
	}

	// data check
	orig, _ := json.Marshal(doc)
	cmp, _ := json.Marshal(loaded)
	if string(orig) != string(cmp) {
		t.Errorf("Loaded Document differs. Expected: %s, Obtained: %s", orig, cmp)
	}

	// cleanup
	os.Remove(test_path)
}

func TestDeleteDocument(t *testing.T) {
	// load into file new element
	db.Store = make(map[string]any)

	doc := map[string]any{"name": "Clare", "age": 24}
	id := "delete-id"

	db.Store[id] = doc
	db.SaveData(test_path)

	// retrieve from file and check if is still there
	db.Store = make(map[string]any)
	db.LoadData(test_path)
	delete(db.Store, id)
	if _, found := db.Store[id]; found {
		t.Fatal("Document not deleted")
	}

	os.Remove(test_path)
}

func TestInsertNastierObj(t *testing.T) {
	// init store
	db.Store = map[string]interface{}{}

	// nasty element
	doc := map[string]interface{}{
		"name":    "frank",
		"surname": "white",
		"age":     41,
		"items": []interface{}{
			"music", "sport",
		},
		"workRevenues": map[string]interface{}{
			"freelance": 800,
			"9-5":       2000,
		},
	}
	id := "nasty-id"

	// storing on file
	db.Store[id] = doc
	err := db.SaveData(test_path)
	if err != nil {
		t.Fatalf("Error saving: %v", err)
	}

	// clean store
	db.Store = map[string]interface{}{}

	// load from file
	err = db.LoadData(test_path)
	if err != nil {
		t.Fatalf("Error retrieving: %v", err)
	}

	loaded := db.Store[id]

	// data check
	orig, _ := json.Marshal(doc)
	cmp, _ := json.Marshal(loaded)
	if string(orig) != string(cmp) {
		t.Errorf("Loaded Document differs. Expected: %s, Obtained: %s", orig, cmp)
	}

	// cleanup
	os.Remove(test_path)
}

func TestFindWhereEq(t *testing.T) {
	db.Store = map[string]interface{}{
		"1": map[string]interface{}{"name": "Alice", "age": 30},
		"2": map[string]interface{}{"name": "Bob", "age": 25},
		"3": map[string]interface{}{"name": "Charlie", "age": 30},
	}
	query := map[string]interface{}{"name": "Alice"}
	results := db.Find(query)
	if len(results) != 1 || results[0]["name"] != "Alice" {
		t.Errorf("Expected 1 result with name 'Alice', got %v results", results)
	}
}

func TestFindWhereLte(t *testing.T) {
	db.Store = map[string]interface{}{
		"1": map[string]interface{}{"name": "Alice", "age": 30},
		"2": map[string]interface{}{"name": "Bob", "age": 25},
		"3": map[string]interface{}{"name": "Charlie", "age": 20},
	}
	query := map[string]interface{}{"age": map[string]interface{}{"$lte": 25}}
	results := db.Find(query)
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	if results[0]["age"] != 25 && results[1]["age"] != 20 {
		t.Errorf("Expected ages 25 and 20, got %v", results)
	}
}

func TestFindWhereGte(t *testing.T) {
	db.Store = map[string]interface{}{
		"1": map[string]interface{}{"name": "Alice", "age": 30},
		"2": map[string]interface{}{"name": "Bob", "age": 25},
		"3": map[string]interface{}{"name": "Charlie", "age": 20},
	}
	query := map[string]interface{}{"age": map[string]interface{}{"$gte": 25}}
	results := db.Find(query)
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	if (results[0]["age"] != 30 || results[1]["age"] != 25) && (results[0]["age"] != 25 || results[1]["age"] != 30) {
		t.Errorf("Expected ages 30 and 25, got %v", results)
	}
}

func TestFindWithMultipleConditions(t *testing.T) {
	db.Store = map[string]interface{}{
		"1": map[string]interface{}{"name": "Alice", "age": 30},
		"2": map[string]interface{}{"name": "Bob", "age": 25},
		"3": map[string]interface{}{"name": "Charlie", "age": 30},
	}
	query := map[string]interface{}{
		"name": "Alice",
		"age":  map[string]interface{}{"$gte": 20},
	}
	results := db.Find(query)
	if len(results) != 1 || results[0]["name"] != "Alice" {
		t.Errorf("Expected 1 result with name 'Alice', got %v results", results)
	}
}
