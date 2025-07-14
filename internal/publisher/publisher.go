package publisher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type patchPayload struct {
	Body []struct {
		Value  string `json:"value"`
		Format string `json:"format"`
	} `json:"body"`
	Type []struct {
		TargetID string `json:"target_id"`
	} `json:"type"`
}

// Publish sends the generated HTML to the configured web page.
func Publish(html, url, user, pass string) error {
	payload := patchPayload{
		Body: []struct {
			Value  string `json:"value"`
			Format string `json:"format"`
		}{
			{Value: html, Format: "2"}, // Assuming "2" is "Full HTML"
		},
		Type: []struct {
			TargetID string `json:"target_id"`
		}{
			{TargetID: "page"},
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s?_format=json", url), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.SetBasicAuth(user, pass)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send PATCH request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	return nil
}
