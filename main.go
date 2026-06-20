package main

import (
	"bufio"
	"fmt"
	"github.com/Pooya-Oladazimi/pokedex/repl"
	"os"
)

func main() {
	buffer := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		ok := buffer.Scan()
		if !ok {
			break
		}
		userInput := buffer.Text()
		tokens := repl.CleanInput(userInput)
		if len(tokens) == 0 {
			continue
		}
		fmt.Printf("Your command was: %s\n", tokens[0])
	}
}
