# `markpdf`

markpdf is a command-line tool to add bookmarks to PDFs. Thanks to the human
friendly format used to specify bookmarks, invoking markpdf can often be more
convenient than equivalent workflows using applications such as Adobe Acrobat.

Below is an example of the bookmark format; a
[complete explanation is available later in this document](#bookmark-format).

```
p1. Cover
p2. Contents
p3. Topic 1. Stoichiometric Relationships
    p3. 1.1: Particulate nature of matter
    p7. 1.2: Molar volume of a gas and calculations
# ...
```

## Usage

Clone this repository and run `go build`. Then simply run the resulting binary,
passing in three files: your bookmarks, input PDF, and output.

```
usage: markpdf bookmarks.txt input.pdf output.pdf [-r]
  -r    replace existing bookmarks (default true)
```

## Bookmark Format

markpdf accepts a sequence of bookmarks separated by newlines. Lines that only
contain whitespace or begin with a `#` will be ignored. A single bookmark has
the following format:

```
p<page_number>. <title>
```

For example, `p1. Introduction` specifies a bookmark titled `Introduction`
pointing to page 1.

Bookmarks may be arbitrarily nested using indentation to create a hierarchy. For
example:

```
p1. Data structures
    p1. Stacks and queues
        p3. Deques
    p5. Circular lists
```

corresponds to a bookmark hierarchy with one top-level bookmark, "Data structures",
containing two immediate children, "Stacks and queues" and "Circular lists". In
addition, "Stacks and queues"
has a child called "Deques".

### Custom starting page numbers

When adding bookmarks according to a table of contents, it is often convenient
to set a custom starting page, such that `p1` refers not to page 1 of the PDF
but rather to, say, page 7, the first page of content. To do this, use the `set`
command. For example:

```
p1. Cover
p4. Table of Contents

set 7
p1. Mechanics
```

The above creates three bookmarks: one to page 1 of the PDF titled "Cover", one
to page 4 of the PDF titled "Table of Contents", and another to page 7 of the
PDF titled "Mechanics".

# Contributing

Contributions in the form of issue reports or pull requests are welcomed. Those
hacking on the markpdf codebase may find it informative to read through
[CONTRIBUTING.md](./CONTRIBUTING.md), which explains how the project is
structured.

# Author

markpdf is authored and maintained by [Joe L.](https://github.com/jo3-l/) and is
made available under the terms of the [MIT license](./LICENSE.md).
