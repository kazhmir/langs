package re

import (
	"sort"
	"sync"
)

type concSlice struct {
	i  int
	dt []string
	mu sync.Mutex
}

func (a *concSlice) Append(s string) {
	a.mu.Lock()
	if a.i >= len(a.dt) {
		a.dt = append(a.dt, make([]string, 10)...)
	}
	a.dt[a.i] = s
	a.i++
	a.mu.Unlock()
}

func (a *concSlice) Slice() []string {
	return a.dt[:a.i]
}

func exclude(a, b []rune) []rune {
	notInA := []rune{}
	mA := map[rune]int{}
	for _, c := range a {
		mA[c] = 0
	}
	mB := map[rune]int{}
	for _, c := range b {
		mB[c] = 0
		if _, ok := mA[c]; !ok { // add everything in b that is not in a
			notInA = append(notInA, c)
		}
	}
	return notInA
}

func intersect(a, b []rune) []rune {
	output := []rune{}
	beep := map[rune]int{}
	for _, c := range a {
		beep[c] = 0
	}
	for _, c := range b {
		if _, ok := beep[c]; ok {
			output = append(output, c)
		}
	}
	return output
}

/*
	Removes a from b, if any
*/
func rmIfAny(a, b []rune) []rune {
	out := []rune{}
	for _, cA := range a {
		for _, cB := range b {
			if cA == cB {
				goto nextRune
			}
		}
		out = append(out, cA)
	nextRune:
	}
	return out
}

func addIfAny(a, b []rune) []rune {
	a = append(a, b...)
	rSet := runeSlice(a)
	rSet = rmDuplicates(rSet)
	sort.Sort(&rSet)
	return rSet
}

type runeSlice []rune

func (rn *runeSlice) Len() int {
	return len(*rn)
}

func (rn *runeSlice) Less(i, j int) bool {
	return (*rn)[i] < (*rn)[j]
}

func (rn *runeSlice) Swap(i, j int) {
	c := (*rn)[i]
	(*rn)[i] = (*rn)[j]
	(*rn)[j] = c
}

func rmDuplicates(s runeSlice) runeSlice {
	dic := map[rune]struct{}{}
	for _, c := range s {
		dic[c] = struct{}{}
	}
	output := make(runeSlice, len(dic))
	var i int
	for c := range dic {
		output[i] = c
		i++
	}
	return output
}

func rmTr(slice []transition, indexes map[int]int) (out []transition) {
	for i, tr := range slice {
		if _, ok := indexes[i]; !ok {
			out = append(out, tr)
		}
	}
	return out
}
