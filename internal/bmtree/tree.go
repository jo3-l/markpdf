package bmtree

// A Tree represents a bookmark hierarchy.
type Tree struct{ TopLevel []*Bookmark }

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

	visitAll(t.TopLevel)
}

type Bookmark struct {
	Page     int
	Title    string
	Children []*Bookmark
}
