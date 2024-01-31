package nvim

import (
	"fmt"
	"image/color"
	"reflect"

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
		// Nvim is done redrawing the screen. For an implementation that renders
		// to an internal buffer, this is the time to display the redrawn parts
		// to the user.
		// No additional entries
		n.Refresh()

	case "grid_resize":
		// Additional entries: grid, width, height

	case "default_colors_set":
		// The RGB values will always be valid colors, by default. If no colors
		// have been set, they will default to black and white, depending on
		// 'background'. By setting the ext_termcolors option, instead -1 will
		// be used for unset colors. This is mostly useful for a TUI
		// implementation, where using the terminal builtin ("ANSI") defaults
		// are expected.
		// Note: Unlike the corresponding ui-grid-old events, the screen is not
		// always cleared after sending this event. The UI must repaint the
		// screen with changed background color itself.
		// Additional entries: rgb_fg, rgb_bg, rgb_sp, cterm_fg, cterm_bg
		defaultHL.Fg = extractRGBAFromEvent(event, 0)
		defaultHL.Bg = extractRGBAFromEvent(event, 1)
		defaultHL.Special = extractRGBAFromEvent(event, 2)
		// cterm_fg, cterm_bg are ignored
		n.Refresh()

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
		// Move the cursor to position (row, col). Currently, the same cursor is
		// used to define the position for text insertion and the visible
		// cursor. However, only the last cursor position, after processing the
		// entire array in the "redraw" event, is intended to be a visible cursor
		// position.
		// Additional entries: row, col
		pos := event[1].([]interface{})
		row, _ := pos[0].(int64)
		col, _ := pos[1].(int64)
		n.cursorRow = int(row)
		n.cursorCol = int(col)

	// Events to set the default colors
	// Additional entries: color
	case "update_fg":
		defaultHL.Fg = extractRGBAFromEvent(event, 0)
	case "update_bg":
		defaultHL.Bg = extractRGBAFromEvent(event, 0)
	case "update_sp":
		defaultHL.Special = extractRGBAFromEvent(event, 0)
	case "highlight_set":
		// Set the attributes that the next text put on the grid will have.
		// Additional entries: attrs which is a dictionary
		m := event[1].([]interface{})[0].(map[string]interface{})
		newHL := highlight{
			Fg:      defaultHL.Fg,
			Bg:      defaultHL.Bg,
			Special: defaultHL.Special,
		}
		setHLFromMap(m, &defaultHL)
		n.hl = newHL

	case "put":
		// The (utf-8 encoded) string text is put at the cursor position (and
		// the cursor is advanced), with the highlights as set by the last
		// highlight_set update.
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
		fmt.Println("Unknown event type: ", event[0])
	}
}

// Expects a map which defines the attributes for highlighting etc. and a target
// to write them to.
func setHLFromMap(personMap map[string]interface{}, target *highlight) {
	targetValue := reflect.ValueOf(target).Elem()

	for i := 0; i < targetValue.NumField(); i++ {
		field := targetValue.Type().Field(i)
		tag := field.Tag.Get("map")
		if value, ok := personMap[tag]; ok {
			if field.Type == reflect.TypeOf(color.RGBA{}) {
				value, ok = extractRGBA[uint64](value)
				if !ok {
					value, _ = extractRGBA[int64](value)
				}
			}
			targetValue.Field(i).Set(reflect.ValueOf(value))
		}
	}

	if target.Reverse {
		target.Fg, target.Bg = target.Bg, target.Fg
	}
}

// A helper to wrap the extraction of RGBA colors from nvim events
func extractRGBAFromEvent(event []interface{}, pos int) color.RGBA {
	entry := event[1].([]interface{})[pos]
	c, ok := extractRGBA[uint64](entry)
	if !ok {
		c, _ = extractRGBA[int64](entry)
	}
	return c
}

// Constraint for the color in nvim events
type NvimColor interface {
	uint64 | int64
}

// Expects a uint64 or int64 and returns its corresponding color.RGBA
func extractRGBA[T NvimColor](i interface{}) (color.RGBA, bool) {
	n, ok := i.(T)
	if !ok {
		return color.RGBA{}, false
	}

	r := (n >> 16) & 0xFF
	g := (n >> 8) & 0xFF
	b := n & 0xFF
	a := 255
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}, true
}

// Writes a rune to the textgrid
func (n *NeoVim) writeRune(r rune) {
	currRow, currCol := n.cursorRow, n.cursorCol

	// make sure the and columns exist
	for len(n.content.Rows)-1 < currRow {
		n.content.Rows = append(n.content.Rows, widget.TextGridRow{})
	}

	fg, bg := n.hl.Fg, n.hl.Bg
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
