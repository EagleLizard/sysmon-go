package gf

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func Gf() {
	fmt.Println("gf")

	// fmt.Print("goFuncs:\n")
	// goFuncs()
	// fmt.Print("\ngoFuncs2:\n")
	// goFuncs2()
	// fmt.Print("\ngoFuncs3:\n")
	// goFuncs3()
	fmt.Print("\ngoFuncs4:\n")
	goFuncs4()
}

func goFuncs4() {
	testStr := "_a_lot_of_llamas_loll_lazily_"
	strRepN := 10
	str := strings.Repeat(testStr, strRepN)
	charCounts := countCharsSync(str)
	fmt.Printf("\"%s\":\n", str)
	printChars(charCounts)
}

func countCharsSync(str string) map[rune]int {
	charCounts := make(map[rune]int)
	for _, c := range str {
		charCounts[c]++
	}
	return charCounts
}

func printChars(charCountMap map[rune]int) {
	sortedKeys := sortCharMapKeys(charCountMap)
	for _, k := range sortedKeys {
		fmt.Printf("'%s': %d\n", string(k), charCountMap[k])
	}
}

func sortCharMapKeys(charCountMap map[rune]int) []rune {
	type pair struct {
		Key rune
		Val int
	}
	var pairs []pair
	for k, v := range charCountMap {
		pairs = append(pairs, pair{k, v})
	}
	sort.SliceStable(pairs, func(i, j int) bool {
		return pairs[i].Val > pairs[j].Val
	})
	var sortedKeys []rune
	for _, p := range pairs {
		sortedKeys = append(sortedKeys, p.Key)
	}
	return sortedKeys
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
