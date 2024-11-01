package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func CourseClient(method, path string, body io.Reader) (*http.Response, error) {
	// Get course service URL from environment variable
	courseServiceURL := os.Getenv("COURSE_SERVICE_URL")
	if courseServiceURL == "" {
		courseServiceURL = "http://course-service:8082" // Default URL
	}

	// Create new request
	req, err := http.NewRequest(method, courseServiceURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Create HTTP client and send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	return resp, nil
} 