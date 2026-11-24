package editor

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func rowTexts(text string, m ansi.Method) []string {
	gs := Graphemes(text)
	var out []string
	for _, r := range Layout(gs, textWidth, m) {
		out = append(out, strings.Join(gs[r.Start:r.End], ""))
	}
	return out
}

func assertRows(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %d rows %q, want %d rows %q", len(got), got, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("row %d is %q, want %q", i, got[i], want[i])
		}
	}
}

func TestLayoutOfEmptyTextIsOneEmptyRow(t *testing.T) {
	rows := Layout(nil, textWidth, ansi.GraphemeWidth)
	if len(rows) != 1 || rows[0].Start != 0 || rows[0].End != 0 {
		t.Errorf("got rows %+v, want one empty row", rows)
	}
}

func TestLayoutBreaksAtANewLine(t *testing.T) {
	assertRows(t, rowTexts("one\ntwo\n", ansi.GraphemeWidth), []string{"one", "two", ""})
}

func TestLayoutWrapsAtASpaceAndKeepsTheSpaceAtTheEndOfItsRow(t *testing.T) {
	text := strings.Repeat("word ", 11) + "word end"
	got := rowTexts(text, ansi.GraphemeWidth)
	assertRows(t, got, []string{strings.Repeat("word ", 11), "word end"})
	if strings.Join(got, "") != text {
		t.Errorf("rows join to %q, want the original text", strings.Join(got, ""))
	}
}

func TestLayoutBreaksAWordWiderThanTheRow(t *testing.T) {
	got := rowTexts(strings.Repeat("x", 100), ansi.GraphemeWidth)
	assertRows(t, got, []string{strings.Repeat("x", 57), strings.Repeat("x", 43)})
}

func TestLayoutMovesAWideCharacterToTheNextRowWhole(t *testing.T) {
	got := rowTexts(strings.Repeat("合", 40), ansi.GraphemeWidth)
	assertRows(t, got, []string{strings.Repeat("合", 28), strings.Repeat("合", 12)})
}

func TestLayoutFillsARowExactlyToTheWidth(t *testing.T) {
	text := strings.Repeat("a", 57) + "bbb"
	rows := Layout(Graphemes(text), textWidth, ansi.GraphemeWidth)
	assertRows(t, rowTexts(text, ansi.GraphemeWidth), []string{strings.Repeat("a", 57), "bbb"})
	if rows[0].EndCol != 57 {
		t.Errorf("first row ends at column %d, want 57", rows[0].EndCol)
	}
	if rows[1].Start != rows[0].End {
		t.Errorf("second row starts at %d, want it to start where the first ends, at %d", rows[1].Start, rows[0].End)
	}
}

func TestLayoutRecordsTheColumnOfEveryCharacter(t *testing.T) {
	rows := Layout(Graphemes("a合b"), textWidth, ansi.GraphemeWidth)
	want := []int{0, 1, 3}
	for i, col := range want {
		if rows[0].Cols[i] != col {
			t.Errorf("character %d is at column %d, want %d", i, rows[0].Cols[i], col)
		}
	}
	if rows[0].EndCol != 4 {
		t.Errorf("row ends at column %d, want 4", rows[0].EndCol)
	}
}

func TestLayoutRowsStartWhereThePreviousEnded(t *testing.T) {
	text := strings.Repeat("alpha beta gamma\n", 30)
	gs := Graphemes(text)
	prevEnd := 0
	for i, r := range Layout(gs, 20, ansi.GraphemeWidth) {
		if i > 0 && r.Start != prevEnd && r.Start != prevEnd+1 {
			t.Fatalf("row %d starts at %d, previous row ended at %d", i, r.Start, prevEnd)
		}
		if r.EndCol > 20 {
			t.Errorf("row %d is %d columns wide, want at most 20", i, r.EndCol)
		}
		prevEnd = r.End
	}
}

func TestHeartWithVariationSelectorIsTwoColumnsByGraphemesAndOneByCodePoints(t *testing.T) {
	heart := "❤️"
	if got := ansi.GraphemeWidth.StringWidth(heart); got != 2 {
		t.Fatalf("grapheme width is %d, want 2", got)
	}
	if got := ansi.WcWidth.StringWidth(heart); got != 1 {
		t.Fatalf("code point width is %d, want 1", got)
	}

	byGraphemes := Layout(Graphemes(heart+"a"), textWidth, ansi.GraphemeWidth)
	if got := byGraphemes[0].Cols[1]; got != 2 {
		t.Errorf("by graphemes, the character after the heart is at column %d, want 2", got)
	}
	byCodePoints := Layout(Graphemes(heart+"a"), textWidth, ansi.WcWidth)
	if got := byCodePoints[0].Cols[1]; got != 1 {
		t.Errorf("by code points, the character after the heart is at column %d, want 1", got)
	}
}

func TestLayoutWrapsHeartsDifferentlyByMethod(t *testing.T) {
	text := strings.Repeat("❤️", 40)
	byGraphemes := rowTexts(text, ansi.GraphemeWidth)
	if got := len(Graphemes(byGraphemes[0])); got != 28 {
		t.Errorf("by graphemes, the first row holds %d hearts, want 28", got)
	}
	byCodePoints := rowTexts(text, ansi.WcWidth)
	if got := len(byCodePoints); got != 1 {
		t.Errorf("by code points, got %d rows, want 1", got)
	}
}
