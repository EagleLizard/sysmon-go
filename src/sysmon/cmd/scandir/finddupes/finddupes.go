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
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/EagleLizard/sysmon-go/src/sysmon/cmd/scandir/scandirutil"
	"github.com/EagleLizard/sysmon-go/src/util/chron"
)

const gfhCoef = float32(1) / 700

// const maxRunningHashFns = 64
var maxRunningHashFns int32

func init() {
	numCpu := runtime.NumCPU()
	fmt.Printf("NumCPU(): %d\n", numCpu)
	// maxRunningHashFns = int32(numCpu)
	maxRunningHashFns = 64
	// maxRunningHashFns = 256
	fmt.Printf("maxRunningHashFns: %d\n", maxRunningHashFns)
}

func FindDupes(filesDataFilePath string) {
	fmt.Printf("filesDataFilePath: %s\n", filesDataFilePath)
	sizeMap := getPossibleDupeSizes(filesDataFilePath)
	possibleDupeCount := 0
	for currKey := range sizeMap {
		currFileCount := sizeMap[currKey]
		possibleDupeCount += currFileCount
	}
	fmt.Printf("Possible dupes: %d\n", possibleDupeCount)
	sw := chron.Start()
	hashCountMap := getFileHashes(filesDataFilePath, sizeMap, possibleDupeCount)
	elapsed := sw.Stop()
	fmt.Printf("getFileHashes() took: %s\n", elapsed)
	fmt.Printf("hashCountMap size: %d\n", len(hashCountMap))
}

/*
	59768
*/

func getFileHashes(filesDataFilePath string, sizeMap map[int]int, possibleDupeCount int) map[string]int {
	gfhMod := gfhCoef * float32(possibleDupeCount)
	fmt.Printf("gfhMod: %v\n", gfhMod)
	filesDataFile, err := os.Open(filesDataFilePath)
	if err != nil {
		panic(err)
	}
	defer filesDataFile.Close()
	scanDirOutDir := scandirutil.GetScanDirOutDirPath()
	hashFileName := "0_hashes.txt"
	hashFilePath := filepath.Join(scanDirOutDir, hashFileName)

	var hashMu sync.Mutex
	var runningHashCount atomic.Int32
	w, err := os.Create(hashFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	lineCount := 0
	hashCountMap := make(map[string]int)

	sc := bufio.NewScanner(filesDataFile)

	gfhSw := chron.Start()

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
			for runningHashCount.Load() > maxRunningHashFns {
				time.Sleep(1 * time.Millisecond)
			}
			runningHashCount.Add(1)
			go func() {
				defer func() {
					runningHashCount.Add(-1)
					if float32(gfhSw.Current().Milliseconds()) > gfhMod {
						// fmt.Print(".")
						fmt.Print("⸱")
						gfhSw.Reset()
					}
				}()
				hashStr, err := getFileHashTrunc(currPath)
				if err != nil {
					if errors.Is(err, os.ErrPermission) || errors.Is(err, os.ErrNotExist) {
						return
					} else {
						panic(err)
					}
				}
				hashMu.Lock()
				hashCountMap[hashStr]++
				w.Write([]byte(fmt.Sprintf("%x %d %s\n", hashStr, size, currPath)))
				hashMu.Unlock()
			}()
		}
	}
	for runningHashCount.Load() > 0 {
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Print("\n")
	return hashCountMap
}

func getFileHashTrunc(filePath string) (string, error) {
	hashStr, err := getFileHash(filePath)
	if err != nil {
		return "", err
	}
	/*
	   approx. 1 collision every 1 trillion (1e12) documents
	     see: https://stackoverflow.com/a/22156338/4677252
	*/
	hashStr = hashStr[:5]
	return hashStr, nil
}

func getFileHash(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
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
	return string(hSum), nil
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
