package main

import (
	"fmt"
	"time"
	"flag"
	"os/exec"
	"os"
	"io/ioutil"
	"strings"
)

var wait = flag.Duration("t", 200*time.Millisecond, "Time between evals")

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Takes one argument (the file)")
		os.Exit(1)
	}
	filedata, err := ioutil.ReadFile(args[0])
	if err != nil {
		panic(err)
	}
	strdata := strings.Split(string(filedata), "\n")
	grid := StrToGrid(strdata)
	m := &Machine{
		G: grid,
		Stack: make([]byte, 10),
		StackPtr: -1,
	}

	Run(m)
}

func StrToGrid(strs []string) Grid {
	out := Grid{}
	for i := range strs {
		out = append(out, []byte(strs[i]))
	}
	return out
}

type Grid [][]byte

type Machine struct {
	G        Grid
	DPointer Vector
	IPointer Vector
	Stack    []byte
	StackPtr int
}

func Normalize(grid *Grid) {
	max := 0
	for _, row := range *grid {
		if max < len(row) {
			max = len(row)
		}
	}
	for i, row := range *grid {
		if len(row) < max {
			diff := max - len(row)
			(*grid)[i] = append((*grid)[i], noopPadding(diff)...)
		}
	}
}

func noopPadding(length int) []byte {
	out := make([]byte, length)
	for i := range out {
		out[i] = ' '
	}
	return out
}

func Display(m *Machine) {
	ClearGrid(m.G)
	fmt.Print("\033[A")
	fmt.Print("\033[A")
	fmt.Printf("DP: %s, IP: %s                            \n",
		string(Read(m.DPointer, &m.G)), string(Read(m.IPointer, &m.G)))
	PrintStack(m.Stack, m.StackPtr)
	for y, row := range m.G {
		for x, b := range row {
			ColorPointers(m.DPointer, m.IPointer, x, y)
			fmt.Printf("%s", string(b))
			fmt.Print("\033[0m")
		}
		fmt.Println()
	}
}

func ColorPointers(D, I Vector, x, y int) {
	if x == D.X && y == D.Y && x == I.X && y == I.Y {
		fmt.Print("\033[0;41m")
	} else if x == D.X && y == D.Y {
		fmt.Print("\033[0;44m")
	} else if x == I.X && y == I.Y {
		fmt.Print("\033[0;42m")
	}
}

func PrintStack(s []byte, ptr int) {
	for i, b := range s {
		if i == ptr {
			fmt.Print("\033[0;42m")
			fmt.Printf("%X", b)
			fmt.Print("\033[0m")
			fmt.Print(" ")
		} else if i > ptr {
			fmt.Print("  ")
		} else {
			fmt.Printf("%X ", b)
		}
	}
	fmt.Println()
}

func ClearGrid(g Grid) {
	for i := 0; i < len(g); i++ {
		fmt.Print("\033[A")
	}
}

type Vector struct {
	X, Y int
	Ori  Dir
}

type Dir byte

const (
	Right Dir = iota
	Down
	Left
	Up 
)

func Run(m *Machine) {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	m.IPointer = Vector{X: 0, Y: 0, Ori: Right}
	m.DPointer = Vector{X: 0, Y: 0, Ori: Right}
	Normalize(&m.G)
	Display(m)
	for {
		instr := Read(m.IPointer, &m.G)
		if Eval(m, instr) {
			fmt.Println("Halt Condition: Collided with the edge!")
			break
		}
		time.Sleep(*wait)
		Display(m)
	}
}

func Read(v Vector, g *Grid) byte {
	if v.Y >= len(*g) {
		*g = append(*g, make([]byte, 0))
		Normalize(g)
	}
	if v.X >= len((*g)[0]) {
		(*g)[0] = append((*g)[0], ' ')
		Normalize(g)
	}
	return (*g)[v.Y][v.X]
}

func Write(v Vector, g *Grid, data byte) {
	if v.Y >= len(*g) {
		*g = append(*g, make([]byte, 0))
		Normalize(g)
	}
	if v.X >= len((*g)[0]) {
		(*g)[0] = append((*g)[0], ' ')
		Normalize(g)
	}
	(*g)[v.Y][v.X] = data
}

var CharToHex = map[byte]byte{
	'0': 0x0, '1': 0x1, '2': 0x2, '3': 0x3,
	'4': 0x4, '5': 0x5, '6': 0x6, '7': 0x7,
	'8': 0x8, '9': 0x9, 'A': 0xA, 'B': 0xB,
	'C': 0xC, 'D': 0xD, 'E': 0xE, 'F': 0xF,
}

var ByteToDir = map[byte]Dir{
	0: Right, 1: Down, 2: Left, 3: Up,
}

func Eval(m *Machine, instr byte) bool {
	switch instr {
	case '-':
		a := Pop(m)
		b := Pop(m)
		Push(m, b-a)
	case '+':
		a := Pop(m)
		b := Pop(m)
		Push(m, b+a)
	case '*':
		a := Pop(m)
		b := Pop(m)
		Push(m, b*a)
	case '/':
		a := Pop(m)
		b := Pop(m)
		Push(m, b/a)

	case 'w':
		val := Pop(m)
		Write(m.DPointer, &m.G, val)
	case 'r':
		Push(m, Read(m.DPointer, &m.G))
	case ':':
		dir := Pop(m)
		if dir >= 0 && dir <= 3 {
			m.DPointer.Ori = ByteToDir[dir]
		}
	case 'i':
		Move(&m.DPointer)

	case '>':
		m.IPointer.Ori = Right
	case '!':
		m.IPointer.Ori = Down
	case '<':
		m.IPointer.Ori = Left
	case '^':
		m.IPointer.Ori = Up

	case 'l':
		c := Pop(m)
		t := Top(m)
		if t < c {
			Rotate(m)
		}
	case 'm':
		c := Pop(m)
		t := Top(m)
		if t > c {
			Rotate(m)
		}
	case '=':
		c := Pop(m)
		t := Top(m)
		if t == c {
			Rotate(m)
		}
	case '~':
		c := Pop(m)
		t := Top(m)
		if t != c {
			Rotate(m)
		}
	default:
		if (instr >= '0' && instr <= '9') || (instr >= 'A' && instr <= 'F') {
			Push(m, CharToHex[instr])
		}
	}
	Move(&m.IPointer)
	if m.DPointer.X < 0 || m.DPointer.Y < 0 ||
		m.IPointer.X < 0 || m.IPointer.Y < 0 {
		return true
	}
	return false
}

func Move(v *Vector) {
	switch v.Ori {
	case Up:
		v.Y--
	case Down:
		v.Y++
	case Left:
		v.X--
	case Right:
		v.X++
	}
}

func Push(m *Machine, data byte) {
	m.StackPtr++
	if m.StackPtr >= len(m.Stack) {
		m.Stack = append(m.Stack, make([]byte, len(m.Stack))...)
	}
	m.Stack[m.StackPtr] = data
}

func Pop(m *Machine) byte {
	if m.StackPtr < 0 {
		m.StackPtr = 0
	}
	val := m.Stack[m.StackPtr]
	m.Stack[m.StackPtr] = 0
	m.StackPtr--
	return val
}

func Top(m *Machine) byte {
	if m.StackPtr < 0 {
		m.StackPtr = 0
	}
	return m.Stack[m.StackPtr]
}

func Rotate(m *Machine) {
	m.IPointer.Ori = Dir((m.IPointer.Ori  +  1)%4)
}





