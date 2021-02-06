package main

import (
	"fmt"
	"runtime"
)

func printMem() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v KiB", bToKb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v KiB", bToKb(m.TotalAlloc))
	fmt.Printf("\tSys = %v KiB", bToKb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToKb(b uint64) uint64 {
	return b / 1024
}
