package finddupes

import (
	"bufio"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/EagleLizard/sysmon-go/src/sysmon/cmd/scandir/finddupes/hashinfo"
	"github.com/EagleLizard/sysmon-go/src/sysmon/cmd/scandir/scandirutil"
	"github.com/EagleLizard/sysmon-go/src/util/chron"
	"github.com/EagleLizard/sysmon-go/src/util/clicolors"
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
	gfhRes := getFileHashes(filesDataFilePath, sizeMap, possibleDupeCount)
	elapsed := sw.Stop()
	fmt.Printf("getFileHashes() took: %s\n", clicolors.Chartreuse_light(elapsed))
	fmt.Printf("hashCountMap size: %d\n", len(gfhRes.hashCountMap))

	gfdSw := chron.Start()
	gfdRes := getFileDupes(gfhRes.hashFilePath, gfhRes.hashCountMap)
	gfdElapsed := gfdSw.Stop()
	fmt.Printf("getFileDupes() took: %s\n", clicolors.Peach(gfdElapsed))
	fmt.Printf("totalDupeCount: %s\n", clicolors.Yellow_light(gfdRes.totalDupeCount))

	sortSw := chron.Start()
	sortDuplicates(gfdRes.dupesFilePath, gfdRes.totalDupeCount)
	sortElapsed := sortSw.Stop()
	fmt.Printf("sortDuplicates() took: %v\n", clicolors.Pink(sortElapsed))
}

func sortDuplicates(dupeFilePath string, totalDupeCount int) {
	tmpDirPath := filepath.Join(
		scandirutil.GetScanDirOutDirPath(),
		scandirutil.TmpDirName,
	)
	err := os.RemoveAll(tmpDirPath)
	if err != nil {
		log.Fatal(err)
	}
	os.MkdirAll(tmpDirPath, 0755)
	wtSw := chron.Start()
	writeTmpDupeSortChunks(dupeFilePath, tmpDirPath, totalDupeCount)
	wtElapsed := wtSw.Stop()
	fmt.Printf("writeTmpDupeSortChunks() took: %v\n", clicolors.Cyan(wtElapsed))

	stSw := chron.Start()
	sortTmpDupChunks(tmpDirPath)
	stElapsed := stSw.Stop()
	fmt.Printf("sortTmpDupChunks() took: %v\n", clicolors.Cyan(stElapsed))
}

func sortTmpDupChunks(tmpDirPath string) {
	dupesFmtFileName := "z1_dupes_fmt.txt"
	dupesFmtFilePath := filepath.Join(
		scandirutil.GetScanDirOutDirPath(),
		dupesFmtFileName,
	)

	dirEntries, err := os.ReadDir(tmpDirPath)
	if err != nil {
		log.Fatal(err)
	}
	tmpFilePaths := []string{}
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			panic(fmt.Sprintf("Unexpected DirEntry in sd tmp dir: %s", dirEntry.Name()))
		}
		tmpFileRx := regexp.MustCompile(`^[0-9]+\.txt$`)
		if !tmpFileRx.Match([]byte(dirEntry.Name())) {
			panic(fmt.Sprintf("sd tmp File with invalid name: %s", dirEntry.Name()))
		}
		tmpFilePaths = append(tmpFilePaths, filepath.Join(tmpDirPath, dirEntry.Name()))
	}

	sortFileCounter := 0

	for len(tmpFilePaths) > 1 {
		faPath := tmpFilePaths[0]
		fbPath := tmpFilePaths[1]
		tmpFilePaths = tmpFilePaths[2:]
		isLastSort := len(tmpFilePaths) == 0
		sortFileName := fmt.Sprintf("a_%d.txt", sortFileCounter)
		sortFileCounter++

		var sortFilePath string
		if isLastSort {
			sortFilePath = dupesFmtFilePath
		} else {
			sortFilePath = filepath.Join(tmpDirPath, sortFileName)
		}

		fa, err := os.Open(faPath)
		if err != nil {
			log.Fatal(err)
		}
		fb, err := os.Open(fbPath)
		if err != nil {
			log.Fatal(err)
		}
		w, err := os.Create(sortFilePath)
		if err != nil {
			log.Fatal(err)
		}
		scA := bufio.NewScanner(fa)
		scB := bufio.NewScanner(fb)
		scARes := scA.Scan()
		scBRes := scB.Scan()

		for scARes || scBRes {
			var aHashInfo hashinfo.FileHashInfo
			var bHashInfo hashinfo.FileHashInfo

			writeA := false
			writeB := false

			if scARes {
				lineA := scA.Text()
				aHashInfo = hashinfo.ParseHashInfo(lineA)
			}
			if scBRes {
				lineB := scB.Text()
				bHashInfo = hashinfo.ParseHashInfo(lineB)
			}
			if aHashInfo.Size > bHashInfo.Size {
				writeA = true
			} else if aHashInfo.Size < bHashInfo.Size {
				writeB = true
			} else {
				if aHashInfo.Size > 0 && bHashInfo.Size > 0 {
					if aHashInfo.Hash > bHashInfo.Hash {
						writeA = true
					} else if aHashInfo.Hash < bHashInfo.Hash {
						writeB = true
					} else {
						writeA = true
						writeB = true
					}
				} else if aHashInfo.Size > 0 {
					writeA = true
				} else if bHashInfo.Size > 0 {
					writeB = true
				}
			}
			if writeA {
				w.Write([]byte(fmt.Sprintf("%s %d %s\n", aHashInfo.Hash, aHashInfo.Size, aHashInfo.FilePath)))
				scARes = scA.Scan()
			}
			if writeB {
				w.Write([]byte(fmt.Sprintf("%s %d %s\n", bHashInfo.Hash, bHashInfo.Size, bHashInfo.FilePath)))
				scBRes = scB.Scan()
			}
		}
		fa.Close()
		fb.Close()
		w.Close()
		err = os.Remove(faPath)
		if err != nil {
			log.Fatal(err)
		}
		err = os.Remove(fbPath)
		if err != nil {
			log.Fatal(err)
		}
		tmpFilePaths = append(tmpFilePaths, sortFilePath)
	}
}

func writeTmpDupeSortChunks(dupeFilePath string, tmpDirPath string, totalDupeCount int) {
	chunkFileSize := 250

	dupesFile, err := os.Open(dupeFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer dupesFile.Close()

	currDupeLines := []string{}
	tmpFileCounter := 0

	// sw := chron.Start()
	// pctSw := chron.Start()

	_writeTmpFile := func() {
		type lineSizeRec struct {
			hash string
			size int
			line string
		}
		tmpFileName := fmt.Sprintf("%d.txt", tmpFileCounter)
		tmpFilePath := filepath.Join(tmpDirPath, tmpFileName)
		tmpFileCounter++
		lineSizeRecs := []lineSizeRec{}
		for _, line := range currDupeLines {
			lineRx := regexp.MustCompile("^(?P<fileHash>[a-f0-9]+) (?P<fileSize>[0-9]+) .*$")
			rxMatch := lineRx.FindStringSubmatch(line)
			rxRes := make(map[string]string)
			for i, name := range lineRx.SubexpNames() {
				if i != 0 && name != "" {
					rxRes[name] = rxMatch[i]
				}
			}
			fileHash := rxRes["fileHash"]
			fileSizeStr := rxRes["fileSize"]
			fileSize, err := strconv.Atoi(fileSizeStr)
			if err != nil {
				log.Fatal(err)
			}
			lineSizeRecs = append(lineSizeRecs, lineSizeRec{
				hash: fileHash,
				size: fileSize,
				line: line,
			})
		}
		sort.SliceStable(lineSizeRecs, func(i int, j int) bool {
			return lineSizeRecs[i].size > lineSizeRecs[j].size
		})
		w, err := os.Create(tmpFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer w.Close()
		for _, currRec := range lineSizeRecs {
			w.Write([]byte(fmt.Sprintf("%s\n", currRec.line)))
		}
		currDupeLines = []string{}
	}

	sc := bufio.NewScanner(dupesFile)
	for sc.Scan() {
		line := sc.Text()
		currDupeLines = append(currDupeLines, line)
		if len(currDupeLines) >= chunkFileSize {
			/*
				write to tmp file and clear currDupeLines
			*/
			_writeTmpFile()
		}
	}
	if len(currDupeLines) > 0 {
		_writeTmpFile()
	}
}

type getFileDupesRes struct {
	dupesFilePath  string
	totalDupeCount int
}

func getFileDupes(hashFilePath string, hashCountMap map[string]int) getFileDupesRes {
	hashFile, err := os.Open(hashFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer hashFile.Close()

	dupesFileName := "z1_dupes.txt"
	dupesFilePath := filepath.Join(
		scandirutil.GetScanDirOutDirPath(),
		dupesFileName,
	)
	w, err := os.Create(dupesFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	totalDupeCount := 0

	sc := bufio.NewScanner(hashFile)
	for sc.Scan() {
		line := sc.Text()
		lineRx := regexp.MustCompile("^(?P<fileHash>[a-f0-9]+) [0-9]+ .*$")
		rxMatch := lineRx.FindStringSubmatch(line)
		rxRes := make(map[string]string)
		for i, name := range lineRx.SubexpNames() {
			if i != 0 && name != "" {
				rxRes[name] = rxMatch[i]
			}
		}
		hashCount := hashCountMap[rxRes["fileHash"]]
		if hashCount > 1 {
			totalDupeCount += hashCount
			w.Write([]byte(fmt.Sprintf("%s\n", line)))
		}
	}
	return getFileDupesRes{
		dupesFilePath:  dupesFilePath,
		totalDupeCount: totalDupeCount,
	}
}

type getFileHashesRes struct {
	hashFilePath string
	hashCountMap map[string]int
}

func getFileHashes(filesDataFilePath string, sizeMap map[int]int, possibleDupeCount int) getFileHashesRes {
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
	finishedHashCount := 0

	sc := bufio.NewScanner(filesDataFile)

	gfhSw := chron.Start()
	pctSw := chron.Start()

	var gfhProgressMu sync.Mutex
	gfhProgressFn := func() {
		/*
			need a mutex otherwise the threads sometimes print
				all at once when the hash write mutex unlocks
		*/
		gfhProgressMu.Lock()
		defer gfhProgressMu.Unlock()

		if float32(gfhSw.Current().Milliseconds()) > gfhMod {
			// fmt.Print(".")
			fmt.Print("⸱")
			gfhSw.Reset()
		}
		if float32(pctSw.Current().Milliseconds()) > (gfhMod * 8) {
			currPct := int(math.Round((float64(finishedHashCount) / float64(possibleDupeCount)) * 100))
			fmt.Printf("%d", currPct)
			pctSw.Reset()
		}
	}

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
		if size > 0 && sizeMap[size] > 1 {
			for runningHashCount.Load() > maxRunningHashFns {
				time.Sleep(1 * time.Millisecond)
			}
			runningHashCount.Add(1)
			go func() {
				defer runningHashCount.Add(-1)

				hashStr, err := getFileHashTrunc(currPath)
				if err != nil {
					if errors.Is(err, os.ErrPermission) || errors.Is(err, os.ErrNotExist) {
						return
					} else {
						panic(err)
					}
				}
				hashMu.Lock()
				// fmt.Printf("%x\n", hashStr)
				hashCountMap[hashStr]++
				finishedHashCount++
				w.Write([]byte(fmt.Sprintf("%s %d %s\n", hashStr, size, currPath)))
				hashMu.Unlock()
				gfhProgressFn()
			}()
		}
	}
	for runningHashCount.Load() > 0 {
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Print("\n")
	return getFileHashesRes{
		hashFilePath: hashFilePath,
		hashCountMap: hashCountMap,
	}
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
	hashStr = fmt.Sprintf("%x", hashStr[:5])
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
