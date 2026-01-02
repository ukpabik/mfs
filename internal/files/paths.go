package files

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"path/filepath"
	"strings"
)

const HashSegmentSize = 3
const NumSegments = 3

// TransformPath converts a user-provided path into a hashed directory structure.
// This distributes files across multiple directories to avoid single large directories.
//
// Input: "myfile.txt"
// Output: "abc/def/ghi/myfile.txt" (where abc/def/ghi come from sha256(input))
//
// Returns an error if the path is empty, invalid, or contains path traversal attempts.
func TransformPath(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", errors.New("empty path")
	}

	base := filepath.Base(path)
	if base == "." || base == ".." || base == "" {
		return "", errors.New("invalid file name")
	}

	if strings.ContainsRune(base, filepath.Separator) || strings.Contains(base, "/") || strings.Contains(base, "\\") {
		return "", errors.New("invalid file name")
	}

	sum := sha256.Sum256([]byte(path))
	hexSum := hex.EncodeToString(sum[:])

	segments := make([]string, 0, NumSegments+1)
	for i := range NumSegments {
		start := i * HashSegmentSize
		end := start + HashSegmentSize
		segments = append(segments, hexSum[start:end])
	}
	segments = append(segments, base)

	rel := filepath.Join(segments...)

	if filepath.IsAbs(rel) {
		return "", errors.New("invalid transformed path")
	}
	clean := filepath.Clean(rel)
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", errors.New("invalid transformed path")
	}

	return clean, nil
}
