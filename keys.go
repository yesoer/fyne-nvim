package nvim

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

// These are the modifier keys fyne supports, which are also sent to TypedKey
var neovimModifierMap = map[fyne.KeyModifier]string{
	// base modifiers
	fyne.KeyModifierShift:   "S-",
	fyne.KeyModifierAlt:     "A-",
	fyne.KeyModifierControl: "C-",
	fyne.KeyModifierSuper:   "M-",
	// all possible combinations of the above
	fyne.KeyModifierShift | fyne.KeyModifierAlt:                                                   "S-A-",
	fyne.KeyModifierShift | fyne.KeyModifierControl:                                               "S-C-",
	fyne.KeyModifierShift | fyne.KeyModifierSuper:                                                 "S-M-",
	fyne.KeyModifierAlt | fyne.KeyModifierControl:                                                 "A-C-",
	fyne.KeyModifierAlt | fyne.KeyModifierSuper:                                                   "A-M-",
	fyne.KeyModifierControl | fyne.KeyModifierSuper:                                               "C-M-",
	fyne.KeyModifierShift | fyne.KeyModifierAlt | fyne.KeyModifierControl:                         "S-A-C-",
	fyne.KeyModifierShift | fyne.KeyModifierAlt | fyne.KeyModifierSuper:                           "S-A-M-",
	fyne.KeyModifierShift | fyne.KeyModifierControl | fyne.KeyModifierSuper:                       "S-C-M-",
	fyne.KeyModifierAlt | fyne.KeyModifierControl | fyne.KeyModifierSuper:                         "A-C-M-",
	fyne.KeyModifierShift | fyne.KeyModifierAlt | fyne.KeyModifierControl | fyne.KeyModifierSuper: "S-A-C-M-",
}
