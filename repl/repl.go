package repl

import (
	"fmt"
	"github.com/Pooya-Oladazimi/pokedex/poke"
	"github.com/Pooya-Oladazimi/pokedex/pokecache"
	"os"
	"strings"
)

func CleanInput(text string) []string {
	text = strings.ToLower(text)
	return strings.Fields(text)
}

type CliCommand struct {
	Name        string
	Description string
	Callback    func(c *Config) error
}

type Config struct {
	Cache    *pokecache.Cache
	Next     string
	Previous string
}

func Map(c *Config) error {
	apiResponse, cached, err := poke.FetchPokeLocation(c.Next, c.Cache)
	if err != nil {
		return err
	}
	if cached {
		fmt.Println("---Cached response---")
	}
	for _, loc := range apiResponse.Results {
		fmt.Println(loc.Name)
	}
	c.Previous = apiResponse.Previous
	c.Next = apiResponse.Next
	return nil
}

func Mapb(c *Config) error {
	if c.Previous == "" {
		return fmt.Errorf("There is no previous page.")
	}
	apiResponse, cached, err := poke.FetchPokeLocation(c.Previous, c.Cache)
	if err != nil {
		return err
	}
	if cached {
		fmt.Println("---Cached response---")
	}
	for _, loc := range apiResponse.Results {
		fmt.Println(loc.Name)
	}
	c.Next = apiResponse.Next
	c.Previous = apiResponse.Previous
	return nil
}

func CommndExit(c *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func CommandHelp(c *Config) error {
	fmt.Printf(`Welcome to the Pokedex!
Usage:

map: Get the next 20 Poke's locations
mapb: Get the previous 20 Poke's locations
help: Displays a help message
exit: Exit the Pokedex
`)
	return nil
}
