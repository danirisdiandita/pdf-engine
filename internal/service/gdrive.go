package service

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/danirisdiandita/pdf-engine/internal/config"
	"github.com/danirisdiandita/pdf-engine/internal/utils"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
)

func DownloadGdriveByProjectId(config *config.Config, fileID string) (string, bool, bool, bool, error) {
	filePath := fmt.Sprintf("./tmp/%s", uuid.New().String())
	if _, err := os.Stat("./tmp"); os.IsNotExist(err) {
		if err := os.MkdirAll("./tmp", 0755); err != nil {
			return "", false, false, false, fmt.Errorf("failed to create directory ./tmp: %v", err)
		}
	}

	url := fmt.Sprintf("https://drive.usercontent.google.com/download?id=%s&confirm=xxx", fileID)
	resp, err := http.Get(url)
	if err != nil {
		return "", false, false, false, fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", false, false, false, fmt.Errorf("failed to download file: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false, false, false, fmt.Errorf("failed to read response body: %v", err)
	}

	if err := os.WriteFile(filePath, body, 0644); err != nil {
		return "", false, false, false, fmt.Errorf("failed to write file to disk: %v", err)
	}

	if utils.IsText(filePath) {
		return filePath, false, false, true, nil
	}

	buf, _ := os.ReadFile(filePath)

	kind, _ := filetype.Match(buf)
	if kind == filetype.Unknown {
		return "", false, false, false, fmt.Errorf("unsupported file type")
	}

	availableAudioExtensions := []string{"mp3", "ogg", "flac", "aac", "aiff", "m4a"} // wma not supported
	availableVideoExtensions := []string{"webm"}
	availableDocumentExtension := []string{"pdf", "txt"}
	isValidAudio := false
	isValidVideo := false
	isValidDocument := false

	// check if file is audio

	for _, ext := range availableAudioExtensions {
		if kind.Extension == ext {
			isValidAudio = true
			break
		}
	}

	// check if file is video
	for _, ext := range availableVideoExtensions {
		if kind.Extension == ext {
			isValidVideo = true
			break
		}
	}

	// check if file is document
	for _, ext := range availableDocumentExtension {
		if kind.Extension == ext {
			isValidDocument = true
			break
		}
	}

	if !isValidAudio && !isValidVideo && !isValidDocument {
		return "", false, false, false, fmt.Errorf("unsupported file type")
	}

	return filePath, isValidAudio, isValidVideo, isValidDocument, nil
}
