# `markpdf`

markpdf is a command-line tool to add bookmarks to PDFs. Thanks to the human
friendly format used to specify bookmarks, invoking markpdf can often be more
convenient than equivalent workflows using applications such as Adobe Acrobat.

Below is an example of the bookmark format; a
[complete explanation is available later in this document](#bookmark-format).

```
Cover, 1
Contents, 2
Topic 1. Stoichiometric Relationships, 3
    1.1: Particulate nature of matter, 3
    1.2: Molar volume of a gas and calculations, 7
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
<title>, <page_number>
```

For example, `Introduction, 1` specifies a bookmark titled `Introduction`
pointing to page 1.

Bookmarks may be arbitrarily nested using indentation to create a hierarchy. For
example:

```
Data structures, 1
    Stacks and queues, 1
        Deques, 3
    Circular lists, 5
```

corresponds to a bookmark hierarchy with one top-level bookmark, "Data structures",
containing two immediate children, "Stacks and queues" and "Circular lists". In
addition, "Stacks and queues"
has a child called "Deques".

### Custom starting page numbers

When adding bookmarks according to a table of contents, it is often convenient
to set a custom starting page, such that `1` refers not to page 1 of the PDF
but rather to, say, page 7, the first page of content. To do this, use the `set`
command. For example:

```
Cover, 1
Table of Contents, 4

set 7
Mechanics, 1
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
