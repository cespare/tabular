// Package tabular implements a simple tabular formatter for text.
//
// The package is an alternative to text/tabwriter for some use cases.
// In comparison with text/tabwriter, this package:
//
//   - Is oriented around rows and cells, not tab-delimited text
//   - Inserts padding between cells, not after each cell
//   - Does not count padding toward min-width
//   - Does not insert any padding after the last cell of each row
//   - Allows for per-cell right-alignment
//   - Omits several lesser-used features of tabwriter
//   - Attempts to guess the width of multibyte code points
//   - Ignores ANSI CSI sequences for width calculations
package tabular

import (
	"cmp"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strings"

	"github.com/mattn/go-runewidth"
)

// Options configure a [Buffer].
type Options struct {
	// MinWidth is the minimum cell width (not including padding).
	MinWidth int
	// Padding is the number of padding characters between each cell.
	Padding int
	// PadChar is the character to use for padding.
	// The default, if PadChar is zero, is to use space.
	PadChar byte
	// AlignRight controls whether all cells are right-aligned by default.
	AlignRight bool
}

// A Buffer stores rows of text and prints them as a table.
//
// To calculate the width of a table cell, it strips out any ANSI CSI sequences
// from the text and then uses go-runewidth to guess the width of the resulting
// text.
type Buffer struct {
	opts Options
	buf  []byte
	rows [][]cell
}

type cell struct {
	wb    int  // width in bytes
	wc    int  // width in visible cells
	right bool // whether to right-align
}

// New constructs a [Buffer] with options.
func New(opts Options) *Buffer {
	opts.PadChar = cmp.Or(opts.PadChar, ' ')
	return &Buffer{opts: opts}
}

// Right marks a value passed to [Buffer.AddRow] for right alignment.
func Right(v any) any {
	return right{v}
}

type right struct{ v any }

func (r right) String() string {
	return fmt.Sprint(r.v)
}

// Left marks a value passed to [Buffer.AddRow] for left alignment.
func Left(v any) any {
	return left{v}
}

type left struct{ v any }

func (l left) String() string {
	return fmt.Sprint(l.v)
}

// AddRow adds a row of values to the buffer.
//
// Each value is turned into a string using the same formatting as fmt.Sprint.
func (b *Buffer) AddRow(vs ...any) {
	row := make([]cell, len(vs))
	for i, v := range vs {
		c := cell{right: b.opts.AlignRight}
		if r, ok := v.(right); ok {
			v = r.v
			c.right = true
		}
		if l, ok := v.(left); ok {
			v = l.v
			c.right = false
		}
		s := fmt.Sprint(v)
		c.wb = len(s)
		c.wc = cellWidth(s)
		row[i] = c
		b.buf = append(b.buf, s...)
	}
	b.rows = append(b.rows, row)
}

var csiRegexp = regexp.MustCompile(`\x1b\[[\x30-\x3f]*[\x20-\x2f]*[\x40-\x7e]`)

func cellWidth(s string) int {
	// Strip out all ANSI CSI sequences. In this context, they are typically
	// used for styling and coloring text.
	s = csiRegexp.ReplaceAllString(s, "")
	return runewidth.StringWidth(s)
}

// WriteTo writes the buffered rows as a text table.
func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	var widths []int
	for _, row := range b.rows {
		for i, c := range row {
			if i < len(widths) {
				widths[i] = max(widths[i], c.wc)
			} else {
				widths = append(widths, c.wc)
			}
		}
	}
	for i, w := range widths {
		widths[i] = max(w, b.opts.MinWidth)
	}
	maxPad := max(slices.Max(widths), b.opts.Padding)
	padBuf := strings.Repeat(string(b.opts.PadChar), maxPad)

	var i int
	var line []byte
	var written int64
	for _, row := range b.rows {
		line = line[:0]
		for j, c := range row {
			if j > 0 {
				line = append(line, padBuf[:b.opts.Padding]...)
			}
			width := widths[j]
			if c.right {
				line = append(line, padBuf[:width-c.wc]...)
			}
			line = append(line, b.buf[i:i+c.wb]...)
			i += c.wb
			if !c.right && j < len(row)-1 {
				line = append(line, padBuf[:width-c.wc]...)
			}
		}
		line = append(line, '\n')
		n, err := w.Write(line)
		written += int64(n)
		if err != nil {
			return written, err
		}
	}
	return written, nil
}
