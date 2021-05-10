package main

import (
	"encoding/json"
	"net/http"
)

func unsortedStrings(w http.ResponseWriter, r *http.Request) {
	collection := map[string]string{
		"one":   "one",
		"two":   "two",
		"three": "three",
		"four":  "four",
	}
	mapStrings := []string{}

	for _, v := range collection {
		mapStrings = append(mapStrings, v)
	}

	b, _ := json.Marshal(mapStrings)

	w.Write(b)
}
