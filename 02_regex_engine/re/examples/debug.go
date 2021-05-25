package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/kazhmir/ata"
)

var filePath = flag.String("f", "", "Use a file as input.")

func main() {
	flag.Parse()
	args := flag.Args()
	pattern := args[0]
	var input string
	if *filePath != "" {
		s, err := ioutil.ReadFile(*filePath)
		if err != nil {
			panic(err)
		}
		input = string(s)
	} else {
		input = args[1]
	}
	err := ata.Debug(pattern, input)
	if err != nil {
		fmt.Println(err)
	}
}
