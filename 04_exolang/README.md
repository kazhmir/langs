# SGTM

SGTM may mean Sounds Good To Me, or Stack-Grid Turing Machine, you choose.

The included interpreter only works on linux, and was only tested on Fedora
with Go 1.16.4.

## Description

### Syntax

```
^><! 0123
-+*/ 4567
rw:i 89AB 
lm=~ CDEF 
```

Rows are terminated with a newline character.

### Semantics

The language has a stack and a 2D tape (grid), 
a instruction vector and a data vector,
both the program and data live in the same tape.

The instruction vector has two components,
a value and a orientation. The value is always
incremented in the direction of the orientation.

Both vectors start at position (0,0) with orientation >.

The machine Halts if it collides with the walls (when x or y is negative).

Hexadecimal (0-F) numbers push the respective value into the stack.

Arithmetic operators (*, /, +, -) pop two values from 
the stack and pushes the result (`ab-` means `a-b`).

'^', '<', '>' and '!' sets the orientation of the instruction vector
('^' is up, '!' is down, '<' is left, '>' is right).

'r' reads a value from the tape and pushes it in the stack,
'w' pops a value from the stack and writes it to the tape.

':' pops a value from the stack and sets it as the orientation of the data pointer 
(0 = >, 1 = !, 2 = <, 3 = ^),
'i' increments the data pointer.

'l', 'm', '=' and '\~' all pop one value from the stack 
and compares it to the next value on the top of the stack 
('l' is less, 'm' is more, '=' is equal, '~' is different).
If the condition is true, the direction is rotated
90° clockwise, otherwise, the direction of the instruction
pointer stays the same.

All other chars are noop.

Have fun!

# Examples

Run with `go run . -t <time> <file>`, where
`<time>` is the time between evaluations and
`<file>` is the file name.

Eg: `go run . -t 200ms ./programs/hello.sgtm`

### Hello World

```
>1:i0: 89*w i 7A*1-w i 7A*6+w i 7A*6+w i 8A*1-w i 84*w i !
             ^w+8*A6 i w+6*A7 i w+2*A8 i w-1*A8 i w+7*A8 <
```

### Self-copy

```
> r1:iiw3:i      !
^ i:0iii:3wii:1r <
```

