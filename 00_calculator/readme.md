# calc

Just a simple textual calculator, i implemented it as an exercise to learn how to generate AST from EBNF grammars using recursive descent parsing. I will briefly explain what i did so that if you have the same questions i had this can serve you well. The lexer is complete overkill for this particular task, but it is [Rob Pike's design](https://youtu.be/HxaD_trXwRE).

If you spot an error or want to ask something, feel free to raise an issue.

First, your grammar should be in LL(1) form, meaning each Non-Terminal has to have a single obvious derivation given the symbol read, it needs to be left-factored, right-recursive and unambiguous. (See Engineering a Compiler: Section 3.3)

The grammar is:

```ebnf
Expr ::= Term {("+" | "-") Term}

Term ::= Power {( "*" | "/" | "%" ) Power}

Power ::= Unary { "^" Unary}

Unary ::= Factor ["!"]

Factor ::= "(" Expr ")"
	| SigNum

SigNum ::= [("+" | "-")] Num

 - - - - - - - This is taken care by the lexer

Num ::= {digits} ["." {digits}]

digits ::= '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9'
```

And the precedence table, from lowest to highest, is:

| Operators | Associativity |
| :-------: | :-----------: |
| +, -      | left to right |
| \*, /, %  | left to right |
| ^         | right to left |
| !         | unary         |

Expressions inside parentheses have the highest precedence.

The precedence is directly encoded in the grammar and each row has it's own procedure. Associativity is expressed inside each procedure by using iteration (left-to-right) or recursion (right-to-left).

That means the operators `+, -`, and all that are left-to-right associative will have a loop in the form: 

```
function Procedure() ASTNode {
	symbols <- ["+", "-"]
	last <- NextProcedure()
	while word exists in symbols { // iterative
		parent <- NewASTNode(word)
		parent.NewLeaf(last)
		NextWord()
		parent.NewLeaf(NextProcedure())
		last <- parent
	}
	return last
}
```

In contrast to "^" and all right-to-left associative that will have the form:

```
function Procedure() ASTNode {
	symbols <- ["^"]
	last <- NextProcedure()
	if word exists in symbols {
		parent <- NewASTNode(word)
		parent.NewLeaf(last)
		NextWord()
		parent.NewLeaf(Procedure()) // recursive
		return parent
	}
	return last
}
```

This only solves associativity with basic expressions, things get complicated if you have two types of associativity in the same precedence level. Context sensitive and non-LL(1) languages will need different modifications to the algorithm to work properly (See Engineering a Compiler: Section 3.5.3). [It](https://gcc.gnu.org/gcc-4.1/changes.html) [seems](https://gcc.gnu.org/legacy-ml/gcc/2005-03/msg00742.html) that recursive-descent parsers generate a lot of confusion about which grammars they're able to parse, but as usual with hand-written things, it can be extensively modified to fit your need.

Complete pseudocode for recursive-descent parsers of the classic expression grammar (with much additional explanation) can be found in:
 - Engineering a Compiler: Figure 3.10 (Only Syntax validation, does not generate AST)
 - https://www.engr.mun.ca/~theo/Misc/exp_parsing.htm#classic (Generates AST and shows how to handle associativity)
