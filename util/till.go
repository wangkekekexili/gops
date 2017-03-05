package util

// till is SMS service.

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type tillRequest struct {
	Phone []string `json:"phone"`
	Text  string   `json:"text"`
}

// SendWarningSMS returns true if it sends a warning SMS successfully.
func SendWarningSMS() bool {
	// Get till url and target phone number.
	tillURL := os.Getenv("TILL_URL")
	phone := os.Getenv("TILL_TARGET")
	if tillURL == "" || phone == "" {
		return false
	}

	// Build a request and post it.
	req := &tillRequest{
		Phone: []string{phone},
		Text:  "failed to ping database - computer is offline?",
	}
	reqBytes, _ := json.Marshal(req)
	reader := bytes.NewReader(reqBytes)
	if _, err := http.Post(tillURL, "application/json", reader); err != nil {
		return false
	}

	return true
}
