package easypdf

import (
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

type BookmarkTree struct{ TopLevel []*Bookmark }

func (t *BookmarkTree) Count() int {
	n := 0
	t.Inspect(func(b *Bookmark) {
		n += 1
	})
	return n
}

// Inspect traverses the bookmark tree in depth-first order. For each top-level
// bookmark b, Inspect calls f(b), then invokes f recursively for each child of
// b.
func (t *BookmarkTree) Inspect(f func(*Bookmark)) {
	var visitAll func([]*Bookmark)
	visitAll = func(bookmarks []*Bookmark) {
		for _, b := range bookmarks {
			f(b)
			visitAll(b.Children)
		}
	}

	visitAll(t.TopLevel)
}

type Bookmark struct {
	Page     int
	Title    string
	Children []*Bookmark
}

func EditBookmarks(pdfFile *os.File, outputFilename string, replace bool, bookmarks *BookmarkTree) error {
	tmpFile, err := os.Create(pdfFile.Name() + ".tmp")
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	err = api.AddBookmarks(pdfFile, tmpFile, convertBookmarks(bookmarks), replace, nil)
	if err != nil {
		return err
	}

	return os.Rename(tmpFile.Name(), outputFilename)
}

func convertBookmarks(bookmarks *BookmarkTree) []pdfcpu.Bookmark {
	var convertAll func([]*Bookmark) []pdfcpu.Bookmark
	convertAll = func(bookmarks []*Bookmark) []pdfcpu.Bookmark {
		if len(bookmarks) == 0 {
			// HACK: pdfcpu panics given a bookmark with empty but non-nil
			// children slice, so ensure we always pass nil. See
			// https://github.com/pdfcpu/pdfcpu/issues/669.
			return nil
		}

		converted := make([]pdfcpu.Bookmark, len(bookmarks))
		for i, b := range bookmarks {
			converted[i] = pdfcpu.Bookmark{
				PageFrom: b.Page,
				Title:    b.Title,
				Children: convertAll(b.Children),
			}
		}
		return converted
	}

	return convertAll(bookmarks.TopLevel)
}
