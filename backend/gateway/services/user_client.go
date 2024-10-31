package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

var userServiceURL = os.Getenv("USER_SERVICE_URL")

func UserClient(method, endpoint string, data interface{}) (*http.Response, error) {
	url := userServiceURL + endpoint
	var body []byte
	if data != nil {
		body, _ = json.Marshal(data)
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	return client.Do(req)
}
