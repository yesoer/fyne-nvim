package nvim

import (
	"fmt"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/neovim/go-client/nvim"
)

// According to https://neovim.io/doc/user/options.html it seems that
// neovim is supposed to have at least 12 columns and 1 row
// In practice the neovim client didn't work with less than 13x1
const MIN_ROWS = 1
const MIN_COLS = 13

// As long as multigrid isn't used there will only be one grid
const GLOBAL_GRID = 1

// Declare conformity with the widget interface
var _ fyne.Widget = (*NeoVim)(nil)

// Declare conformity with the shortcut interface
// So that we can receive and handle shortcut events, which includes modifiers
// For support of other shortcuts add fyne.ShortCutHandler
var _ fyne.Shortcutable = (*NeoVim)(nil)

// Declare conformity with the focusable interface
// So that we can receive and handle text input events
var _ fyne.Focusable = (*NeoVim)(nil)

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

	Altfont interface{} `map:"altfont"` // TODO : implement
	Blend   interface{} `map:"blend"`   // TODO : implement
}

var defaultHL = highlight{
	Fg:      color.RGBA{255, 255, 255, 0},
	Bg:      color.RGBA{0, 0, 0, 255},
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
	hl                   map[int]highlight // the highlight table used by ext_hlstate
}

// Create a new NeoVim widget
func New() *NeoVim {
	neovim := &NeoVim{}
	neovim.hl = make(map[int]highlight)

	tgrid := widget.NewTextGrid()
	neovim.content = tgrid

	neovim.ExtendBaseWidget(neovim)
	err := neovim.startNeovim()
	if err != nil {
		fmt.Println("Error starting neovim: ", err)
	}

	return neovim
}

// Helper to estimate the size of a cell in the textgrid
func guessCellSize() fyne.Size {
	cell := canvas.NewText("M", color.White)
	cell.TextStyle.Monospace = true

	min := cell.MinSize()
	return fyne.NewSize(float32(math.Round(float64(min.Width))), float32(math.Round(float64(min.Height))))
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

	// tell nvim we want to draw the screen (using the new line based API)
	uiOpt := make(map[string]any)
	uiOpt["ext_hlstate"] = true  // detailed highlight state
	uiOpt["ext_linegrid"] = true // new line based grid events
	uiOpt["ext_multigrid"] = false
	err = nvimInstance.AttachUI(MIN_COLS, MIN_ROWS, uiOpt)
	if err != nil {
		fmt.Println("Error attaching UI: ", err)
		return err
	}

	nvimInstance.RegisterHandler("redraw", func(events ...[]interface{}) {
		for _, event := range events {
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
	n.resizeGrid(s)
}

// Resizes the neovim internal grid
func (n *NeoVim) resizeGrid(s fyne.Size) {
	cellSize := guessCellSize()
	rowsCnt := int(s.Height / cellSize.Height)
	colsCnt := int(s.Width / cellSize.Width)

	err := n.engine.TryResizeUIGrid(GLOBAL_GRID, colsCnt, rowsCnt)
	if err != nil {
		fmt.Println("Error resizing grid: ", err)
	}
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
	cellSize := guessCellSize()
	minWidth := cellSize.Width * MIN_COLS
	minHeight := cellSize.Height * MIN_ROWS
	return fyne.NewSize(minWidth, minHeight)
}

// Refresh implements fyne.WidgetRenderer
// The Refresh() method is triggered when the widget this renderer draws has
// changed or if the theme is altered
func (r *render) Refresh() {
	r.refreshCursor()
	r.content.Refresh()
}

// refreshCursor draws the cursor
func (r *render) refreshCursor() {
	cellStyle := &widget.CustomTextGridStyle{
		FGColor: color.RGBA{200, 200, 200, 180},
		BGColor: color.RGBA{255, 255, 255, 180},
	}

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
}

// Objects implements fyne.WidgetRenderer
func (r *render) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.content}
}

// Destroy implements fyne.WidgetRenderer
// Is called when this renderer is no longer needed so it should clear any
// resources that would otherwise leak
func (r *render) Destroy() {
	r.engine.Close()
}
