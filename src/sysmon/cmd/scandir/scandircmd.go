package scandir

import (
	"fmt"

	"github.com/EagleLizard/sysmon-go/src/lib/argv"
)

func ScanDirCmd(pargv argv.ParsedArgv) {
	fmt.Println("ScanDirCmd()")
	dirs := pargv.Args
	fmt.Println("Scanning:")
	for _, currDir := range dirs {
		/*
			TODO: make this async
		*/
		fmt.Printf("%s\n", currDir)

	}
}
