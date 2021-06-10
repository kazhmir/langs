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
	m := re.BuildOne("ab|a[^]*c", act)
	fmt.Println(m)
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	m = re.BuildOne("[^]*c", act)
	fmt.Println(m)
	//
	//	err := m.RunStr("abcda")
	//	fmt.Printf("%#v\n", out)
	//	if err != nil {
	//		panic(err)
	//	}
}
