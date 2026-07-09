package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/iamASR27/pokedex/internal/pokeapi"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{
		client:  pokeapi.NewClient(5 * time.Minute),
		pokedex: make(map[string]Pokemon),
	}
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		cleanedInput := cleanInput(input)
		if len(cleanedInput) == 0 {
			continue
		}

		commandName := cleanedInput[0]
		cmd, ok := commands[commandName]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		var args []string
		if len(cleanedInput) > 1 {
			args = cleanedInput[1:]
		}

		if err := cmd.callback(cfg, args); err != nil {
			fmt.Println("Error:", err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
	}
}
