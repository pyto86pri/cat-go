package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func readFile(fileName string) <-chan *Line {
	ch := make(chan *Line)
	go func() {
		defer close(ch)
		f, err := os.Open(fileName)
		defer f.Close()

		if err != nil {
			ch <- &Line{err: errors.Wrap(err, "can't find file")}
			return
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			ch <- &Line{raw: scanner.Text()}
		}
		if err := scanner.Err(); err != nil {
			ch <- &Line{err: errors.Wrap(err, "can't scan file")}
			return
		}
	}()
	return ch
}

var squeezeBlank = flag.Bool("s", false, "")

func doSqueezeBlank(ctx context.Context, lch <-chan *Line) <-chan *Line {
	lb := ""
	ch := make(chan *Line)
	go func() {
		defer close(ch)
		for l := range lch {
			if l.err == nil {
				if lb == "" && l.raw == "" {
					continue
				}
				lb = l.raw
			}
			select {
			case <-ctx.Done():
				return
			case ch <- l:
			}
		}
	}()
	return ch
}

var number = flag.Bool("n", false, "")

func doNumber(ctx context.Context, lch <-chan *Line, n *int) <-chan *Line {
	ch := make(chan *Line)
	go func() {
		defer close(ch)
		for l := range lch {
			if l.err == nil {
				*n++
				l.prefix += fmt.Sprintf("%d: ", *n)
			}
			select {
			case <-ctx.Done():
				return
			case ch <- l:
			}
		}
	}()
	return ch
}

var numberNonblank = flag.Bool("b", false, "")

func doNumberNonblank(ctx context.Context, lch <-chan *Line, n *int) <-chan *Line {
	ch := make(chan *Line)
	go func() {
		defer close(ch)
		for l := range lch {
			if l.err == nil {
				if l.raw != "" {
					*n++
					l.prefix += fmt.Sprintf("%d:", *n)
				}
			}
			select {
			case <-ctx.Done():
				return
			case ch <- l:
			}
		}
	}()
	return ch
}

var showEnds = flag.Bool("e", false, "")

func doShowEnds(ctx context.Context, lch <-chan *Line) <-chan *Line {
	ch := make(chan *Line)
	go func() {
		defer close(ch)
		for l := range lch {
			if l.err == nil {
				l.suffix += "$"
			}
			select {
			case <-ctx.Done():
				return
			case ch <- l:
			}
		}
	}()
	return ch
}

var showTabs = flag.Bool("t", false, "")

func doShowTabs(ctx context.Context, lch <-chan *Line) <-chan *Line {
	ch := make(chan *Line)
	go func() {
		defer close(ch)
		for l := range lch {
			if l.err == nil {
				l.replacer = append(l.replacer, "\t", "^I")
			}
			select {
			case <-ctx.Done():
				return
			case ch <- l:
			}
		}
	}()
	return ch
}

func writeLine(line *Line) {
	fmt.Fprintln(os.Stdout, line)
}

func writeLines(lines <-chan *Line) {
	for l := range lines {
		if l.err != nil {
			log.Fatal(l.err)
		}
		writeLine(l)
	}
}

// Line express line component
type Line struct {
	prefix   string
	raw      string
	suffix   string
	replacer []string
	err      error
}

func (l *Line) String() string {
	replacer := strings.NewReplacer(l.replacer...)
	return fmt.Sprint(l.prefix, replacer.Replace(l.raw), l.suffix)
}

func main() {
	flag.Parse()
	fileNames := flag.Args()
	n := 0
	ctx := context.TODO()

	for _, fn := range fileNames {
		ch := readFile(fn)
		if *squeezeBlank {
			ch = doSqueezeBlank(ctx, ch)
		}
		if *showTabs {
			ch = doShowTabs(ctx, ch)
		}
		if *number && !*numberNonblank {
			ch = doNumber(ctx, ch, &n)
		}
		if *numberNonblank {
			ch = doNumberNonblank(ctx, ch, &n)
		}
		if *showEnds {
			ch = doShowEnds(ctx, ch)
		}
		writeLines(ch)
	}

}
