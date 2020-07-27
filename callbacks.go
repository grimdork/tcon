package tscreen

import "github.com/gdamore/tcell"

type (
	// OnFunc is the signature for callbacks without any input.
	OnFunc = func()
	// ObRuneFunc is the callback for keypresses corresponding to a printable character, including space.
	// Use this to filter or modify the input. Return the rune to display as-is, 0 to skip it, or modify the
	// command buffer directly.
	OnRuneFunc = func(rune) rune
	// OnCommandFunc is the callback to run when a command has been entered and return has been pressed.
	OnCommandFunc = func(string)
	// OnCtrlFunc is the callback for ctrl+<key> combinations without other modifiers.
	OnCtrlFunc = func(tcell.Key)
)

// SetRuneFunc sets an input filter callback for command input.
func (s *Screen) SetRuneFunc(fn OnRuneFunc) {
	s.OnRune = fn
}

// SetCommandFunc sets a command handler.
func (s *Screen) SetCommandFunc(fn OnCommandFunc) {
	s.OnCommand = fn
}

// SetTabFunc sets a TAB callback.
func (s *Screen) SetTabFunc(fn OnFunc) {
	s.OnTab = fn
}

// SetEscFunc sets a callback for escape.
func (s *Screen) SetEscFunc(fn OnFunc) {
	s.OnEsc = fn
}

// SetEscFunc sets a callback for control keys.
func (s *Screen) SetCtrlFunc(fn OnCtrlFunc) {
	s.OnCtrl = fn
}
