package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

var authServiceURL = os.Getenv("AUTH_SERVICE_URL")

func AuthClient(method, endpoint string, data interface{}) (*http.Response, error) {
	url := authServiceURL + endpoint
	body, _ := json.Marshal(data)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	return client.Do(req)
}
