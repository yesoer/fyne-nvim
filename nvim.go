package nvim

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/neovim/go-client/nvim"
)

// Declare conformity with the widget interface
var _ fyne.Widget = (*NeoVim)(nil)

// TODO : Other interfaces we might want to implement :
// - validatable

// Declare conformity with the shortcut interface
// So that we can receive and handle shortcut events, which includes modifiers
// For support of other shortcuts add fyne.ShortCutHandler
var _ fyne.Shortcutable = (*NeoVim)(nil)

// Declare conformity with the focusable interface
// So that we can receive and handle text input events
var _ fyne.Focusable = (*NeoVim)(nil)

// Colorscheme
type highlight struct {
	Fg      color.RGBA `map:"foreground"`
	Bg      color.RGBA `map:"background"`
	Special color.RGBA `map:"special"` // color to use for underlines
	Reverse bool       `map:"reverse"` // trigger switch of fg and bg

	// text styles
	Italic        bool `map:"italic"`
	Bold          bool `map:"bold"`
	Strikethrough bool `map:"strikethrough"`

	// underline styles which all use the special color
	Underline   bool `map:"underline"`
	Undercurl   bool `map:"undercurl"`
	Underdouble bool `map:"underdouble"`
	Underdotted bool `map:"underdotted"`
	Underdashed bool `map:"underdashed"`
}

var defaultHL = highlight{
	Fg:      color.RGBA{0, 0, 0, 0},
	Bg:      color.RGBA{255, 255, 255, 255},
	Special: color.RGBA{0, 0, 0, 255},
}

type NeoVim struct {
	// Widget requirements
	widget.BaseWidget

	// Additional fields
	// It is standard in a Fyne widget to export the fields which define
	// behaviour (just like the primitives defined in the canvas package).
	content              *widget.TextGrid
	cursorRow, cursorCol int
	engine               *nvim.Nvim
	hl                   highlight // the color scheme for the next "put" event
}

// Create a new NeoVim widget
func New() *NeoVim {
	neovim := &NeoVim{}

	tgrid := widget.NewTextGrid()
	neovim.content = tgrid

	neovim.ExtendBaseWidget(neovim)
	err := neovim.startNeovim()
	if err != nil {
		fmt.Println("Error starting neovim: ", err)
	}

	return neovim
}

// Helper to start neovim
func (n *NeoVim) startNeovim() error {
	// start neovim
	// --embed to use stdin/stdout as a msgpack-RPC channel
	opt := nvim.ChildProcessArgs("--embed")
	nvimInstance, err := nvim.NewChildProcess(opt)
	if err != nil {
		return err
	}

	// tell nvim we want to draw the screen
	uiOpt := make(map[string]any)
	err = nvimInstance.AttachUI(100, 100, uiOpt)
	if err != nil {
		return err
	}

	nvimInstance.RegisterHandler("redraw", func(events ...[]interface{}) {
		for _, event := range events {
			fmt.Println("Event ", event)
			n.HandleNvimEvent(event)
		}
	})

	n.engine = nvimInstance

	return nil
}

// Override resize to adjust the textgrid
func (n *NeoVim) Resize(s fyne.Size) {
	n.BaseWidget.Resize(s) // must be included
	n.content.Resize(s)
}

// CreateRenderer implements fyne.Widget
func (n *NeoVim) CreateRenderer() fyne.WidgetRenderer {
	return &render{n}
}

// FocusGained implements fyne.Focusable
// FocusGained is a hook called by the focus handling logic after this object gained the focus.
func (n *NeoVim) FocusGained() {
	n.Refresh()
}

// FocusGained implements fyne.Focusable
// FocusLost is a hook called by the focus handling logic after this object lost the focus.
func (n *NeoVim) FocusLost() {
	n.Refresh()
}

// FocusGained implements fyne.Focusable
// TypedRune is a hook called by the input handling logic on text input events if this object is focused.
func (n *NeoVim) TypedRune(r rune) {
	n.engine.Input(string(r))
}

// FocusGained implements fyne.Focusable
// TypedKey is a hook called by the input handling logic on key events if this object is focused.
func (n *NeoVim) TypedKey(e *fyne.KeyEvent) {
	n.engine.Input(neovimKeyMap[e.Name])
}

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

// Declare conformity with the widget renderer interface
var _ fyne.WidgetRenderer = (*render)(nil)

type render struct {
	*NeoVim
}

// Layout implements fyne.WidgetRenderer
func (r *render) Layout(s fyne.Size) {
	r.content.Resize(s)
}

// MinSize implements fyne.WidgetRenderer
func (r *render) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)
}

// Refresh implements fyne.WidgetRenderer
// The Refresh() method is triggered when the widget this renderer draws has
// changed or if the theme is altered
func (r *render) Refresh() {
	// draw the cursor by inverting fore- and background
	fg, bg := color.Black, color.White
	cellStyle := &widget.CustomTextGridStyle{FGColor: fg, BGColor: bg}
	currentRune := ' '
	if r.cursorRow >= 0 && r.cursorRow < len(r.content.Rows) &&
		r.cursorCol >= 0 && r.cursorCol < len(r.content.Rows[r.cursorRow].Cells) {
		currentRune = r.content.Rows[r.cursorRow].Cells[r.cursorCol].Rune
	}
	cursorCell := widget.TextGridCell{
		Rune:  currentRune,
		Style: cellStyle,
	}
	r.content.SetCell(r.cursorRow, r.cursorCol, cursorCell)

	r.content.Refresh()
}

// Objects implements fyne.WidgetRenderer
func (r *render) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.content}
}

// Destroy implements fyne.WidgetRenderer
// Is called when this renderer is no longer needed so it should clear any
// resources that would otherwise leak
func (r *render) Destroy() {
}
