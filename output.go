package nvim

import "fyne.io/fyne/v2/widget"

// Substitutes all runes with ' '
func (n *NeoVim) ClearGrid() {
	for i := range n.content.Rows {
		for j := range n.content.Rows[i].Cells {
			n.content.Rows[i].Cells[j].Rune = ' '
		}
	}
}

// Moves the displayed text up/down/lef/right
func (n *NeoVim) ScrollGrid(top, bot, left, right, rows int) {
	if rows > 0 {
		// Scroll down
		for row := top; row < bot-rows; row++ {
			for col := left; col < right; col++ {
				cell := n.content.Rows[row+rows].Cells[col]
				n.content.Rows[row].Cells[col] = cell
			}
		}
	} else {
		// Scroll up, start at bot-1 to skip the status line
		for row := bot - 1; row > top+(-rows); row-- {
			for col := left; col < right; col++ {
				cell := n.content.Rows[row+rows].Cells[col]
				n.content.Rows[row].Cells[col] = cell
			}
		}
	}
}

// Recovers the previous locations colors on horizontal movement and updates the
// cursor position
func (n *NeoVim) MoveGridCursor(oldRow, oldCol, newRow, newCol int) {
	// recover the previous locations colors on horizontal movement
	if oldRow == newRow {
		r := n.content.Rows[oldRow].Cells[oldCol].Rune
		cellStyle := &widget.CustomTextGridStyle{
			FGColor: n.cursorCellFg,
			BGColor: n.cursorCellBg,
		}
		recoveredCell := widget.TextGridCell{Rune: r, Style: cellStyle}
		n.content.SetCell(oldRow, oldCol, recoveredCell)
	}

	newCell := n.content.Rows[newRow].Cells[newCol]
	n.cursorCellFg = newCell.Style.TextColor()
	n.cursorCellBg = newCell.Style.BackgroundColor()
	n.cursorRow = newRow
	n.cursorCol = newCol
}

// Writes a line of text (as defined by neovims ui events) to the textgrid
func (n *NeoVim) WriteGridLine(row, col int, cells []interface{}) {
	lastHL_id := 0
	for _, cell := range cells {
		cell := cell.([]interface{})
		s := cell[0].(string)
		r := rune(s[0])

		if len(cell) > 1 {
			lastHL_id, _ = intOrUintToInt(cell[1])
		}

		repeat := 1
		if len(cell) > 2 {
			repeat, _ = intOrUintToInt(cell[2])
		}

		for i := 0; i < repeat; i++ {
			n.writeRune(row, col, r, lastHL_id)
			if len(s) > 1 {
				n.writeRune(row, col, ' ', lastHL_id)
			}
			col++
		}
	}
}

// Changes the size of the textgrid, creating or removing rows and columns as
// needed
// TODO: One may consider only downsizing rows/cols if the difference is
// significant
func (n *NeoVim) ChangeVisualGridSize(targetRow, targetCol int) {
	// remove rows
	if targetRow < len(n.content.Rows) {
		n.content.Rows = n.content.Rows[:targetRow]
	}

	cellStyle := gridStyleFromHL(defaultHL)

	for currRow := 0; currRow < targetRow; currRow++ {
		// append new row if needed
		if currRow > len(n.content.Rows)-1 {
			n.content.Rows = append(n.content.Rows, widget.TextGridRow{})
		}

		// remove columns
		numCols := len(n.content.Rows[currRow].Cells)
		if numCols > targetCol {
			n.content.Rows[currRow].Cells =
				n.content.Rows[currRow].Cells[:targetCol-1]
		}

		// append new columns if needed
		for len(n.content.Rows[currRow].Cells) < targetCol {
			newCell := widget.TextGridCell{
				Rune:  ' ',
				Style: cellStyle,
			}
			n.content.Rows[currRow].Cells =
				append(n.content.Rows[currRow].Cells, newCell)
		}
	}
}

// Writes a rune to the textgrid
func (n *NeoVim) writeRune(row int, col int, r rune, hl_id int) {

	hl, ok := n.hl[hl_id]
	if !ok {
		hl = defaultHL
	}

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
