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
