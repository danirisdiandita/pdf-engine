package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/danirisdiandita/pdf-engine/internal/config"
)

func Transcribe(config *config.Config, filePath string, lang string) ([]interface{}, string, error) {
	// Create a new file for the request
	file, err := os.Open(filePath)
	if err != nil {
		return nil, lang, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Create a buffer to store our request body as multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add the audio file
	part, err := writer.CreateFormFile("audio", filepath.Base(filePath))
	if err != nil {
		return nil, lang, fmt.Errorf("failed to create form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, lang, fmt.Errorf("failed to copy file: %v", err)
	}

	// Add the language parameter if provided
	if lang != "" {
		err = writer.WriteField("language", lang)
		if err != nil {
			return nil, lang, fmt.Errorf("failed to write form field: %v", err)
		}
	}

	// Close the writer
	err = writer.Close()
	if err != nil {
		return nil, lang, fmt.Errorf("failed to close writer: %v", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", config.DeepInfraUrl, &requestBody)
	if err != nil {
		return nil, lang, fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "bearer "+config.DeepInfraApiKey)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, lang, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, lang, fmt.Errorf("bad status: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse the response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, lang, fmt.Errorf("failed to decode response: %v", err)
	}

	fmt.Println("result", result)

	// Extract the segments
	segments, ok := result["segments"].([]interface{})
	if !ok {
		return nil, lang, fmt.Errorf("transcription segments not found in response")
	}

	if lang == "" {
		langCode, ok := GetLanguageCode(result["language"].(string))
		if !ok {
			return nil, lang, fmt.Errorf("language not found in response")
		}
		lang = langCode
	}

	return segments, lang, nil
}
