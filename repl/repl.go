package repl

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strings"

	"github.com/Pooya-Oladazimi/pokedex/poke"
	"github.com/Pooya-Oladazimi/pokedex/pokecache"
)

func CleanInput(text string) []string {
	text = strings.ToLower(text)
	return strings.Fields(text)
}

type CliCommand struct {
	Name        string
	Description string
	Callback    func(c *Config, p []Parameter) error
}

type Parameter struct {
	Name  string
	Value string
}

type Config struct {
	Cache    *pokecache.Cache
	Next     string
	Previous string
}

func Map(c *Config, params []Parameter) error {
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

func Mapb(c *Config, params []Parameter) error {
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

func Explore(c *Config, params []Parameter) error {
	if len(params) == 0 {
		return fmt.Errorf("explore command needs a location parameter.")
	}
	location := params[0].Value
	pokemons, cached, err := poke.FetchPokemonsInLocation(location, c.Cache)
	if err != nil {
		return err
	}
	if cached {
		fmt.Println("---Cached response---")
	}
	fmt.Println("Exploring pastoria-city-area...")
	fmt.Println("Found Pokemon:")
	for _, pokE := range pokemons {
		fmt.Printf("  - %s\n", pokE.Pokemon.Name)
	}

	return nil
}

func Catch(c *Config, params []Parameter) error {
	if len(params) == 0 {
		return fmt.Errorf("catch command needs a pokemon's name paramter.")
	}
	pokemonName := params[0].Value
	pokemon, err := poke.FetchPokemon(pokemonName, c.Cache)
	if err != nil {
		return fmt.Errorf("Not able to fetch the pokemon: %s", err)
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	upper := pokemon.BaseExperience / 10
	shot := rand.IntN(upper + 1)
	if shot == upper {
		fmt.Printf("%s was caught!\n", pokemonName)
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}

func CommndExit(c *Config, params []Parameter) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func CommandHelp(c *Config, params []Parameter) error {
	fmt.Printf(`Welcome to the Pokedex!
Usage:

map: Get the next 20 Poke's locations
mapb: Get the previous 20 Poke's locations
explore <location>: list all the Pokemons in a location
catch <pokemon_name>: catch a pokemon
help: Displays a help message
exit: Exit the Pokedex
`)
	return nil
}
