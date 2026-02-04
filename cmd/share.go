/*
Copyright © 2026 Sebastiaan van Vliet <sebastiaan.van.vliet@hotmail.nl>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/sebaswvv/lan-share/internal/server"

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
			var err error
			filePath, err = selectFile()
			if err != nil {
				log.Fatalf("Error selecting file: %v", err)
			}
		} else {
			filePath = args[0]
		}

		if err := validateFile(filePath); err != nil {
			log.Fatalf("Error: %v", err)
		}

		fmt.Printf("Sharing file: %s\n", filePath)

		localIP := getLocalIP()
		srv := setupServer(filePath)

		displayServerInfo(localIP, port, "download")
		runServerWithGracefulShutdown(srv, nil)
	},
}

func selectFile() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	pterm.DefaultSection.Println("📂 File Selection")
	fmt.Println()

	for {
		// get items in current directory
		entries, err := os.ReadDir(currentDir)
		if err != nil {
			return "", fmt.Errorf("error reading directory: %w", err)
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
			return "", fmt.Errorf("no files or folders found")
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
			return "", fmt.Errorf("error selecting: %w", err)
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
			return filepath.Join(currentDir, entry.Name()), nil
		}
	}
}

func validateFile(filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file '%s' does not exist", filePath)
		}
		return fmt.Errorf("unable to access file '%s': %w", filePath, err)
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a file", filePath)
	}

	return nil
}

func setupServer(filePath string) *server.Server {
	fileHandler := server.NewFileHandler(filePath)
	mux := fileHandler.SetupRoutes()
	return server.New(port, mux)
}

func init() {
	rootCmd.AddCommand(shareCmd)

	// add port flag
	shareCmd.Flags().StringVarP(&port, "port", "p", server.DefaultPort, "Port to run the server on")
}
