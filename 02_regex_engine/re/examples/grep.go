package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/kazhmir/ata"
)

func main() {
	found := false
	find := func(*ata.Match) bool {
		found = true
		return false
	}
	scanner := bufio.NewScanner(os.Stdin)
	m := ata.BuildOne(os.Args[1], find)
	for scanner.Scan() {
		line := scanner.Text()
		err := m.RunStr(line)
		if err != nil {
			fmt.Println("err: ", err)
			continue
		}
		if found {
			fmt.Fprintln(os.Stdout, line)
			found = false
		}
	}
}
