package gf

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func Gf() {
	fmt.Println("gf")
	fmt.Print("goFuncs:\n")
	goFuncs()
	fmt.Print("\ngoFuncs2:\n")
	goFuncs2()
	fmt.Print("\ngoFuncs3:\n")
	goFuncs3()
}

func goFuncs3() {
	const testStr = "etc~"
	const maxRunningFns = len(testStr)

	rCh := make(chan rune)

	_gFn := func(c rune, ch chan rune) {
		// time.Sleep(time.Duration(c) * time.Millisecond)
		time.Sleep(1 * time.Millisecond)
		ch <- c
	}

	strRepN := 10
	str := strings.Repeat(testStr, strRepN)
	runningFns := 0

	for _, c := range str {
		runningFns++
		go _gFn(c, rCh)
	}
	for c := range rCh {
		fmt.Print(string(c))
		runningFns--
		if runningFns < 1 {
			close(rCh)
			break
		}
	}
}

func goFuncs2() {
	const testStr = "abc123"
	const maxRunningFns = len(testStr)

	var runningFns atomic.Int64
	var resStrMu sync.Mutex
	resStr := ""

	_gfFn := func(c rune) {
		resStrMu.Lock()
		defer resStrMu.Unlock()

		resStr += string(c)
		// time.Sleep(time.Duration(c * 1e3))
		fmt.Print(string(c))
	}

	strRepN := 10
	str := strings.Repeat(testStr, strRepN)

	for _, c := range str {
		if runningFns.Load() >= int64(maxRunningFns) {
			// time.Sleep(10 * 1e3)
			time.Sleep(100 * time.Millisecond)
		}
		runningFns.Add(1)
		go func() {
			defer func() {
				runningFns.Add(-1)
			}()
			_gfFn(c)
		}()
	}
	for runningFns.Load() > 0 {
		time.Sleep(10 * 1e6)
	}
	fmt.Print("\n")
	fmt.Println(resStr)
}

func goFuncs() {
	var wg sync.WaitGroup
	_print := func(c rune) {
		time.Sleep(time.Duration(c) * 1000)
		fmt.Print(string(c))
	}
	testStr := "etc_I_am_a_teapot"
	str := strings.Repeat(testStr, 10)
	for _, c := range str {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// time.Sleep(waitNs)
			_print(c)
		}()
		// time.Sleep(1e3)
	}
	wg.Wait()
}
