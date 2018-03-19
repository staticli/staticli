package lib

import (
	"encoding/json"
	"net/http"
	"time"
)

// GetJson does a HTTP request to a specified URL, then applies data to your interface
func GetJson(url string, target interface{}) error {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}