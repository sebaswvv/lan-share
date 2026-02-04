/*
Copyright © 2026 Sebastiaan van Vliet <sebastiaan.van.vliet@hotmail.nl>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
	qrterminal "github.com/mdp/qrterminal/v3"
	"github.com/sebaswvv/lan-share/internal/server"
)

// getLocalIP retrieves the local IP address
func getLocalIP() string {
	localIP, err := server.GetLocalIP()
	if err != nil {
		log.Printf("Warning: could not determine local IP: %v", err)
		return "localhost"
	}
	return localIP
}

// displayServerInfo shows server connection information with QR code
func displayServerInfo(localIP, port, mode string) {
	green := color.New(color.FgGreen, color.Bold)
	cyan := color.New(color.FgCyan, color.Bold)
	magenta := color.New(color.FgMagenta, color.Bold)
	yellow := color.New(color.FgYellow)

	url := fmt.Sprintf("http://%s:%s", localIP, port)

	fmt.Println()
	green.Println("✓ Server started successfully!")
	fmt.Println()

	if mode == "upload" {
		magenta.Println("📱 Scan QR code to upload files:")
	} else {
		magenta.Println("📱 Scan QR code or use the URL below:")
	}

	fmt.Println()
	qrterminal.GenerateHalfBlock(url, qrterminal.L, os.Stdout)
	fmt.Println()
	magenta.Print("🌐  URL: ")
	cyan.Println(url)
	fmt.Println()

	if mode == "upload" {
		yellow.Println("📥 Waiting for uploads... Press Ctrl+C to stop")
	} else {
		yellow.Println("📡 Waiting for connections... Press Ctrl+C to stop")
	}

	fmt.Println()
}

// runServerWithGracefulShutdown runs the server with signal handling
func runServerWithGracefulShutdown(srv *server.Server, onShutdown func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-sigChan

	ctx, cancel := context.WithTimeout(context.Background(), server.ShutdownTimeout)
	defer cancel()

	if onShutdown != nil {
		onShutdown()
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	red := color.New(color.FgRed, color.Bold)
	red.Println("\n🛑 Server stopped.")
}
