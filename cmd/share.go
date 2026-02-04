/*
Copyright © 2026 Sebastiaan van Vliet <sebastiaan.van.vliet@hotmail.nl>
*/
package cmd

import (
	"context"
	"fmt"
	"lanshare/internal/server"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fatih/color"
	qrterminal "github.com/mdp/qrterminal/v3"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var port string

// shareCmd represents the share command
var shareCmd = &cobra.Command{
	Use:   "share [file]",
	Short: "Share a file over the local network",
	Long: `Share a file over the local network. 
Provide the path to the file you want to share as an argument.
The file must exist and cannot be a directory.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var filePath string
		if len(args) == 0 {
			filePath = selectFile()
		} else {
			filePath = args[0]
		}

		validateFile(filePath)
		fmt.Printf("Sharing file: %s\n", filePath)

		localIP := getLocalIP()
		srv := setupServer(filePath, localIP)

		displayServerInfo(localIP)
		runServer(srv)
	},
}

func selectFile() string {
	currentDir, _ := os.Getwd()

	pterm.DefaultSection.Println("📂 File Selection")
	fmt.Println()

	for {
		// get items in current directory
		entries, err := os.ReadDir(currentDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading directory: %v\n", err)
			os.Exit(1)
		}

		var options []string
		var items []os.DirEntry

		// add parent directory option
		parentDir := filepath.Dir(currentDir)
		if currentDir != parentDir {
			options = append(options, "📁 .. (parent directory)")
			items = append(items, nil)
		}

		// add folders first, then files
		for _, entry := range entries {
			if entry.IsDir() {
				options = append(options, "📁 "+entry.Name()+"/")
				items = append(items, entry)
			}
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				options = append(options, "📄 "+entry.Name())
				items = append(items, entry)
			}
		}

		if len(options) == 0 {
			fmt.Println("No files or folders found")
			os.Exit(1)
		}

		// show current directory
		pterm.DefaultBasicText.Printf("Current: %s\n\n", currentDir)

		// create interactive select
		selected, err := pterm.DefaultInteractiveSelect.
			WithOptions(options).
			WithDefaultText("Select a file or folder (↑/↓ to navigate, Enter to select)").
			WithMaxHeight(15).
			Show()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error selecting: %v\n", err)
			os.Exit(1)
		}

		// find selected index
		selectedIndex := -1
		for i, opt := range options {
			if opt == selected {
				selectedIndex = i
				break
			}
		}

		if selectedIndex == -1 {
			continue
		}

		// handle parent directory
		if selectedIndex == 0 && currentDir != filepath.Dir(currentDir) && items[0] == nil {
			currentDir = filepath.Dir(currentDir)
			continue
		}

		entry := items[selectedIndex]
		if entry == nil {
			continue
		}

		if entry.IsDir() {
			// navigate into directory
			currentDir = filepath.Join(currentDir, entry.Name())
		} else {
			// selected a file
			return filepath.Join(currentDir, entry.Name())
		}
	}
}

func validateFile(filePath string) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: file '%s' does not exist\n", filePath)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Error: unable to access file '%s': %v\n", filePath, err)
		os.Exit(1)
	}

	if fileInfo.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: '%s' is a directory, not a file\n", filePath)
		os.Exit(1)
	}
}

func getLocalIP() string {
	localIP, err := server.GetLocalIP()
	if err != nil {
		log.Printf("Warning: could not determine local IP: %v", err)
		return "localhost"
	}
	return localIP
}

func setupServer(filePath, localIP string) *server.Server {
	fileHandler := server.NewFileHandler(filePath)
	mux := fileHandler.SetupRoutes()
	return server.New(port, mux)
}

func displayServerInfo(localIP string) {
	green := color.New(color.FgGreen, color.Bold)
	cyan := color.New(color.FgCyan, color.Bold)
	magenta := color.New(color.FgMagenta, color.Bold)
	yellow := color.New(color.FgYellow)

	url := fmt.Sprintf("http://%s:%s", localIP, port)

	fmt.Println()
	green.Println("✓ Server started successfully!")
	fmt.Println()
	magenta.Println("📱 Scan QR code or use the URL below:")
	fmt.Println()
	qrterminal.GenerateHalfBlock(url, qrterminal.L, os.Stdout)
	fmt.Println()
	magenta.Print("🌐  URL: ")
	cyan.Println(url)
	fmt.Println()
	yellow.Println("📡 Waiting for connections... Press Ctrl+C to stop")
	fmt.Println()
}

func runServer(srv *server.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-sigChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	red := color.New(color.FgRed, color.Bold)
	red.Println("\n🛑 Server stopped.")
}

func init() {
	rootCmd.AddCommand(shareCmd)

	// add port flag
	shareCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to run the server on")
}
