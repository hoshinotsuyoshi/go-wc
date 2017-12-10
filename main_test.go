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

func ExampleCounter_Show_lines() {
	c := main.Counter{}
	opts := main.FlagOptions{}
	opts.PrintLines = true
	reader := strings.NewReader("あいう\ndef\n")
	c.Count(reader)
	c.Show(&opts, "filename")
	// Output:
	// 2 filename
}

func ExampleCounter_Show_bytes() {
	c := main.Counter{}
	opts := main.FlagOptions{}
	opts.PrintBytes = true
	reader := strings.NewReader("あいう\ndef\n")
	c.Count(reader)
	c.Show(&opts, "filename")
	// Output:
	// 14 filename
}
