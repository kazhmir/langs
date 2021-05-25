package re

/*This file has the types used in the automata.go part of the package*/

import (
	"fmt"
	"sort"
)

type state struct {
	i     int
	trans []transition
	act   Action
}

func (st *state) move(c rune) *state {
	for _, tr := range st.trans {
		if tr.set.Contains(c) {
			return tr.next
		}
	}
	// this means an error state,
	// the Run function will deal with this by restarting the automata
	return nil
}

func (st *state) String() string {
	format := "S%v{"
	values := []interface{}{st.i}

	if len(st.trans) > 0 {
		format += "t: %v, "
		values = append(values, st.trans)
	}
	if st.act != nil {
		format += "act: %v"
		values = append(values, st.act)
	}
	format += "}"
	return fmt.Sprintf(format, values...)
}

func (st *state) Enum(prev *map[*state]int) {
	(*prev)[st] = len(*prev)
	st.i = (*prev)[st]
	for j := 0; j < len(st.trans); j++ {
		if n := st.trans[j].next; n != nil {
			if _, ok := (*prev)[n]; !ok { // prevent infinite loop
				n.Enum(prev)
			}
		}
	}
}

func (st *state) addEmptyTr(next *state) {
	st.trans = append(st.trans, transition{set: nil, next: next})
}

func (st *state) addTr(set *Set, next *state) {
	tr := transition{set: set, next: next}
	st.trans = append(st.trans, tr)
}

type Set struct {
	Items   []rune
	Negated bool
}

func (s *Set) String() string {
	var val string
	if s.Negated {
		val = "not"
	}
	return val + `"` + string(s.Items) + `"`
}

/*
tr.Negated defines if it's negated or not. Which means if true,
it'll invert the result of the binary search. This is why
here we always return tr.Negated or !tr.Negated. Contains becomes not contains.
*/
func (s *Set) Contains(c rune) bool {
	bot, top := 0, len(s.Items)-1
	if len(s.Items) == 0 || c < s.Items[bot] || c > s.Items[top] {
		return s.Negated
	}
	var curr rune
	var mid int
	for bot <= top {
		mid = int((top + bot) / 2)
		curr = s.Items[mid]
		if c == curr {
			return !s.Negated
		}
		if c < curr {
			top = mid - 1
		}
		if c > curr {
			bot = mid + 1
		}
	}
	return s.Negated
}

func (s *Set) rm(other *Set) {
	switch true {
	case s.Negated && other.Negated:
		s.Items = exclude(s.Items, other.Items)
		s.Negated = false
	case s.Negated && !other.Negated:
		s.Items = addIfAny(s.Items, other.Items)
	case !s.Negated && other.Negated:
		s.Items = intersect(s.Items, other.Items)
	case !s.Negated && !other.Negated:
		s.Items = rmIfAny(s.Items, other.Items)
	}
}

func (s *Set) intersect(other *Set) *Set {
	out := &Set{Items: []rune{}}
	switch true {
	case s.Negated && other.Negated:
		out.Items = addIfAny(s.Items, other.Items)
		out.Negated = true
	case s.Negated && !other.Negated:
		out.Items = exclude(s.Items, other.Items)
	case !s.Negated && other.Negated:
		out.Items = exclude(other.Items, s.Items)
	case !s.Negated && !other.Negated:
		out.Items = intersect(s.Items, other.Items)
	}
	return out
}

// TODO: this code doesn't work if the Other is negated
func (s *Set) add(other *Set) {
	if s.Negated {
		s.Items = rmIfAny(other.Items, s.Items)
	} else {
		s.Items = append(s.Items, other.Items...)
		rSet := runeSlice(s.Items)
		rSet = rmDuplicates(rSet)
		sort.Sort(&rSet)
		s.Items = rSet
	}
}

type transition struct {
	set  *Set
	next *state
}

func (tr transition) String() string {
	set := "ε"
	next := "nil"
	if tr.next != nil {
		next = "S" + fmt.Sprint(tr.next.i)
	}
	if tr.set != nil {
		set = tr.set.String()
	}
	return fmt.Sprintf("{%s -> %s}", set, next)
}

type automaton struct {
	start, acc *state
}

func NewAtmt(s *Set) *automaton {
	start := &state{}
	acc := &state{}
	start.addTr(s, acc)
	return &automaton{start, acc}
}

func prettyPrint(states *map[*state]int) string {
	order := make([]*state, len(*states))
	for k, v := range *states {
		order[v] = k
	}
	output := ""
	for _, st := range order {
		output += fmt.Sprintf("%v\n", st)
	}
	return output
}
