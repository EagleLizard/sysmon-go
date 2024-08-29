package argv

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/EagleLizard/sysmon-go/src/lib/argv/argvparser"
	"github.com/EagleLizard/sysmon-go/src/lib/argv/argvtoken"
)

const FLAG_ASSIGNMENT_DELIM = "="

func ParseArgv(args []string) {
	tokenStack := []argvtoken.ArgvToken{}

	argParser := getArgParser(args)
	var currToken *argvtoken.ArgvToken
	for {
		currToken = argParser()
		if currToken == nil {
			break
		}
		// fmt.Printf("kind: %v\nval: %v\n\n", currToken.Kind, currToken.Val)
		if currToken.Kind != argvtoken.ARG {
			/*
				consume current stack
			*/
			tokenStack = tokenStack[:0]
		}
		tokenStack = append(tokenStack, *currToken)
		fmt.Printf("%v\n", tokenStack)
	}
	for i, arg := range args {
		fmt.Printf("%v %v\n", i, arg)
	}
}

func getArgParser(_args []string) func() *argvtoken.ArgvToken {
	parseState := argvparser.INIT
	pos := 0
	args := _args[1:]
	var next func() *argvtoken.ArgvToken
	next = func() *argvtoken.ArgvToken {
		if pos >= len(args) {
			return nil
		}
		currArg := args[pos]
		switch parseState {
		case argvparser.INIT:
			if pos == 0 && isCmdStr(currArg) {
				parseState = argvparser.CMD
			} else if isFlagArg(currArg) {
				parseState = argvparser.FLAG
			} else {
				parseState = argvparser.ARG
			}
		case argvparser.CMD:
			res := &argvtoken.ArgvToken{
				Kind: argvtoken.CMD,
				Val:  currArg,
			}
			pos++
			parseState = argvparser.INIT
			return res
		case argvparser.FLAG:
			hasAssignment := isAssignment(currArg)
			if !hasAssignment {
				res := &argvtoken.ArgvToken{
					Kind: argvtoken.FLAG,
					Val:  currArg,
				}
				pos++
				parseState = argvparser.INIT
				return res
			}
			assignmentParts := strings.Split(currArg, FLAG_ASSIGNMENT_DELIM)
			if len(assignmentParts) != 2 {
				panic(fmt.Sprintf("Invalid flag assignment: %s", currArg))
			}
			lhs := assignmentParts[0]
			rhs := assignmentParts[1]
			args = args[:len(args)-1]
			args = append(args, lhs, rhs)
		case argvparser.ARG:
			res := &argvtoken.ArgvToken{
				Kind: argvtoken.ARG,
				Val:  currArg,
			}
			pos++
			parseState = argvparser.INIT
			return res
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
