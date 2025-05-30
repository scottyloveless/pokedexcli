package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/scottyloveless/pokedexcli/internal"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	go internal.NewCache(time.Second * 90)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		words := cleanInput(scanner.Text())
		if value, exists := supportedCommands[words[0]]; exists {
			err := value.callback(internal.ApiConfig)
			if err != nil {
				fmt.Println(err)
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
	callback    func(internal.Config) error
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
	}
}

func commandExit(config internal.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config internal.Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println(" ")
	for _, cmd := range supportedCommands {
		fmt.Printf("%v: %v\n", cmd.name, cmd.description)
	}
	return nil
}

func mapForward(config internal.Config) error {
	pmap, err := internal.ApiLocationFetch(config.Next)
	if err != nil {
		return err
	}

	for _, m := range pmap.Results {
		fmt.Printf("%v\n", m.Name)
	}
	return nil
}

func mapBack(config internal.Config) error {
	if config.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	pmap, err := internal.ApiLocationFetch(config.Previous)
	if err != nil {
		return err
	}

	for _, m := range pmap.Results {
		fmt.Printf("%v\n", m.Name)
	}
	return nil
}
