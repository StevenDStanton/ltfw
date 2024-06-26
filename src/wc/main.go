package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/StevenDStanton/cli-tools-for-windows/common"
)

type lineData struct {
	lineCount  int
	wordCount  int
	charCount  int
	byteCount  int
	maxLineLen int
	name       string
	err        string
}

type lineSize struct {
	lineCount  int
	wordCount  int
	charCount  int
	byteCount  int
	maxLineLen int
}

type flags struct {
	printBytes    bool
	printChars    bool
	printLines    bool
	maxLineLength bool
	printWords    bool
	helpFlag      bool
	versionFlag   bool
	debug         bool
}

var (
	fileNames       []string
	totalBytes      = 0
	totalChars      = 0
	totalLines      = 0
	totalWords      = 0
	totalMaxLineLen = 0
	cmdFlags        flags
)

const (
	tool    = "wc"
	version = "1.0.11"
)

func main() {

	parseArgs()

	if cmdFlags.debug {

		fmt.Printf(`Debugging
printBytes: %t
printChars: %t
printLines: %t
maxLineLength: %t
printWords: %t
helpFlag: %t
versionFlag: %t
IsAllFalse: %t

`, cmdFlags.printBytes, cmdFlags.printChars, cmdFlags.printLines, cmdFlags.maxLineLength, cmdFlags.printWords, cmdFlags.helpFlag, cmdFlags.versionFlag, allFlagsFalse())
	}

	if cmdFlags.versionFlag {
		versionInformation := common.PrintVersion(tool, version)
		fmt.Println(versionInformation)
		os.Exit(0)
	}

	if cmdFlags.helpFlag {
		printHelp()
	}

	lines := []lineData{}
	if len(fileNames) == 0 {
		lines = append(lines, parseUserInput())
	} else {
		for _, filename := range fileNames {
			if filename == "-" {
				lines = append(lines, parseUserInput())
				continue
			}
			lines = append(lines, readFile(filename))
		}
	}

	if len(fileNames) > 1 {
		lines = append(lines, lineData{totalLines, totalWords, totalChars, totalBytes, totalMaxLineLen, "total", ""})
	}

	lineSize := getMaxWidths(lines)

	printLinesToConsole(lines, lineSize)

}

func allFlagsFalse() bool {
	return !cmdFlags.printBytes &&
		!cmdFlags.printChars &&
		!cmdFlags.printLines &&
		!cmdFlags.maxLineLength &&
		!cmdFlags.printWords
}

func parseArgs() {
	//I am aware of the flags package. However as I am trying to replicate how wc works on linux it proved to limited for my needs.
	args := os.Args[1:]
	processingFlags := true
	nullSeparatesFileNames := false

	for _, arg := range args {

		if strings.HasPrefix(arg, "--files0-from=") {
			fmt.Println("Not Implemented due to Windows limitations")
			os.Exit(1)
		}

		if !processingFlags && !nullSeparatesFileNames {
			fileNames = append(fileNames, arg)
			continue
		}
		if arg == "--" {
			processingFlags = false
			continue
		}

		if strings.HasPrefix(arg, "--") {
			// Handle long options
			switch arg {
			case "--bytes":
				cmdFlags.printBytes = true
			case "--chars":
				cmdFlags.printChars = true
			case "--lines":
				cmdFlags.printLines = true
			case "--max-line-length":
				cmdFlags.maxLineLength = true
			case "--words":
				cmdFlags.printWords = true
			case "--help":
				cmdFlags.helpFlag = true
			case "--version":
				cmdFlags.versionFlag = true
			default:
				fmt.Printf("wc: unrecognized option %s\n", arg)
				fmt.Println("Try 'wc --help' for more information.")
				os.Exit(1)
			}

			continue
		}

		if strings.HasPrefix(arg, "-") && arg != "-" {
			for _, flag := range arg {
				switch flag {
				case '-':
					continue
				case 'c':
					cmdFlags.printBytes = true
				case 'm':
					cmdFlags.printChars = true
				case 'l':
					cmdFlags.printLines = true
				case 'L':
					cmdFlags.maxLineLength = true
				case 'w':
					cmdFlags.printWords = true
				case 'd':
					cmdFlags.debug = true
				default:
					fmt.Printf("wc: invalid option -- %s\n", string(flag))
					fmt.Println("Try 'wc --help' for more information.")
					os.Exit(1)
				}
			}
			continue
		}

		if !nullSeparatesFileNames {
			fileNames = append(fileNames, arg)
		}

	}

}

func readFile(fileName string) lineData {
	file, err := os.Open(fileName)
	fileData := lineData{0, 0, 0, 0, 0, fileName, ""}
	if err != nil {
		fileData.err = "No Such File or Directory;"
		return fileData
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(scanLinesWithNewlines)
	for scanner.Scan() {
		line := scanner.Bytes()
		parseText(&fileData, line)
	}

	totalBytes += fileData.byteCount
	totalChars += fileData.charCount
	totalLines += fileData.lineCount
	totalWords += fileData.wordCount
	if fileData.maxLineLen > totalMaxLineLen {
		totalMaxLineLen = fileData.maxLineLen
	}

	return fileData
}

func scanLinesWithNewlines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		return i + 1, data[:i+1], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

func parseUserInput() lineData {
	reader := bufio.NewReader(os.Stdin)
	userData := lineData{0, 0, 0, 0, 0, "-", ""}

	for {
		userInput, err := reader.ReadBytes('\x00')
		if err != nil {
			if err == io.EOF {
				break
			}
			userData.err = "Error Reading Standard Input"
			return userData
		}
		parseText(&userData, userInput)
	}

	totalBytes += userData.byteCount
	totalChars += userData.charCount
	totalLines += userData.lineCount
	totalWords += userData.wordCount
	if userData.maxLineLen > totalMaxLineLen {
		totalMaxLineLen = userData.maxLineLen
	}

	return userData
}

func parseText(fileData *lineData, line []byte) {
	//Updated to use built in functions instead of the loop previously used
	fileData.lineCount++
	fileData.byteCount += len(line)
	fileData.wordCount += len(bytes.Fields(line))
	fileData.charCount += utf8.RuneCount(line)

	if len(line) > fileData.maxLineLen {
		fileData.maxLineLen = len(line)
	}
}

func printLinesToConsole(lineData []lineData, lineSize lineSize) {

	for _, line := range lineData {
		if line.err != "" {
			fmt.Printf("wc: %s: %s\n", line.name, line.err)
			continue
		}
		if cmdFlags.printLines || allFlagsFalse() {
			fmt.Printf("%*d ", lineSize.lineCount, line.lineCount)
		}

		if cmdFlags.printWords || allFlagsFalse() {
			fmt.Printf("%*d ", lineSize.wordCount, line.wordCount)
		}

		if cmdFlags.printChars {
			fmt.Printf("%*d ", lineSize.charCount, line.charCount)
		}

		if cmdFlags.printBytes || allFlagsFalse() {
			fmt.Printf("%*d ", lineSize.byteCount, line.byteCount)
		}

		if cmdFlags.maxLineLength {
			fmt.Printf("%*d ", lineSize.maxLineLen, line.maxLineLen)
		}

		fmt.Printf("%s \n", line.name)
	}
}

func getMaxWidths(lines []lineData) lineSize {
	var maxWidths lineSize

	for _, line := range lines {
		maxWidths.lineCount = max(maxWidths.lineCount, len(fmt.Sprintf("%d", line.lineCount)))
		maxWidths.wordCount = max(maxWidths.wordCount, len(fmt.Sprintf("%d", line.wordCount)))
		maxWidths.charCount = max(maxWidths.charCount, len(fmt.Sprintf("%d", line.charCount)))
		maxWidths.byteCount = max(maxWidths.byteCount, len(fmt.Sprintf("%d", line.byteCount)))
		maxWidths.maxLineLen = max(maxWidths.maxLineLen, len(fmt.Sprintf("%d", line.maxLineLen)))
	}

	return maxWidths
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func printHelp() {
	help := `Usage: wc [OPTION]... [FILE]...

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
-d                     Enables Debugging

--files0-from=F        Has not been included in the Windows version due to issues with null terminators in Windows. 

--help                 display this help and exit
--version              output version information and exit`

	fmt.Println(help)
	os.Exit(0)
}
