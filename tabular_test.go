package tabular

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPadding(t *testing.T) {
	b := New(Options{Padding: 2, PadChar: '.'})
	b.AddRow("this", "is", Right("a"), "test")
	b.AddRow(1, 2, Right(true), false)
	testOutput(t, b, `
this..is.....a..test
1.....2...true..false
`)
}

func TestMultiByte(t *testing.T) {
	b := New(Options{Padding: 2, PadChar: '.'})
	b.AddRow("⌘", "liberté", Right("égalité"), "☃")
	b.AddRow(1, 2, Right(true), "fraternité")
	testOutput(t, b, `
⌘..liberté..égalité..☃
1..2...........true..fraternité
`)
}

func TestRightAlign(t *testing.T) {
	b := New(Options{Padding: 2, PadChar: '.', AlignRight: true})
	b.AddRow("this", Right("is"), Left("a"), "test")
	b.AddRow(1, Left(2), Right(true), false)
	testOutput(t, b, `
this..is..a......test
...1..2...true..false
`)
}

func TestMinWidth(t *testing.T) {
	b := New(Options{MinWidth: 5, PadChar: '.'})
	b.AddRow("this", "is", Right("a"), "test")
	b.AddRow(1, 2, Right(true), false)
	testOutput(t, b, `
this.is.......atest
1....2.....truefalse
`)
}

func TestMinWidthAndPadding(t *testing.T) {
	b := New(Options{MinWidth: 5, Padding: 1, PadChar: '.'})
	b.AddRow("this", "is", Right("a"), "test")
	b.AddRow(1, 2, Right(true), false)
	testOutput(t, b, `
this..is........a.test
1.....2......true.false
`)
}

func TestMismatchedRows(t *testing.T) {
	b := New(Options{MinWidth: 3, Padding: 2, PadChar: '.'})
	b.AddRow("this", "is", Right("a"), "test")
	b.AddRow(1, Right(2), true, false, "extra")
	b.AddRow(" xyzabc blah blah")
	b.AddRow(Right(0.0), 9, 56, 12.34)
	testOutput(t, b, `
this...............is......a..test
1....................2..true..false..extra
 xyzabc blah blah
................0..9....56....12.34
`)
}

func TestEmptyCells(t *testing.T) {
	b := New(Options{MinWidth: 3, Padding: 2, PadChar: '.'})
	b.AddRow("this", "", Right("a"), "test")
	b.AddRow("", 2, "", false)
	testOutput(t, b, `
this.........a..test
......2.........false
`)
}

func TestDoubleRight(t *testing.T) {
	b := New(Options{Padding: 2, PadChar: ' '})
	b.AddRow("this", "is", Right(Right("a")), "test")
	b.AddRow(1, 2, Right(Right(Right(true))), false)
	testOutput(t, b, `
this  is     a  test
1     2   true  false
`)
}

func TestAddRowAfterWrite(t *testing.T) {
	b := New(Options{Padding: 2, PadChar: '.'})
	b.AddRow("this", "is", Right("a"), "test")
	b.AddRow(1, 2, Right(true), false)
	testOutput(t, b, `
this..is.....a..test
1.....2...true..false
`)
	b.AddRow(2, Right("x"), 1839834, "yes")
	testOutput(t, b, `
this..is........a..test
1.....2......true..false
2......x..1839834..yes
`)
}

func testOutput(t *testing.T, w *Buffer, want string) {
	t.Helper()
	want = strings.TrimPrefix(want, "\n")
	var buf bytes.Buffer
	w.WriteTo(&buf)
	got := buf.String()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("wrong output (-got, +want):\n%s", diff)
	}
}
