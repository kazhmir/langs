package re

import (
	//"fmt"
	"sort"
	"strconv"
)

// Thompson's construction
func createAtmt(n *node) *automaton {
	if n == nil { // empty set
		return NewAtmt(&Set{Items: []rune{}})
	}
	out := &automaton{start: &state{}, acc: &state{}}
	if n.tp == set {
		out.start.addTr(n.set, out.acc)
		return out
	}

	atmts := []*automaton{}
	for _, child := range n.children {
		atmts = append(atmts, createAtmt(child))
	}
	switch n.tp {
	case or:
		for _, atmt := range atmts {
			atmt.acc.addEmptyTr(out.acc)
			out.start.addEmptyTr(atmt.start)
		}
	case and:
		curr := out.start
		for _, atmt := range atmts {
			curr.addEmptyTr(atmt.start)
			curr = atmt.acc
		}
		curr.addEmptyTr(out.acc)
	case star: // guaranteed to have a single leaf
		newSt := &state{}
		atmts[0].acc.addEmptyTr(atmts[0].start) // ε transition from acc to start
		atmts[0].acc.addEmptyTr(newSt)          // acc to new
		out.start.addEmptyTr(atmts[0].start)    // out.start to start
		out.start.addEmptyTr(newSt)             // out.start to new
		newSt.addEmptyTr(out.acc)               // new to out.acc
	}
	return out
}

/*
	Returns a UNORDERED slice of states reachable through ε transitions
	from the given state.
	ε transitions are regarded as nil slices (which are illegal in the
	context of this package)
*/
func eFind(st *state) []*state {
	out := make([]*state, 10)
	m := make(map[*state]int, 10)
	return underEFind(st, &m, &out, 0)
}

/* underEFind is a recursive function that keeps track of the number
of items inside the *prev map. Each *state in the map has it's index
as the value, so in the end this *state can be added to the out slice.
*/
func underEFind(st *state, prev *map[*state]int, out *[]*state, iter int) []*state {
	start := iter
	(*prev)[st] = iter
	iter++
	for _, tr := range st.trans {
		if tr.set == nil { // ε transition
			if _, ok := (*prev)[tr.next]; !ok {
				for _, sta := range underEFind(tr.next, prev, out, iter) {
					(*prev)[sta] = iter
					iter++
				}
			}
		}
	}

	for sta, i := range *prev {
		if iter > len(*out) {
			(*out) = append(*out, make([]*state, 10)...)
		}
		(*out)[i] = sta
	}
	return (*out)[start:len(*prev)]
}

//Joins all states from the list into the given state.
func fuse(output *state, states []*state) *state {
	if len(states) == 0 {
		return nil
	}
	trs := []transition{}

	for _, st := range states {
		if st.act != nil { // if any state is an accepting state
			output.act = st.act // the output will also be an accepting state
		}
		for _, tr := range st.trans {
			trs = append(trs, tr)
		}
	}
	output.trans = trs
	return output
}

// returns a state containing all other states reachable through ε transitions
func eClosure(starters ...*state) (*state, string) {
	eC := []*state{}
	for _, st := range starters {
		eC = append(eC, eFind(st)...)
	}
	id := createID(eC)
	out := fuse(&state{}, eC)
	out.trans = rmEmpty(out.trans)
	return out, id
}

/*Creating an ID is needed to keep track of which combinations
of states we have already made.*/
func createID(in []*state) string {
	mp := map[string]int{}
	for _, st := range in {
		mp[strconv.Itoa(st.i)] = 0
	}
	lst := make([]string, len(mp))
	i := 0
	for s := range mp {
		lst[i] = s
		i++
	}

	sort.Strings(lst)

	out := ""
	for i := range lst {
		out += lst[i]
	}
	return out
}

/*removes ε and impossible transitions*/
func rmEmpty(trs []transition) (out []transition) {
	for _, tr := range trs {
		if tr.set != nil && len(tr.set.Items) > 0 {
			out = append(out, tr)
		}
	}
	return out
}

/*
	This function depends on all states being correctly enumerated,
	so it can safely generate ids for each computed state,
	use state.enum function to enumerate.
*/
func powerSet(prev *map[string]*state, starters ...*state) *state {
	eSt, id := eClosure(starters...)
	if st, ok := (*prev)[id]; ok { // if already computed the state
		return st
	}
	(*prev)[id] = eSt

	newTrans := intersectAll(eSt.trans)
	eSt.trans = []transition{}
	for set, states := range newTrans {
		tr := transition{
			set:  set,
			next: powerSet(prev, states...),
		}
		eSt.trans = append(eSt.trans, tr)
	}

	return eSt
}

/*
	This is an heuristic approach to get all the intersections between n sets,
*/
func intersectAll(trans []transition) map[*Set][]*state {
	out := map[*Set][]*state{}
	for done := true; done; { // we will iterate until there are no more intersections
		done = false
		for i := range trans {
			finalSect := trans[i].set
			toRemove := []int{i}
			states := []*state{trans[i].next}
			for j, jTr := range trans {
				if i != j { // not itself
					sect := finalSect.intersect(jTr.set)
					if len(sect.Items) > 0 {
						finalSect = sect
						toRemove = append(toRemove, j)
						states = append(states, jTr.next)
						done = true
					}
				}
			}
			if len(toRemove) > 1 { // contains something more than the initial set
				for _, x := range toRemove {
					trans[x].set.rm(finalSect)
				}
			}
			out[finalSect] = states
		}
	}
	for i := range trans {
		if len(trans[i].set.Items) > 0 {
			out[trans[i].set] = []*state{trans[i].next}
		}
	}
	return out
}
