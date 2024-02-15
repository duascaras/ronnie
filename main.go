package main

import (
	"fmt"
	"os"
	"os/user"
	"ronnie/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello, %s!\n", user.Name)
	fmt.Printf("This is Ronnie. Feel free to type in commands.\n")
	repl.Start(os.Stdin, os.Stdout)
}
