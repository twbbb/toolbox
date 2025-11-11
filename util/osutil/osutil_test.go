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
			name:     "PDF file with UTF-8 BOM",
			buffer:   []byte{0xEF, 0xBB, 0xBF, 0x25, 0x50, 0x44, 0x46, 0x2D, 0x31, 0x2E, 0x35},
			expected: "pdf",
			hasError: false,
		},
		{
			name:     "PDF file with UTF-16 LE BOM",
			buffer:   []byte{0xFF, 0xFE, 0x25, 0x50, 0x44, 0x46, 0x2D, 0x31, 0x2E, 0x35},
			expected: "pdf",
			hasError: false,
		},
		{
			name:     "PDF file with UTF-16 BE BOM",
			buffer:   []byte{0xFE, 0xFF, 0x25, 0x50, 0x44, 0x46, 0x2D, 0x31, 0x2E, 0x35},
			expected: "pdf",
			hasError: false,
		},
		{
			name:     "PDF file with UTF-32 LE BOM",
			buffer:   []byte{0xFF, 0xFE, 0x00, 0x00, 0x25, 0x50, 0x44, 0x46, 0x2D, 0x31, 0x2E, 0x35},
			expected: "pdf",
			hasError: false,
		},
		{
			name:     "PDF file with UTF-32 BE BOM",
			buffer:   []byte{0x00, 0x00, 0xFE, 0xFF, 0x25, 0x50, 0x44, 0x46, 0x2D, 0x31, 0x2E, 0x35},
			expected: "pdf",
			hasError: false,
		},
		{
			name:     "PDF file with garbage data prefix",
			buffer:   []byte{0x75, 0xab, 0x5a, 0xfd, 0xaa, 0x69, 0x96, 0x27, 0x1a, 0xb6, 0x2a, 0x27, 0xfe, 0x97, 0x5f, 0xfd, 0xb6, 0xac, 0x7b, 0xae, 0x3f, 0x25, 0x50, 0x44, 0x46, 0x2D, 0x31, 0x2E, 0x33, 0x0d, 0x0a},
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
			name:     "Text file (UTF-8 without BOM)",
			buffer:   []byte("Hello, World! This is a text file.\nSecond line.\nThird line."),
			expected: "txt",
			hasError: false,
		},
		{
			name:     "JSON file",
			buffer:   []byte(`{"name": "test", "value": 123}`),
			expected: "json",
			hasError: false,
		},
		{
			name:     "XML file",
			buffer:   []byte(`<?xml version="1.0" encoding="UTF-8"?><root></root>`),
			expected: "xml",
			hasError: false,
		},
		{
			name:     "HTML file",
			buffer:   []byte(`<!DOCTYPE html><html><head></head><body></body></html>`),
			expected: "html",
			hasError: false,
		},
		{
			name:     "CSV file",
			buffer:   []byte("name,age,city\nJohn,25,New York\nJane,30,London"),
			expected: "csv",
			hasError: false,
		},
		{
			name:     "ZIP file",
			buffer:   []byte{0x50, 0x4B, 0x03, 0x04, 0x14, 0x00, 0x08, 0x00},
			expected: "zip",
			hasError: false,
		},
		{
			name:     "RAR file",
			buffer:   []byte{0x52, 0x61, 0x72, 0x21, 0x1A, 0x07, 0x00, 0x00},
			expected: "rar",
			hasError: false,
		},
		{
			name:     "7Z file",
			buffer:   []byte{0x37, 0x7A, 0xBC, 0xAF, 0x27, 0x1C, 0x00, 0x00},
			expected: "7z",
			hasError: false,
		},
		{
			name:     "GZ file",
			buffer:   []byte{0x1F, 0x8B, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00},
			expected: "gz",
			hasError: false,
		},
		{
			name:     "TIFF file (little endian)",
			buffer:   []byte{0x49, 0x49, 0x2A, 0x00, 0x08, 0x00, 0x00, 0x00},
			expected: "tiff",
			hasError: false,
		},
		{
			name:     "TIFF file (big endian)",
			buffer:   []byte{0x4D, 0x4D, 0x00, 0x2A, 0x00, 0x00, 0x00, 0x08},
			expected: "tiff",
			hasError: false,
		},
		{
			name:     "MP3 file (ID3v2)",
			buffer:   []byte("ID3\x03\x00\x00\x00\x00\x00"),
			expected: "mp3",
			hasError: false,
		},
		{
			name:     "WAV file",
			buffer:   []byte{0x52, 0x49, 0x46, 0x46, 0x00, 0x00, 0x00, 0x00, 0x57, 0x41, 0x56, 0x45},
			expected: "wav",
			hasError: false,
		},
		{
			name:     "MP4 file",
			buffer:   []byte{0x00, 0x00, 0x00, 0x18, 0x66, 0x74, 0x79, 0x70},
			expected: "mp4",
			hasError: false,
		},
		{
			name:     "AVI file",
			buffer:   []byte{0x52, 0x49, 0x46, 0x46, 0x00, 0x00, 0x00, 0x00, 0x41, 0x56, 0x49, 0x20},
			expected: "avi",
			hasError: false,
		},
		{
			name:     "Buffer too short",
			buffer:   []byte{0x25, 0x50, 0x44},
			expected: "",
			hasError: true,
		},
		{
			name:     "Buffer too short after BOM removal",
			buffer:   []byte{0xEF, 0xBB, 0xBF, 0x25},
			expected: "",
			hasError: true,
		},
		{
			name:     "Unknown file type",
			buffer:   []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
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
