package handler

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/danirisdiandita/pdf-engine/internal/config"
	"github.com/danirisdiandita/pdf-engine/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Generate(c *gin.Context) {
	var payload model.PDFRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request payload"})
		return
	}

	Config := config.Load()

	content_ := payload.Content
	type_ := payload.Type

	directory := ""

	if Config.AppEnv == "production" {
		directory = "/tmp"
	} else {
		directory = "./tmp"
	}

	inputFilePath := directory + "/" + uuid.New().String() + ".md"

	// save content_ (as string) to a file in directory
	os.WriteFile(inputFilePath, []byte(content_), 0644)

	filePath := directory + "/" + uuid.New().String() + ".pdf"
	if type_ == model.PDFTypeNote {
		cmd := exec.Command("pandoc", inputFilePath, "-o",
			filePath,
			"-V",
			"geometry:margin=1in",
			"-V", "mainfont=Helvetica")
		cmd.Run()

		// read file as byte and make response
		file, err := os.ReadFile(filePath)

		// delete input and output file
		if err := os.Remove(inputFilePath); err != nil && !os.IsNotExist(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("Warning: could not remove file %s: %v\n", inputFilePath, err)})
			return
		}
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("Warning: could not remove file %s: %v\n", filePath, err)})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
		c.Data(http.StatusOK, "application/pdf", file)
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "result": payload, "language": ""})
}
