package filededuper

import (
	"crypto/sha512"
	"io"
	"os"
	"path/filepath"
)

// Deduper returns deduplication closure function for files (SHA512 content hash).
func Deduper(cb func(dupePath string) error) func(path string) error {
	hasher := sha512.New()
	hashBytes := make([]byte, 0, hasher.Size())

	seen := map[string]struct{}{}
	dupe := map[string]string{}

	return func(path string) error {
		if !filepath.IsAbs(path) {
			var err error
			path, err = filepath.Abs(path)
			if err != nil {
				return err
			}
		}

		stats, err := os.Stat(path)
		if err != nil {
			return err
		}
		if !stats.Mode().IsRegular() {
			return nil
		}

		hasher.Reset()
		hashBytes = hashBytes[:0]

		if _, ok := seen[path]; ok {
			return nil
		}
		seen[path] = struct{}{}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(hasher, f); err != nil {
			return err
		}

		hashBytes = hasher.Sum(hashBytes)
		original, isDupe := dupe[string(hashBytes)]

		if !isDupe {
			dupe[string(hashBytes)] = path
			return nil
		}

		dupePath := path
		if len(filepath.Base(original)) > len(filepath.Base(path)) {
			dupePath = original
			dupe[string(hashBytes)] = path
		}
		return cb(dupePath)
	}
}
