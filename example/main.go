package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/grimdork/tscreen"
)

func main() {
	s, err := tscreen.New()
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
		}
	})

	s.SetCommandFunc(func(cmd string) {
		s.AddText(cmd)
		s.PL()
	})

	defer s.Close()
	s.Run()
}
