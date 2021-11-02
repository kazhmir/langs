package main

import (
	S "air/parser/state"
	P "air/parser"
	"fmt"
	"flag"
	"io/ioutil"
	"os"
)

var file = flag.String("f", "", "-f filename")

func main() {
	flag.Parse()
	source := ""
	if *file != "" {
		contents, err := ioutil.ReadFile(*file)
		if err != nil {
			panic(err)
		}
		source = string(contents)
	} else {
		fmt.Println("Expected a file: use -f filename")
		os.Exit(1)
	}
	st := S.NewState(source, *file)
	root := P.Program(st)
	fmt.Println(root)
}
