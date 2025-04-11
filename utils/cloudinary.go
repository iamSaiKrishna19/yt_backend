package utils

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// CloudinaryService handles all Cloudinary operations
type CloudinaryService struct {
	cld *cloudinary.Cloudinary
}

// NewCloudinaryService creates a new Cloudinary service instance
func NewCloudinaryService() (*CloudinaryService, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return nil, err
	}

	return &CloudinaryService{cld: cld}, nil
}

// SaveToTempFile saves a multipart file to a temporary file and returns the file path
func SaveToTempFile(file *multipart.FileHeader) (string, error) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "upload-*"+filepath.Ext(file.Filename))
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Copy the file content to the temporary file
	if _, err := io.Copy(tempFile, src); err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}

// HandleUpload handles the complete process of saving and uploading files to Cloudinary
func HandleUpload(ctx context.Context, file *multipart.FileHeader, folder string, resourceType string) (string, error) {
	// Initialize Cloudinary service
	cloudinaryService, err := NewCloudinaryService()
	if err != nil {
		return "", err
	}

	// Save to temporary file
	tempFilePath, err := SaveToTempFile(file)
	if err != nil {
		return "", err
	}

	// Open the file
	fileHandle, err := os.Open(tempFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer fileHandle.Close()
	defer os.Remove(tempFilePath) // Clean up the temporary file

	// Upload options
	useFilename := true
	uniqueFilename := true
	uploadParams := uploader.UploadParams{
		Folder:         folder,
		ResourceType:   resourceType,
		UseFilename:    &useFilename,
		UniqueFilename: &uniqueFilename,
	}

	// Upload the file
	result, err := cloudinaryService.cld.Upload.Upload(ctx, fileHandle, uploadParams)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	fmt.Printf("Successfully uploaded file to Cloudinary: %s\n", result.SecureURL)
	return result.SecureURL, nil
}

// HandleImageUpload is a convenience function for uploading images
func HandleImageUpload(ctx context.Context, file *multipart.FileHeader, imageType string) (string, error) {
	return HandleUpload(ctx, file, imageType, "image")
}

// HandleVideoUpload is a convenience function for uploading videos
func HandleVideoUpload(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
	return HandleUpload(ctx, file, folder, "video")
}

// DeleteFromCloudinary deletes a file from Cloudinary using its URL
func DeleteFromCloudinary(ctx context.Context, fileURL string) error {
	// Initialize Cloudinary service
	cloudinaryService, err := NewCloudinaryService()
	if err != nil {
		return err
	}

	// Extract public ID from URL
	// Cloudinary URL format: https://res.cloudinary.com/<cloud_name>/<resource_type>/upload/<public_id>
	// We need to extract the public_id part
	publicID := extractPublicID(fileURL)
	if publicID == "" {
		return fmt.Errorf("invalid Cloudinary URL")
	}

	// Delete the file
	_, err = cloudinaryService.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from Cloudinary: %v", err)
	}

	fmt.Printf("Successfully deleted file from Cloudinary: %s\n", publicID)
	return nil
}

// extractPublicID extracts the public ID from a Cloudinary URL
func extractPublicID(url string) string {
	// Split the URL by '/'
	parts := strings.Split(url, "/")

	// Find the index of "upload"
	uploadIndex := -1
	for i, part := range parts {
		if part == "upload" {
			uploadIndex = i
			break
		}
	}

	// If we found "upload" and there's a part after it
	if uploadIndex != -1 && uploadIndex+1 < len(parts) {
		// Get the public ID (remove file extension)
		publicID := parts[uploadIndex+1]
		return strings.TrimSuffix(publicID, filepath.Ext(publicID))
	}

	return ""
}
