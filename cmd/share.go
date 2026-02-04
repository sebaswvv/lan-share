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
	Use:   "share [file/folder...]",
	Short: "Share files or folders over the local network",
	Long: `Share one or more files or folders over the local network. 
Provide the path(s) to the file(s) or folder(s) you want to share as arguments.
Multiple files/folders will be automatically zipped together.
Folders will be zipped with their contents.`,
	Run: func(cmd *cobra.Command, args []string) {
		var paths []string
		if len(args) == 0 {
			var err error
			path, err := selectFileOrFolder()
			if err != nil {
				log.Fatalf("Error selecting: %v", err)
			}
			paths = []string{path}
		} else {
			paths = args
		}

		if err := validatePaths(paths); err != nil {
			log.Fatalf("Error: %v", err)
		}

		if len(paths) == 1 {
			fmt.Printf("Sharing: %s\n", paths[0])
		} else {
			fmt.Printf("Sharing %d items:\n", len(paths))
			for _, p := range paths {
				fmt.Printf("  - %s\n", p)
			}
		}

		localIP := getLocalIP()
		srv, cleanup := setupServerForPaths(paths)

		displayServerInfo(localIP, port, "download")
		runServerWithGracefulShutdown(srv, cleanup)
	},
}

func selectFileOrFolder() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	pterm.DefaultSection.Println("📂 File/Folder Selection")
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

		// add current folder as shareable option
		options = append(options, "📦 [Share this entire folder]")
		items = append(items, nil)

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
			WithDefaultText("Select a file, folder, or choose to share current folder (↑/↓ to navigate, Enter to select)").
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
		if selected == "📁 .. (parent directory)" {
			currentDir = filepath.Dir(currentDir)
			continue
		}

		// handle share current folder
		if selected == "📦 [Share this entire folder]" {
			return currentDir, nil
		}

		entry := items[selectedIndex]
		if entry == nil {
			continue
		}

		if entry.IsDir() {
			// Ask if user wants to share the folder or navigate into it
			choice, err := pterm.DefaultInteractiveSelect.
				WithOptions([]string{"📦 Share this folder", "📂 Navigate into folder"}).
				WithDefaultText("What would you like to do?").
				Show()
			
			if err != nil {
				return "", fmt.Errorf("error selecting action: %w", err)
			}

			if choice == "📦 Share this folder" {
				return filepath.Join(currentDir, entry.Name()), nil
			}
			// Navigate into directory
			currentDir = filepath.Join(currentDir, entry.Name())
		} else {
			// selected a file
			return filepath.Join(currentDir, entry.Name()), nil
		}
	}
}

func validatePaths(paths []string) error {
	for _, path := range paths {
		_, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("path '%s' does not exist", path)
			}
			return fmt.Errorf("unable to access path '%s': %w", path, err)
		}
	}
	return nil
}

func setupServerForPaths(paths []string) (*server.Server, func()) {
	var filePath string
	var cleanup func()
	
	// Check if we have a single file (not a directory)
	if len(paths) == 1 {
		info, err := os.Stat(paths[0])
		if err == nil && !info.IsDir() {
			// Single file - serve it directly
			filePath = paths[0]
			cleanup = nil
		} else {
			// Single directory or error - create zip
			zipPath, err := server.CreateZipArchive(paths)
			if err != nil {
				log.Fatalf("Error creating archive: %v", err)
			}
			filePath = zipPath
			cleanup = createCleanupFunc(zipPath)
		}
	} else {
		// Multiple paths - create zip
		zipPath, err := server.CreateZipArchive(paths)
		if err != nil {
			log.Fatalf("Error creating archive: %v", err)
		}
		filePath = zipPath
		cleanup = createCleanupFunc(zipPath)
	}
	
	fileHandler := server.NewFileHandler(filePath)
	mux := fileHandler.SetupRoutes()
	return server.New(port, mux), cleanup
}

func createCleanupFunc(zipPath string) func() {
	return func() {
		os.Remove(zipPath)
	}
}

func init() {
	rootCmd.AddCommand(shareCmd)

	// add port flag
	shareCmd.Flags().StringVarP(&port, "port", "p", server.DefaultPort, "Port to run the server on")
}
