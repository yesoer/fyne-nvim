package nvim

import (
	"fmt"
	"image/color"
	"reflect"

	"fyne.io/fyne/v2"
)

// Handles events for the NeoVim instance
// Currently only the following event groups are handled:
// - Global Events
// - Grid Events (line-based)
// For the documentation of the events see:
// https://neovim.io/doc/user/ui.html
// The go client calls this function sequentially for each event, so we don't
// have to worry about preserving order
func (n *NeoVim) HandleNvimEvent(event []interface{}) {
	// fmt.Println("Handling event: ", event)
	// fmt.Println("Handling event: ", event[0])

	for _, e := range event[1:] {
		entries, ok := e.([]interface{})
		if !ok {
			entries = []interface{}{e}
		}

		switch event[0] {

		//------------------------------Global Events-------------------------------

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

		//-------------------------Grid Events (line-based)-------------------------

		case "grid_resize":
			// The grid is resized to width and height cells.
			// Additional entries: grid, width, height

			colsCnt, _ := intOrUintToInt(entries[1])
			rowsCnt, _ := intOrUintToInt(entries[2])
			n.ChangeVisualGridSize(rowsCnt, colsCnt)

			cellSize := guessCellSize()
			s := fyne.NewSize(float32(colsCnt*int(cellSize.Width)),
				float32(rowsCnt*int(cellSize.Height)))

			n.BaseWidget.Resize(s) // must be included
			n.content.Resize(s)

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

			defaultHL.Fg, _ = extractRGBA(entries[0])
			defaultHL.Bg, _ = extractRGBA(entries[1])
			defaultHL.Special, _ = extractRGBA(entries[2])
			// cterm_fg, cterm_bg are ignored
			n.Refresh()

		case "hl_attr_define":
			// Add a new highlight with id to the highlight table. rgb_attr carries
			// information on fore-/background, special color, text attributes and
			// underline styles. cterm_attr is relevant for 256-color terminals so
			// it is ignored. info is used by the ext_hlstate extension to add
			// semantic information.
			// Additional entries: id, rgb_attr, cterm_attr, info

			id, _ := intOrUintToInt(entries[0])

			rgbAttr := entries[1].(map[string]interface{})

			newHL := highlight{
				Fg:      RGBA_SENTINEL,
				Bg:      RGBA_SENTINEL,
				Special: RGBA_SENTINEL,
			}
			setHLFromMap(rgbAttr, &newHL)
			n.hl[id] = newHL

			// Info is ignored since we don't need semantic information as of
			// now

		case "hl_group_set":
			// The built-in highlight group name was set to use the attributes hl_id
			// defined by a previous hl_attr_define call. This event is not needed
			// to render the grids which use attribute ids directly, but is useful
			// for a UI who want to render its own elements with consistent
			// highlighting. For instance a UI using ui-popupmenu events, might use
			// the hl-Pmenu family of builtin highlights.
			// Additional entries: name, hl_id

		case "grid_line":
			// Write row from col_start with cells. Cells is an array of arrays each
			// with 1 to 3 items: [text(, hl_id, repeat)]. The text should be
			// styled with the colorscheme at hl_id in the table. If no hl_id is
			// provided, use the most recent from this call (is always present for
			// the first cell). repeat is a number indicating how many times cell
			// should be placed.
			// The right cell of a double-width char will be represented as the
			// empty string. Double-width chars never use repeat.
			// wrap is a boolean indicating that this line wraps to the next row.
			// When redrawing a line which wraps to the next row, Nvim will emit a
			// grid_line event covering the last column of the line with wrap set
			// to true, followed immediately by a grid_line event starting at the
			// first column of the next row.
			// Additional entries: grid, row, col_start, cells, wrap

			row, _ := intOrUintToInt(entries[1])
			col, _ := intOrUintToInt(entries[2])
			cells := entries[3].([]interface{})
			// wrap := entries[4].(bool) // TODO : use wrap

			n.WriteGridLine(row, col, cells)

		case "grid_clear":
			// Clear a grid
			// Additional entries: grid

			n.ClearGrid()

		case "grid_destroy":
			// Grid will not be used anymore and the UI can free any data associated
			// with it.
			// Additional entries: grid

			n.content = nil

		case "grid_cursor_goto":
			// Makes grid the current grid and row, column the cursor position on
			// this grid. This event will be sent at most once in a redraw batch and
			// indicates the visible cursor position.
			// Additional entries: grid, row, column

			oldRow, oldCol := n.cursorRow, n.cursorCol
			newRow, _ := entries[1].(int64)
			newCol, _ := entries[2].(int64)

			n.MoveGridCursor(oldRow, oldCol, int(newRow), int(newCol))

		case "grid_scroll":
			// Scroll a region of grid. This is semantically unrelated to editor
			// scrolling, rather this is an optimized way to say "copy these
			// screen cells".
			// If rows is bigger than 0, move a rectangle in the SR up,
			// this can happen while scrolling down.
			// If rows is less than zero, move a rectangle in the SR down, this
			// can happen while scrolling up.
			// cols is always zero in this version of Nvim, and reserved for
			// future use.
			// The scrolled-in area will be filled using ui-event-grid_line
			// directly after the scroll event. The UI thus doesn't need to
			// clear this area as part of handling the scroll event.
			// Additional entries: grid, top, bot, left, right, rows, cols

			top, _ := intOrUintToInt(entries[1])
			bot, _ := intOrUintToInt(entries[2])
			left, _ := intOrUintToInt(entries[3])
			right, _ := intOrUintToInt(entries[4])
			rows, _ := intOrUintToInt(entries[5])

			n.ScrollGrid(top, bot, left, right, rows)

		default:
			// Handle unknown entry type
			fmt.Println("Unknown event type: ", event[0])
		}
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
				value, _ = extractRGBA(value)
				if !ok {
					fmt.Println("Unknown type: ", value)
					continue
				}
			}
			targetValue.Field(i).Set(reflect.ValueOf(value))
		}
	}

	if target.Reverse {
		target.Fg, target.Bg = target.Bg, target.Fg
	}
}

func intOrUintToInt(i interface{}) (int, bool) {
	switch i.(type) {
	case uint64:
		return int(i.(uint64)), true
	case int64:
		return int(i.(int64)), true
	default:
		fmt.Println("Unknown type: ", i)
		return 0, false
	}
}

// Expects a uint64 or int64 and returns its corresponding color.RGBA
func extractRGBA(i interface{}) (color.RGBA, bool) {
	n, ok := intOrUintToInt(i)
	if !ok {
		return color.RGBA{}, false
	}

	r := (n >> 16) & 0xFF
	g := (n >> 8) & 0xFF
	b := n & 0xFF
	a := 255
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}, true
}
