package main

import (
	"bufio"
	"fmt"
	"limits/database"
	"limits/processor"
	"limits/processor/deposits"
	"log"
	"os"
)

func main() {
	// init
	db := database.New()
	DepositGetter := deposits.NewDepositGetter(db)
	depositsPutter := deposits.NewDepositPutter(db)

	processor := processor.New(
		db,
		DepositGetter,
		depositsPutter,
	)

	// open input and output files
	source, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer source.Close()

	dest, err := os.Create("./tmp/test_output.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer dest.Close()

	// process input file line by line
	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		res, err := processor.Process(scanner.Bytes())
		if err != nil {
			log.Fatal(err)
		}

		if res == nil {
			continue
		}

		// write result to output file
		_, err = fmt.Fprintln(dest, string(res))
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("All done!")
}
