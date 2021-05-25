package main

import (
	"fmt"
	"re/re"
)

func main() {
	out := []string{}
	act := func(mat *re.Match) bool {
		out = append(out, mat.S)
		return false
	}
	m := re.BuildOne("a", act)
	err := m.RunStr("abcda")
	fmt.Printf("%#v\n", out)
	if err != nil {
		panic(err)
	}
}
