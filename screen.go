package tcon

import (
	"sync"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
)

// Screen controls a terminal's input and output.
type Screen struct {
	sync.WaitGroup
	tcell.Screen

	title, status string
	w, h          int

	titleStyle  tcell.Style
	statusStyle tcell.Style
	whiteStyle  tcell.Style
	textStyle   tcell.Style

	lines   []string
	next    []string
	nothing []rune

	// Output section position
	tx, ty, tmaxy int

	// Input section position and buffer
	cx, cy        int
	cbuf          []rune
	cpos          int
	cmdhistory    []string
	cmdhistorypos int
	insertMode    bool

	// callbacks
	OnRune    OnRuneFunc
	OnCommand OnCommandFunc
	OnTab     OnFunc
	OnEsc     OnFunc
	OnCtrl    OnCtrlFunc

	q chan bool
}

// New returns a pointer to a Screen structure which wraps some tcell functionality.
func New() (*Screen, error) {
	s := &Screen{
		insertMode: true,
		q:          make(chan bool),
	}
	scr, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	s.Screen = scr
	err = s.Init()
	if err != nil {
		return nil, err
	}

	s.Sync()
	s.title = "Untitled"
	s.status = "Enter a command:"

	s.w, s.h = s.Size()
	s.ty = 1
	s.tmaxy = s.h - 3
	s.cx = 1
	s.cy = s.h - 1
	s.cbuf = make([]rune, s.w)
	s.PC()

	encoding.Register()

	s.titleStyle = tcell.StyleDefault.
		Reverse(true).
		Bold(true)

	s.statusStyle = s.titleStyle

	s.whiteStyle = tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorWhite)

	s.textStyle = tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack).
		Bold(false)

	s.nothing = make([]rune, s.w)
	s.ClearText()
	return s, nil
}

// Close the screen and restore the terminal.
func (s *Screen) Close() {
	s.Fini()
}

// Run the input loop.
func (s *Screen) Run() {
	s.Refresh()

	for {
		select {
		case _, ok := <-s.q:
			if !ok {
				return
			}
		default:
		}

		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.w, s.h = ev.Size()
			s.tmaxy = s.h - 3
			s.cy = s.h - 1
			buf := make([]rune, s.w)
			max := len(buf)
			if len(s.cbuf) < max {
				max = len(s.cbuf)
			}
			for i := 0; i < max; i++ {
				buf[i] = s.cbuf[i]
			}
			s.cbuf = buf
			s.Refresh()

		case *tcell.EventKey:
			switch ev.Modifiers() {
			case tcell.ModCtrl:
				s.handleCtrl(ev)

			case tcell.ModNone:
				s.handlePlain(ev)
			}

			s.ShowCursor(s.cx+s.cpos, s.cy)
			s.Show()
		}
	}
}

// Quit signals that the screen should close and the input loop should stop.
func (s *Screen) Quit() {
	close(s.q)
}

// handleCtrl input combinations.
func (s *Screen) handleCtrl(ev *tcell.EventKey) {
	if s.OnCtrl != nil {
		s.OnCtrl(ev.Key())
	}
}

// handlePlain rune input without modifier keys.
func (s *Screen) handlePlain(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEsc:
		if s.OnEsc != nil {
			s.OnEsc()
		}

	case tcell.KeyTAB:
		if s.OnTab != nil {
			s.OnTab()
		}

	case tcell.KeyRune:
		// Filtering.
		var r rune
		if s.OnRune != nil {
			r = s.OnRune(ev.Rune())
			if r == 0 {
				break
			}
		} else {
			r = ev.Rune()
		}

		if s.cpos < len(s.cbuf) {
			tmp := []rune{}
			x := s.cpos
			for x := s.cpos; x < len(s.cbuf); x++ {
				tmp = append(tmp, s.cbuf[x])
			}
			s.PL()
			s.cbuf[s.cpos] = r
			s.cpos++
			if s.insertMode {
				save := s.cpos
				for _, c := range tmp {
					if c == 0 || x >= len(s.cbuf) {
						break
					}
					s.cbuf[s.cpos] = c
					s.cpos++
				}
				s.cpos = save
			}
			s.PC()
		}

	case tcell.KeyBackspace2:
		// KeyBackspace is an alias for Ctrl-H., so we use KeyBackspace2.
		if s.cpos > 0 {
			s.cpos--
			s.cbuf[s.cpos] = 0
			for i := s.cpos + 1; i < len(s.cbuf); i++ {
				if s.cbuf[i] != 0 {
					s.cbuf[i-1] = s.cbuf[i]
					s.cbuf[i] = 0
				}
			}
			s.PC()
		}

	case tcell.KeyDelete:
		if s.cpos < len(s.cbuf) {
			s.cbuf = append(s.cbuf[:s.cpos], s.cbuf[s.cpos+1:]...)
			s.cbuf = append(s.cbuf, 0)
		}
		s.PC()

	case tcell.KeyEnter:
		if s.cbuf[0] == 0 {
			break
		}

		s.AddHistory(string(s.cbuf))
		if s.OnCommand != nil {
			s.OnCommand(string(s.cbuf))
		}
		s.ClearCommand()

	case tcell.KeyUp:
		if len(s.cmdhistory) == 0 || s.cmdhistorypos == 0 {
			break
		}

		if s.cmdhistorypos > 0 {
			s.cmdhistorypos--
			s.fetchHistory()
		}

	case tcell.KeyDown:
		if len(s.cmdhistory) == 0 {
			break
		}

		if s.cmdhistorypos < len(s.cmdhistory)-1 {
			s.cmdhistorypos++
			s.fetchHistory()
		} else {
			s.cmdhistorypos++
			s.ClearCommand()
		}

	case tcell.KeyLeft:
		if s.cpos > 0 {
			s.cpos--
		}

	case tcell.KeyRight:
		if s.cpos < len(s.cbuf) {
			s.cpos++
			if s.cbuf[s.cpos] == 0 && s.cbuf[s.cpos-1] == 0 {
				s.cpos--
			}
		}
	}
}

// Refresh the display.
func (s *Screen) Refresh() {
	s.Clear()
	s.UpdateTitle()
	s.UpdateStatus()
	s.ShowCursor(1, s.h-1)
	s.Show()
	s.PL()
	s.PC()
}

// SetTitle sets the text at the top of the terminal, with white background.
func (s *Screen) SetTitle(t string) {
	s.title = t
	s.UpdateTitle()
}

// UpdateTitle refreshes the title display.
func (s *Screen) UpdateTitle() {
	s.Whiteout(0, 0, s.w)
	s.SetStyle(s.titleStyle)
	s.P(0, 0, s.title)
	s.SetStyle(s.textStyle)
}

// SetStatus sets the status text above the command entry.
func (s *Screen) SetStatus(txt string) {
	s.status = txt
	s.UpdateStatus()
}

// UpdateStatus updates the status text and symbol in front of the command line.
func (s *Screen) UpdateStatus() {
	s.Whiteout(0, s.h-2, s.w)
	s.SetStyle(s.statusStyle)
	s.P(0, s.h-2, s.status)
	s.SetStyle(s.statusStyle)
	s.P(0, s.h-1, ">")
	s.Show()
	s.SetStyle(s.textStyle)
}

// SetInsertMode sets overwrite (false) or insert (true) modes.
func (s *Screen) SetInsertMode(t bool) {
	s.insertMode = t
}
