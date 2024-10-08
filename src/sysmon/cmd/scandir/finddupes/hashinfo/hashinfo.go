package hashinfo

import (
	"log"
	"regexp"
	"strconv"
)

type FileHashInfo struct {
	Hash     string
	Size     int
	FilePath string
}

var lineRx *regexp.Regexp
var regexpNames []string

func init() {
	lineRx = regexp.MustCompile("^(?P<fileHash>[a-f0-9]+) (?P<fileSize>[0-9]+) (?P<filePath>.*)$")
	regexpNames = lineRx.SubexpNames()
}

func ParseHashInfo(line string) FileHashInfo {

	// lineRx := regexp.MustCompile("^(?P<fileHash>[a-f0-9]+) (?P<fileSize>[0-9]+) (?P<filePath>.*)$")
	rxMatch := lineRx.FindStringSubmatch(line)
	rxRes := make(map[string]string)
	for i, name := range regexpNames {
		if i != 0 && name != "" {
			rxRes[name] = rxMatch[i]
		}
	}
	fileHash := rxRes["fileHash"]
	fileSizeStr := rxRes["fileSize"]
	filePath := rxRes["filePath"]

	fileSize, err := strconv.Atoi(fileSizeStr)
	if err != nil {
		log.Fatal(err)
	}
	res := FileHashInfo{
		Hash:     fileHash,
		Size:     fileSize,
		FilePath: filePath,
	}
	return res
}
