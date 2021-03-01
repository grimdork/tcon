package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/gdamore/tcell"
	"github.com/grimdork/tcon"
)

func main() {
	var err error
	s, err := tcon.New()
	s.SetTitle("Demo app")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(2)
	}

	s.SetCtrlFunc(func(k tcell.Key) {
		switch k {
		case tcell.KeyCtrlL:
			s.ClearTextBuffer()
			s.ClearText()

		case tcell.KeyCtrlQ, tcell.KeyCtrlX:
			s.Quit()

		case tcell.KeyCtrlR:
			s.Refresh()

		case tcell.KeyCtrlT:
			for i := 0; i < 100; i++ {
				s.AddText(fmt.Sprintf("%d", i))
			}
			s.PL()

		case tcell.KeyCtrlD:
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			t := fmt.Sprintf("Current memory: %dk All-time usage: %dk Sys: %dk NumGC: %d",
				m.Alloc/1024, m.TotalAlloc/1024, m.Sys/1024, m.NumGC)
			s.AddText(t)
			s.PL()
		}
	})

	s.SetCommandFunc(func(cmd string) {
		s.AddText(cmd)
		s.PL()
	})

	defer s.Close()
	s.Run()
}
