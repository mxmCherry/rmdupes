package filewalker

import (
	"os"
	"path/filepath"
	"strings"
)

// Flags defines Walk behaviour.
type Flags uint8

const (
	// None defines default Walk behaviour.
	None Flags = 0
	// Recursive makes Walk to descend into inner dirs.
	Recursive = 1 << (iota - 1)
	// Hidden makes Walk to process hidden files/dirs.
	Hidden
)

// Walk iterates all the regular files in path.
// Path can be file itself.
// If Recursive flag passed, it will descend into inner dirs.
// If Hidden flag passed, it will not skip hidden (starting with dot) files/dirs.
func Walk(path string, flags Flags, cb func(path string) error) error {
	noRecurse := flags&Recursive == 0
	skipHidden := flags&Hidden == 0

	metDir := false
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if metDir && (noRecurse || skipHidden && strings.HasPrefix(info.Name(), ".")) {
				return filepath.SkipDir
			}
			metDir = true
			return nil
		}

		if !info.Mode().IsRegular() || skipHidden && strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		return cb(path)
	})
}
