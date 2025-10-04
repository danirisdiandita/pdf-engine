package utils

import (
	"github.com/gabriel-vasile/mimetype"
)

func IsText(filePath string) bool {
	mime, err := mimetype.DetectFile(filePath)
	if err != nil {
		return false
	}
	return mime.Is("text/plain")
}
