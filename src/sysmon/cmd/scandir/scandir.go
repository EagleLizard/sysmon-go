package scandir

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ScanDirRes struct {
	FileCount, DirCount int
}

type ScanDirCbParams struct {
	IsDir     bool
	IsSymLink bool
	FullPath  string
	Stats     os.FileInfo
}

/*
return:

	0 to continue
	1 to skip
*/
type ScanDirCb func(scanDirCbParams ScanDirCbParams) int

func ScanDir(dir string, scanDirCb ScanDirCb) ScanDirRes {

	dirFs := os.DirFS(dir)
	dirCount := 0
	fileCount := 0

	fs.WalkDir(dirFs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if !strings.Contains(err.Error(), "permission denied") {
				log.Fatal(err)
			}
		}
		fullPath := filepath.Join(dir, path)
		lstats, err := os.Lstat(fullPath)
		if err != nil {
			fmt.Println(fmt.Errorf("path: %s", fullPath))
			panic(err)
		}
		if d.IsDir() {
			dirCount++
		} else {
			fileCount++
		}
		scandDirCbRes := scanDirCb(ScanDirCbParams{
			IsDir:     d.IsDir(),
			IsSymLink: lstats.Mode()&fs.ModeSymlink != 0,
			FullPath:  fullPath,
			Stats:     lstats,
		})
		if scandDirCbRes == 1 && d.IsDir() {
			return filepath.SkipDir
		}
		return nil
	})
	res := ScanDirRes{
		FileCount: fileCount,
		DirCount:  dirCount,
	}
	return res
}
