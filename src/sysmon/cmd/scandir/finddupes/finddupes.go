package finddupes

import (
	"bufio"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/EagleLizard/sysmon-go/src/sysmon/cmd/scandir/scandirutil"
)

func FindDupes(filesDataFilePath string) {
	fmt.Printf("filesDataFilePath: %s\n", filesDataFilePath)
	sizeMap := getPossibleDupeSizes(filesDataFilePath)
	possibleDupeCount := 0
	for currKey := range sizeMap {
		currFileCount := sizeMap[currKey]
		possibleDupeCount += currFileCount
	}
	fmt.Printf("Possible dupes: %d\n", possibleDupeCount)
	// hashCountMap := getFileHashes(filesDataFilePath, sizeMap)
	// fmt.Printf("hashCountMap size: %d", len(hashCountMap))
}

func getFileHashes(filesDataFilePath string, sizeMap map[int]int) map[string]int {
	filesDataFile, err := os.Open(filesDataFilePath)
	if err != nil {
		panic(err)
	}
	defer filesDataFile.Close()
	scanDirOutDir := scandirutil.GetScanDirOutDirPath()
	hashFileName := "0_hashes.txt"
	hashFilePath := filepath.Join(scanDirOutDir, hashFileName)
	w, err := os.Create(hashFilePath)
	if err != nil {
		log.Fatal(err)
	}

	lineCount := 0
	hashCountMap := make(map[string]int)

	sc := bufio.NewScanner(filesDataFile)
	for sc.Scan() {
		line := sc.Text()
		// fmt.Println(line)
		lineCount++
		delimIdx := strings.Index(line, " ")
		if delimIdx == -1 {
			log.Fatalf("Invalid entry at line %d", lineCount)
		}
		sizeStr := line[:delimIdx]
		currPath := line[delimIdx+1:]
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			log.Fatalf("Invalid size string on line %d:\n%s", lineCount, line)
		}
		if sizeMap[size] > 1 {
			// hash the file
			f, err := os.Open(currPath)
			if err != nil {
				if errors.Is(err, os.ErrPermission) || errors.Is(err, os.ErrNotExist) {
					fmt.Println(err)
					continue
				} else {
					panic(err)
				}
			}
			defer f.Close()
			h := sha1.New()
			fr := bufio.NewReader(f)
			buf := make([]byte, 1*1024)
			for {
				n, err := fr.Read(buf)
				if err != nil {
					if !errors.Is(err, io.EOF) {
						log.Fatal(err)
					}
					break
				}
				if n != 0 {
					h.Write(buf[:n])
				}
			}
			hSum := h.Sum(nil)
			hashCountMap[string(hSum)]++
			w.Write([]byte(fmt.Sprintf("%x %d %s\n", hSum, size, currPath)))
		}
	}
	return hashCountMap
}

func getPossibleDupeSizes(filesDataFilePath string) map[int]int {
	f, err := os.Open(filesDataFilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	lineCount := 0
	sizeMap := make(map[int]int)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		lineCount++
		// fmt.Println(line)
		delimIdx := strings.Index(line, " ")
		if delimIdx == -1 {
			log.Fatalf("Invalid entry at line %d", lineCount)
		}
		sizeStr := line[:delimIdx]
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			log.Fatalf("Invalid size string on line %d:\n%s", lineCount, line)
		}
		sizeMap[size]++
	}
	for key := range sizeMap {
		if sizeMap[key] < 2 {
			delete(sizeMap, key)
		}
	}
	return sizeMap
}
