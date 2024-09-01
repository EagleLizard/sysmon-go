package scandir

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/EagleLizard/sysmon-go/src/constants"
	"github.com/EagleLizard/sysmon-go/src/lib/argv"
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
		lineBytes := []byte(fmt.Sprintf("%s\n", params.FullPath))
		if params.IsDir {
			dirsWriter.Write(lineBytes)
		} else {
			filesWriter.Write(lineBytes)
		}
		pathCount++
		if pathCount%progressMod == 0 {
			fmt.Print(".")
		}
	}

	dirs := pargv.Args
	fmt.Println("Scanning:")
	startTime := time.Now()
	for _, currDir := range dirs {
		/*
			TODO: make this async
		*/
		fmt.Printf("%s\n", currDir)
		ScanDir(currDir, scanDirCb)
		fmt.Print("\n")
	}
	endTime := time.Since(startTime)
	fmt.Printf("Scan took: %s\n", endTime)
}

func initScanDir() string {
	// outDataDirPath := filepath.Join(constants.BaseDir(), constants.OutDataDirName)
	// scanDirOutDirPath := filepath.Join(outDataDirPath, constants.ScanDirOutDirName)
	scanDirOutDirPath := GetScanDirOutDirPath()
	// os.Mkdir(outDataDirPath, 0755)
	os.MkdirAll(scanDirOutDirPath, 0755)
	return scanDirOutDirPath
}

func GetScanDirOutDirPath() string {
	outDataDirPath := filepath.Join(constants.BaseDir(), constants.OutDataDirName)
	scanDirOutDirPath := filepath.Join(outDataDirPath, constants.ScanDirOutDirName)
	return scanDirOutDirPath
}
