package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/neovim/go-client/nvim"
)

// Declare conformity with the widget interface
var _ fyne.Widget = (*NeoVim)(nil)

// Other interfaces we might want to implement :
// - shortcutable
// - validatable

// Declare conformity with the focusable interface
// So that we can receive and handle text input events
var _ fyne.Focusable = (*NeoVim)(nil)

type NeoVim struct {
	// Widget requirements
	widget.BaseWidget

	// Additional fields
	// It is standard in a Fyne widget to export the fields which define
	// behaviour (just like the primitives defined in the canvas package).
	content *widget.TextGrid
}

// Create a new NeoVim widget
func New() *NeoVim {
	neovim := &NeoVim{}

	tgrid := widget.NewTextGrid()
	neovim.content = tgrid

	neovim.ExtendBaseWidget(neovim)
	err := startNeovim()
	if err != nil {
		fmt.Println("Error starting neovim: ", err)
	}

	return neovim
}

// Helper to start neovim
func startNeovim() error {
	opt := nvim.ChildProcessArgs("--embed")
	nvimInstance, err := nvim.NewChildProcess(opt)
	if err != nil {
		return err
	}

	uiOpt := make(map[string]any)
	err = nvimInstance.AttachUI(100, 100, uiOpt)
	if err != nil {
		return err
	}

	nvimInstance.RegisterHandler("redraw", func(events ...[]interface{}) {
		for _, event := range events {
			fmt.Println("Event ", event)
		}
	})

	return nil
}

func (n *NeoVim) SetText(text string) {
	n.content.SetText(text)
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
	fmt.Println("Focus gained")
	n.Refresh()
}

// FocusGained implements fyne.Focusable
// FocusLost is a hook called by the focus handling logic after this object lost the focus.
func (n *NeoVim) FocusLost() {
	fmt.Println("Focus lost")
	n.Refresh()
}

// FocusGained implements fyne.Focusable
// TypedRune is a hook called by the input handling logic on text input events if this object is focused.
func (n *NeoVim) TypedRune(r rune) {
	fmt.Println("Typed rune: ", r)
	// TODO : possibly should send this to neovim, let it handle it and
	// execute this code in the redraw handler
	// TODO : dynamic (cursor) position and color state in the widget
	row, col := 0, 0
	fg, bg := color.White, color.Black
	cellStyle := &widget.CustomTextGridStyle{FGColor: fg, BGColor: bg}
	n.content.SetCell(row, col, widget.TextGridCell{Rune: r, Style: cellStyle})
}

// FocusGained implements fyne.Focusable
// TypedKey is a hook called by the input handling logic on key events if this object is focused.
func (n *NeoVim) TypedKey(e *fyne.KeyEvent) {
	fmt.Println("Typed key: ", e)
	// TODO : handle
}

// Declare conformity with the widget renderer interface
var _ fyne.WidgetRenderer = (*render)(nil)

type render struct {
	*NeoVim
}

// Layout implements fyne.WidgetRenderer
func (r *render) Layout(fyne.Size) {
	r.content.Resize(r.MinSize())
}

// MinSize implements fyne.WidgetRenderer
func (r *render) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)
}

// Refresh implements fyne.WidgetRenderer
// The Refresh() method is triggered when the widget this renderer draws has
// changed or if the theme is altered
func (r *render) Refresh() {
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
