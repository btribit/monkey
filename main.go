package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}

// Output:
// $ go run main.go
// Hello jason! This is the Monkey programming language!
// Feel free to type in commands
// >> let add = fn(x, y) { x + y; };
// {Type:LET Literal:let}
// {Type:IDENT Literal:add}
// ..
