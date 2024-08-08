package constants

import (
	"log"
	"path/filepath"
	"runtime"
)

const (
	OutDataDirName    = "out_data"
	ScanDirOutDirName = "scandir"
)

func BaseDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Error getting BaseDir")
	}
	baseDir := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	return baseDir
}
