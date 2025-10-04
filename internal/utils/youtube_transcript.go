package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/danirisdiandita/pdf-engine/internal/model"
	"golang.org/x/net/html"
)

// GetTranscriptFromYoutube extracts captions/transcript from a YouTube video
func GetTranscriptFromYoutube(url string, lang string) (*model.TranscriptResponse, error) {
	// Fetch the YouTube page
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check if captions exist
	if !strings.Contains(string(data), "captionTracks") {
		return nil, errors.New(model.ErrorCaptionsNotFound)
	}

	// Extract caption tracks using regex
	regex := regexp.MustCompile(`"captionTracks":(\[.*?\])`)
	matches := regex.FindStringSubmatch(string(data))
	if len(matches) < 2 {
		return nil, errors.New(model.ErrorCaptionsNotFound)
	}

	// Parse the JSON
	var captionData struct {
		CaptionTracks []struct {
			BaseURL string `json:"baseUrl"`
			Name    struct {
				SimpleText string `json:"simpleText"`
			} `json:"name"`
			VssID          string `json:"vssId"`
			LanguageCode   string `json:"languageCode"`
			Kind           string `json:"kind"`
			IsTranslatable bool   `json:"isTranslatable"`
			TrackName      string `json:"trackName"`
		} `json:"captionTracks"`
	}

	err = json.Unmarshal([]byte(fmt.Sprintf(`{"captionTracks":%s}`, matches[1])), &captionData)
	if err != nil {
		return nil, err
	}

	if len(captionData.CaptionTracks) == 0 {
		return nil, errors.New(model.ErrorCaptionsNotFound)
	}

	// Use default language if none specified
	langCode := lang
	if langCode == "" {
		langCode = captionData.CaptionTracks[0].LanguageCode
	}

	// Find the appropriate track
	var selectedTrack int = -1

	for i, t := range captionData.CaptionTracks {
		if t.LanguageCode == langCode {
			selectedTrack = i
			break
		}
	}

	// If requested language not found, use default
	if selectedTrack == -1 {
		langCode = captionData.CaptionTracks[0].LanguageCode
		selectedTrack = 0
	}

	// Fetch the transcript XML
	transcriptResp, err := http.Get(captionData.CaptionTracks[selectedTrack].BaseURL)
	if err != nil {
		return nil, err
	}
	defer transcriptResp.Body.Close()

	transcriptData, err := io.ReadAll(transcriptResp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the XML
	transcript := string(transcriptData)
	transcript = strings.Replace(transcript, `<?xml version="1.0" encoding="utf-8" ?><transcript>`, "", 1)
	transcript = strings.Replace(transcript, "</transcript>", "", 1)

	// Process transcript lines
	lines := []model.TranscriptLine{}
	for _, line := range strings.Split(transcript, "</text>") {
		if line == "" || strings.TrimSpace(line) == "" {
			continue
		}

		// Extract start time
		startRegex := regexp.MustCompile(`start="([\d.]+)"`)
		startMatches := startRegex.FindStringSubmatch(line)
		if len(startMatches) < 2 {
			continue
		}
		start := startMatches[1]

		// Extract duration
		durRegex := regexp.MustCompile(`dur="([\d.]+)"`)
		durMatches := durRegex.FindStringSubmatch(line)
		if len(durMatches) < 2 {
			continue
		}
		dur := durMatches[1]

		// Extract and clean text
		textRegex := regexp.MustCompile(`<text.+>(.*?)$`)
		textMatches := textRegex.FindStringSubmatch(line)
		if len(textMatches) < 2 {
			continue
		}

		htmlText := textMatches[1]
		htmlText = strings.ReplaceAll(htmlText, "&amp;", "&")

		// Decode HTML entities and strip HTML tags
		decodedText := html.UnescapeString(htmlText)
		decodedText = stripHTMLTags(decodedText)

		lines = append(lines, model.TranscriptLine{
			Start: start,
			Dur:   dur,
			Text:  decodedText,
		})
	}

	// Process into chunks based on maximum words per range
	maxWordsPerTranscriptionRange := 150
	if envMax, exists := os.LookupEnv("MAX_WORDS_PER_TRANSCRIPTION_RANGE"); exists {
		if val, err := strconv.Atoi(envMax); err == nil {
			maxWordsPerTranscriptionRange = val
		}
	}

	transcript_ := []model.TranscriptRange{}
	chunk := struct {
		Start float64
		End   float64
		Text  string
	}{
		Start: 0,
		End:   0,
		Text:  "",
	}

	for i := 0; i < len(lines); i++ {
		start, _ := strconv.ParseFloat(lines[i].Start, 64)
		dur, _ := strconv.ParseFloat(lines[i].Dur, 64)
		text := lines[i].Text

		// Check if current chunk exceeds word limit
		totalWords := len(strings.Fields(chunk.Text))
		if totalWords >= maxWordsPerTranscriptionRange && chunk.Text != "" {
			// Fix end time
			if chunk.End > 0 {
				chunk.End = min(chunk.End, start-0.001)
			}

			chunk.Text = strings.TrimSpace(chunk.Text)
			startStr := fmt.Sprintf("%.3f", chunk.Start)
			endStr := fmt.Sprintf("%.3f", chunk.End)

			transcript_ = append(transcript_, model.TranscriptRange{
				Start: startStr,
				End:   endStr,
				Text:  chunk.Text,
			})

			chunk.Text = ""
			chunk.Start = start
			chunk.End = start + dur
		}

		chunk.End = start + dur
		chunk.Text += " " + text
	}

	// Add the final chunk
	if chunk.Text != "" {
		chunk.Text = strings.TrimSpace(chunk.Text)
		startStr := fmt.Sprintf("%.3f", chunk.Start)
		endStr := fmt.Sprintf("%.3f", chunk.End)

		transcript_ = append(transcript_, model.TranscriptRange{
			Start: startStr,
			End:   endStr,
			Text:  chunk.Text,
		})
	}

	return &model.TranscriptResponse{
		Data:         transcript_,
		DetectedLang: langCode,
		InputLang:    lang,
	}, nil
}

// Helper function to strip HTML tags
func stripHTMLTags(s string) string {
	var result strings.Builder
	inTag := false

	for _, r := range s {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// Helper function for min of two float64s
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
