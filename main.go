// nanotocsv is a program to convert ndj files (TSVs under the hood) to CSVs.
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func main() {

	const (
		NumberOfFields = 3 // Set number of fields expected in input file
		RowsToSkip     = 1 // Set number of rows at top of file to skip
	)

	// Check to make sure there's only input file
	if len(os.Args) != 2 {
		// Print usage if not
		fmt.Println("Usage: nanotocsv file.ndj")
		os.Exit(0)
	}

	// Open file for reading
	ndjFile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("nanotocsv: opening file: %s\n", err)
	}
	defer ndjFile.Close()

	// Create tsv reader
	r := csv.NewReader(ndjFile)
	r.Comma = '\t'
	r.FieldsPerRecord = NumberOfFields

	// Skip rows at top of file
	for i := 0; i < RowsToSkip; i++ {
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
	f, err := os.Create(os.Args[1] + ".csv")
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
