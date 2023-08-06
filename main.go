package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jo3-l/markpdf/internal/bmtree"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: markpdf bookmarks.txt input.pdf output.pdf [-r]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

var replaceFlag = flag.Bool("r", true, "replace existing bookmarks")

func main() {
	log.SetFlags(0)
	log.SetPrefix("markpdf: ")

	flag.Usage = usage
	flag.Parse()
	if flag.NArg() < 3 {
		log.Print("not enough arguments")
		fmt.Fprintln(os.Stderr)
		usage()
	}

	bookmarkFile, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatalln(err)
	}
	defer bookmarkFile.Close()

	bookmarks, err := bmtree.ParseReader(bookmarkFile)
	if err != nil {
		log.Fatalln("could not parse bookmarks:", err)
	}

	errs := RunChecks(AllChecks, bookmarks)
	if len(errs) > 0 {
		log.Fatalf("bookmarks failed %d checks:\n%s", len(errs), errors.Join(errs...))
	}

	fmt.Printf("adding %d bookmarks to PDF ... ", bookmarks.Count())
	err = editBookmarks(flag.Arg(1), flag.Arg(2), *replaceFlag, bookmarks.ToPdfCpu())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("ok")
}

func editBookmarks(inputFilename, outputFilename string, replace bool, bookmarks []pdfcpu.Bookmark) error {
	inFile, err := os.Open(inputFilename)
	if err != nil {
		return err
	}
	defer inFile.Close()

	tmpFile, err := os.Create(inputFilename + ".tmp")
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	err = api.AddBookmarks(inFile, tmpFile, bookmarks, replace, nil)
	if err != nil {
		return fmt.Errorf("could not bookmarks: %s", err)
	}

	return os.Rename(tmpFile.Name(), outputFilename)
}
