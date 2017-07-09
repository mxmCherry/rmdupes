package filededuper_test

import (
	. "github.com/mxmCherry/rmdupes/internal/filededuper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deduper", func() {
	var deduper func(string) error
	var dupes []string

	BeforeEach(func() {
		dupes = nil
		deduper = Deduper(func(dupePath string) error {
			dupes = append(dupes, dupePath)
			return nil
		})
	})

	It("should not report files with the same path", func() {
		Expect(deduper("testdata/content1.bin")).To(Succeed())
		Expect(deduper("testdata/content1.bin")).To(Succeed())
		Expect(dupes).To(BeEmpty())
	})

	It("should not report files with different content", func() {
		Expect(deduper("testdata/content1.bin")).To(Succeed())
		Expect(deduper("testdata/content2.bin")).To(Succeed())
		Expect(dupes).To(BeEmpty())
	})

	It("should report files with same content", func() {
		Expect(deduper("testdata/content1.bin")).To(Succeed())
		Expect(deduper("testdata/content1-copy.bin")).To(Succeed())

		Expect(dupes).To(HaveLen(1))
		Expect(dupes[0]).To(HaveSuffix("testdata/content1-copy.bin"))
	})

	It("should skip non-regular items (like dirs)", func() {
		Expect(deduper("testdata")).To(Succeed())
		Expect(deduper("testdata")).To(Succeed())
		Expect(dupes).To(BeEmpty())
	})
})
