package re

var tests = []struct {
	re      string
	tokens  []token
	root    *node
	input   string
	matches []string
}{
	{"a",
		[]token{
			{val: 'a', tp: char},
			{val: eof, tp: end},
		},
		&node{
			set: &Set{Items: []rune{'a'}},
			tp:  set,
		},
		"abAbaa",
		[]string{"a", "a", "a"},
	},
	{"(a)",
		[]token{
			{val: '(', tp: ope},
			{val: 'a', tp: char},
			{val: ')', tp: ope},
			{val: eof, tp: end},
		},
		&node{
			set: &Set{Items: []rune{'a'}},
			tp:  set,
		},
		"a",
		[]string{"a"},
	},
	{"abc",
		[]token{
			{val: 'a', tp: char},
			{val: 'b', tp: char},
			{val: 'c', tp: char},
			{val: eof, tp: end},
		},
		&node{
			tp: and,
			children: []*node{
				{set: &Set{Items: []rune("a")}, tp: set},
				{set: &Set{Items: []rune("b")}, tp: set},
				{set: &Set{Items: []rune("c")}, tp: set},
			},
		},
		"aaaaaabcaaabac",
		[]string{"abc"},
	},
	{`\s\t\n`,
		[]token{
			{val: ' ', tp: char},
			{val: '\t', tp: char},
			{val: '\n', tp: char},
			{val: eof, tp: end},
		},

		&node{
			tp: and,
			children: []*node{
				{set: &Set{Items: []rune(" ")}, tp: set},
				{set: &Set{Items: []rune("\t")}, tp: set},
				{set: &Set{Items: []rune("\n")}, tp: set},
			},
		},
		"aaaa\ta\n aaa \t\n",
		[]string{" \t\n"},
	},
	{"a|b|c",
		[]token{
			{val: 'a', tp: char},
			{val: '|', tp: ope},
			{val: 'b', tp: char},
			{val: '|', tp: ope},
			{val: 'c', tp: char},
			{val: eof, tp: end},
		},
		&node{
			tp: or,
			children: []*node{
				{set: &Set{Items: []rune("a")}, tp: set},
				{set: &Set{Items: []rune("b")}, tp: set},
				{set: &Set{Items: []rune("c")}, tp: set},
			},
		},
		"azzzbczzzbza",
		[]string{"a", "b", "c", "b", "a"},
	},
	{"ab|ac",
		[]token{
			{val: 'a', tp: char},
			{val: 'b', tp: char},
			{val: '|', tp: ope},
			{val: 'a', tp: char},
			{val: 'c', tp: char},
			{val: eof, tp: end},
		},
		&node{
			tp: or,
			children: []*node{
				{
					tp: and,
					children: []*node{
						{set: &Set{Items: []rune("a")}, tp: set},
						{set: &Set{Items: []rune("b")}, tp: set},
					},
				},
				{
					tp: and,
					children: []*node{
						{set: &Set{Items: []rune("a")}, tp: set},
						{set: &Set{Items: []rune("c")}, tp: set},
					},
				},
			},
		},
		"abacacbaaab",
		[]string{"ab", "ac", "ac", "ab"},
	},
	{"a*", []token{
		{val: 'a', tp: char},
		{val: '*', tp: ope},
		{val: eof, tp: end},
	},
		&node{
			tp: star,
			children: []*node{
				{set: &Set{Items: []rune{'a'}}, tp: set},
			},
		},
		"aaaabc",
		[]string{"aaaa", "", ""},
	},
	{"ab*",
		[]token{
			{val: 'a', tp: char},
			{val: 'b', tp: char},
			{val: '*', tp: ope},
			{val: eof, tp: end},
		},
		&node{
			tp: and,
			children: []*node{
				{set: &Set{Items: []rune{'a'}}, tp: set},
				{
					tp: star,
					children: []*node{
						{set: &Set{Items: []rune{'b'}}, tp: set},
					},
				},
			},
		},
		"aabbabaaa",
		[]string{"aab", "ab", "aaa"},
	},
	{"a*b*",
		[]token{
			{val: 'a', tp: char},
			{val: '*', tp: ope},
			{val: 'b', tp: char},
			{val: '*', tp: ope},
			{val: eof, tp: end},
		},
		&node{
			tp: and,
			children: []*node{
				{
					tp: star,
					children: []*node{
						{set: &Set{Items: []rune{'a'}}, tp: set},
					},
				},
				{
					tp: star,
					children: []*node{
						{set: &Set{Items: []rune{'b'}}, tp: set},
					},
				},
			},
		},
		"aabbabaaa",
		[]string{"aabb", "ab", "aaa"},
	},
	{"a|b*",
		[]token{
			{val: 'a', tp: char},
			{val: '|', tp: ope},
			{val: 'b', tp: char},
			{val: '*', tp: ope},
			{val: eof, tp: end},
		},
		&node{
			tp: or,
			children: []*node{
				{set: &Set{Items: []rune{'a'}}, tp: set},
				{
					tp: star,
					children: []*node{
						{set: &Set{Items: []rune{'b'}}, tp: set},
					},
				},
			},
		},
		"aabbabaaa",
		[]string{"a", "a", "bb", "a", "b", "a", "a", "a"},
	},
	{"(a|b)*|c", []token{
		{val: '(', tp: ope},
		{val: 'a', tp: char},
		{val: '|', tp: ope},
		{val: 'b', tp: char},
		{val: ')', tp: ope},
		{val: '*', tp: ope},
		{val: '|', tp: ope},
		{val: 'c', tp: char},
		{val: eof, tp: end},
	},
		&node{
			tp: or,
			children: []*node{
				{
					tp: star,
					children: []*node{
						{
							tp: or,
							children: []*node{
								{set: &Set{Items: []rune{'a'}}, tp: set},
								{set: &Set{Items: []rune{'b'}}, tp: set},
							},
						},
					},
				},
				{set: &Set{Items: []rune{'c'}}, tp: set},
			},
		},
		"acaacbbcbcc",
		[]string{"a", "c", "aa", "c", "bb", "c", "b", "c", "c"},
	},
	{"[a-z]",
		[]token{
			{val: '[', tp: ope},
			{val: 'a', tp: char},
			{val: '-', tp: ope},
			{val: 'z', tp: char},
			{val: ']', tp: ope},
			{val: eof, tp: end},
		},
		&node{set: &Set{Items: []rune("abcdefghijklmnopqrstuvwxyz")}, tp: set},
		"abcdZABC",
		[]string{"a", "b", "c", "d"},
	},
	{`[\t\n\s]`,
		[]token{
			{val: '[', tp: ope},
			{val: '\t', tp: char},
			{val: '\n', tp: char},
			{val: ' ', tp: char},
			{val: ']', tp: ope},
			{val: eof, tp: end},
		},
		&node{set: &Set{Items: []rune("\n\t ")}, tp: set},
		"aaaa\taaaa\n a a \n",
		[]string{"\t", "\n", " ", " ", " ", "\n"},
	},
	{"[a-z]|[A-Z]",
		[]token{
			{val: '[', tp: ope},
			{val: 'a', tp: char},
			{val: '-', tp: ope},
			{val: 'z', tp: char},
			{val: ']', tp: ope},
			{val: '|', tp: ope},
			{val: '[', tp: ope},
			{val: 'A', tp: char},
			{val: '-', tp: ope},
			{val: 'Z', tp: char},
			{val: ']', tp: ope},
			{val: eof, tp: end},
		},
		&node{
			tp: or,
			children: []*node{
				{set: &Set{Items: []rune("abcdefghijklmnopqrstuvwxyz")}, tp: set},
				{set: &Set{Items: []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")}, tp: set},
			},
		},
		"aAbBcC",
		[]string{"a", "A", "b", "B", "c", "C"},
	},
	{"[A-Z][a-z]*",
		[]token{
			{val: '[', tp: ope},
			{val: 'A', tp: char},
			{val: '-', tp: ope},
			{val: 'Z', tp: char},
			{val: ']', tp: ope},
			{val: '[', tp: ope},
			{val: 'a', tp: char},
			{val: '-', tp: ope},
			{val: 'z', tp: char},
			{val: ']', tp: ope},
			{val: '*', tp: ope},
			{val: eof, tp: end},
		},
		&node{
			tp: and,
			children: []*node{
				{set: &Set{Items: []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")}, tp: set},
				{
					tp: star,
					children: []*node{
						{set: &Set{Items: []rune("abcdefghijklmnopqrstuvwxyz")}, tp: set},
					},
				},
			},
		},
		"aAbBcC",
		[]string{"a", "A", "b", "B", "c", "C"},
	},
	{`\s[a-z ]*`,
		[]token{
			{val: ' ', tp: char},
			{val: '[', tp: ope},
			{val: 'a', tp: char},
			{val: '-', tp: ope},
			{val: 'z', tp: char},
			{val: ' ', tp: char},
			{val: ']', tp: ope},
			{val: '*', tp: ope},
			{val: eof, tp: end},
		},
		&node{
			tp: and,
			children: []*node{
				{set: &Set{Items: []rune(" ")}, tp: set},
				{
					tp: star,
					children: []*node{
						{set: &Set{Items: []rune(" abcdefghijklmnopqrstuvwxyz")}, tp: set},
					},
				},
			},
		},
		"  abc , abc a",
		[]string{"  abc ", " abc a"},
	},
	{"(a)|(c)",
		[]token{
			{val: '(', tp: ope},
			{val: 'a', tp: char},
			{val: ')', tp: ope},
			{val: '|', tp: ope},
			{val: '(', tp: ope},
			{val: 'c', tp: char},
			{val: ')', tp: ope},
			{val: eof, tp: end},
		},
		&node{
			tp: or,
			children: []*node{
				{set: &Set{Items: []rune("a")}, tp: set},
				{set: &Set{Items: []rune("c")}, tp: set},
			},
		},
		"aaacc",
		[]string{"a", "a", "a", "c", "c"},
	},
	{"a[^a-z]*",
		[]token{
			{val: 'a', tp: char},
			{val: '[', tp: ope},
			{val: '^', tp: ope},
			{val: 'a', tp: char},
			{val: '-', tp: ope},
			{val: 'z', tp: char},
			{val: ']', tp: ope},
			{val: '*', tp: ope},
			{val: eof, tp: end},
		},
		&node{
			tp: and,
			children: []*node{
				{set: &Set{Items: []rune("a")}, tp: set},
				{
					tp: star,
					children: []*node{
						{set: &Set{Items: []rune("abcdefghijklmnopqrstuvwxyz"), Negated: true}, tp: set},
					},
				},
			},
		},
		"aaaZZZ",
		[]string{"a", "a", "aZZZ"},
	},
	{`a[a-z^]*`,
		[]token{
			{val: 'a', tp: char},
			{val: '[', tp: ope},
			{val: 'a', tp: char},
			{val: '-', tp: ope},
			{val: 'z', tp: char},
			{val: '^', tp: char},
			{val: ']', tp: ope},
			{val: '*', tp: ope},
			{val: eof, tp: end},
		},
		&node{
			tp: and,
			children: []*node{
				{set: &Set{Items: []rune("a")}, tp: set},
				{
					tp: star,
					children: []*node{
						{set: &Set{Items: []rune("^abcdefghijklmnopqrstuvwxyz")}, tp: set},
					},
				},
			},
		},
		"aaaZZZ^^",
		[]string{"a", "a", "aZZZ^^"},
	},
	{`a|\e`,
		[]token{
			{val: 'a', tp: char},
			{val: '|', tp: ope},
			{val: 0, tp: empty},
			{val: eof, tp: end},
		},
		&node{
			tp: or,
			children: []*node{
				{set: &Set{Items: []rune("a")}, tp: set},
				{tp: emptyStr},
			},
		},
		"abbaab",
		[]string{"a", "", "", "a", "a", ""},
	},
	{`a|[]`,
		[]token{
			{val: 'a', tp: char},
			{val: '|', tp: ope},
			{val: '[', tp: ope},
			{val: ']', tp: ope},
			{val: eof, tp: end},
		},
		&node{
			tp: or,
			children: []*node{
				{set: &Set{Items: []rune("a")}, tp: set},
				{set: &Set{Items: []rune{}}, tp: set},
			},
		},
		"abbaab",
		[]string{"a", "a", "a"},
	},
	{`[a-z]*\s\*`,
		[]token{
			{val: '[', tp: ope},
			{val: 'a', tp: char},
			{val: '-', tp: ope},
			{val: 'z', tp: char},
			{val: ']', tp: ope},
			{val: '*', tp: ope},
			{val: ' ', tp: char},
			{val: '*', tp: char},
			{val: eof, tp: end},
		},
		&node{
			tp: and,
			children: []*node{
				{
					tp: star,
					children: []*node{
						{set: &Set{Items: []rune("abcdefghijklmnopqrstuvwxyz")}, tp: set},
					},
				},
				{set: &Set{Items: []rune(" ")}, tp: set},
				{set: &Set{Items: []rune("*")}, tp: set},
			},
		},
		" *ab *baab",
		[]string{" *", "ab *"},
	},
	{"([A-Z]|[a-z])*",
		[]token{
			{val: '(', tp: ope},
			{val: '[', tp: ope},
			{val: 'A', tp: char},
			{val: '-', tp: ope},
			{val: 'Z', tp: char},
			{val: ']', tp: ope},
			{val: '|', tp: ope},
			{val: '[', tp: ope},
			{val: 'a', tp: char},
			{val: '-', tp: ope},
			{val: 'z', tp: char},
			{val: ']', tp: ope},
			{val: ')', tp: ope},
			{val: '*', tp: ope},
			{val: eof, tp: end},
		},
		&node{
			tp: star,
			children: []*node{
				{
					tp: or,
					children: []*node{
						{set: &Set{Items: []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")}, tp: set},
						{set: &Set{Items: []rune("abcdefghijklmnopqrstuvwxyz")}, tp: set},
					},
				},
			},
		},
		"AaBbCc_+,",
		[]string{"A", "a", "B", "b", "C", "c"},
	},
	{"\uFFFF",
		[]token{
			{val: '\uFFFF', tp: char},
			{val: eof, tp: end},
		},
		&node{
			set: &Set{Items: []rune("\uFFFF")},
			tp:  set,
		},
		"\uFFFFabc",
		[]string{"\uFFFF"},
	},
	{"\uFFFF[\u0000-\u0004]", []token{
		{val: '\uFFFF', tp: char},
		{val: '[', tp: ope},
		{val: '\u0000', tp: char},
		{val: '-', tp: ope},
		{val: '\u0004', tp: char},
		{val: ']', tp: ope},
		{val: eof, tp: end},
	},
		&node{
			tp: and,
			children: []*node{
				{set: &Set{Items: []rune("\uFFFF")}, tp: set},
				{set: &Set{Items: []rune("\u0000\u0001\u0002\u0003\u0004")}, tp: set},
			},
		},
		"\uFFFF\u0003abc",
		[]string{"\uFFFF\u0003"},
	},
	{"\uFFFF\\s",
		[]token{
			{val: '\uFFFF', tp: char},
			{val: ' ', tp: char},
			{val: eof, tp: end},
		},
		&node{
			tp: and,
			children: []*node{
				{set: &Set{Items: []rune("\uFFFF")}, tp: set},
				{set: &Set{Items: []rune(" ")}, tp: set},
			},
		},
		"\uFFFF abc",
		[]string{"\uFFFF "},
	},
}

var expandTests = []struct {
	set []rune
	ans []rune
}{
	{[]rune("a-z"), []rune("abcdefghijklmnopqrstuvwxyz")},
	{[]rune("A-Z"), []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")},
	{[]rune("\u0000-\u0005"), []rune("\u0000\u0001\u0002\u0003\u0004\u0005")},
	// must output in unicode order
	{[]rune("A-Za-z"), []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")},
	{[]rune("A-Z0-9"), []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")},
	{[]rune("A-Za-z0-9"), []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")},
	{[]rune("a-z_"), []rune("_abcdefghijklmnopqrstuvwxyz")},
	// '-' at the end is a considered a terminal
	{[]rune("a-z-"), []rune("-abcdefghijklmnopqrstuvwxyz")},
	// must remove duplicates
	{[]rune("aaaa"), []rune("a")},
	{[]rune("abc"), []rune("abc")},
	{[]rune("a-bc"), []rune("abc")},
	{[]rune(""), []rune("")},
	{[]rune("a-za-z"), []rune("abcdefghijklmnopqrstuvwxyz")},
	// escapes
	{[]rune(`\s\n\t`), []rune("\t\n ")},
}
