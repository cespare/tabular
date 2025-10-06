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
//   - Omits several lesser-used features
package tabular

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

// Options configure a [Buffer].
type Options struct {
	MinWidth   int  // Minimum cell width (not including padding).
	Padding    int  // Padding between each cell.
	PadChar    byte // The character to use for padding.
	AlignRight bool // Align cells to the right by default.
}

// A Buffer stores rows of text and prints them as a table.
// It assumes that each Unicode code point has a width of 1.
type Buffer struct {
	opts Options
	buf  []byte
	rows [][]cell
}

type cell struct {
	wb    int  // width in bytes
	wr    int  // width in runes
	right bool // whether to right-align
}

// New constructs a [Buffer] with options.
func New(opts Options) *Buffer {
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
		c.wr = utf8.RuneCountInString(s)
		row[i] = c
		b.buf = append(b.buf, s...)
	}
	b.rows = append(b.rows, row)
}

// WriteTo writes the buffered rows as a text table.
func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	var widths []int
	for _, row := range b.rows {
		for i, c := range row {
			if i < len(widths) {
				if c.wr > widths[i] {
					widths[i] = c.wr
				}
			} else {
				widths = append(widths, c.wr)
			}
		}
	}
	for i, n := range widths {
		if n < b.opts.MinWidth {
			widths[i] = b.opts.MinWidth
		}
	}
	var maxPad int
	for _, n := range widths {
		if n > maxPad {
			maxPad = n
		}
	}
	if b.opts.Padding > maxPad {
		maxPad = b.opts.Padding
	}
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
				line = append(line, padBuf[:width-c.wr]...)
			}
			line = append(line, b.buf[i:i+c.wb]...)
			i += c.wb
			if !c.right && j < len(row)-1 {
				line = append(line, padBuf[:width-c.wr]...)
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
