package main

import (
	"bufio"
	"fmt"
	"github.com/kazhmir/ata"
	"os"
)

/*FindAllChan returns a machine, the matching strings
are send through the given channel.
You have to call machine.Run (tipically in a new goroutine)
on the input to start sending matches through the out channel.
*/
func FindAllChan(pattern string, out chan string) *ata.Machine {
	get := func(mat *ata.Match) bool {
		out <- mat.S
		return false
	}
	m := ata.BuildOne(pattern, get)
	return m
}

/*If you want to print the whole line just like grep,
use the pattern: "\N*<yourpattern>\N*\n"
*/
func main() {
	nOfFiles := len(os.Args[2:])
	pattern := os.Args[1]
	matches := make(chan string, nOfFiles)
	signal := make(chan struct{})
	m := FindAllChan(pattern, matches)
	for _, filename := range os.Args[2:] {
		f, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		go func() {
			err := m.Run(bufio.NewReader(f))
			if err != nil {
				fmt.Println(err)
			}
			signal <- struct{}{}
		}()
	}

	for i := nOfFiles; i > 0; {
		select {
		case match := <-matches:
			fmt.Println(match)
		case <-signal:
			i--
		}
	}
	close(matches)
	for match := range matches { // empty buffer
		fmt.Println(match)
	}
}
