package main

import (
	"fmt"
	"re/re"
)

func main() {
	act := func(){}
	m := re.BuildOne("ab|a[^]*c", act)
	fmt.Println(m)
	m = re.BuildOne("[^]*c", act)
	fmt.Println(m)
	//
	//	err := m.RunStr("abcda")
	//	fmt.Printf("%#v\n", out)
	//	if err != nil {
	//		panic(err)
	//	}
}
