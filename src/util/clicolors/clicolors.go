package clicolors

import "fmt"

type ColorFormatter func(val interface{}) string

var Yellow_light = Rgb(199, 196, 62)
var Yellow_yellow = Rgb(255, 255, 0)
var Chartreuse = Rgb(127, 255, 0)
var Chartreuse_light = Rgb(190, 255, 125)
var Peach = Rgb(255, 197, 109)
var Pink = Rgb(247, 173, 209)
var Cyan = Rgb(142, 250, 253)

func Rgb(r int, g int, b int) ColorFormatter {
	return func(val any) string {
		valStr := fmt.Sprintf("%v", val)
		return fmt.Sprintf("\x1B[38;2;%d;%d;%dm%s\x1B[39m", r, g, b, valStr)
	}
}
