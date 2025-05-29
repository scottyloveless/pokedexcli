package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {

		fmt.Print("Pokedex > ")
		scanner.Scan()
		words := cleanInput(scanner.Text())
		if value, exists := supportedCommands[words[0]]; exists {
			err := value.callback()
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
	callback    func() error
}

type config struct {
	Next     string `json:"next"`
	Previous any    `json:"previous"`
}

var supportedCommands map[string]cliCommand

func init() {
	supportedCommands = map[string]cliCommand{
		"map": {
			name:        "map",
			description: "Displays a list of maps",
			callback:    fetchMap,
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

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println(" ")
	for _, cmd := range supportedCommands {
		fmt.Printf("%v: %v\n", cmd.name, cmd.description)
	}
	return nil
}

func fetchMap() error {
	return nil
}
