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

	for _, t := range st.trans {
		format += "t: %v, "
		values = append(values, t.String())
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
	st.trans = append(st.trans, transition{epsilon: true, next: next})
}

func (st *state) addTr(set Set, next *state) {
	tr := transition{set: set, next: next}
	st.trans = append(st.trans, tr)
}

type Set struct {
	Items   map[rune]struct{}
	Negated bool
}

func NewSet(s string, n bool) *Set {
	out := &Set{
		Items:   make(map[rune]struct{}, len(s)),
		Negated: n,
	}
	for _, r := range s {
		out.Items[r] = struct{}{}
	}
	return out
}

func (s *Set) String() string {
	var val string
	if s.Negated {
		val = "not"
	}
	items := runeSlice{}
	for c := range s.Items {
		items = append(items, c)
	}
	sort.Sort(&items)
	return fmt.Sprintf("%v\"%v\"", val, string(items))
}

/*
tr.Negated defines if it's negated or not. Which means if true,
it'll invert the result of the binary search. This is why
here we always return tr.Negated or !tr.Negated. Contains becomes not contains.
*/
func (s *Set) Contains(c rune) bool {
	_, ok := s.Items[c]
	if ok {
		return !s.Negated
	}
	return s.Negated
}

func (s *Set) IsNotEmpty() bool {
	if len(s.Items) > 0 {
		return true
	}
	return s.Negated
	//return s.Negated || len(s.Items) > 0
}

func (s *Set) rm(other Set) {
	switch true {
	case s.Negated && other.Negated:
		s.Items = setDifference(other.Items, s.Items) // [^a] - [^b] = [b]
		s.Negated = false
	case s.Negated && !other.Negated:
		s.Items = union(s.Items, other.Items) // [^a] - [b] = [^ab]
	case !s.Negated && other.Negated:
		s.Items = intersection(s.Items, other.Items) // [a] - [^b] = []
	case !s.Negated && !other.Negated:
		s.Items = setDifference(s.Items, other.Items) // [a] - [b] = [a]
	}
}

func (s Set) intersect(other Set) *Set {
	out := NewSet("", false)
	switch true {
	case s.Negated && other.Negated:
		out.Items = union(s.Items, other.Items) // [^a] ∩ [^b] = [^ab]
		out.Negated = true
	case s.Negated && !other.Negated:
		out.Items = setDifference(other.Items, s.Items) // [^a] ∩ [b] = [b]
	case !s.Negated && other.Negated:
		out.Items = setDifference(s.Items, other.Items) // [a] ∩ [^b] = [a]
	case !s.Negated && !other.Negated:
		out.Items = intersection(s.Items, other.Items) // [a] ∩ [b] = []
	}
	return out
}

func (s *Set) add(other Set) {
	switch true {
	case s.Negated && other.Negated:
		s.Items = intersection(s.Items, other.Items) // [^a] ∪ [^b] = [^];	[^a] ∪ [^a] = [^a]; 	[^a] ∪ [^ab] = [^a];
		s.Negated = true
	case s.Negated && !other.Negated:
		s.Items = setDifference(other.Items, s.Items) // [^a] ∪ [b] = [^a];	 [^a] ∪ [a] = [^];
		s.Negated = true
	case !s.Negated && other.Negated:
		s.Items = setDifference(s.Items, other.Items) // [a] ∪ [^b] = [^b];	 [a] ∪ [^a] = [^];
		s.Negated = true
	case !s.Negated && !other.Negated:
		s.Items = union(s.Items, other.Items) // [a] ∪ [b] = [ab];	 [a] ∪ [a] = [a];
	}
}

type transition struct {
	set  Set
	epsilon bool
	next *state
}

func (tr transition) String() string {
	set := "ε"
	next := "nil"
	if tr.next != nil {
		next = "S" + fmt.Sprint(tr.next.i)
	}
	if !tr.epsilon {
		set = tr.set.String()
	}
	return fmt.Sprintf("{%s -> %s}", set, next)
}

type automaton struct {
	start, acc *state
}

func NewAtmt(s Set) *automaton {
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
