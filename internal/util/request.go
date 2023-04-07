package util

import (
	"encoding/json"
	"net/http"
)

func GetJSON(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func GetStatusCode(url string) (int, error) {
	r, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	defer r.Body.Close()

	return r.StatusCode, nil
}
