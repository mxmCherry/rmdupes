/*

Command rmdupes removes duplicated files by content hash.

Usage:

	rmdupes             # get help
	rmdupes path/to/dir # remove duplicated files in path/to/dir

*/
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mxmCherry/rmdupes/internal/filededuper"
	"github.com/mxmCherry/rmdupes/internal/filewalker"
)

var flags struct {
	hash          string
	printOnly     bool
	recursive, r  bool
	processHidden bool
}

func init() {
	flag.BoolVar(&flags.printOnly, "print", false, "Only print files instead of removing")
	flag.BoolVar(&flags.recursive, "recursive", false, "Process files recursively (descend into inner dirs)")
	flag.BoolVar(&flags.r, "r", false, "Process files recursively (descend into inner dirs)")
	flag.BoolVar(&flags.processHidden, "hidden", false, "Process hidden (starting with dot) files/directories")
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

	walkerFlags := filewalker.None
	if flags.recursive || flags.r {
		walkerFlags |= filewalker.Recursive
	}
	if flags.processHidden {
		walkerFlags |= filewalker.Hidden
	}

	deduper := filededuper.Deduper(func(dupePath string) error {
		fmt.Fprintln(os.Stdout, dupePath)
		if flags.printOnly {
			return nil
		}
		return os.Remove(dupePath)
	})

	for _, path := range paths {
		path, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		if err := filewalker.Walk(path, walkerFlags, deduper); err != nil {
			return err
		}
	}
	return nil
}
