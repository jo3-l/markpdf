package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jo3-l/markpdf/internal/bmtreesyntax"
	"github.com/jo3-l/markpdf/internal/easypdf"
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

	bookmarks, err := bmtreesyntax.ParseReader(bookmarkFile)
	if err != nil {
		log.Fatalln("could not parse bookmarks:", err)
	}

	pdfFile, err := os.Open(flag.Arg(1))
	if err != nil {
		log.Fatalln(err)
	}
	defer pdfFile.Close()

	pdfInfo, err := easypdf.ExtractPDFInfo(pdfFile)
	if err != nil {
		log.Fatalln("could not extract pdf info:", err)
	}

	errs := RunChecks(AllChecks, bookmarks, pdfInfo)
	if len(errs) > 0 {
		log.Fatalf("bookmarks failed %d checks:\n%s\n", len(errs), errors.Join(errs...))
	}

	fmt.Printf("adding %d bookmarks to PDF ... ", bookmarks.Count())
	if err = easypdf.EditBookmarks(pdfFile, flag.Arg(2), *replaceFlag, bookmarks); err != nil {
		fmt.Println()
		log.Fatalln("could not edit bookmarks:", err)
	}
	fmt.Println("ok")
}
