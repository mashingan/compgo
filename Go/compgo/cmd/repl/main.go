package main

import (
	"bufio"
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
		evl := interp.Eval(prg)
		if evl != nil {
			fmt.Println(evl.Inspect())
		}
	}
}

func printParserErrors(o io.Writer, errs []string) {
	for _, e := range errs {
		io.WriteString(o, "\t"+e+"\n")
	}
}
