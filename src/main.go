package main

import (
	"os"

	"github.com/EagleLizard/sysmon-go/src/lib/argv"
	"github.com/EagleLizard/sysmon-go/src/sysmon"
)

func main() {
	parsedArgv := argv.ParseArgv(os.Args)
	sysmon.SysmonMain(parsedArgv)
}

// func getPossibleDupes(filesPath string) string {
// 	f, err := os.OpenFile(filesPath, os.O_RDONLY, os.ModePerm)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer f.Close()
// 	sizesFilePath := filepath.Join(scandir.GetScanDirOutDirPath(), "0_sizes.txt")
// 	sizesWriter, err := os.Create(sizesFilePath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer sizesWriter.Close()
// 	fScanner := bufio.NewScanner(f)

// 	for fScanner.Scan() {
// 		line := fScanner.Text()
// 		fStat, err := os.Lstat(line)
// 		notExist := errors.Is(err, fs.ErrNotExist)
// 		if err != nil && !notExist {
// 			fmt.Printf("%+v", err)
// 			log.Fatal(err)
// 		}
// 		if !notExist {
// 			sizeLine := fmt.Sprint(fStat.Size()) + " " + line
// 			// fmt.Println(fmt.Sprint(fStat.Size()) + " " + line)
// 			sizesWriter.Write([]byte(sizeLine + "\n"))
// 		} else {
// 			log.Println("File not found: " + line)
// 		}
// 	}
// 	return sizesFilePath
// }
