# arith. expr. compiler

It compiles arithmetic expressions into NASM x64 assembly. A recursive descent parser produces AST, the AST is converted into non-destructive 3 address code, then a local allocator works with the 3 address code, allocates the registers and generates NASM x64 assembly.
