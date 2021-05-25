package main

/*Tells the linker where the program starts
and sets the base pointer to the stack pointer*/
const Header = `
	section .bss
buff:	resb 	24		; 24 byte buffer

	global _start
	section .text
_start:
	mov 	r12, [rsp]	; argc
convert_loop:
	mov 	r13, r12
	imul 	r13, 8
	add 	r13, rsp 	; r13 == rsp+8*argc
	mov 	r14, [r13]
	push 	r14

	call 	atoi
	mov 	[r13], rax	; substitute address of string for int
	add 	rsp, 8		; cleans stack
	dec 	r12
	cmp 	r12, 0		; stop if there are no more arguments
	jne 	convert_loop	; otherwise continue loop

	mov	rbp, rsp
`

/*Makes an exit syscall with whatever value is in the
rdi register, we can use this to return values smaller than 255*/
const Tail = `
exit:
	call 	itoa		; converts the result to string
	mov 	rdx, rax	; returns size of string in rax
	mov 	rax, 1		; write syscall
	mov 	rdi, 1		; file == stdout
	mov 	rsi, buff	; string is in buffer
	syscall

	mov 	rax, 60
	syscall

; atoi takes one argument:
;	start address of a string
; and returns one result in rax:
;	the integer
atoi:
	push 	rbp
	mov	rbp, rsp
	xor 	rax, rax	; sets rax to 0
	xor 	rbx, rbx	; sets rbx to 0
	mov 	rcx, 1		; rcx is the signal (starts as positive)
	mov 	rdx, [rbp+16] 	; gets address of string

	mov 	bl, [rdx]	; gets char into bl
	cmp 	bl, 0		; if empty
	je 	ret_atoi 	; then return
	cmp 	bl, '-' 	; if minus sign
	je 	minus		; then jump to minus
	jmp 	atoi_loop	; otherwise go to loop
minus:
	mov 	rcx, -1		; sets the signal as negative
	inc 	rdx
	mov 	bl, [rdx]	; gets char into bl again
	
atoi_loop:
	sub 	bl, 48		; converts char to integer
	add 	rax, rbx	; adds to result

	inc 	rdx		; increments pointer
	xor 	rbx, rbx 	; zeroes rbx
	mov 	bl, [rdx]	; gets char into bl
	cmp 	bl, 0	 	; if end of string
	je 	ret_atoi	; then return

	imul 	rax, 10		; otherwise shift decimal digit
	jmp 	atoi_loop	; and loop again
	
ret_atoi:
	imul 	rax, rcx	; applies signal
	mov 	rsp, rbp
	pop 	rbp
	ret
	

; itoa takes one argument:
;	a 64 bit integer
; and returns one result in rax:
;	the size of the string
; buff will contain the string
itoa:
	push 	rbp
	mov	rbp, rsp
	mov 	rax, [rbp+16] 	; gets integer
	mov 	r8, buff	; moves address of buffer to r8

	cmp 	rax, 0 		; compares integer with zero
	jg	plus		; if greater than zero then go to plus
	mov 	cl, '-'		; otherwise set cl to minus sign
	imul 	rax, -1		; method only works with positive integers
	jmp 	itoa_loop	; and jump to loop
plus:
	mov 	cl, '+'		; set cl to plus sign
itoa_loop:
	inc 	r8
	xor 	rdx, rdx
	mov 	rbx, 10
	div	rbx
	
	add 	rdx, 48
	mov 	[r8], dl
	cmp 	rax, 0
	jne 	itoa_loop
	
	inc 	r8
	mov 	[r8], cl	  	; puts sign at the end
	mov 	[r8+1], byte(10)	; moves	end of string to the end

	mov 	rbx, r8		
	add 	rbx, 2		; append newline, and increment so it's inclusive
	sub 	rbx, buff	; computes total size of string
	mov 	rax, rbx	; put size into return register (rax)

	mov 	rbx, buff	; rbx contains the start of string
				; r8 contains the last char

reverse:			; reverses string in place
	mov 	al, [r8]	; gets char at one end
	mov 	cl, [rbx]	; gets char of other end
	mov 	[r8], cl	; swap one char
	mov 	[rbx], al	; with the other
	dec 	r8		; decrements the end
	inc 	rbx		; increments the beginning
	cmp 	r8, rbx		; they should meet at the middle
	jg	reverse		; if r8 is bigger than rbx continue

ret_itoa:
	mov 	rsp, rbp
	pop 	rbp
	ret
`

var OpToASM = map[Operator]string{
	SUB: "sub",
	ADD: "add",
	MUL: "imul",
	DIV: "idiv",
}

/*Note that we skip the special purpose registers,
there are 14 in total, but the register allocator
can work with only 3 by pushing values to the stack*/
var x64Reg = []string{
	"rax", "rbx", "rcx", "rdx", "rsi", "rdi",
	"r8", "r9", "r10", "r11", "r12", "r13", "r14", "r15",
}
