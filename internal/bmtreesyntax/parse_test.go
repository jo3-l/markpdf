package bmtreesyntax

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/jo3-l/markpdf/internal/easypdf"
)

var updateFlag = flag.Bool("update", false, "update or create snapshots")

func TestParseOK(t *testing.T) {
	paths, err := filepath.Glob(filepath.Join("testdata", "ok", "*.input"))
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range paths {
		_, filename := filepath.Split(path)
		testname := strings.TrimSuffix(filename, filepath.Ext(filename))

		t.Run(testname, func(t *testing.T) {
			f, err := os.Open(path)
			if err != nil {
				t.Fatal("error reading input file:", err)
			}
			defer f.Close()

			bookmarks, err := ParseReader(f)
			if err != nil {
				t.Fatal("unexpected error:", err)
			}
			got := dumpBookmarks(bookmarks)

			matchSnapshot(t, filepath.Join("testdata", "ok", testname+".snap"), got)
		})
	}
}

func TestParseErr(t *testing.T) {
	paths, err := filepath.Glob(filepath.Join("testdata", "err", "*.input"))
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range paths {
		_, filename := filepath.Split(path)
		testname := strings.TrimSuffix(filename, filepath.Ext(filename))

		t.Run(testname, func(t *testing.T) {
			f, err := os.Open(path)
			if err != nil {
				t.Fatal("error reading input file:", err)
			}
			defer f.Close()

			_, err = ParseReader(f)
			if err == nil {
				t.Fatal("expected error, got none")
			}

			var perr ParseError
			if !errors.As(err, &perr) {
				t.Fatal("expected parse error")
			}

			matchSnapshot(t, filepath.Join("testdata", "err", testname+".snap"), perr.Error())
		})
	}
}

func matchSnapshot(t *testing.T, snapshotFile string, got string) {
	t.Helper()

	want, err := os.ReadFile(snapshotFile)
	if err != nil {
		if *updateFlag && os.IsNotExist(err) {
			err = os.WriteFile(snapshotFile, []byte(got), 0644)
			if err != nil {
				t.Fatal("error writing snapshot:", err)
			}
			t.Fatalf("created snapshot:\n%s\n", got)
		}

		t.Fatal("error reading snapshot:", err)
	}

	if normalizeLineEndings(got) != normalizeLineEndings(string(want)) {
		if *updateFlag {
			err = os.WriteFile(snapshotFile, []byte(got), 0644)
			if err != nil {
				t.Fatal("error updating snapshot:", err)
			}
			t.Fatalf("updated snapshot.\nprevious:\n%s\n\nnew:\n%s\n", want, got)
		}

		t.Fatalf("snapshot not equal.\ngot:\n%s\n\nwant:\n%s\n", got, want)
	}
}

func normalizeLineEndings(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

func dumpBookmarks(bookmarks *easypdf.BookmarkTree) string {
	var sb strings.Builder
	for i, b := range bookmarks.TopLevel {
		if i > 0 {
			sb.WriteByte('\n')
		}
		dumpBookmark(&sb, b, 0)
	}
	return sb.String()
}

func dumpBookmark(sb *strings.Builder, b *easypdf.Bookmark, depth int) {
	for i := 0; i < depth; i++ {
		sb.WriteByte('\t')
	}

	sb.WriteString("- ")
	sb.WriteString(strconv.Quote(b.Title))
	sb.WriteString(" @ p")
	sb.WriteString(strconv.Itoa(b.Page))
	for _, c := range b.Children {
		sb.WriteByte('\n')
		dumpBookmark(sb, c, depth+1)
	}
}
