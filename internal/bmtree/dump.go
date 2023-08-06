package bmtree

import (
	"strconv"
	"strings"
)

func dump(bookmarks *Tree) string {
	var sb strings.Builder
	for i, b := range bookmarks.toplevel {
		if i > 0 {
			sb.WriteByte('\n')
		}
		dumpBookmarkTo(&sb, b, 0)
	}
	return sb.String()
}

func dumpBookmarkTo(sb *strings.Builder, b *Bookmark, depth int) {
	for i := 0; i < depth; i++ {
		sb.WriteByte('\t')
	}

	sb.WriteString("- ")
	sb.WriteString(strconv.Quote(b.Title))
	sb.WriteString(" @ p")
	sb.WriteString(strconv.Itoa(b.Page))
	for _, c := range b.Children {
		sb.WriteByte('\n')
		dumpBookmarkTo(sb, c, depth+1)
	}
}
