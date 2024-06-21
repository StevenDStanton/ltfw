package main

import (
	"flag"
	"log"
	"os"

	"github.com/StevenDStanton/ltfw/common"
)

var (
	inputFileName = flag.String("i", "", "Input File Name")
	developerFlag = flag.Bool("d", false, "In Development Flag to block accidental use")
)

const (
	tool    = "cut"
	version = "v0.0.1"
)

func main() {

	common.PrintVersion(tool, version)
	if !*developerFlag {
		log.Fatalln("This program is in development and not ready for usage")
	}

	readFile()
}

func readFile() {
	file, err := os.Open(*inputFileName)
	if err != nil {
		log.Fatalf("Unable to read file, %v", err)
	}
	defer file.Close()

	// scanner := bufio.NewScanner(file)
	// for scanner.Scan() {
	// 	chunk := scanner.Text()
	// }

}
