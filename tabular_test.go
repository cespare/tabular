package tabular

import (
	"strings"
	"testing"

	"github.com/pkg/diff"
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

func TestMultibyte(t *testing.T) {
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

func TestWriteToResets(t *testing.T) {
	b := New(Options{Padding: 2, PadChar: '.'})
	b.AddRow("this", "is", Right("a"), "test")
	b.AddRow(1, 2, Right(true), false)
	testOutput(t, b, `
this..is.....a..test
1.....2...true..false
`)
	b.AddRow(2, Right("x"), 1839834, "yes")
	testOutput(t, b, `
2..x..1839834..yes
`)
}

func TestMultibyteWidth(t *testing.T) {
	b := New(Options{Padding: 2, PadChar: '.'})
	b.AddRow("hello", "world", "!")
	b.AddRow(Right("☺"), "你好", "x")
	testOutput(t, b, `
hello..world..!
....☺..你好...x
`)
}

func TestCSI(t *testing.T) {
	b := New(Options{Padding: 2, PadChar: '.'})
	b.AddRow("\x1b[1mabc\x1b[0m", "def", "ghi")
	b.AddRow("jkl", "\x1b[32mmno\x1b[0m", "pqr")
	testOutput(t, b, "\x1b[1mabc\x1b[0m..def..ghi\njkl..\x1b[32mmno\x1b[0m..pqr")
}

func TestDefaultOptions(t *testing.T) {
	b := New(Options{Padding: 2})
	b.AddRow("a", "bbbbbb", "c")
	b.AddRow("dddd", "e", "fff")
	testOutput(t, b, `
a     bbbbbb  c
dddd  e       fff
`)
	b = New(Options{})
	b.AddRow("a", "bb")
	b.AddRow("ccc", "d")
	testOutput(t, b, `
a  bb
cccd
`)
}

func testOutput(t *testing.T, w *Buffer, want string) {
	t.Helper()
	want = strings.TrimPrefix(want, "\n")
	var buf strings.Builder
	w.WriteTo(&buf)
	got := buf.String()
	if got == want {
		return
	}
	var diffBuf strings.Builder
	if err := diff.Text("got", "want", got, want, &diffBuf); err != nil {
		panic(err)
	}
	result := diffBuf.String()
	if result == "--- got\n+++ want\n" {
		return
	}
	t.Error(result)
}
