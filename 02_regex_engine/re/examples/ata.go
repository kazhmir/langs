package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/kazhmir/ata"
)

var cpuProfile = flag.String("ppcpu", "", "Write cpu profile to file")
var p = flag.Bool("p", false, "Prints matches to terminal")
var my = flag.Bool("my", false, "My regular expression engine")
var std = flag.Bool("std", false, "Standard regular expression engine")

func main() {
	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	nonFlags := flag.Args()

	if len(nonFlags) < 2 {
		fmt.Println("Wrong")
	}
	content, err := ioutil.ReadFile(nonFlags[1])
	if err != nil {
		panic(err)
	}
	sCon := string(content)

	if *p {
		if *std {
			StdPrintAndMeasure(&sCon, nonFlags[0])
		}
		if *my {
			MyPrintAndMeasure(&sCon, nonFlags[0])
		}
	} else {
		if *std {
			StdCountAndMeasure(&sCon, nonFlags[0])
		}
		if *my {
			MyCountAndMeasure(&sCon, nonFlags[0])
		}
	}
}

func StdPrintAndMeasure(input *string, pattern string) {
	var stdRe *regexp.Regexp
	measure(func() {
		stdRe = regexp.MustCompile(pattern)
	}, "stdReCompile:")
	measure(func() {
		fmt.Printf("stdRe: %#v\n", stdRe.FindAllString(*input, -1))
	}, "stdReRun:")
}

func MyPrintAndMeasure(input *string, pattern string) {
	var myRe *ata.Machine
	out := make([]string, 0)
	act := func(mat *ata.Match) bool {
		out = append(out, mat.S)
		return false
	}
	txt := strings.NewReader(*input)
	measure(func() {
		myRe = ata.BuildOne(pattern, act)
	}, "myReCompile:")
	measure(func() {
		err := myRe.Run(txt)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("MyRe: %#v\n", out)
	}, "myReRun:")
}

func StdCountAndMeasure(input *string, pattern string) {
	var stdRe *regexp.Regexp
	measure(func() {
		stdRe = regexp.MustCompile(pattern)
	}, "stdReCompile:")
	measure(func() {
		fmt.Println("stdRe:", len(stdRe.FindAllString(*input, -1)))
	}, "stdReRun:")
}

func MyCountAndMeasure(input *string, pattern string) {
	var myRe *ata.Machine
	out := 0
	act := func(mat *ata.Match) bool {
		out++
		return false
	}
	txt := strings.NewReader(*input)
	measure(func() {
		myRe = ata.BuildOne(pattern, act)
	}, "myReCompile:")
	measure(func() {
		err := myRe.Run(txt)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("MyRe: %#v\n", out)
	}, "myReRun:")
}

func measure(a func(), header string) {
	tStart := time.Now()
	a()
	tEnd := time.Now()
	fmt.Println(header, tEnd.Sub(tStart))
}
