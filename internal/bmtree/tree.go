package bmtree

import "github.com/pdfcpu/pdfcpu/pkg/pdfcpu"

// A Tree represents a bookmark hierarchy.
type Tree struct{ toplevel []*Bookmark }

func (t *Tree) Count() int {
	n := 0
	t.Inspect(func(b *Bookmark) {
		n += 1
	})
	return n
}

// Inspect traverses the bookmark hierarchy in depth-first order. For each
// top-level bookmark b, Inspect calls f(b), then invokes f recursively for each
// of the non-nil children of b.
func (t *Tree) Inspect(f func(*Bookmark)) {
	var visitAll func([]*Bookmark)
	visitAll = func(bookmarks []*Bookmark) {
		for _, b := range bookmarks {
			f(b)
			visitAll(b.Children)
		}
	}

	visitAll(t.toplevel)
}

func (t *Tree) ToPdfCpu() []pdfcpu.Bookmark {
	var convertAll func([]*Bookmark) []pdfcpu.Bookmark
	convertAll = func(bookmarks []*Bookmark) []pdfcpu.Bookmark {
		if len(bookmarks) == 0 {
			// pdfcpu panics given a bookmark with empty but non-nil children
			// slice, so just always use nil. See
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

	return convertAll(t.toplevel)
}

type Bookmark struct {
	Page     int
	Title    string
	Children []*Bookmark
}
