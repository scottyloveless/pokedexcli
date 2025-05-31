package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/scottyloveless/pokedexcli/internal"
	"github.com/scottyloveless/pokedexcli/internal/pokecache"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	pokecache.GlobalCache = pokecache.NewCache(time.Second * 5)
	internal.GlobalConfig = internal.Config{
		Next:     nil,
		Previous: nil,
		Cache: pokecache.Cache{
			Map:   map[string]pokecache.CacheEntry{},
			Mutex: &sync.Mutex{},
		},
	}
	internal.GlobalPokedex = internal.Pokedex{
		Inventory: map[string]internal.Pokemon{},
		Mutex:     &sync.Mutex{},
	}
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		words := cleanInput(scanner.Text())
		if value, exists := supportedCommands[words[0]]; exists {
			if len(words) == 1 {
				err := value.callback(internal.GlobalConfig, "")
				if err != nil {
					fmt.Println(err)
				}

			}
			if len(words) == 2 {
				err := value.callback(internal.GlobalConfig, words[1])
				if err != nil {
					fmt.Println(err)
				}
			}

		} else {
			fmt.Println("Unknown command")
		}

	}
}

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	parts := strings.Fields(lowered)
	return parts
}

type cliCommand struct {
	name        string
	description string
	callback    func(internal.Config, string) error
}

var supportedCommands map[string]cliCommand

func init() {
	supportedCommands = map[string]cliCommand{
		"map": {
			name:        "map",
			description: "Displays a list of next 20 maps",
			callback:    mapForward,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays a list of previous 20 maps",
			callback:    mapBack,
		},

		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"explore": {
			name:        "explore",
			description: "Enter location name to display list of Pokemon",
			callback:    explore,
		},
		"catch": {
			name:        "catch",
			description: "Attemp to catch a pokemon",
			callback:    catch,
		},
		"inspect": {
			name:        "inspect",
			description: "View details of caught Pokemon",
			callback:    inspect,
		},
	}
}

func commandExit(config internal.Config, str string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config internal.Config, str string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println(" ")
	for _, cmd := range supportedCommands {
		fmt.Printf("%v: %v\n", cmd.name, cmd.description)
	}
	return nil
}

func mapForward(config internal.Config, str string) error {
	pmap, err := internal.ApiLocationFetch()
	if err != nil {
		return err
	}

	for _, m := range pmap.Results {
		fmt.Printf("%v\n", m.Name)
	}
	return nil
}

func mapBack(config internal.Config, str string) error {
	if config.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	pmap, err := internal.ApiLocationFetchBack()
	if err != nil {
		return err
	}

	for _, m := range pmap.Results {
		fmt.Printf("%v\n", m.Name)
	}
	return nil
}

func explore(config internal.Config, str string) error {
	pexplore, err := internal.ApiLocationFetchExplore(str)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %v...\n", str)

	for _, e := range pexplore.PokemonEncounters {
		fmt.Printf("%v\n", e.Pokemon.Name)
	}
	return nil
}

func catch(config internal.Config, str string) error {
	poke_info, err := internal.ApiPokemonData(str)
	if err != nil {
		return err
	}

	baseExp := float64(poke_info.BaseExperience)
	maxExp := 300.0

	baseChance := 20.0
	successChance := (maxExp / baseExp) * baseChance

	fmt.Printf("Throwing a Pokeball at %v...\n", str)
	success := rollChance(successChance)
	if success {
		fmt.Printf("%v was caught!\n", str)
		internal.GlobalPokedex.Mutex.Lock()
		internal.GlobalPokedex.Inventory[str] = poke_info
		internal.GlobalPokedex.Mutex.Unlock()
		// fmt.Printf("Pokedex contains %v", internal.GlobalPokedex.Inventory[str].Name)
	} else {
		fmt.Printf("%v escaped!\n", str)
	}

	return nil
}

func rollChance(percentage float64) bool {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	roll := rand.Float64() * 100

	return roll <= percentage
}

func inspect(config internal.Config, str string) error {
	v, exists := internal.GlobalPokedex.Inventory[str]
	if !exists {
		fmt.Println("you have not caught that pokemon")
	} else {
		fmt.Printf("Name: %v\n", v.Name)
		fmt.Printf("Height: %v\n", v.Height)
		fmt.Printf("Weight: %v\n", v.Weight)
		fmt.Println("Stats:")
		for _, s := range v.Stats {
			fmt.Printf(" -%v: %v\n", s.Stat.Name, s.BaseStat)
		}
		fmt.Println("Types:")
		for _, t := range v.Types {
			fmt.Printf(" - %v\n", t.Type.Name)
		}
	}
	return nil
}
