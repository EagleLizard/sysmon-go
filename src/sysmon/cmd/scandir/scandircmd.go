package scandir

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/EagleLizard/sysmon-go/src/lib/argv"
	"github.com/EagleLizard/sysmon-go/src/sysmon/cmd/scandir/finddupes"
	"github.com/EagleLizard/sysmon-go/src/sysmon/cmd/scandir/scandiropts"
	"github.com/EagleLizard/sysmon-go/src/sysmon/cmd/scandir/scandirutil"
	"github.com/EagleLizard/sysmon-go/src/util/chron"
	"github.com/EagleLizard/sysmon-go/src/util/clicolors"
)

const progressMod = 1e4

func ScanDirCmd(pargv argv.ParsedArgv) {
	fmt.Println("ScanDirCmd()")

	scanDirOutDir := initScanDir()

	sdOpts := scandiropts.GetScanDirOpts(pargv)

	filesFileName := "0_files.txt"
	filesPath := filepath.Join(scanDirOutDir, filesFileName)
	dirsFileName := "0_dirs.txt"
	dirsPath := filepath.Join(scanDirOutDir, dirsFileName)

	filesWriter, err := os.Create(filesPath)
	if err != nil {
		log.Fatal(err)
	}
	defer filesWriter.Close()
	dirsWriter, err := os.Create(dirsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer dirsWriter.Close()

	pathCount := 0

	scanDirCb := func(params ScanDirCbParams) int {
		if params.IsSymLink {
			return 0
		}
		exclude := false
		if len(sdOpts.Exclude) > 0 {
			for _, excludeDirPath := range sdOpts.Exclude {
				if strings.Contains(params.FullPath, excludeDirPath) {
					exclude = true
				}
			}
		}
		if params.IsDir {
			if exclude {
				return 1
			}
			dirsWriter.Write([]byte(fmt.Sprintf("%s\n", params.FullPath)))
		} else {
			if params.Stats == nil {
				log.Fatalf("No stats for file: %s", params.FullPath)
			}
			if exclude {
				return 0
			}
			filesWriter.Write([]byte(fmt.Sprintf("%d %s\n", params.Stats.Size(), params.FullPath)))
		}
		pathCount++
		if pathCount%progressMod == 0 {
			fmt.Print(".")
		}
		return 0
	}

	dirs := pargv.Args
	fmt.Println("Scanning:")
	sw := chron.Start()
	totalFileCount := 0
	totalDirCount := 0
	for _, currDir := range dirs {
		/*
			TODO: make this async
		*/
		fmt.Printf("%s\n", currDir)
		sdRes := ScanDir(currDir, scanDirCb)
		fmt.Print("\n")
		totalFileCount += sdRes.FileCount
		totalDirCount += sdRes.DirCount
	}
	elapsed := sw.Stop()
	fmt.Printf("# files: %d\n", totalFileCount)
	fmt.Printf("# dirs: %d\n", totalDirCount)
	fmt.Printf("Scan took: %s\n", elapsed)
	if sdOpts.FindDuplicates {
		fdSw := chron.Start()
		finddupes.FindDupes(filesPath)
		fdElapsed := fdSw.Stop()
		fmt.Printf("findDupes() took: %s\n", clicolors.Chartreuse_light(fdElapsed))
	}
}

func initScanDir() string {
	// outDataDirPath := filepath.Join(constants.BaseDir(), constants.OutDataDirName)
	// scanDirOutDirPath := filepath.Join(outDataDirPath, constants.ScanDirOutDirName)
	scanDirOutDirPath := scandirutil.GetScanDirOutDirPath()
	// os.Mkdir(outDataDirPath, 0755)
	os.MkdirAll(scanDirOutDirPath, 0755)
	return scanDirOutDirPath
}
