package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/DonBigBon/parser-backend/internal/filehandler"
	"github.com/DonBigBon/parser-backend/internal/parser"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("document")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filePath, err := filehandler.SaveUploadedFile(file, handler.Filename)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	codeData, err := parser.ParseDocument(filePath)
	if err != nil {
		http.Error(w, "Error parsing document: "+err.Error(), http.StatusInternalServerError)
		return
	}

	csvFiles, err := filehandler.GenerateCSV(codeData)
	if err != nil {
		http.Error(w, "Error generating CSV files", http.StatusInternalServerError)
		return
	}

	sqlDump, err := filehandler.GenerateSQLDump(codeData)
	if err != nil {
		http.Error(w, "Error generating SQL dump", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":  "File processed successfully",
		"csvFiles": csvFiles,
		"sqlDump":  sqlDump,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filename := r.URL.Query().Get("file")
	if filename == "" {
		http.Error(w, "File parameter is required", http.StatusBadRequest)
		return
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var contentType string
	ext := filepath.Ext(filename)
	switch ext {
	case ".csv":
		contentType = "text/csv"
	case ".xlsx":
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".sql":
		contentType = "application/sql"
	default:
		contentType = "application/octet-stream"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filename))

	io.Copy(w, file)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Загрузка документа</title>
		</head>
		<body>
			<h1>Загрузка документа для парсинга</h1>
			<form method="post" action="/upload" enctype="multipart/form-data">
				<input type="file" name="document" accept=".docx,.doc,.txt,.rtf" required />
				<button type="submit">Загрузить и обработать</button>
			</form>
		</body>
		</html>
	`)
}
