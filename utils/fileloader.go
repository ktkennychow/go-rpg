package utils

// LoadFile loads a file using the platform-specific implementation
func LoadFile(path string) ([]byte, error) {
	return loadFileNative(path)
}
