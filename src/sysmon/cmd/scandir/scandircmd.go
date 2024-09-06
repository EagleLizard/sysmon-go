package scandir

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/EagleLizard/sysmon-go/src/lib/argv"
	"github.com/EagleLizard/sysmon-go/src/sysmon/cmd/scandir/finddupes"
	"github.com/EagleLizard/sysmon-go/src/sysmon/cmd/scandir/scandirutil"
	"github.com/EagleLizard/sysmon-go/src/util/chron"
)

const progressMod = 1e4

func ScanDirCmd(pargv argv.ParsedArgv) {
	fmt.Println("ScanDirCmd()")

	scanDirOutDir := initScanDir()
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

	scanDirCb := func(params ScanDirCbParams) {
		if params.IsSymLink {
			return
		}
		if params.IsDir {
			dirsWriter.Write([]byte(fmt.Sprintf("%s\n", params.FullPath)))
		} else {
			if params.Stats == nil {
				log.Fatalf("No stats for file: %s", params.FullPath)
			}
			filesWriter.Write([]byte(fmt.Sprintf("%d %s\n", params.Stats.Size(), params.FullPath)))
		}
		pathCount++
		if pathCount%progressMod == 0 {
			fmt.Print(".")
		}
	}

	dirs := pargv.Args
	fmt.Println("Scanning:")
	sw := chron.Start()
	for _, currDir := range dirs {
		/*
			TODO: make this async
		*/
		fmt.Printf("%s\n", currDir)
		ScanDir(currDir, scanDirCb)
		fmt.Print("\n")
	}
	elapsed := sw.Stop()
	fmt.Printf("Scan took: %s\n", elapsed)
	finddupes.FindDupes(filesPath)
}

func initScanDir() string {
	// outDataDirPath := filepath.Join(constants.BaseDir(), constants.OutDataDirName)
	// scanDirOutDirPath := filepath.Join(outDataDirPath, constants.ScanDirOutDirName)
	scanDirOutDirPath := scandirutil.GetScanDirOutDirPath()
	// os.Mkdir(outDataDirPath, 0755)
	os.MkdirAll(scanDirOutDirPath, 0755)
	return scanDirOutDirPath
}
