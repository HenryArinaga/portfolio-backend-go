package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	UploadDir     = "uploads/images"
	MaxUploadSize = 10 << 20 // 10 MB
)

func init() {
	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create upload directory: %v", err))
	}
}

// SaveImage saves an uploaded image and returns the URL path
func SaveImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Validate file size
	if header.Size > MaxUploadSize {
		return "", fmt.Errorf("file too large: max size is 10MB")
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !isValidImageExt(ext) {
		return "", fmt.Errorf("invalid file type: only jpg, jpeg, png, gif, webp allowed")
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d%s", timestamp, ext)
	filepath := filepath.Join(UploadDir, filename)

	// Create destination file
	dst, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer dst.Close()

	// Copy uploaded file to destination
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	// Return URL path (relative to server root)
	return "/" + filepath, nil
}

// DeleteImage removes an image file from storage
func DeleteImage(imageURL string) error {
	if imageURL == "" {
		return nil
	}

	// Remove leading slash to get filesystem path
	path := strings.TrimPrefix(imageURL, "/")
	
	// Only delete if it's in our upload directory
	if !strings.HasPrefix(path, UploadDir) {
		return nil
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete image: %v", err)
	}

	return nil
}

func isValidImageExt(ext string) bool {
	validExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	return validExts[ext]
}
