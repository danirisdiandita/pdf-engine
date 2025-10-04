package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danirisdiandita/pdf-engine/internal/config"
	"github.com/danirisdiandita/pdf-engine/internal/model"
)

func UpdateNoteStatus(config *config.Config, req model.TranscribeRequest, status string) error {

	var url string
	if req.WebhookURL != "" {
		url = req.WebhookURL + "/api/uploader/update-note-status"
	} else {
		url = config.PlatformUrl + "/api/uploader/update-note-status"
	}

	payload := map[string]interface{}{
		"note_id":  req.NoteID,
		"user_id":  req.UserID,
		"status":   status,
		"language": req.Language,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+config.PlatformApiKey)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer response.Body.Close()

	return nil
}

func UpdateNoteTranscriptionCompleted(config *config.Config, req model.TranscribeRequest, transcription []TranscriptionOutput, language string, processingTime int) error {
	var url string
	if req.WebhookURL != "" {
		url = req.WebhookURL + "/api/uploader/update-note-status"
	} else {
		url = config.PlatformUrl + "/api/uploader/update-note-status"
	}

	transcriptJSON, err := json.Marshal(transcription)
	if err != nil {
		return fmt.Errorf("failed to marshal transcription: %v", err)
	}

	payload := map[string]interface{}{
		"note_id":  req.NoteID,
		"user_id":  req.UserID,
		"status":   "TRANSCRIPTION_COMPLETED",
		"language": language,
		"transcription_processing_time_in_seconds": processingTime,
		"transcript":       string(transcriptJSON),
		"progress_percent": 100,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+config.PlatformApiKey)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer response.Body.Close()

	return nil
}

func UpdateNoteGdriveStatus(config *config.Config, req model.TranscribeGdriveRequest, status string) error {
	var url string
	if req.WebhookURL != "" {
		url = req.WebhookURL + "/api/uploader/update-note-status"
	} else {
		url = config.PlatformUrl + "/api/uploader/update-note-status"
	}

	payload := map[string]interface{}{
		"note_id":  req.NoteID,
		"user_id":  req.UserID,
		"status":   status,
		"language": req.Language,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+config.PlatformApiKey)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer response.Body.Close()

	return nil
}

func UpdateNoteTranscriptionGdriveCompleted(config *config.Config, req model.TranscribeGdriveRequest, transcription []TranscriptionOutput, language string, processingTime int, fileKey string, fileType string) error {

	if fileType == "audio" {
		var url string
		if req.WebhookURL != "" {
			url = req.WebhookURL + "/api/uploader/update-note-status"
		} else {
			url = config.PlatformUrl + "/api/uploader/update-note-status"
		}

		transcriptJSON, err := json.Marshal(transcription)
		if err != nil {
			return fmt.Errorf("failed to marshal transcription: %v", err)
		}

		payload := map[string]interface{}{
			"note_id":  req.NoteID,
			"user_id":  req.UserID,
			"status":   "TRANSCRIPTION_COMPLETED",
			"language": language,
			"transcription_processing_time_in_seconds": processingTime,
			"transcript":       string(transcriptJSON),
			"progress_percent": 100,
			"audio_url_path":   fileKey,
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %v", err)
		}

		request, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to create request: %v", err)
		}

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+config.PlatformApiKey)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			return fmt.Errorf("failed to send request: %v", err)
		}
		defer response.Body.Close()

		return nil
	} else if fileType == "pdf" {
		var url string
		if req.WebhookURL != "" {
			url = req.WebhookURL + "/api/document/pdf-callback"
		} else {
			url = config.PlatformUrl + "/api/document/pdf-callback"
		}

		payload := map[string]interface{}{
			"note_id":  req.NoteID,
			"user_id":  req.UserID,
			"status":   "TRANSCRIPTION_COMPLETED",
			"language": language,
			"transcription_processing_time_in_seconds": processingTime,
			"progress_percent":                         100,
			"pdf_url_path":                             fileKey,
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %v", err)
		}

		request, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to create request: %v", err)
		}

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+config.PlatformApiKey)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			return fmt.Errorf("failed to send request: %v", err)
		}
		defer response.Body.Close()

		return nil
	} else if fileType == "txt" {
		var url string
		if req.WebhookURL != "" {
			url = req.WebhookURL + "/api/document/txt-callback"
		} else {
			url = config.PlatformUrl + "/api/document/txt-callback"
		}

		payload := map[string]interface{}{
			"note_id":  req.NoteID,
			"user_id":  req.UserID,
			"status":   "TRANSCRIPTION_COMPLETED",
			"language": language,
			"transcription_processing_time_in_seconds": processingTime,
			"progress_percent":                         100,
			"txt_url_path":                             fileKey,
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %v", err)
		}

		request, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to create request: %v", err)
		}

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+config.PlatformApiKey)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			return fmt.Errorf("failed to send request: %v", err)
		}
		defer response.Body.Close()

		return nil
	}

	return nil
}
