package bmtreesyntax

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/jo3-l/markpdf/internal/easypdf"
)

func ParseReader(r io.Reader) (*easypdf.BookmarkTree, error) {
	sc := bufio.NewScanner(r)
	state := newParseState()
	for sc.Scan() {
		if err := state.process(sc.Text()); err != nil {
			return nil, err
		}
		state.lineno++
	}

	if sc.Err() != nil {
		return nil, fmt.Errorf("io error: %w", sc.Err())
	}
	return state.bookmarks, nil
}

func newParseState() *parseState {
	return &parseState{lineno: 1, pageOffset: 1, bookmarks: &easypdf.BookmarkTree{}}
}

type parseState struct {
	lineno     int // one-based line number, for error reporting
	pageOffset int // page n refers to n+pageOffset-1 page of PDF

	bookmarks *easypdf.BookmarkTree
	parent    []*easypdf.Bookmark // stack of parents for nested bookmarks. invariant: len(parent) is the current nesting depth
	prev      *easypdf.Bookmark   // last bookmark processed
}

type ParseError struct {
	lineno int
	msg    string
}

func (p ParseError) Error() string {
	return fmt.Sprintf("line %d: %s", p.lineno, p.msg)
}

var bookmarkRe = regexp.MustCompile(`^(\s*)(.+),\s+(-?\d+)`)

func (p *parseState) process(line string) error {
	trim := strings.TrimSpace(line)

	if trim == "" || strings.HasPrefix(trim, "#") {
		// comment or blank line; skip it
		return nil
	}

	if s, ok := strings.CutPrefix(trim, "set "); ok {
		n, err := strconv.Atoi(s)
		if err != nil {
			return p.errorf("invalid page number %q passed to set command", s)
		}
		p.pageOffset = n
		return nil
	}

	// Don't use trimmed version of line; we need the leading indentation.
	m := bookmarkRe.FindStringSubmatch(line)
	if m == nil {
		return p.errorf("invalid bookmark %q (should start with title, then comma, then page number, e.g. 'Introduction, 1')", line)
	}

	page, err := strconv.Atoi(m[3])
	if err != nil {
		return p.errorf("invalid page number %q", m[3])
	}

	depth, err := p.indentDepth(m[1])
	if err != nil {
		return err
	}
	return p.insertBookmark(depth, p.pageOffset+page-1, m[2])
}

func (p *parseState) insertBookmark(depth int, page int, title string) error {
	if depth > len(p.parent) {
		// Current bookmark is nested; ensure there's a parent bookmark and add
		// it to the stack.
		if depth-len(p.parent) > 1 {
			// invalid:
			//   a
			//   .  .  b
			return p.errorf("bookmark nested too deep")
		} else if p.prev == nil {
			// invalid:
			//   .  a
			return p.errorf("invalid nested bookmark without parent")
		} else {
			// Consider:
			//   a
			//   .  b
			// The bookmark immediately preceding b (that is, a) is the parent
			// of b.
			p.parent = append(p.parent, p.prev)
		}
	} else if depth < len(p.parent) {
		// We're moving out of a nested position, so trim the parent stack.
		p.parent = p.parent[:depth]
	}

	// Attach the bookmark to its parent.
	b := &easypdf.Bookmark{Page: page, Title: title}
	if len(p.parent) > 0 {
		par := p.parent[len(p.parent)-1]
		par.Children = append(par.Children, b)
	} else {
		p.bookmarks.TopLevel = append(p.bookmarks.TopLevel, b)
	}

	p.prev = b
	return nil
}

const indentWidth = 4

func (p *parseState) indentDepth(indent string) (int, error) {
	depth := 0
	for indent != "" {
		switch indent[0] {
		case '\t':
			depth++
			indent = indent[1:]

		case ' ':
			numSpace := 0
			for indent != "" && indent[0] == ' ' {
				numSpace++
				indent = indent[1:]
			}

			if numSpace%indentWidth != 0 {
				return 0, p.errorf("invalid number of spaces in indentation (must be multiple of %d)", indentWidth)
			}
			depth += numSpace / indentWidth

		default:
			return 0, p.errorf("unrecognized space character %q in indentation", indent[0])
		}
	}

	return depth, nil
}

func (p *parseState) errorf(format string, a ...any) error {
	return ParseError{p.lineno, fmt.Sprintf(format, a...)}
}
