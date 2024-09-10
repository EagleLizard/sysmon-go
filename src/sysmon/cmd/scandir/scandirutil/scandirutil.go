package scandirutil

import (
	"path/filepath"

	"github.com/EagleLizard/sysmon-go/src/constants"
)

const TmpDirName = "fd_tmp"

func GetScanDirOutDirPath() string {
	outDataDirPath := filepath.Join(constants.BaseDir(), constants.OutDataDirName)
	scanDirOutDirPath := filepath.Join(outDataDirPath, constants.ScanDirOutDirName)
	return scanDirOutDirPath
}
