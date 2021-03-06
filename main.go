package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/tomocy/monkey/repl"
)

func main() {
	sayHelloToUser()
	repl.Start(os.Stdin, os.Stdout)
}

func sayHelloToUser() {
	user, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	fmt.Println("Feel free to type in commands")
}
