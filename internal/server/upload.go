/*
Copyright © 2026 Sebastiaan van Vliet <sebastiaan.van.vliet@hotmail.nl>
*/
package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

// UploadHandler manages file upload requests
type UploadHandler struct {
	savePath       string
	pendingUploads chan *PendingUpload
}

// PendingUpload represents a file waiting for approval
type PendingUpload struct {
	Filename string
	Filesize int64
	TempPath string
	Response chan bool
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler() *UploadHandler {
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Warning: could not get working directory, using temp: %v", err)
		cwd = os.TempDir()
	}
	return &UploadHandler{
		savePath:       cwd,
		pendingUploads: make(chan *PendingUpload, PendingUploadBufferSize),
	}
}

// ServeUploadPage serves the upload page
func (h *UploadHandler) ServeUploadPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := GenerateUploadHTML()
	w.Write([]byte(html))
}

// sanitizeFilename removes path traversal attempts and dangerous characters
func sanitizeFilename(filename string) (string, error) {
	// Get just the base filename, removing any path components
	filename = filepath.Base(filename)

	// Check for empty filename after cleaning
	if filename == "" || filename == "." || filename == ".." {
		return "", fmt.Errorf("invalid filename")
	}

	// Remove any remaining path separators
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")

	return filename, nil
}

// HandleUpload processes file uploads
func (h *UploadHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// parse multipart form with size limit
	err := r.ParseMultipartForm(MaxUploadSize)
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Error parsing form or file too large", http.StatusBadRequest)
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

	// Sanitize filename for security
	filename, err = sanitizeFilename(filename)
	if err != nil {
		log.Printf("Invalid filename: %v", err)
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

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
		progressbar.OptionSetWidth(ProgressBarWidth),
		progressbar.OptionThrottle(ProgressBarThrottle),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprintf(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(ProgressBarSpinnerType),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
	)

	// copy to temp file with progress and context cancellation
	ctx := r.Context()
	done := make(chan error, 1)

	go func() {
		_, err := io.Copy(io.MultiWriter(tempFile, bar), file)
		done <- err
	}()

	var copyErr error
	select {
	case <-ctx.Done():
		tempFile.Close()
		os.Remove(tempPath)
		log.Printf("Upload cancelled by client")
		return
	case copyErr = <-done:
		tempFile.Close()
	}

	if copyErr != nil {
		os.Remove(tempPath)
		log.Printf("Error saving file: %v", copyErr)
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

// SetupRoutes sets up the HTTP routes for uploads
func (h *UploadHandler) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.ServeUploadPage)
	mux.HandleFunc("/upload", h.HandleUpload)
	return mux
}

// ProcessUploads handles pending upload approvals with context cancellation
func (h *UploadHandler) ProcessUploads(ctx context.Context) {
	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen, color.Bold)
	red := color.New(color.FgRed, color.Bold)

	for {
		select {
		case <-ctx.Done():
			// Shutdown requested, reject any pending uploads
			for {
				select {
				case pending := <-h.pendingUploads:
					os.Remove(pending.TempPath)
					pending.Response <- false
				default:
					return
				}
			}
		case pending := <-h.pendingUploads:
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

				// Try to rename (move) the file
				err := os.Rename(pending.TempPath, destPath)
				if err != nil {
					// If rename fails (different filesystems), copy instead
					if copyErr := copyFile(pending.TempPath, destPath); copyErr != nil {
						log.Printf("Error saving file: %v", copyErr)
						red.Printf("❌ Error saving file: %v\n", copyErr)
						os.Remove(pending.TempPath)
						pending.Response <- false
						fmt.Println()
						continue
					}
					// Remove temp file after successful copy
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
}

// copyFile copies a file from src to dst with proper error handling
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination: %w", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// Ensure data is written to disk
	if err := dstFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	return nil
}
