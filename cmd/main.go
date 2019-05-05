package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	gokmv "github.com/positiveblue/go-kmv"
)

func getScanner(fileName string) *bufio.Scanner {
	if fileName != "" {
		f, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}
		return bufio.NewScanner(f)
	} else {
		return bufio.NewScanner(os.Stdin)
	}
}

func fmtMessage(estimator *gokmv.KMV) string {
	distinct := estimator.EstimateCardinality()
	total := estimator.ElementsAdded()
	size := estimator.Size()
	return fmt.Sprintf("%d %d %d", distinct, total, size)
}

func main() {
	sizePtr := flag.Int("size", 64, "initial size for the kmv data structure")
	fileNamePtr := flag.String("filename", "", "File name to process (otherwhise will read from StdIn")

	flag.Parse()

	scanner := getScanner(*fileNamePtr)
	kmv := gokmv.NewKMV(*sizePtr)
	for scanner.Scan() {
		for _, word := range strings.Fields(scanner.Text()) {
			kmv.InsertString(word)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	} else {
		fmt.Println(fmtMessage(kmv))
	}
}
