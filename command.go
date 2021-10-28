package tcon

// PC prints the command buffer.
func (s *Screen) PC() {
	i := s.cx
	s.SetStyle(s.textStyle)
	s.P(s.cx, s.cy, string(s.cbuf))
	s.Show()
	s.ShowCursor(i, s.cy)
	s.Show()
}

// ClearCommand clears the command buffer and updates the display.
func (s *Screen) ClearCommand() {
	for i := 0; i < len(s.cbuf); i++ {
		s.cbuf[i] = 0
	}
	s.cpos = 0
	s.PC()
}

// AddHistory adds a string to the command history.
func (s *Screen) AddHistory(cmd string) {
	if len(s.cmdhistory) > 0 && cmd == s.cmdhistory[len(s.cmdhistory)-1] {
		return
	}

	if len(s.cmdhistory) > 1000 {
		s.cmdhistory = s.cmdhistory[1:]
	}

	s.cmdhistory = append(s.cmdhistory, cmd)
	s.cmdhistorypos = len(s.cmdhistory)
}

// fetchHistory sets the current command to the current history point.
func (s *Screen) fetchHistory() {
	s.cpos = 0
	if s.cmdhistorypos >= len(s.cmdhistory) {
		s.cmdhistorypos = len(s.cmdhistory) - 1
	}
	for i, r := range s.cmdhistory[s.cmdhistorypos] {
		if i >= len(s.cbuf) {
			break
		}
		s.cbuf[i] = r
		if r != 0 {
			s.cpos++
		}
	}
	s.PC()
}
