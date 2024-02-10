package nvim

import (
	"fmt"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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

// The sentinel value for Fg, Bg and Special to indicate that the coloris not
// set i.e. the default color should be used
var RGBA_SENTINEL = color.RGBA{255, 255, 255, 0}

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

	Altfont interface{} `map:"altfont"`
	Blend   interface{} `map:"blend"`
}

var defaultHL = highlight{
	Fg:      color.RGBA{255, 255, 255, 255},
	Bg:      color.RGBA{0, 0, 0, 255},
	Special: color.RGBA{0, 0, 0, 255},
}

// Declare conformity with the widget interface
var _ fyne.Widget = (*NeoVim)(nil)

type NeoVim struct {
	// Widget requirements
	widget.BaseWidget

	// Additional fields
	// It is standard in a Fyne widget to export the fields which define
	// behaviour (just like the primitives defined in the canvas package).
	content                    *widget.TextGrid
	cursorRow, cursorCol       int
	cursorCellFg, cursorCellBg color.Color // store color of the underlying cell
	engine                     *nvim.Nvim
	hl                         map[int]highlight // the highlight table used by ext_hlstate
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

// Helper to estimate the size of a cell in the textgrid
func guessCellSize() fyne.Size {
	cell := canvas.NewText("M", color.White)
	cell.TextStyle.Monospace = true

	min := cell.MinSize()
	return fyne.NewSize(float32(math.Round(float64(min.Width))), float32(math.Round(float64(min.Height))))
}
