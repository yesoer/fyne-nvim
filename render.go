package nvim

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

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
