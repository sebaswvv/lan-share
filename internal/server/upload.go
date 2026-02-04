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

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

// uploadHandler manages file upload requests
type UploadHandler struct {
	savePath       string
	pendingUploads chan *PendingUpload
}

// pendingUpload represents a file waiting for approval
type PendingUpload struct {
	Filename string
	Filesize int64
	TempPath string
	Response chan bool
}

// newUploadHandler creates a new upload handler
func NewUploadHandler() *UploadHandler {
	cwd, _ := os.Getwd()
	return &UploadHandler{
		savePath:       cwd,
		pendingUploads: make(chan *PendingUpload, 10),
	}
}

// serveUploadPage serves the upload page
func (h *UploadHandler) ServeUploadPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := GenerateUploadHTML()
	w.Write([]byte(html))
}

// handleUpload processes file uploads
func (h *UploadHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// parse multipart form (max 100MB)
	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file: %v", err)
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := header.Filename
	filesize := header.Size

	yellow := color.New(color.FgYellow, color.Bold)

	fmt.Println()
	yellow.Printf("📤 Incoming file: %s (%.2f MB)\n", filename, float64(filesize)/(1024*1024))

	// save to temp file first
	tempFile, err := os.CreateTemp("", "lanshare-*")
	if err != nil {
		log.Printf("Error creating temp file: %v", err)
		http.Error(w, "Error processing file", http.StatusInternalServerError)
		return
	}
	tempPath := tempFile.Name()

	// create progress bar for receiving
	bar := progressbar.NewOptions64(
		filesize,
		progressbar.OptionSetDescription(fmt.Sprintf("📥 Receiving %s", filename)),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(40),
		progressbar.OptionThrottle(65*1000000),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprintf(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
	)

	// copy to temp file with progress
	_, err = io.Copy(io.MultiWriter(tempFile, bar), file)
	tempFile.Close()

	if err != nil {
		os.Remove(tempPath)
		log.Printf("Error saving file: %v", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// send for approval
	pending := &PendingUpload{
		Filename: filename,
		Filesize: filesize,
		TempPath: tempPath,
		Response: make(chan bool),
	}

	h.pendingUploads <- pending

	// wait for approval
	accepted := <-pending.Response

	if accepted {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Upload Successful</title>
	<style>
		body {
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
			display: flex;
			align-items: center;
			justify-content: center;
			min-height: 100vh;
			margin: 0;
			background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		}
		.container {
			text-align: center;
			background: white;
			padding: 48px;
			border-radius: 24px;
			box-shadow: 0 20px 60px rgba(0,0,0,0.3);
		}
		.success-icon {
			font-size: 64px;
			margin-bottom: 24px;
		}
		h1 {
			color: #2d3748;
			margin-bottom: 16px;
		}
		p {
			color: #718096;
			margin-bottom: 32px;
		}
		button {
			background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
			color: white;
			border: none;
			padding: 16px 32px;
			font-size: 16px;
			border-radius: 12px;
			cursor: pointer;
		}
	</style>
</head>
<body>
	<div class="container">
		<div class="success-icon">✅</div>
		<h1>Upload Accepted!</h1>
		<p>Your file has been accepted and saved.</p>
		<button onclick="window.location.href='/'">Upload Another File</button>
	</div>
</body>
</html>`))
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Upload Rejected</title>
	<style>
		body {
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
			display: flex;
			align-items: center;
			justify-content: center;
			min-height: 100vh;
			margin: 0;
			background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		}
		.container {
			text-align: center;
			background: white;
			padding: 48px;
			border-radius: 24px;
			box-shadow: 0 20px 60px rgba(0,0,0,0.3);
		}
		.icon {
			font-size: 64px;
			margin-bottom: 24px;
		}
		h1 {
			color: #2d3748;
			margin-bottom: 16px;
		}
		p {
			color: #718096;
			margin-bottom: 32px;
		}
		button {
			background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
			color: white;
			border: none;
			padding: 16px 32px;
			font-size: 16px;
			border-radius: 12px;
			cursor: pointer;
		}
	</style>
</head>
<body>
	<div class="container">
		<div class="icon">❌</div>
		<h1>Upload Rejected</h1>
		<p>The file was rejected by the receiver.</p>
		<button onclick="window.location.href='/'">Try Again</button>
	</div>
</body>
</html>`))
	}
}

// setupRoutes sets up the HTTP routes for uploads
func (h *UploadHandler) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.ServeUploadPage)
	mux.HandleFunc("/upload", h.HandleUpload)
	return mux
}

// processUploads handles pending upload approvals
func (h *UploadHandler) ProcessUploads() {
	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen, color.Bold)
	red := color.New(color.FgRed, color.Bold)

	for pending := range h.pendingUploads {
		fmt.Println()
		cyan.Printf("📋 File: %s (%.2f MB)\n", pending.Filename, float64(pending.Filesize)/(1024*1024))
		fmt.Print("Accept this file? (y/n): ")

		var response string
		fmt.Scanln(&response)

		accepted := response == "y" || response == "Y" || response == "yes" || response == "Yes"

		if accepted {
			// move from temp to final location
			destPath := filepath.Join(h.savePath, pending.Filename)

			// check if file exists, append number if needed
			counter := 1
			for {
				if _, err := os.Stat(destPath); os.IsNotExist(err) {
					break
				}
				ext := filepath.Ext(pending.Filename)
				nameWithoutExt := pending.Filename[:len(pending.Filename)-len(ext)]
				destPath = filepath.Join(h.savePath, fmt.Sprintf("%s_%d%s", nameWithoutExt, counter, ext))
				counter++
			}

			err := os.Rename(pending.TempPath, destPath)
			if err != nil {
				// if rename fails, try copy
				srcFile, _ := os.Open(pending.TempPath)
				dstFile, _ := os.Create(destPath)
				io.Copy(dstFile, srcFile)
				srcFile.Close()
				dstFile.Close()
				os.Remove(pending.TempPath)
			}

			green.Printf("✅ File saved: %s\n", destPath)
			pending.Response <- true
		} else {
			os.Remove(pending.TempPath)
			red.Println("❌ File rejected and deleted")
			pending.Response <- false
		}
		fmt.Println()
	}
}
