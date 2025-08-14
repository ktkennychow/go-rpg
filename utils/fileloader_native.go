//go:build !wasm
// +build !wasm

package utils

import "os"

func loadFileNative(path string) ([]byte, error) {
	return os.ReadFile(path)
}

