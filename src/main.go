package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/EagleLizard/sysmon-go/src/constants"
	"github.com/EagleLizard/sysmon-go/src/lib/cmd/scandir"
)

func main() {
	etc := "hello"
	fmt.Printf("%s\n", etc)
	args := os.Args
	fmt.Printf("args: %s\n", strings.Join(args, ", "))
	wd := constants.BaseDir()
	dirArg := args[1]
	fmt.Printf("dirArg: %v", dirArg)
	var isRelative = true
	if strings.HasPrefix(dirArg, "/") {
		isRelative = false
	}
	baseDir := ""
	if isRelative {
		baseDir = wd
	}
	dir := filepath.Join(baseDir, dirArg)
	fmt.Printf("%s\n", dir)

	scanDirRes := scandir.ScanDir(dir)
	fmt.Printf("# files: %v"+"\n", scanDirRes.FileCount)
	fmt.Printf("# dirs: %v"+"\n", scanDirRes.DirCount)
	getPossibleDupes(scanDirRes.FilesFilePath)
}

func getPossibleDupes(filesPath string) string {
	f, err := os.OpenFile(filesPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	sizesFilePath := filepath.Join(scandir.GetScanDirOutDirPath(), "0_sizes.txt")
	sizesWriter, err := os.Create(sizesFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer sizesWriter.Close()
	fScanner := bufio.NewScanner(f)

	for fScanner.Scan() {
		line := fScanner.Text()
		fStat, err := os.Lstat(line)
		notExist := errors.Is(err, fs.ErrNotExist)
		if err != nil && !notExist {
			fmt.Printf("%+v", err)
			log.Fatal(err)
		}
		if !notExist {
			sizeLine := fmt.Sprint(fStat.Size()) + " " + line
			// fmt.Println(fmt.Sprint(fStat.Size()) + " " + line)
			sizesWriter.Write([]byte(sizeLine + "\n"))
		} else {
			log.Println("File not found: " + line)
		}
	}
	return sizesFilePath
}
