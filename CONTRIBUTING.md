# Project Structure

`markpdf` consists of three main components:

- [`./internal/bmtreesyntax`](./internal/bmtreesyntax), a parser that converts textual representations of bookmarks — the syntax of which is explained in [README.md](./README.md) — to structured `BookmarkTree`s;
- [`./internal/easypdf`](./internal/easypdf/), a facade around the relevant parts of the underlying PDF used (https://github.com/pdfcpu/pdfcpu);
- and the command-line interface ([`./main.go`](./main.go)) which glues everything together: it parses the input bookmarks using `bmtreesyntax`, checks for validity — page numbers in correct range and so on ([`./check.go`](./main.go)), and finally creates an edited version of the PDF using `easypdf`.

## Testing

The only component of `markpdf` that is unit tested currently is `bmtreesyntax`, which uses a form of snapshot testing (also known as golden testing.) Specifically, each input file in [`./internal/bmtreesyntax/testdata/ok`](./internal/bmtreesyntax/testdata/ok) is parsed and the resulting `BookmarkTree` structure serialized. The test passes if and only if this serialized output matches the reference output stored alongside the input file (the snapshot.) Similar can be said for [`./internal/bmtreesyntax/testdata/err`](./internal/bmtreesyntax/testdata/err), the difference being that the files in that directory are checked for errors, as the directory name suggests.

If you have made a change to the code that affects the outputs of some of the tests, but you are confident that the new output is correct — for example, perhaps you have fixed a bug — you can forcibly update all snapshots so that they match the current output using

```
go test -update
```

Adding a new test is equally simple: simply add a new input file in `testdata/ok`, run `go test -update` to generate a snapshot of the output, and ensure that everything looks correct.

---

See https://eli.thegreenplace.net/2022/file-driven-testing-in-go/ for a more detailed explanation to the topic of golden/snapshot testing.
