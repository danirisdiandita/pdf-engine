package service

import (
	"strings"

	"strconv"

	"github.com/danirisdiandita/pdf-engine/internal/config"
)

type TranscriptionOutput struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text"`
}

func GroupSegment(config *config.Config, segments []interface{}) ([]TranscriptionOutput, error) {
	// Initialize variables
	var transcriptionOutput []TranscriptionOutput
	var currentStart, currentEnd float64
	currentStart = -1
	currentText := ""
	currentWordTotalCount := 0
	maxWordsPerRange, err := strconv.Atoi(config.MaxWordsPerTranscriptionRange)
	if err != nil {
		maxWordsPerRange = 150 // Default value
	}

	// Iterate through segments
	for _, segment := range segments {
		segmentMap, ok := segment.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract segment properties
		text, _ := segmentMap["text"].(string)
		start, _ := segmentMap["start"].(float64)
		end, _ := segmentMap["end"].(float64)

		// Count words
		totalWords := len(strings.Fields(text))

		// Initialize start time if needed
		if currentStart < 0 {
			currentStart = start
		}
		currentEnd = end

		// Append text
		currentText += text + " "
		currentWordTotalCount += totalWords

		// Check if word count exceeds threshold
		if currentWordTotalCount > maxWordsPerRange {
			transcriptionOutput = append(transcriptionOutput, TranscriptionOutput{
				Start: currentStart,
				End:   currentEnd,
				Text:  strings.TrimSpace(currentText),
			})

			// Reset variables for next grouping
			currentStart = -1
			currentEnd = 0
			currentText = ""
			currentWordTotalCount = 0
		}
	}

	// Add the last segment group if exists
	if currentStart >= 0 {
		transcriptionOutput = append(transcriptionOutput, TranscriptionOutput{
			Start: currentStart,
			End:   currentEnd,
			Text:  strings.TrimSpace(currentText),
		})
	}

	return transcriptionOutput, nil
}
