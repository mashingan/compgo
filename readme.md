# Compgo

Learning building interpreter and compiler from book [Writing an Interpreter In Go](interpreter-book)
and [Writing A Compiler In Go](compiler-book).

While this is following the method and code from the books. There are several differences
for the exercises such as:

1. The code organization is flat for each book. Only has submodules for both interpreter (`interp`) and compiler (`comp`).
2. Overall naming variable is different from book and very likely is not too consistent.
3. This implementation immediately using modules instead of `.env` in interpreter book for adjusting `GOPATH`.
4. Single line re-definition of functions are avoided and immediately using that's single line code only. Re-definition is definitely good for general reading and comprehension but since this is exercise so we're using it as it is.
4. Lexer support utf-8 both in string and in identifier.
5. Lexer using map token instead of map keywords to check whether scanning is stopping for next token or continue.
6. Some parser utilized goroutines to produces the AST, e.g. Infix expression.
7. While the compiler part doesn't add the macros, the interpreter module still has the macro part from interpreter book [lost chapter](https://interpreterbook.com/lost/).
8. The builtin function has additional support for string instead of array only.
9. The compiler doesn't have scope and each time entering the function scope it will flatly compile and readjust itself during walking the instructions.
10. Vm stack has different implementation by appending and deleting last element instead of allocating fixed stack size in book.
11. The builtin functions from [the compiler book](compiler-book) copy-pasting from the builtin module but this implementation literally re-use the builtin functions from the interpreter module, only mapping the identifier. Only need to export the builtin map functions from the `interp` module.
12. New objects defined for the compiler module is defined only there without changing the definition of objects in interpreter.
13. Most of infix operations (parsing/compiling/eval/vm running) is using map to unify operator definition and applying instead of switching which operator it is.

## Impression

Overally both books are enjoyable. I admit not actually read the sentence and immediately looked for the codes but I also read intermittently for something that I need a proper comprehension about the code itself.  
Mostly the code is self-explained, apart from the Go itself is easily readable.

Both books actually the journeys, the [interpreter book](interpreter-book) is very light hearted and very endearing to read while the [compiler book](compiler-book) shifted to more serious but still has the light-hearted nuances.

Since I did something different with the book (for the exercises), I often found myself fixing on non-existing bugs in books' code that I wrote myself.  
For example out-of-bound error because my stack doesn't allocate fixed size stack but appending/slicing the latest element, so if I missed some logic for popping and/or pushing the stack, I often finding myself either got empty stack or accessing index more than stack itself.
Because of this the implementation of stack also returning error to ensure the code itself won't crash but properly checking whether I accessed the stack properly or not.

[interpreter-book]: https://interpreterbook.com
[compiler-book]: https://compilerbook.com