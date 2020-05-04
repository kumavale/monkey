package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	if len(os.Args) < 2 {
		user, err := user.Current()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
		fmt.Printf("Feel free to type in commands\n")
		fmt.Printf("Use exit() to exit\n")
		repl.Start(os.Stdin, os.Stdout)
	} else {
		b, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			fmt.Println(os.Stderr, err)
			os.Exit(1)
		}

		run(string(b))
	}
}

func run(input string) {
	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			io.WriteString(os.Stdout, "\t"+msg+"\n")
		}
		return
	}

	evaluator.DefineMacros(program, macroEnv)
	expanded := evaluator.ExpandMacros(program, macroEnv)

	evaluated := evaluator.Eval(expanded, env)
	if evaluated != nil {
		io.WriteString(os.Stdout, evaluated.Inspect())
		io.WriteString(os.Stdout, "\n")
	}
}
