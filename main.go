// this is main package.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
	"unicode"
	"unicode/utf8"
)

const version = "0.0.1"

// FlagOptions remembers flag options.
type FlagOptions struct {
	PrintLines bool
	PrintBytes bool
	PrintWords bool
	PrintChars bool
}

// Counter remembers lines, words, bytes, and chars.
type Counter struct {
	lines int
	words int
	bytes int
	chars int
	mux   sync.Mutex
}

// Count adds lines, words, byte and chars.
func (c *Counter) Count(r io.Reader) (bool, error) {
	reader := bufio.NewReader(r)
	var wg sync.WaitGroup
	for {
		p := make([]byte, 4*1024)
		n, err := reader.Read(p)
		if n == 0 {
			break
		}
		c.AddBytes(n)

		wg.Add(1)
		go func() {
			var localCounter = &Counter{}
			bytesRead := p[:n]
			inField := false
			for i := 0; i < len(bytesRead); {
				r, size := utf8.DecodeRune(bytesRead[i:])
				wasInField := inField
				inField = !unicode.IsSpace(r)
				if inField && !wasInField {
					localCounter.words++
				}
				if r == '\n' {
					localCounter.lines++
				}
				localCounter.chars++
				i += size
			}
			c.Add(localCounter)
			wg.Done()
		}()

		if err == io.EOF {
			break
		}

		// fix word count between the read buffer
		next, err := reader.Peek(1)
		if err != nil && err != io.EOF {
			return false, err
		}
		if !unicode.IsSpace(rune(p[n-1 : n][0])) && !unicode.IsSpace(rune(next[0])) {
			c.AddWords(-1)
		}
	}
	wg.Wait()
	return true, nil
}

// Show shows flag options.
func (c *Counter) Show(opts FlagOptions, filename string) {
	if opts.PrintLines {
		fmt.Printf(" %7d", c.lines)
	}
	if opts.PrintWords {
		fmt.Printf(" %7d", c.words)
	}
	if opts.PrintBytes {
		fmt.Printf(" %7d", c.bytes)
	}
	if opts.PrintChars {
		fmt.Printf(" %7d", c.chars)
	}
	fmt.Printf(" %s\n", filename)
}

// Add adds lines, words, byte and chars.
func (c *Counter) Add(src *Counter) {
	c.mux.Lock()
	c.lines += src.lines
	c.bytes += src.bytes
	c.words += src.words
	c.chars += src.chars
	c.mux.Unlock()
}

// AddLines add lines.
func (c *Counter) AddLines(n int) {
	c.mux.Lock()
	c.lines += n
	c.mux.Unlock()
}

// AddBytes add bytes.
func (c *Counter) AddBytes(n int) {
	c.mux.Lock()
	c.bytes += n
	c.mux.Unlock()
}

// AddWords add words.
func (c *Counter) AddWords(n int) {
	c.mux.Lock()
	c.words += n
	c.mux.Unlock()
}

func parseFlagOptions() FlagOptions {
	var opts = FlagOptions{false, false, false, false}

	flag.BoolVar(&opts.PrintLines, "l", false, "print lines")
	flag.BoolVar(&opts.PrintBytes, "c", false, "print bytes")
	flag.BoolVar(&opts.PrintWords, "w", false, "print words")
	flag.BoolVar(&opts.PrintChars, "m", false, "print chars")
	flag.Parse()

	if opts.PrintChars {
		opts.PrintBytes = false
	}

	if !opts.PrintLines && !opts.PrintBytes && !opts.PrintWords && !opts.PrintChars {
		opts.PrintLines = true
		opts.PrintBytes = true
		opts.PrintWords = true
	}

	return opts
}

// Execute is main function.
func Execute(stdin io.Reader, stdout io.Writer, stderr io.Writer, opts FlagOptions) int {

	var totalCount = &Counter{}

	filenames := flag.Args()
	if len(filenames) == 0 {
		var c = &Counter{}
		_, err := c.Count(stdin)
		if err != nil {
			fmt.Fprintln(stderr, "stdin: count: ", err)
			return 1
		}
		c.Show(opts, "")
		return 0
	}

	for _, filename := range filenames {
		var c = &Counter{}
		fp, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(stderr, "%s: open: %s\n", filename, err)
			continue
		}
		_, err = c.Count(fp)
		if err != nil {
			fmt.Fprintf(stderr, "%s: count: %s\n", filename, err)
			continue
		}
		totalCount.Add(c)
		c.Show(opts, filename)
		fp.Close()
	}

	if len(filenames) > 1 {
		totalCount.Show(opts, "total")
	}

	return 0
}

func main() {
	opts := parseFlagOptions()
	os.Exit(Execute(os.Stdin, os.Stdout, os.Stderr, opts))
}
