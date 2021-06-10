package re

import (
	//"fmt"
	"log"
	"strings"
)

/*Action is a function that receives a pointer to the machine Head
Actions are called on accepting states, according to syntax
*/
type Action func()

/*Machine is the automaton generated by the Build function
It contains structural data regarding the underlying automaton
*/
type Machine struct {
	Start   *state
	Pattern string
	Syntax  map[string]Action
}

func (m *Machine) String() string {
	dt := &map[*state]int{}
	m.Start.Enum(dt)
	return m.Pattern + "\n" + prettyPrint(dt) + "\n"
}

func (m *Machine) FullyMatches(s string) bool {
    currState := m.Start
    for _, r := range s {
	next := currState.move(r)
	if next == nil { // move returns nil when rune is not accepted
    		return false
	}
	currState = next
    }
    if currState.act != nil { // if ended up in a matching state
        return true
    }
    return false
}

/*BuildOne returns a machine for a single pattern and action.
 */
func BuildOne(pattern string, act Action) *Machine {
	start := compile(pattern, act).start
	return &Machine{
		Start:   start,
		Pattern: pattern,
		Syntax:  map[string]Action{pattern: act},
	}
}

/*Build creates a machine with many patterns and actions.
The regexes are built one by one and joined through alternation '|'.
*/
func Build(syntax map[string]Action) *Machine {
	atmts := make([]*automaton, len(syntax))
	patterns := make([]string, len(syntax))
	i := 0
	//this can be paralelized
	for re, act := range syntax {
		patterns[i] = re
		atmts[i] = compile(re, act)
		i++
	}
	final := &automaton{}
	for _, atmt := range atmts { // joining through alternation
		atmt.acc.addEmptyTr(final.acc)
		final.start.addEmptyTr(atmt.start)
	}
	final.start.Enum(&map[*state]int{})
	start := powerSet(&map[string]*state{}, final.start)
	return &Machine{
		Start:   start,
		Pattern: strings.Join(patterns, "|"),
		Syntax:  syntax,
	}
}

func compile(pattern string, act Action) *automaton {
	if act == nil {
		log.Fatal("Action cannot be nil")
	}
	tokens := lexString(pattern)
	p := &parser{}
	root := p.run(tokens)
	//fmt.Println(root)
	atmt := createAtmt(root)
	atmt.acc.act = act
	mp := map[*state]int{}
	atmt.start.Enum(&mp)
	//fmt.Println("thomps:", prettyPrint(&mp))

	atmt.start = powerSet(&map[string]*state{}, atmt.start)
	mp = map[*state]int{}
	atmt.start.Enum(&mp)
	//fmt.Println("powerset:", prettyPrint(&mp))

	return atmt
}
