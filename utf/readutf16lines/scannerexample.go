package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// Creates a scanner similar to os.Open() but decodes the file as UTF-16.
// Useful when reading data from MS-Windows systems that generate UTF-16BE
// files, but will do the right thing if other BOMs are found.
func NewScannerUTF16(filename string) (io.Reader, error) {

	// Read the file into a []byte:
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// Make an tranformer that converts MS-Win default to UTF8:
	win16be := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	// Make a transformer that is like win16be, but abides by BOM:
	utf16bom := unicode.BOMOverride(win16be.NewDecoder())

	// This technique is recommended by the W3C for use in HTML 5:
	// "For compatibility with deployed content, the byte order
	// mark (also known as BOM) is considered more authoritative
	// than anything else." http://www.w3.org/TR/encoding/#specification-hooks

	// Make a Reader that uses utf16bom:
	unicodeReader := transform.NewReader(file, utf16bom)
	return unicodeReader, nil
}

func main() {

	s, err := NewScannerUTF16("inputfile.txt")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(s)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // Println will add back the final '\n'
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading inputfile:", err)
	}

}