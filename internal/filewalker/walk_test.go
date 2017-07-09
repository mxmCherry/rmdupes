package filewalker_test

import (
	. "github.com/mxmCherry/rmdupes/internal/filewalker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Walk", func() {
	list := func(path string, flags Flags) (res []string) {
		err := Walk(path, flags, func(path string) error {
			res = append(res, path)
			return nil
		})
		Expect(err).NotTo(HaveOccurred())
		return res
	}

	It("should iterate files", func() {
		Expect(list("testdata", None)).To(ConsistOf(
			"testdata/file.bin",
		))
	})

	It("should recursively iterate files", func() {
		Expect(list("testdata", Recursive)).To(ConsistOf(
			"testdata/file.bin",
			"testdata/subdir/file.bin",
		))
	})

	It("should iterate hidden files", func() {
		Expect(list("testdata", Hidden)).To(ConsistOf(
			"testdata/file.bin",
			"testdata/.file.bin",
		))
	})

	It("should recursively iterate hidden files", func() {
		Expect(list("testdata", Recursive|Hidden)).To(ConsistOf(
			"testdata/file.bin",
			"testdata/.file.bin",
			"testdata/subdir/file.bin",
			"testdata/subdir/.hidden/.file.bin",
		))
	})
})
