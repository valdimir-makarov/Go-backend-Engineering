package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/kafka"
)

type FileHander struct {
	fileProd *kafka.KafkaProducer
}

func NewFileHandler(producer *kafka.KafkaProducer) *FileHander {
	return &FileHander{fileProd: producer}

}

func (fl *FileHander) SendFileHandler(w http.ResponseWriter, r *http.Request) {
	senderId := r.URL.Query().Get("sender_id")
	receiverId := r.URL.Query().Get("receiver_id")

	if senderId == "" || receiverId == "" {
		http.Error(w, "sender_id and receiver_id required", http.StatusBadRequest)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File error", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadPath := "./uploads"
	os.MkdirAll(uploadPath, os.ModePerm)
	filePath := fmt.Sprintf("%s/%s", uploadPath, handler.Filename)

	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	io.Copy(out, file)
	fl.fileProd.SendFileUpLoadEvent(senderId, receiverId, handler.Filename)
	// Create download URL (in production use actual domain)
	downloadURL := fmt.Sprintf("http://localhost:8080/uploads/%s", handler.Filename)
	fmt.Printf(downloadURL)
}
