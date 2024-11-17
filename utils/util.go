package utils

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
)

// DecompressRequest checks if the request is gzip-compressed and decompresses it.
// It returns the decompressed body as a byte slice or an error.
func DecompressRequest(r *http.Request) ([]byte, error) {
	var body []byte
	var err error

	// Check if the request is gzip-compressed
	if r.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			return nil, err
		}
		defer gz.Close()

		body, err = io.ReadAll(gz)
		if err != nil {
			return nil, err
		}
	} else {
		// Log a warning that the request is not gzip-compressed
		log.Println("Warning: Request is not gzip-compressed")

		// Handle plain body as fallback
		body, err = io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
	}

	return body, nil
}
