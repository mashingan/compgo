package main

import (
	"compgo/comp"
	"compgo/interp"
	"flag"
	"log"
	"time"
)

var engine = flag.String("engine", "vm", "use vm or val")

var input = `
let fibonacci = fn(x) {
	if (x == 0) { return 0 }
	else {
		if (x == 1) { return 1; }
		else { fibonacci(x - 1) + fibonacci(x - 2) }
	}
};
fibonacci(35);
`

func main() {
	flag.Parse()
	var (
		dur time.Duration
		res interp.Object
	)
	p := interp.NewParser(interp.NewLexer(input))
	prg := p.ParseProgram()
	if *engine == "vm" {
		compiler := comp.New()
		if err := compiler.Compile(prg); err != nil {
			log.Printf("compiler error: %s", err)
			return
		}
		mcn := comp.NewVm(compiler.Bytecode())
		start := time.Now()
		if err := mcn.Run(); err != nil {
			log.Printf("vm run error: %s", err)
			return
		}
		dur = time.Since(start)
		res = mcn.LastPop()
	} else {
		env := interp.NewEnvironment()
		start := time.Now()
		res = interp.Eval(prg, env)
		dur = time.Since(start)
	}

	log.Printf("engine=%s, result=%s, duration=%s\n",
		*engine, res.Inspect(), dur)
}
