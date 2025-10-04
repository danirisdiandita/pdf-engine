package service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/danirisdiandita/pdf-engine/internal/config"
	"github.com/danirisdiandita/pdf-engine/internal/model"
)

func ProcessTranscription(req model.TranscribeRequest) ([]interface{}, string, error) {
	filePath, err := DownloadFile(config.Load(), req.AudioPathKey)
	if err != nil {
		fmt.Printf("Failed to download file: %v\n", err)
		return nil, "", err
	}
	segments, lang, err := Transcribe(config.Load(), filePath, req.Language)
	if err != nil {
		fmt.Printf("Failed to transcribe: %v\n", err)
		return nil, "", err
	}

	// delete file
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		fmt.Printf("Warning: could not remove file %s: %v\n", filePath, err)
	}
	return segments, lang, nil
}

func ProcessGdriveTranscription(req model.TranscribeGdriveRequest) ([]interface{}, string, string, string, error) {
	filePath, isValidAudio, isValidVideo, isValidDocument, err := DownloadGdriveByProjectId(config.Load(), req.FileID)
	if err != nil {
		fmt.Printf("Failed to download file: %v\n", err)
		return nil, "", "", "", err
	}

	if isValidAudio || isValidVideo {
		segments, lang, err := Transcribe(config.Load(), filePath, req.Language)
		if err != nil {
			fmt.Printf("Failed to transcribe: %v\n", err)
			return nil, "", "", "", err
		}

		// upload audio to s3 before delete file

		// extract fileName from filePath
		fileName := filepath.Base(filePath)

		fileKey := fmt.Sprintf("%s/%s", "uploads", fileName)

		if err := UploadFile(config.Load(), filePath, fileKey); err != nil {
			fmt.Printf("Failed to upload file: %v\n", err)
			return nil, "", "", "", err
		}

		// delete file
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			fmt.Printf("Warning: could not remove file %s: %v\n", filePath, err)
		}
		return segments, lang, fileKey, "audio", nil
	} else if isValidDocument {
		fileName := filepath.Base(filePath)
		fileKey := fmt.Sprintf("%s/%s", "uploads", fileName)

		if err := UploadFile(config.Load(), filePath, fileKey); err != nil {
			fmt.Printf("Failed to upload file: %v\n", err)
			return nil, "", "", "", err
		}

		// delete file
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			fmt.Printf("Warning: could not remove file %s: %v\n", filePath, err)
		}
		return nil, req.Language, fileKey, "pdf", nil
	}

	return nil, "", "", "", nil
}
