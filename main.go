/*

Command rmdupes removes duplicated files by content hash.

Usage:

	rmdupes             # get help
	rmdupes path/to/dir # remove duplicated files in path/to/dir

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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var flags struct {
	hash          string
	printOnly     bool
	processHidden bool
	recursive, r  bool
}

func init() {
	flag.StringVar(&flags.hash, "hash", "sha512", "Hash to detect duplicated files (md5, sha1, sha256, sha512)")
	flag.BoolVar(&flags.printOnly, "print", false, "Only print files instead of removing")
	flag.BoolVar(&flags.processHidden, "hidden", false, "Process hidden (starting with dot) files/directories")
	flag.BoolVar(&flags.recursive, "recursive", false, "Process files recursively (descend into inner dirs)")
	flag.BoolVar(&flags.r, "r", false, "Process files recursively (descend into inner dirs)")
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

func run() error {
	paths := flag.Args()
	if len(paths) == 0 {
		flag.PrintDefaults()
		return nil
	}

	hasher := hasherByName(flags.hash)
	if hasher == nil {
		return fmt.Errorf("unknown hash: " + flags.hash + "; try these: md5, sha1, sha256, sha512")
	}

	walker := walkFiles
	if flags.recursive || flags.r {
		walker = walkFilesRecursive
	}

	onFile := deduper(hasher, flags.printOnly)

	for _, path := range paths {
		path, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		if err := walker(path, flags.processHidden, onFile); err != nil {
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

func walkFilesRecursive(path string, processHidden bool, cb func(path string) error) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !processHidden && strings.HasPrefix(info.Name(), ".") {
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

func walkFiles(path string, processHidden bool, cb func(path string) error) error {
	stats, err := os.Stat(path)
	if err != nil {
		return err
	}

	if stats.Mode().IsRegular() {
		return cb(path)
	}
	if !stats.IsDir() {
		return nil
	}

	fis, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		if !fi.Mode().IsRegular() {
			continue
		}
		if !processHidden && strings.HasPrefix(fi.Name(), ".") {
			continue
		}

		fip := filepath.Join(path, fi.Name())
		if err := cb(fip); err != nil {
			return err
		}
	}
	return nil
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
