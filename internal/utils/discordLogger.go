package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danirisdiandita/pdf-engine/internal/config"
)

func DiscordLogger(message string) bool {
	payload := map[string]interface{}{
		"content": nil,
		"embeds": []map[string]interface{}{
			{
				"title":       "Message Logger",
				"description": fmt.Sprintf("**Message:** ```%s```", message),
				"color":       5814783,
			},
		},
		"attachments": []interface{}{},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return false
	}

	cfg := config.Load()
	resp, err := http.Post(
		cfg.DiscordWebhookURL,
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		fmt.Printf("Error sending message to Discord: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Printf("Failed to send message to Discord: %d\n", resp.StatusCode)
		return false
	}

	return true
}
