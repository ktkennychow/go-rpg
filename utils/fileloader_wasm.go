//go:build wasm
// +build wasm

package utils

import (
	"fmt"
	"io"
	"net/http"
)

func loadFileNative(path string) ([]byte, error) {
	// In WASM, we use HTTP requests instead of filesystem access
	resp, err := http.Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch %s: status %d", path, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body for %s: %w", path, err)
	}

	return data, nil
}
