package argv

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/EagleLizard/sysmon-go/src/lib/argv/argvparser"
	"github.com/EagleLizard/sysmon-go/src/lib/argv/argvtoken"
)

const FLAG_ASSIGNMENT_DELIM = "="

type ArgvFlag struct {
	Flag     string
	FlagOpts []string
}

type ParsedArgv struct {
	Cmd  string
	Args []string
	Opts []ArgvFlag
}

func ParseArgv(args []string) ParsedArgv {
	var cmdToken *argvtoken.ArgvToken
	cmdArgs := []string{}
	flags := []ArgvFlag{}
	tokenStack := []argvtoken.ArgvToken{}

	argParser := getArgParser(args)
	var currToken argvtoken.ArgvToken

	consumeCmdOrFlag := func() {
		var token *argvtoken.ArgvToken
		argTokens := []argvtoken.ArgvToken{}
		for len(tokenStack) > 0 {
			token = &tokenStack[len(tokenStack)-1]
			tokenStack = tokenStack[:len(tokenStack)-1]
			if token.Kind == argvtoken.ARG {
				argTokens = append(argTokens, *token)
			}
		}
		if token == nil {
			return
		}
		if token.Kind == argvtoken.CMD {
			if cmdToken != nil {
				panic("cmd token encountered, but cmdToken already set")
			}
			cmdToken = &argvtoken.ArgvToken{
				Kind: token.Kind,
				Val:  token.Val,
			}
			for len(argTokens) > 0 {
				argToken := argTokens[len(argTokens)-1]
				argTokens = argTokens[:len(argTokens)-1]
				cmdArgs = append(cmdArgs, argToken.Val)
			}
		} else {
			if cmdToken == nil {
				panic(fmt.Sprintf("Unexpected Token: cmd not set. Expected %v, received: %v", argvtoken.CMD, token.Kind))
			}
			flagOpts := []string{}
			for len(argTokens) > 0 {
				argToken := argTokens[len(argTokens)-1]
				argTokens = argTokens[:len(argTokens)-1]
				flagOpts = append(flagOpts, argToken.Val)
			}
			flags = append(flags, ArgvFlag{
				token.Val,
				flagOpts,
			})
		}
	}

	for {
		currToken = argParser()
		if currToken.Kind != argvtoken.ARG {
			/*
				consume current stack
			*/
			consumeCmdOrFlag()
		}
		tokenStack = append(tokenStack, currToken)
		if currToken.Kind == argvtoken.END {
			break
		}
	}
	res := ParsedArgv{
		Cmd:  cmdToken.Val,
		Args: cmdArgs,
		Opts: flags,
	}
	return res
}

func getArgParser(_args []string) func() argvtoken.ArgvToken {
	parseState := argvparser.INIT
	pos := 0
	args := _args[1:]
	var next func() argvtoken.ArgvToken
	next = func() argvtoken.ArgvToken {
		if pos >= len(args) {
			return argvtoken.ArgvToken{
				Kind: argvtoken.END,
				Val:  "",
			}
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
			pos++
			parseState = argvparser.INIT
			return argvtoken.ArgvToken{
				Kind: argvtoken.CMD,
				Val:  currArg,
			}
		case argvparser.FLAG:
			hasAssignment := isAssignment(currArg)
			if !hasAssignment {
				pos++
				parseState = argvparser.INIT
				return argvtoken.ArgvToken{
					Kind: argvtoken.FLAG,
					Val:  currArg,
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
		case argvparser.ARG:
			pos++
			parseState = argvparser.INIT
			return argvtoken.ArgvToken{
				Kind: argvtoken.ARG,
				Val:  currArg,
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
