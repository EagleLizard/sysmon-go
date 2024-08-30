package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/EagleLizard/sysmon-go/src/lib/argv"
)

func main() {
	etc := "hello"
	fmt.Printf("%s\n", etc)
	args := os.Args
	fmt.Printf("args: %s\n", strings.Join(args, ", "))
	parsedArgv := argv.ParseArgv(os.Args)
	fmt.Printf("cmd: %s\n", parsedArgv.Cmd)
	fmt.Printf("cmdArgs: %v\n", parsedArgv.Args)
	fmt.Print("opts:\n")
	for _, opt := range parsedArgv.Opts {
		fmt.Printf("%s = %v\n", opt.Flag, opt.FlagOpts)
		// if len(opt.FlagOpts) > 0 {
		// 	fmt.Printf("flagOpts: %v\n", opt.FlagOpts)
		// }
	}
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
