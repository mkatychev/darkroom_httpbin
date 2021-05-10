package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func returnCollection(ex string) map[string]interface{} {
	var collection map[string]interface{}
	switch ex {
	case "1":
		collection = map[string]interface{}{
			"one":   "one",
			"two":   "two",
			"three": "three",
			"four":  "four",
		}
	case "2":
		object := map[string]string{
			"key": "value",
		}
		collection = map[string]interface{}{
			"one":   "one",
			"two":   object,
			"three": "three",
			"four":  "four",
		}
	case "3":
		object1 := map[string]string{
			"key": "value",
		}
		object2 := map[string]string{
			"key": "other_value",
		}
		collection = map[string]interface{}{
			"one":   "one",
			"two":   object1,
			"three": object2,
			"four":  "four",
		}

	}
	return collection
}

func noOrderArray(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	mapStrings := []interface{}{}
	println(parts)

	for _, v := range returnCollection(parts[len(parts)-1]) {
		mapStrings = append(mapStrings, v)
	}

	b, _ := json.Marshal(mapStrings)

	w.Write(b)
}

func unsortedWithObject(w http.ResponseWriter, r *http.Request) {
	object := map[string]string{
		"key": "value",
	}
	collection := map[string]interface{}{
		"one":   "one",
		"two":   object,
		"three": "three",
		"four":  "four",
	}
	mapStrings := []interface{}{}

	for _, v := range collection {
		mapStrings = append(mapStrings, v)
	}

	b, _ := json.Marshal(mapStrings)

	w.Write(b)
}
