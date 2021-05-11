package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func returnCollection(ex string) interface{} {
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
	case "4":
		object1 := map[string]interface{}{
			"key": "value",
			"other_key": map[string]interface{}{
				"other_value":   "stuff",
				"desired_value": "here",
			},
		}
		collection = map[string]interface{}{
			"one": "one",
			"two": object1,
		}

	}

	returnArray := false
	for _, v := range []string{"1", "2", "3"} {
		if ex == v {
			returnArray = true
		}
	}
	if returnArray {
		mapStrings := []interface{}{}
		for _, v := range collection {
			mapStrings = append(mapStrings, v)
		}
		return mapStrings
	}
	return collection
}

func examples(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	collection := returnCollection(parts[len(parts)-1])

	b, _ := json.Marshal(collection)

	w.Write(b)
}
