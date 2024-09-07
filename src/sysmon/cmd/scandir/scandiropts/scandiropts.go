package scandiropts

import (
	"fmt"

	"github.com/EagleLizard/sysmon-go/src/lib/argv"
)

type ScanDirOpts struct {
	FindDuplicates bool
	Exclude        []string
}

func GetScanDirOpts(parsedArgv argv.ParsedArgv) ScanDirOpts {
	sdOpts := ScanDirOpts{}
	for _, argvOpt := range parsedArgv.Opts {
		fmt.Printf("%+v\n", argvOpt)
		switch argvOpt.Flag {
		case "--find-duplicates":
			fallthrough
		case "-d":
			if sdOpts.FindDuplicates {
				panic("FindDuplicates flag already set")
			}
			sdOpts.FindDuplicates = true
		case "--exclude":
			fallthrough
		case "-ex":
			if sdOpts.Exclude != nil {
				panic("Exclude flag already set")
			}
			excludeArr := []string{}
			excludeArr = append(excludeArr, argvOpt.FlagOpts...)
			sdOpts.Exclude = excludeArr
		}
	}
	return sdOpts
}
