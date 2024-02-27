package nvim

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// Declare conformity with the tappable interface
// So that we can focus when tapped
var _ fyne.Tappable = (*NeoVim)(nil)

// Declare conformity with the focusable interface
// So that we can receive and handle text input events
var _ fyne.Focusable = (*NeoVim)(nil)

// Tapped implements fyne.Tappable
// Tapped makes sure we ask for focus if user taps us.
func (n *NeoVim) Tapped(ev *fyne.PointEvent) {
	fyne.CurrentApp().Driver().CanvasForObject(n).Focus(n)
}

// FocusGained implements fyne.Focusable
// FocusGained is a hook called by the focus handling logic after this object
// gained the focus.
func (n *NeoVim) FocusGained() {
	n.Refresh()
}

// FocusGained implements fyne.Focusable
// FocusLost is a hook called by the focus handling logic after this object lost
// the focus.
func (n *NeoVim) FocusLost() {
	n.Refresh()
}

// FocusGained implements fyne.Focusable
// TypedRune is a hook called by the input handling logic on text input events
// if this object is focused.
func (n *NeoVim) TypedRune(r rune) {
	n.engine.Input(string(r))
}

// FocusGained implements fyne.Focusable
// TypedKey is a hook called by the input handling logic on key events if this
// object is focused.
func (n *NeoVim) TypedKey(e *fyne.KeyEvent) {
	n.engine.Input(neovimKeyMap[e.Name])
}

// Declare conformity with the shortcut interface
// So that we can receive and handle shortcut events, which includes modifiers
// For support of other shortcuts add fyne.ShortCutHandler
var _ fyne.Shortcutable = (*NeoVim)(nil)

// TypedShortcut implements fyne.Shortcutable
// TypedShortcut handle the registered shortcut
// TODO : There are other shortcuts e.g. SelectAll (Cmd+A)
func (n *NeoVim) TypedShortcut(s fyne.Shortcut) {
	if ds, ok := s.(*desktop.CustomShortcut); ok {

		char := ds.KeyName[0]
		if ds.Key() == fyne.KeySpace {
			char = ' '
		} else if ds.Key() == "@" {
			char = '@'
		}

		modifiers := neovimModifierMap[ds.Modifier]
		n.engine.Input("<" + modifiers + string(char) + ">")
	}
}
