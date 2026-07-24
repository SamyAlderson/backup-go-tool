// Package utils contains helper functions for the backup-go-tool project.
package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GetExecutablePath returns the path to the executable.
func GetExecutablePath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	return filepath.Dir(exePath), nil
}

// GetHomeDir returns the user's home directory.
func GetHomeDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return homeDir, nil
}

// CheckDirExists checks if a directory exists.
func CheckDirExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist", path)
	}
	return nil
}

// CreateDir creates a directory if it does not exist.
func CreateDir(path string) error {
	if err := CheckDirExists(path); err != nil {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}
	}
	return nil
}

// CheckFileExists checks if a file exists.
func CheckFileExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist", path)
	}
	return nil
}

// CopyFile copies a file from one location to another.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("failed to copy file from %s to %s: %w", src, dst, err)
	}

	return nil
}

// GetFileSize returns the size of a file in bytes.
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("failed to get file size %s: %w", path, err)
	}
	return info.Size(), nil
}

// GetFileExtension returns the file extension.
func GetFileExtension(path string) string {
	return filepath.Ext(path)
}

// IsAbsolutePath checks if a path is absolute.
func IsAbsolutePath(path string) bool {
	return strings.HasPrefix(path, "/")
}

// LogError logs an error with a custom message.
func LogError(err error, message string) {
	log.Printf("Error: %s - %s", message, err)
}