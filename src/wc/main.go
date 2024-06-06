package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	printBytes         = flag.Bool("bytes", false, "print the byte counts")
	printBytesShort    = flag.Bool("c", false, "print the byte counts")
	printChars         = flag.Bool("chars", false, "print the character counts")
	printCharsShort    = flag.Bool("m", false, "print the character counts")
	printLines         = flag.Bool("lines", false, "print the newline counts")
	printLinesShort    = flag.Bool("l", false, "print the newline counts")
	maxLineLength      = flag.Bool("max-line-length", false, "print the length of the longest line")
	maxLineLengthShort = flag.Bool("L", false, "print the length of the longest line")
	printWords         = flag.Bool("words", false, "print the word counts")
	printWordsShort    = flag.Bool("w", false, "print the word counts")
	helpFlag           = flag.Bool("help", false, "display this help and exit")
	versionFlag        = flag.Bool("version", false, "output version information and exit")
	aboutFlag          = flag.Bool("about", false, "display information about the program and exit")
)

func main() {
	const version = "0.5.2"
	flag.Parse()
	args := flag.Args()

	if *versionFlag {
		fmt.Printf("go Version %s\nCopyright 2024 The Simple Dev\nLicense MIT - No Warranty\n\nWritten By Steven Stanton\nReverse Engineered by RTFM", version)
		os.Exit(0)
	}

	if *helpFlag {
		help := `Usage: wc [OPTION]... [FILE]...
Multiple Files: wc [OPTION]... --files0-from=F

Prints new line, word, and byte counts for each FILE and a total line if more file is specified.
A word is a non-zero length sequence of characters delimited by white space.

When no FILE or when FILE is -, read from standard input.

The options below may be used to select which counts are displayed. Always in the following order:
newLine, word, character, byte, maximum line length.

-c, --bytes            print the byte counts
-m, --chars            print the character counts
-l, --lines            print the newline counts
-L, --max-line-length  print the length of the longest line
-w, --words            print the word counts

--files0-from=F        read input from the files specified by NUL-terminated names in file F; If F is - then read names from standard input

--help                 display this help and exit
--version              output version information and exit
--about 			  display information about the program and exit`

		fmt.Println(help)
		os.Exit(0)
	}

	if *aboutFlag {
		about := `This is a simple implementation of the wc command in Go.

This program has been reversed engineered from the GNU Coreutils wc program using only documentation and observed behavior in a clean room environment.

Author:         Steven Stanton
License:        MIT - No Warranty
Author Github:  https//github.com/StevenDStanton
Project Github: https://github.com/StevemStanton/ltfw

Part of my Linux Tools for Windows (ltfw) project.
`
		fmt.Println(about)
		os.Exit(0)
	}

	filename := ""
	fileData := []byte{}

	if len(args) > 0 {
		filename = args[0]
	}

	if filename != "" {
		var err error
		fileData, err = os.ReadFile(filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	charCount := 0
	newLineCount := 0
	maxLineLen := 0
	wordCount := 0

	stringData := string(fileData)
	inWord := false
	currentLineLen := 0
	for _, char := range stringData {
		charCount++

		if char == '\n' {
			newLineCount++
			if currentLineLen > maxLineLen {
				maxLineLen = currentLineLen
			}
			currentLineLen = 0
		} else {
			currentLineLen++
		}

		if inWord && (char == ' ' || char == '\n' || char == '\r' || char == '\t') {
			inWord = false
		}

		if !inWord && (char != ' ' && char != '\n' && char != '\r' && char != '\t') {
			inWord = true
			wordCount++
		}
	}

	if *printLines || *printLinesShort || allFlagsFalse() {
		fmt.Printf("%d ", newLineCount)
	}

	if *printWords || *printWordsShort || allFlagsFalse() {
		fmt.Printf("%d ", wordCount)
	}

	if *printChars || *printCharsShort {
		fmt.Printf("%d ", charCount)
	}

	if *printBytes || *printBytesShort || allFlagsFalse() {
		fmt.Printf("%d ", len(fileData))
	}

	if *maxLineLength || *maxLineLengthShort {
		fmt.Printf("%d ", maxLineLen)
	}

	if filename != "" || allFlagsFalse() {
		fmt.Printf("%s ", filename)
	}

}

func allFlagsFalse() bool {
	return !*printBytes && !*printBytesShort &&
		!*printChars && !*printCharsShort &&
		!*printLines && !*printLinesShort &&
		!*maxLineLength && !*maxLineLengthShort &&
		!*printWords && !*printWordsShort &&
		!*helpFlag && !*versionFlag && !*aboutFlag
}
