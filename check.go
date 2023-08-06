package main

import (
	"errors"
	"fmt"

	"github.com/jo3-l/markpdf/internal/bmtree"
)

type Check interface {
	AppendErrors([]error, *bmtree.Tree) []error
}

var AllChecks = []Check{&NonEmptyCheck{}, &MonotonicallyIncreasingPageNumsCheck{}, &UniqueTitlesCheck{}}

func RunChecks(checks []Check, bookmarks *bmtree.Tree) []error {
	var errs []error
	for _, c := range checks {
		errs = c.AppendErrors(errs, bookmarks)
	}
	return errs
}

// pdfcpu doesn't support calling AddBookmarks with an empty slice; see
// https://github.com/pdfcpu/pdfcpu/blob/a9afcfe683880972fbbb576e12ef74688005ed3a/pkg/api/bookmark.go#L40.
type NonEmptyCheck struct{}

func (*NonEmptyCheck) AppendErrors(errs []error, bookmarks *bmtree.Tree) []error {
	if bookmarks.Count() == 0 {
		return append(errs, errors.New("no bookmarks specified"))
	}
	return errs
}

// pdfcpu requires page numbers of bookmarks to be montonically increasing; see
// https://github.com/pdfcpu/pdfcpu/issues/376.
type MonotonicallyIncreasingPageNumsCheck struct{}

func (*MonotonicallyIncreasingPageNumsCheck) AppendErrors(errs []error, bookmarks *bmtree.Tree) []error {
	var prev *bmtree.Bookmark
	bookmarks.Inspect(func(b *bmtree.Bookmark) {
		if prev != nil && b.Page < prev.Page {
			err := fmt.Errorf("bookmark %q (pg %d) appears after bookmark %q (pg %d) but has lower page number; this is not supported",
				b.Title, b.Page,
				prev.Title, prev.Page)
			errs = append(errs, err)
		}
		prev = b
	})
	return errs
}

// pdfcpu has a bug where AddBookmarks will panic given duplicate titles; see
// https://github.com/pdfcpu/pdfcpu/issues/664.
type UniqueTitlesCheck struct{}

func (*UniqueTitlesCheck) AppendErrors(errs []error, bookmarks *bmtree.Tree) []error {
	count := make(map[string]int)
	bookmarks.Inspect(func(b *bmtree.Bookmark) {
		count[b.Title] += 1
	})

	for title, occurrences := range count {
		if occurrences > 1 {
			errs = append(errs, fmt.Errorf("bookmark title %q is duplicated %d times; titles must be unique", title, occurrences))
		}
	}
	return errs
}
