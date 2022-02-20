package tabular_test

import (
	"os"

	"github.com/cespare/tabular"
)

func ExampleBuffer() {
	b := tabular.New(tabular.Options{Padding: 2, PadChar: ' ', AlignRight: true})
	b.AddRow("x", "x²", "x³")
	for _, x := range []float64{0, 4, 8, 12} {
		b.AddRow(x, x*x, x*x*x)
	}
	b.WriteTo(os.Stdout)

	// Output:
	//  x   x²    x³
	//  0    0     0
	//  4   16    64
	//  8   64   512
	// 12  144  1728
}
