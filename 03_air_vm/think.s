; if a {...}
; if a = 0 {...}
cmp a, 0
je block

; if a < 1 {...}
cmp a, 1
jl block

; if a >= 1 {...}
cmp a, 1
jge block

; if a != 1 {...}
cmp a, 1
jne block

; if a = 0 or a = 1 {...}
cmp a, 0
je block
cmp a, 1
je block

; if a = 0 and a = 1 {...}
cmp a, 0
jne else
cmp a, 1
jne else
jmp block
else:
	...
	jmp overBlock
block:
	...
overBlock:
	...

; if a > 0 and a < 0 and a = 0 {...}
cmp a, 0
jle else
cmp a, 0
jge else
cmp a, 0
jne else
jmp block
else:
	...
	jmp overBlock
block:
	...
overBlock:
	...

; if a = 1 and a > 0 or a < 2 {...}
cmp a, 2
jl block
cmp a, 1
jne else
cmp a, 0
jge else
jmp block
else:
	...
	jmp overBlock
block:
	...
overBlock:
	...
