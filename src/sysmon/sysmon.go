package sysmon

import (
	"fmt"

	"github.com/EagleLizard/sysmon-go/src/lib/argv"
	"github.com/EagleLizard/sysmon-go/src/sysmon/cmd/gf"
	"github.com/EagleLizard/sysmon-go/src/sysmon/cmd/scandir"
)

func SysmonMain(parsedArgv argv.ParsedArgv) {
	fmt.Printf("cmd: %s\n", parsedArgv.Cmd)
	switch parsedArgv.Cmd {
	case "scandir":
		fallthrough
	case "sd":
		scandir.ScanDirCmd(parsedArgv)
	case "gf":
		gf.Gf()
	default:
		fmt.Printf("cmd not supported: \"%s\"\n", parsedArgv.Cmd)
	}
}
