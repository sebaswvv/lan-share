/*
Copyright © 2026 Sebastiaan van Vliet <sebastiaan.van.vliet@hotmail.nl>
*/
package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

// fileHandler manages file sharing requests
type FileHandler struct {
	FilePath string
	FileName string
}

// newFileHandler creates a new file handler
func NewFileHandler(filePath string) *FileHandler {
	return &FileHandler{
		FilePath: filePath,
		FileName: filepath.Base(filePath),
	}
}

// serveHomePage serves the main page with the download button
func (h *FileHandler) ServeHomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := GenerateHTML(h.FileName)
	w.Write([]byte(html))
}

// serveDownload handles file download requests
func (h *FileHandler) ServeDownload(w http.ResponseWriter, r *http.Request) {
	log.Printf("Download request from %s", r.RemoteAddr)

	// open the file
	file, err := os.Open(h.FilePath)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// get file info for size
	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("Error getting file info: %v", err)
		http.Error(w, "Error getting file info", http.StatusInternalServerError)
		return
	}

	// set headers for download
	w.Header().Set("Content-Disposition", "attachment; filename=\""+h.FileName+"\"")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprint(fileInfo.Size()))

	// create progress bar
	bar := progressbar.NewOptions64(
		fileInfo.Size(),
		progressbar.OptionSetDescription(fmt.Sprintf("📤 Sending %s", h.FileName)),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(40),
		progressbar.OptionThrottle(65*1000000), // 65ms
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprintf(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
	)

	// stream the file with progress tracking
	_, err = io.Copy(io.MultiWriter(w, bar), file)
	if err != nil {
		log.Printf("Error streaming file: %v", err)
		return
	}

	log.Printf("File successfully downloaded by %s", r.RemoteAddr)
}

// setupRoutes sets up the HTTP routes
func (h *FileHandler) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.ServeHomePage)
	mux.HandleFunc("/download", h.ServeDownload)
	return mux
}
