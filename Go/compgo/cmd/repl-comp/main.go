package main

import (
	"bufio"
	"compgo/comp"
	"compgo/interp"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
)

const (
	Prompt = ">> "
)

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Welcome %s to REPL!\n\n", user.Username)
	scanner := bufio.NewScanner(os.Stdin)
	// env := interp.NewEnvironment()
	menv := interp.NewEnvironment()
	cmpiler := comp.New()
	var mcn *comp.Vm
	globs := make([]interp.Object, comp.GlobalSize)
	constants := []interp.Object{}
	symbolsTable := comp.NewSymbolTable()
	for i, b := range comp.Builtins {
		symbolsTable.DefineBuiltin(i, b.Name)
	}
	cmpiler.SetSymbolTable(symbolsTable)
	cmpiler.SetConstants(constants)
	for {
		fmt.Print(Prompt)
		scn := scanner.Scan()
		if !scn {
			return
		}
		line := scanner.Text()
		p := interp.NewParser(interp.NewLexer(line))
		prg := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(os.Stdout, p.Errors())
			continue
		}
		interp.DefineMacros(prg, menv)
		mobj := interp.ExpandMacros(prg, menv)
		cmpiler.Instructions = comp.Instructions{}
		err := cmpiler.Compile(mobj)
		if err != nil {
			fmt.Printf("Compilation failed:\n%s\n", err)
			continue
		}
		if mcn == nil {
			mcn = comp.NewVm(cmpiler.Bytecode())
			mcn.SetGlobals(globs)
		} else {
			b := cmpiler.Bytecode()
			mcn.SetFrame(comp.NewFrame(&comp.CompiledFunction{
				Instructions: b.Instructions}, 0))
			mcn.Stack = []interp.Object{}
			mcn.SetConstants(b.Constants)
			mcn.SetGlobals(globs)
		}
		err = mcn.Run()
		if err != nil {
			fmt.Printf("Executing bytecode failed:\n%s\n", err)
			continue
		}

		stacktop := mcn.LastPop()
		fmt.Println(stacktop.Inspect())
	}
}

func printParserErrors(o io.Writer, errs []string) {
	for _, e := range errs {
		io.WriteString(o, "\t"+e+"\n")
	}
}
