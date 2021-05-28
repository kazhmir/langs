package re

/*
returns a set that is the intersection between 'a' and 'b'
*/
func intersection(a, b map[rune]struct{}) map[rune]struct{} {
	out := map[rune]struct{}{}
	for c := range b {
		if _, ok := a[c]; ok {
			out[c] = struct{}{}
		}
	}
	return out
}

/*
returns a set that is 'a' without the items in 'b', if any
*/
func setDifference(a, b map[rune]struct{}) map[rune]struct{} {
	out := map[rune]struct{}{}
	for c := range a {
		if _, ok := b[c]; !ok {
			out[c] = struct{}{}
		}
	}
	return out
}

/*
returns the union of both sets
*/
func union(a, b map[rune]struct{}) map[rune]struct{} {
	out := map[rune]struct{}{}
	for c := range a {
		out[c] = struct{}{}
	}
	for c := range b {
		out[c] = struct{}{}
	}
	return out
}

func rmTr(slice []transition, indexes map[int]int) (out []transition) {
	for i, tr := range slice {
		if _, ok := indexes[i]; !ok {
			out = append(out, tr)
		}
	}
	return out
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
