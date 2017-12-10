package main_test

import (
	"fmt"
	"github.com/hoshinotsuyoshi/go-wc"
	"strings"
)

func ExampleCounter() {
	c := main.Counter{}
	reader := strings.NewReader("abc\n")
	val, _ := c.Count(reader)
	fmt.Println(val)
	// Output:
	// true
}

func ExampleCounter_Show() {
	c := main.Counter{}
	opts := main.FlagOptions{}
	opts.PrintLines = true
	reader := strings.NewReader("abc\ndef\n")
	c.Count(reader)
	c.Show(&opts, "filename")
	// Output:
	// 2 filename
}
