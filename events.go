package main

import (
	"image/color"

	"fyne.io/fyne/v2/widget"
)

func (n *NeoVim) HandleNvimEvent(event []interface{}) {
	switch event[0] {
	case "set_title":
		// Additional entries: title

	case "set_icon":
		// Additional entries: icon

	case "mode_info_set":
		// Additional entries: cursor_style_enabled, mode_info

	case "option_set":
		// Additional entries: name, value

	case "chdir":
		// Additional entries: path

	case "mode_change":
		// Additional entries: mode, mode_idx

	case "mouse_on":
		// No additional entries

	case "mouse_off":
		// No additional entries

	case "busy_start":
		// No additional entries

	case "busy_stop":
		// No additional entries

	case "suspend":
		// No additional entries

	case "update_menu":
		// No additional entries

	case "bell":
		// No additional entries

	case "visual_bell":
		// No additional entries

	case "flush":
		// No additional entries

	case "grid_resize":
		// Additional entries: grid, width, height

	case "default_colors_set":
		// Additional entries: rgb_fg, rgb_bg, rgb_sp, cterm_fg, cterm_bg

	case "hl_attr_define":
		// Additional entries: id, rgb_attr, cterm_attr, info

	case "hl_group_set":
		// Additional entries: name, hl_id

	case "grid_line":
		// Additional entries: grid, row, col_start, cells, wrap

	case "grid_clear":
		// Additional entries: grid

	case "grid_destroy":
		// Additional entries: grid

	case "grid_cursor_goto":
		// Additional entries: grid, row, column

	case "grid_scroll":
		// Additional entries: grid, top, bot, left, right, rows, cols

	case "resize":
		// Additional entries: width, height

	case "clear":
		// No additional entries

	case "eol_clear":
		// No additional entries

	case "cursor_goto":
		// Additional entries: row, col
		pos := event[1].([]interface{})
		row, _ := pos[0].(int64)
		col, _ := pos[1].(int64)
		n.cursorRow = int(row)
		n.cursorCol = int(col)

	case "update_fg":
		// Additional entries: color

	case "update_bg":
		// Additional entries: color

	case "update_sp":
		// Additional entries: color

	case "highlight_set":
		// Additional entries: attrs

	case "put":
		// Additional entries: text
		for _, s := range event[1:] {
			r := s.([]interface{})[0].(string)
			// TODO : can there be more than one rune in r?
			n.writeRune(rune(r[0]))
		}

	case "set_scroll_region":
		// Additional entries: top, bot, left, right

	case "scroll":
		// Additional entries: count

	case "win_pos":
		// Additional entries: grid, win, start_row, start_col, width, height

	case "win_float_pos":
		// Additional entries: grid, win, anchor, anchor_grid, anchor_row, anchor_col, focusable

	case "win_external_pos":
		// Additional entries: grid, win

	case "win_hide":
		// Additional entries: grid

	case "win_close":
		// Additional entries: grid

	case "msg_set_pos":
		// Additional entries: grid, row, scrolled, sep_char

	case "win_viewport":
		// Additional entries: grid, win, topline, botline, curline, curcol, line_count, scroll_delta

	case "win_extmark":
		// Additional entries: grid, win, ns_id, mark_id, row, col

	case "popupmenu_show":
		// Additional entries: items, selected, row, col, grid

	case "popupmenu_select":
		// Additional entries: selected

	case "popupmenu_hide":
		// No additional entries

	case "tabline_update":
		// Additional entries: curtab, tabs, curbuf, buffers

	case "cmdline_show":
		// Additional entries: content, pos, firstc, prompt, indent, level

	case "cmdline_pos":
		// Additional entries: pos, level

	case "cmdline_special_char":
		// Additional entries: c, shift, level

	case "cmdline_hide":
		// No additional entries

	case "cmdline_block_show":
		// Additional entries: lines

	case "cmdline_block_append":
		// Additional entries: line

	case "cmdline_block_hide":
		// No additional entries

	case "msg_show":
		// Additional entries: kind, content, replace_last

	case "msg_clear":
		// No additional entries

	case "msg_showmode":
		// Additional entries: content

	case "msg_showcmd":
		// Additional entries: content

	case "msg_ruler":
		// Additional entries: content

	case "msg_history_show":
		// Additional entries: entries

	case "msg_history_clear":
		// No additional entries

	default:
		// Handle unknown entry type
	}
}

// Writes a rune to the textgrid
func (n *NeoVim) writeRune(r rune) {
	currRow, currCol := n.cursorRow, n.cursorCol

	// make sure the and columns exist
	for len(n.content.Rows)-1 < currRow {
		n.content.Rows = append(n.content.Rows, widget.TextGridRow{})
	}

	fg, bg := color.White, color.Black
	cellStyle := &widget.CustomTextGridStyle{FGColor: fg, BGColor: bg}

	for len(n.content.Rows[currRow].Cells)-1 < currCol {
		newCell := widget.TextGridCell{
			Rune:  ' ',
			Style: cellStyle,
		}
		n.content.Rows[currRow].Cells = append(n.content.Rows[currRow].Cells, newCell)
	}

	n.content.SetCell(currRow, currCol, widget.TextGridCell{Rune: r, Style: cellStyle})

	n.cursorCol++
}
