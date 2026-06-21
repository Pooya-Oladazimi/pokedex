package main

import (
	"bufio"
	"fmt"
	"github.com/Pooya-Oladazimi/pokedex/poke"
	"github.com/Pooya-Oladazimi/pokedex/repl"
	"os"
)

func main() {
	buffer := bufio.NewScanner(os.Stdin)
	commands := make(map[string]repl.CliCommand)
	commands["exit"] = repl.CliCommand{
		Name:        "exist",
		Description: "Exit the Pokedex",
		Callback:    repl.CommndExit,
	}
	commands["help"] = repl.CliCommand{
		Name:        "help",
		Description: "Pokedex manual",
		Callback:    repl.CommandHelp,
	}
	commands["map"] = repl.CliCommand{
		Name:        "map",
		Description: "Fetch the next page of poke's locations",
		Callback:    repl.Map,
	}
	commands["mapb"] = repl.CliCommand{
		Name:        "mapb",
		Description: "Fetch the previous page of poke's locations",
		Callback:    repl.Mapb,
	}
	config := repl.Config{
		Next:     poke.PokeLocationUrlFirstPage,
		Previous: "",
	}
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
		cmd, ok := commands[tokens[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := cmd.Callback(&config)
		if err != nil {
			fmt.Printf("Command exited with error: %v\n", err)
			continue
		}
	}
}
