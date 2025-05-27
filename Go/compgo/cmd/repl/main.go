package main

import (
	"bufio"
	"compgo/interp"
	"fmt"
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
		l := interp.NewLexer(line)
		for tok := l.NextToken(); tok.Type != interp.Eof; tok = l.NextToken() {
			fmt.Printf("%#v\n", tok)
		}
	}
}
