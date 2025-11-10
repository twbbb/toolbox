package osutil

import (
	"os"
	"testing"
)

func TestGuessFileType(t *testing.T) {
	tests := []struct {
		name     string
		buffer   []byte
		expected string
		hasError bool
	}{
		{
			name:     "PDF file",
			buffer:   []byte{0x25, 0x50, 0x44, 0x46, 0x2D, 0x31, 0x2E, 0x35},
			expected: "pdf",
			hasError: false,
		},
		{
			name:     "JPEG file",
			buffer:   []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46},
			expected: "jpg",
			hasError: false,
		},
		{
			name:     "PNG file",
			buffer:   []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
			expected: "png",
			hasError: false,
		},
		{
			name:     "GIF file",
			buffer:   []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x01, 0x00},
			expected: "gif",
			hasError: false,
		},
		{
			name:     "BMP file",
			buffer:   []byte{0x42, 0x4D, 0x36, 0x04, 0x00, 0x00, 0x00, 0x00},
			expected: "bmp",
			hasError: false,
		},
		{
			name:     "WEBP file",
			buffer:   []byte{0x52, 0x49, 0x46, 0x46, 0x00, 0x00, 0x00, 0x00, 0x57, 0x45, 0x42, 0x50},
			expected: "webp",
			hasError: false,
		},
		{
			name:     "Buffer too short",
			buffer:   []byte{0x25, 0x50, 0x44},
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GuessFileType(tt.buffer)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %s, got %s", tt.expected, result)
				}
			}
		})
	}

	file, _ := os.ReadFile("./R11120b9170d92712261keeerm7")
	result, err := GuessFileType(file)
	t.Logf("err:%v", err)
	t.Logf("result:%s", result)
}
