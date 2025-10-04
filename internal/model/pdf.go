package model

type PDFType string

const (
	PDFTypeNote      PDFType = "note"
	PDFTypeFlashcard PDFType = "flashcard"
	PDFTypeQuiz      PDFType = "quiz"
)

type PDFRequest struct {
	Content  string  `json:"content"`
	Type     PDFType `json:"type"`
	Language string  `json:"lang,omitempty"`
}
