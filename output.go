package nvim

import "fyne.io/fyne/v2/widget"

// Make sure the rows and columns exist, if not create them
func (n *NeoVim) fillGrid(row, col int, hl highlight) {
	for len(n.content.Rows)-1 < row {
		n.content.Rows = append(n.content.Rows, widget.TextGridRow{})
	}

	cellStyle := gridStyleFromHL(hl)

	for len(n.content.Rows[row].Cells)-1 < col {
		newCell := widget.TextGridCell{
			Rune:  ' ',
			Style: cellStyle,
		}
		n.content.Rows[row].Cells = append(n.content.Rows[row].Cells, newCell)
	}
}

// Writes a rune to the textgrid
func (n *NeoVim) writeRune(row int, col int, r rune, hl_id int) {

	hl, ok := n.hl[hl_id]
	if !ok {
		hl = defaultHL
	}
	n.fillGrid(row, col, hl)

	cellStyle := gridStyleFromHL(hl)
	n.content.SetCell(row, col, widget.TextGridCell{Rune: r, Style: cellStyle})
}

func gridStyleFromHL(hl highlight) *widget.CustomTextGridStyle {
	style := widget.CustomTextGridStyle{
		FGColor: hl.Fg,
		BGColor: hl.Bg,
	}

	if style.FGColor == RGBA_SENTINEL {
		style.FGColor = defaultHL.Fg
	}

	if style.BGColor == RGBA_SENTINEL {
		style.BGColor = defaultHL.Bg
	}

	return &style
}
