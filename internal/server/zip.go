/*
Copyright © 2026 Sebastiaan van Vliet <sebastiaan.van.vliet@hotmail.nl>
*/
package server

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CreateZipArchive creates a zip archive from the given paths (files or directories)
// and returns the path to the created zip file
func CreateZipArchive(paths []string) (string, error) {
	// Create a temporary zip file
	tmpFile, err := os.CreateTemp("", "lanshare-*.zip")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	zipWriter := zip.NewWriter(tmpFile)

	// Add all paths to the zip
	for _, path := range paths {
		err := addToZip(zipWriter, path, "")
		if err != nil {
			zipWriter.Close()
			tmpFile.Close()
			os.Remove(tmpPath)
			return "", fmt.Errorf("failed to add '%s' to zip: %w", path, err)
		}
	}

	// Close the zip writer and check for errors
	if err := zipWriter.Close(); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to close zip writer: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to close temp file: %w", err)
	}

	return tmpPath, nil
}

// addToZip adds a file or directory to the zip archive
func addToZip(zipWriter *zip.Writer, sourcePath, baseInZip string) error {
	fileInfo, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return addDirToZip(zipWriter, sourcePath, baseInZip)
	}
	return addFileToZip(zipWriter, sourcePath, baseInZip)
}

// addFileToZip adds a single file to the zip archive
func addFileToZip(zipWriter *zip.Writer, filePath, baseInZip string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Determine the path in the zip
	var pathInZip string
	if baseInZip == "" {
		pathInZip = filepath.Base(filePath)
	} else {
		pathInZip = filepath.Join(baseInZip, filepath.Base(filePath))
	}

	// Clean and validate the path to prevent path traversal
	pathInZip = filepath.Clean(pathInZip)
	// Ensure the path doesn't escape (no leading .. or absolute path)
	if strings.HasPrefix(pathInZip, "..") || filepath.IsAbs(pathInZip) {
		return fmt.Errorf("invalid path in archive: %s", pathInZip)
	}

	// Create the file in the zip
	writer, err := zipWriter.Create(pathInZip)
	if err != nil {
		return err
	}

	// Copy file content to zip
	_, err = io.Copy(writer, file)
	return err
}

// addDirToZip adds a directory recursively to the zip archive
func addDirToZip(zipWriter *zip.Writer, dirPath, baseInZip string) error {
	dirName := filepath.Base(dirPath)
	
	// If baseInZip is empty, use the directory name as the base
	if baseInZip == "" {
		baseInZip = dirName
	} else {
		baseInZip = filepath.Join(baseInZip, dirName)
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())
		
		if entry.IsDir() {
			err = addDirToZip(zipWriter, fullPath, baseInZip)
		} else {
			err = addFileToZip(zipWriter, fullPath, baseInZip)
		}
		
		if err != nil {
			return err
		}
	}

	return nil
}

// GetArchiveName generates a suitable name for the archive based on the input paths
func GetArchiveName(paths []string) string {
	if len(paths) == 1 {
		// Single path - use its name
		name := filepath.Base(paths[0])
		// If it's a directory, use the directory name
		if info, err := os.Stat(paths[0]); err == nil && info.IsDir() {
			return name + ".zip"
		}
		// If it's a file, replace extension with .zip
		ext := filepath.Ext(name)
		if ext != "" {
			return strings.TrimSuffix(name, ext) + ".zip"
		}
		return name + ".zip"
	}
	// Multiple paths - use a generic name
	return "shared-files.zip"
}
