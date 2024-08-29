package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	etc := "hello"
	fmt.Printf("%s\n", etc)
	args := os.Args
	fmt.Printf("args: %s\n", strings.Join(args, ", "))
	parseArgs(os.Args)
}

type ArgParserState int

const (
	INIT ArgParserState = iota
	CMD
	FLAG
	ARG
)

const FLAG_ASSIGNMENT_DELIM = "="

func parseArgs(args []string) {
	// var cmd *string
	// parseState := INIT
	// pos := 0
	argParser := getArgParser(args)
	var nextRes *argParserNextRes
	for {
		nextRes = argParser()
		if nextRes == nil {
			break
		}
		fmt.Printf("val: %v\n", nextRes.Val)
	}
	for i, arg := range args {
		fmt.Printf("%v %v\n", i, arg)
	}
}

type argParserNextRes struct {
	Kind ArgParserState
	Val  string
}

func getArgParser(_args []string) func() *argParserNextRes {
	parseState := INIT
	pos := 0
	args := _args[1:]
	var next func() *argParserNextRes
	next = func() *argParserNextRes {
		if pos >= len(args) {
			return nil
		}
		currArg := args[pos]
		switch parseState {
		case INIT:
			if pos == 0 && isCmdStr(currArg) {
				parseState = CMD
			} else if isFlagArg(currArg) {
				parseState = FLAG
			} else {
				parseState = ARG
			}
		case CMD:
			pos++
			parseState = INIT
			return &argParserNextRes{
				parseState,
				currArg,
			}
		case FLAG:
			hasAssignment := isAssignment(currArg)
			if !hasAssignment {
				pos++
				parseState = INIT
				return &argParserNextRes{
					parseState,
					currArg,
				}
			}
			assignmentParts := strings.Split(currArg, FLAG_ASSIGNMENT_DELIM)
			if len(assignmentParts) != 2 {
				panic(fmt.Sprintf("Invalid flag assignment: %s", currArg))
			}
			lhs := assignmentParts[0]
			rhs := assignmentParts[1]
			args = args[:len(args)-1]
			args = append(args, lhs, rhs)
		case ARG:
			pos++
			parseState = INIT
			return &argParserNextRes{
				parseState,
				currArg,
			}
		}
		return next() // advance to next if we didn't return already
	}
	return next
}

func isAssignment(str string) bool {
	return strings.Contains(str, FLAG_ASSIGNMENT_DELIM)
}

var flagRx = regexp.MustCompile("^-{1,2}[a-zA-Z0-9][a-zA-Z0-9-]*=?")

func isFlagArg(str string) bool {
	/*
		-d
		--find-duplicates
		-ex
		--exclude
		-ex=etc
		-ex etc
		-ex etc1 etc2
	*/
	return flagRx.Match([]byte(str))
}

var cmdRx = regexp.MustCompile("^[a-z0-9]+(([a-z0-9]+|-)*[a-z0-9]+)?")

func isCmdStr(str string) bool {
	return cmdRx.Match([]byte(str))
}

// func getPossibleDupes(filesPath string) string {
// 	f, err := os.OpenFile(filesPath, os.O_RDONLY, os.ModePerm)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer f.Close()
// 	sizesFilePath := filepath.Join(scandir.GetScanDirOutDirPath(), "0_sizes.txt")
// 	sizesWriter, err := os.Create(sizesFilePath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer sizesWriter.Close()
// 	fScanner := bufio.NewScanner(f)

// 	for fScanner.Scan() {
// 		line := fScanner.Text()
// 		fStat, err := os.Lstat(line)
// 		notExist := errors.Is(err, fs.ErrNotExist)
// 		if err != nil && !notExist {
// 			fmt.Printf("%+v", err)
// 			log.Fatal(err)
// 		}
// 		if !notExist {
// 			sizeLine := fmt.Sprint(fStat.Size()) + " " + line
// 			// fmt.Println(fmt.Sprint(fStat.Size()) + " " + line)
// 			sizesWriter.Write([]byte(sizeLine + "\n"))
// 		} else {
// 			log.Println("File not found: " + line)
// 		}
// 	}
// 	return sizesFilePath
// }
