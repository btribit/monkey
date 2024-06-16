package main

import (
	"fmt"
	"log"
	"monkey/repl"
	"os"
	"os/user"
	"plugin"
)

func main() {

	loadExtensions("extensions")

	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	// Compile a filename
	if len(os.Args) > 1 && os.Args[1] == "compile" {
		if len(os.Args) < 3 {
			fmt.Println("Please provide a filename to compile")
			return
		}
		filename := os.Args[2]
		repl.CompileFile(filename)
		return
	}

	// Read a command line argument for compiler or evaluator
	if len(os.Args) > 1 && os.Args[1] == "compiler" {
		fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
		fmt.Printf("You are using the Monkey compiler\n")
		repl.StartCompiler(os.Stdin, os.Stdout)
		return
	}

	fmt.Printf("Hello %s! This is the Monkey programming language!\n",
		user.Username)
	fmt.Printf("You are using the Monkey evaluator\n")
	fmt.Printf("Feel free to type in commands\n")
	repl.StartEvaluator(os.Stdin, os.Stdout)
}

// Output:
// $ go run main.go
// Hello jason! This is the Monkey programming language!
// Feel free to type in commands
// >> let add = fn(x, y) { x + y; };
// {Type:LET Literal:let}
// {Type:IDENT Literal:add}
// ..

func loadExtensions(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() && file.Type().IsRegular() && file.Name()[len(file.Name())-3:] == ".so" {
			extPath := dir + "/" + file.Name()
			p, err := plugin.Open(extPath)
			if err != nil {
				log.Printf("Error loading plugin %s: %v", extPath, err)
				continue
			}

			symbol, err := p.Lookup("Register")
			if err != nil {
				log.Printf("Error looking up Register in %s: %v", extPath, err)
				continue
			}

			registerFunc, ok := symbol.(func())
			if !ok {
				log.Printf("Register symbol in %s is not a function", extPath)
				continue
			}

			registerFunc()
		}
	}

	return nil
}
