package tcon

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
)

// Whiteout creates a line of n white blocks at the specified position.
func (s *Screen) Whiteout(x, y, n int) {
	s.SetStyle(s.whiteStyle)
	space := strings.Repeat("\u2588", n)
	s.P(x, y, space)
	s.Show()
	s.SetStyle(s.textStyle)
}

// P prints a string at the desired position.
func (s *Screen) P(x, y int, txt string) {
	for _, r := range txt {
		s.SetContent(x, y, r, nil, tcell.StyleDefault)
		x++
	}
	s.Show()
}

// Printf prints a formatted string to the output buffer.
func (s *Screen) Priintf(format string, args ...interface{}) {
	s.AddText(fmt.Sprintf(format, args...))
	s.PL()
}

// AddText adds text to the output buffer, splitting it if needed.
// TODO: History buffer and PgUp/PgDown.
func (s *Screen) AddText(txt string) {
	lines := WordWrap(txt, s.w)
	if len(lines) > s.tmaxy {
		s.lines[s.tmaxy-1] = ""
		s.lines = lines[:s.tmaxy]
		s.next = lines[s.tmaxy:]
		return
	}

	s.lines = append(s.lines, lines...)
	if len(s.lines) > s.tmaxy {
		s.lines = s.lines[len(lines):]
	}
}

// PL prints all lines in the line buffer.
func (s *Screen) PL() {
	for {
		y := s.ty
		for _, l := range s.lines {
			s.P(s.tx, y, l)
			y++
		}
		s.UpdateTitle()

		if len(s.next) > 0 {
			s.P(s.tx, s.ty+s.tmaxy, "----- Press space to continue -----")
			loop := true
			for loop {
				ev := s.PollEvent()
				switch ev := ev.(type) {
				case *tcell.EventKey:
					switch ev.Rune() {
					case ' ':
						loop = false
						continue
					}
				}
			}

			if len(s.next) > s.tmaxy {
				s.lines = s.next[:s.tmaxy]
				s.next = s.next[s.tmaxy:]
			} else {
				s.lines = s.lines[len(s.next):]
				s.lines = append(s.lines, s.next...)
				s.next = nil
			}
		} else {
			return
		}
	}
}

// ClearText clears the output section.
func (s *Screen) ClearText() {
	s.SetStyle(s.textStyle)
	for i := s.ty; i <= s.tmaxy; i++ {
		s.P(s.tx, i, string(s.nothing))
	}
}

// ClearTextBuffer clears the output buffer.
func (s *Screen) ClearTextBuffer() {
	s.lines = nil
}

// WordWrap returns one or more lines, based on the desired width.
// This will be split at spaces or punctuation as far as possible,
// but if using small widths and long words they will be split and
// a hyphen added.
func WordWrap(txt string, width int) []string {
	lines := []string{}
	if len(txt) <= width {
		lines = append(lines, txt)
		return lines
	}

	var l string
	s := txt
	for s != "" {
		l, s = splitLine(s, width)
		lines = append(lines, l)
	}

	return lines
}

// splitLine into two parts after space or punctuation.
func splitLine(txt string, width int) (string, string) {
	if len(txt) <= width {
		return txt, ""
	}

	i := strings.LastIndexAny(txt[:width], " ,.!?;:…‽")
	if i == -1 || i > width {
		l := txt[:width-1]
		return l + "-", txt[width-1:]
	}

	i++
	return txt[:i], txt[i:]
}
