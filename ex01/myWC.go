package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"unicode/utf8"
)

func countWords(argv string, group *sync.WaitGroup) {
	text, _ := os.ReadFile(argv)
	words := strings.Fields(string(text))
	count := len(words)
	fmt.Printf("%d	%s\n", count, argv)
	group.Done()
}

func countLines(argv string, group *sync.WaitGroup) {
	text, _ := os.ReadFile(argv)
	count := bytes.Count(text, []byte("\n")) + 1
	fmt.Printf("%d	%s\n", count, argv)
	group.Done()
}

func countCharacters(argv string, group *sync.WaitGroup) {
	text, _ := os.ReadFile(argv)
	count := utf8.RuneCount(text)
	fmt.Printf("%d	%s\n", count, argv)
	group.Done()
}

func main() {
	w := flag.Bool("w", false, "count words")
	l := flag.Bool("l", false, "count lines")
	m := flag.Bool("m", false, "count characters")
	flag.Parse()

	if len(os.Args) < 2 || len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "Pass one or more filenames")
		fmt.Fprintln(os.Stderr, "You can also provide '-w', '-l', '-m' flags to specify output")
		os.Exit(1)
	}

	var group sync.WaitGroup
	group.Add(len(flag.Args()))
	if *w || (!*w && !*l && !*m) {
		for i := range flag.Args() {
			go countWords(flag.Arg(i), &group)
		}
	}
	if *l {
		for i := range flag.Args() {
			go countLines(flag.Arg(i), &group)
		}
	}
	if *m {
		for i := range flag.Args() {
			go countCharacters(flag.Arg(i), &group)
		}
	}
	group.Wait()
}
