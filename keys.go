package main

import "fyne.io/fyne/v2"

// These are the keys fyne does not send to TypedRune but to TypedKey
var neovimKeyMap = map[fyne.KeyName]string{
	fyne.KeyEscape:    "<Esc>",
	fyne.KeyReturn:    "<CR>",
	fyne.KeyTab:       "<Tab>",
	fyne.KeyBackspace: "<BS>",
	fyne.KeyInsert:    "<Insert>",
	fyne.KeyDelete:    "<Del>",
	fyne.KeyRight:     "<Right>",
	fyne.KeyLeft:      "<Left>",
	fyne.KeyDown:      "<Down>",
	fyne.KeyUp:        "<Up>",
	fyne.KeyPageUp:    "<PageUp>",
	fyne.KeyPageDown:  "<PageDown>",
	fyne.KeyHome:      "<Home>",
	fyne.KeyEnd:       "<End>",
	fyne.KeyF1:        "<F1>",
	fyne.KeyF2:        "<F2>",
	fyne.KeyF3:        "<F3>",
	fyne.KeyF4:        "<F4>",
	fyne.KeyF5:        "<F5>",
	fyne.KeyF6:        "<F6>",
	fyne.KeyF7:        "<F7>",
	fyne.KeyF8:        "<F8>",
	fyne.KeyF9:        "<F9>",
	fyne.KeyF10:       "<F10>",
	fyne.KeyF11:       "<F11>",
	fyne.KeyF12:       "<F12>",
	fyne.KeyEnter:     "<CR>",
}
