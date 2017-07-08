/*

Command rmdupes removes duplicated files by content hash.

Usage:

	rmdupes --print-only --skip-hidden=false path/to/dir

*/
package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var flags struct {
	hash       string
	printOnly  bool
	skipHidden bool
}

func init() {
	flag.StringVar(&flags.hash, "hash", "sha512", "Hash to detect duplicated files (md5, sha1, sha256, sha512 - default)")
	flag.BoolVar(&flags.printOnly, "print-only", false, "Only print files instead of removing")
	flag.BoolVar(&flags.skipHidden, "skip-hidden", true, "Skip hidden files/directories (starting with dot)")
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

func run() error {
	hasher := hasherByName(flags.hash)
	if hasher == nil {
		return fmt.Errorf("unknown hash: " + flags.hash + ", supported are: md5, sha1, sha256, sha512")
	}

	walker := deduper(hasher, flags.printOnly)
	for _, path := range flag.Args() {
		path, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		if err := walkFiles(path, flags.skipHidden, walker); err != nil {
			return err
		}
	}
	return nil
}

func hasherByName(name string) hash.Hash {
	switch name {
	case "md5":
		return md5.New()
	case "sha1":
		return sha1.New()
	case "sha256":
		return sha256.New()
	case "sha512":
		return sha512.New()
	default:
		return nil
	}
}

func walkFiles(path string, skipHidden bool, cb func(path string) error) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if skipHidden && strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.Mode().IsRegular() {
			return cb(path)
		}
		return nil
	})
}

func deduper(hasher hash.Hash, printOnly bool) func(path string) error {
	seen := map[string]struct{}{}
	dupe := map[string]string{}

	hashBytes := make([]byte, 0, hasher.Size())

	return func(path string) error {
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

		remove := path
		if len(filepath.Base(original)) > len(filepath.Base(path)) {
			remove = original
			dupe[string(hashBytes)] = path
		}

		fmt.Printf("%s\n", remove)

		if printOnly {
			return nil
		}
		return os.Remove(remove)
	}
}
