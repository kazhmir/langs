package re

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestLexString(t *testing.T) {
	fmt.Println("Testing lexer")
	for _, tst := range tests {
		t.Run("  "+tst.re, func(t *testing.T) {
			ans := lexString(tst.re)
			if !reflect.DeepEqual(ans, tst.tokens) {
				t.Errorf("got %v, wanted %v", ans, tst.tokens)
			}
		})
	}
}

func TestParse(t *testing.T) {
	fmt.Println("Testing parser")
	p := parser{}
	for _, tst := range tests {
		name := fmt.Sprintf("%v", tst.tokens)
		t.Run(name, func(t *testing.T) {
			ans := p.run(tst.tokens)
			if !treeEqual(ans, tst.root) {
				t.Errorf("incorrect tree for: %#v\n\n got:\n%v \n\n want:\n%v\n\rpath:\n%v", tst.re, ans, tst.root, p.path)
			}
		})
	}
}

func TestFindAllString(t *testing.T) {
	fmt.Println("Testing FindAllStrings")
	for _, tst := range tests {
		name := fmt.Sprintf("%v:%v", tst.re, tst.input)
		t.Run(name, func(t *testing.T) {
			txt := strings.NewReader(tst.input)
			//out, err := FindAllString(tst.re, txt)
			out := make([]string, 0)
			act := func(mat *Match) bool {
				out = append(out, mat.S)
				return false
			}
			m := BuildOne(tst.re, act)
			err := m.Run(txt)
			if err != nil {
				fmt.Println(err)
			}
			if !reflect.DeepEqual(tst.matches, out) {
				mp := map[*state]int{}
				m.Start.Enum(&mp)
				t.Errorf("%v\n\ngot %#v, wanted %#v", prettyPrint(&mp), out, tst.matches)
				//t.Errorf("got %#v, wanted %#v", out, tst.matches)
			}
		})
	}
}

func treeEqual(a, b *node) bool {
	if a == nil || b == nil {
		return checkNil(a, b)
	}
	if len(a.children) != len(b.children) {
		return false
	}
	for i := 0; i < len(a.children); i++ {
		if !reflect.DeepEqual(a.children[i].set, b.children[i].set) {
			return false
		}
		if !treeEqual(a.children[i], b.children[i]) {
			return false
		}
	}
	return true
}

func checkNil(a, b *node) bool {
	if a == nil && b != nil || a != nil && b == nil {
		return false
	}
	return true
}
