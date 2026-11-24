package editor

import "github.com/charmbracelet/x/ansi"

// Row is one row of text on screen. Start and End are character indexes, End
// exclusive and without the new line that ends a line. Cols holds the column
// of each character, and EndCol the column after the last one.
type Row struct {
	Start, End int
	Cols       []int
	EndCol     int
}

// Layout breaks gs into rows at most width columns wide, measuring each
// character with m. A row breaks after the last space that fits, or inside a
// word that is wider than the row.
func Layout(gs []string, width int, m ansi.Method) []Row {
	var rows []Row
	start := 0
	for {
		end := start
		for end < len(gs) && gs[end] != "\n" {
			end++
		}
		rows = appendLine(rows, gs, start, end, width, m)
		if end == len(gs) {
			return rows
		}
		start = end + 1
	}
}

func appendLine(rows []Row, gs []string, start, end, width int, m ansi.Method) []Row {
	for {
		var cols []int
		col, k := 0, start
		for k < end {
			w := charWidth(gs[k], m)
			if col+w > width && k > start {
				break
			}
			cols = append(cols, col)
			col += w
			k++
		}
		if k == end {
			return append(rows, Row{Start: start, End: end, Cols: cols, EndCol: col})
		}

		rowEnd := k
		for j := k - 1; j >= start; j-- {
			if gs[j] == " " {
				rowEnd = j + 1
				break
			}
		}
		endCol := col
		if rowEnd < k {
			endCol = cols[rowEnd-start]
		}
		rows = append(rows, Row{Start: start, End: rowEnd, Cols: cols[:rowEnd-start], EndCol: endCol})
		start = rowEnd
	}
}

func charWidth(g string, m ansi.Method) int {
	if len(g) == 1 && g[0] >= ' ' && g[0] < 0x7f {
		return 1
	}
	return m.StringWidth(g)
}
