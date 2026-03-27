//
// Simple Go application to read a Seduko defintion from file and write out the solution
// Note: If multiple solutions exist (they shouldn't!) only one will be found
//
// usage:  go-seduko <filename>
//
// Defaults to a sample seduko in the executable directory
//
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	//	maxThreads := flag.Int("t", 1, "maximum number of threads to use")
	flag.Parse()

	if flag.NArg() > 1 {
		flag.Usage()
		return
	}

	fileName := flag.Arg(0)
	if len(fileName) == 0 {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		fileName = filepath.Join(dir, "seduko1.csv")
	}
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	var puzzle seduko
	err = puzzle.init(scanner)
	if err != nil {
		log.Fatalf("Invalid seduko definition supplied: %v\n", err)
		return
	}
	puzzle.print(false)
	err = puzzle.solve()
	if err != nil {
		fmt.Printf("Error solving seduko: %v\n", err)
		return
	}
	puzzle.print(true)
}
