package scandir

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/EagleLizard/sysmon-go/src/constants"
)

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

type ScanDirRes struct {
	DirsFilePath, FilesFilePath string
	FileCount, DirCount         int
}

type ScanDirCbParams struct {
	IsDir     bool
	IsSymLink bool
	FullPath  string
	stats     *os.FileInfo
}

type ScanDirCb func(scanDirCbParams ScanDirCbParams)

func ScanDir(dir string, scanDirCb ScanDirCb) ScanDirRes {
	scanDirOutDir := initScanDir()
	filesFileName := "0_files.txt"
	filesPath := filepath.Join(scanDirOutDir, filesFileName)
	dirsFileName := "0_dirs.txt"
	dirsPath := filepath.Join(scanDirOutDir, dirsFileName)

	dirFs := os.DirFS(dir)
	dirCount := 0
	fileCount := 0

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

	fs.WalkDir(dirFs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if !strings.Contains(err.Error(), "permission denied") {
				log.Fatal(err)
			}
		}
		fullPath := filepath.Join(dir, path)
		if d.IsDir() {
			dirCount++
			dirsWriter.Write([]byte(fullPath + "\n"))
		} else {
			fileCount++
			filesWriter.Write([]byte(fullPath + "\n"))
		}
		// fmt.Println(path)
		return nil
	})
	res := ScanDirRes{
		DirsFilePath:  dirsPath,
		FilesFilePath: filesPath,
		FileCount:     fileCount,
		DirCount:      dirCount,
	}
	return res
}
