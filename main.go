// nanotocsv is a program to convert ndj files (TSVs under the hood) to CSVs.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var numberOfFields int
var rowsToSkip int
var help bool

func init() {
	const (
		fieldsDefault = 147
		fieldsUsage   = "the number of fields (columns) in the input file"
		skipDefault   = 4
		skipUsage     = "the number of rows at the top of the input file to skip"
		helpDefault   = false
		helpUsage     = "show this help message"
	)

	flag.IntVar(&numberOfFields, "fields", fieldsDefault, fieldsUsage)
	flag.IntVar(&numberOfFields, "f", fieldsDefault, fieldsUsage)
	flag.IntVar(&rowsToSkip, "skip", skipDefault, skipUsage)
	flag.IntVar(&rowsToSkip, "s", skipDefault, skipUsage)
	flag.BoolVar(&help, "help", helpDefault, helpUsage)
	flag.BoolVar(&help, "h", helpDefault, helpUsage)
}

func main() {

	usageMessage := "Usage: %s [options] file.ndj\nOptions:\n"

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usageMessage, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	// Check to make sure there's only input file
	if len(flag.Args()) != 1 {
		// Print usage if not
		fmt.Fprintf(flag.CommandLine.Output(), usageMessage, os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Open file for reading
	ndjFile, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Printf("nanotocsv: opening file: %s\n", err)
	}
	defer ndjFile.Close()

	// Create tsv reader
	r := csv.NewReader(ndjFile)
	r.Comma = '\t'
	r.FieldsPerRecord = numberOfFields

	// Skip rows at top of file
	for i := 0; i < rowsToSkip; i++ {
		_, err := r.Read()
		if err == io.EOF {
			break
		}
	}

	// Create variable to store data
	var records [][]string

	// Read in records
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
		records = append(records, record)
	}

	// Create file to write to
	outFile := strings.TrimSuffix(flag.Arg(0), ".ndj") + ".csv"
	f, err := os.Create(outFile)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	// Create csv writer
	w := csv.NewWriter(f)

	// Write data
	w.WriteAll(records)
	if err := w.Error(); err != nil {
		fmt.Println(err)
	}
}
