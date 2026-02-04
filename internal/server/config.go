/*
Copyright © 2026 Sebastiaan van Vliet <sebastiaan.van.vliet@hotmail.nl>
*/
package server

import "time"

const (
	// file size limits
	MaxUploadSize = 100 * 1024 * 1024 // 100 MB

	// HTTP server configuration
	MaxHeaderBytes = 1 * 1024 * 1024 // 1 MB

	// progress bar configuration
	ProgressBarWidth       = 40
	ProgressBarThrottle    = 65 * time.Millisecond
	ProgressBarSpinnerType = 14

	// server defaults
	DefaultPort = "8080"

	// upload configuration
	PendingUploadBufferSize = 10
	ShutdownTimeout         = 5 * time.Second
)
