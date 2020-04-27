package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func readFile(fileName string) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		f, err := os.Open(fileName)
		defer f.Close()

		if err != nil {
			return
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			return
		}
	}()
	return ch
}

var squeezeBlank = flag.Bool("s", false, "")

func doSqueezeBlank(lch <-chan string) <-chan string {
	lb := ""
	ch := make(chan string)
	go func() {
		defer close(ch)
		for l := range lch {
			if lb == "" && l == "" {
				continue
			}
			lb = l
			ch <- l
		}
	}()
	return ch
}

var number = flag.Bool("n", false, "")

func doNumber(lch <-chan string, n *int) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for l := range lch {
			*n++
			ch <- fmt.Sprintf("%d: %s", *n, l)
		}
	}()
	return ch
}

var numberNonblank = flag.Bool("b", false, "")

func doNumberNonblank(lch <-chan string, n *int) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for l := range lch {
			if l == "" {
				ch <- l
			} else {
				*n++
				ch <- fmt.Sprintf("%d: %s", *n, l)
			}
		}
	}()
	return ch
}

var showEnds = flag.Bool("e", false, "")

func doShowEnds(lch <-chan string) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for l := range lch {
			ch <- l + "$"
		}
	}()
	return ch
}

var showTabs = flag.Bool("t", false, "")

func doShowTabs(lch <-chan string) <-chan string {
	ch := make(chan string)
	replacer := strings.NewReplacer("\t", "^I")
	go func() {
		defer close(ch)
		for l := range lch {
			ch <- replacer.Replace(l)
		}
	}()
	return ch
}

func writeLine(line string) {
	fmt.Fprintln(os.Stdout, line)
}

func writeLines(lines <-chan string) {
	for l := range lines {
		writeLine(l)
	}
}

func main() {
	flag.Parse()
	fileNames := flag.Args()
	n := 0

	for _, fn := range fileNames {
		ch := readFile(fn)
		if *squeezeBlank {
			ch = doSqueezeBlank(ch)
		}
		if *showTabs {
			ch = doShowTabs(ch)
		}
		if *number && !*numberNonblank {
			ch = doNumber(ch, &n)
		}
		if *numberNonblank {
			ch = doNumberNonblank(ch, &n)
		}
		if *showEnds {
			ch = doShowEnds(ch)
		}
		writeLines(ch)
	}

}
