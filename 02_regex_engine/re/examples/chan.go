package main

import (
	"fmt"
	"strings"

	"github.com/kazhmir/ata"
)

/*FindAllChan returns a machine, the matching strings
are send through the given channel.
You have to call machine.Run (tipically in a new goroutine)
on the input to start sending matches through the out channel.
*/
func FindAllChan(pattern string, out chan string) *ata.Machine {
	var get ata.Action = func(a *ata.Match) bool {
		out <- a.S
		return false
	}
	syntax := map[string]ata.Action{
		pattern: get,
	}
	m := ata.Build(syntax)
	return m
}

func main() {
	s1 := strings.NewReader("abababababababababababab")

	out := make(chan string)
	m := FindAllChan("a", out)
	go func() {
		err := m.Run(s1)
		if err != nil {
			fmt.Println(err)
		}
		close(out)
	}()
	for match := range out {
		fmt.Println(match)
	}
}
