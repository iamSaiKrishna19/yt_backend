package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// GetVideoDuration gets the duration of a video file in HH:MM:SS format
func GetVideoDuration(file *multipart.FileHeader) (string, error) {
	// Save to temporary file
	tempFilePath, err := SaveToTempFile(file)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}
	defer os.Remove(tempFilePath)

	// Use ffprobe to get video duration
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		tempFilePath,
	)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get video duration: %v", err)
	}

	// Convert string output to float64
	durationStr := strings.TrimSpace(string(output))
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse video duration: %v", err)
	}

	// Convert to HH:MM:SS format
	hours := int(duration) / 3600
	minutes := (int(duration) % 3600) / 60
	seconds := int(duration) % 60

	// Format the duration
	durationFormatted := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	return durationFormatted, nil
}
