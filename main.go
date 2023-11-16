package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"internal/pokeapi"
	"internal/pokecache"
)

var conf = config{
	next:     "https://pokeapi.co/api/v2/location-area/",
	previous: "",
}

var cache pokecache.Cache = pokecache.NewCache(300 * time.Second)

var pokedex map[string]pokeapi.Pokemon = make(map[string]pokeapi.Pokemon)

func main() {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		if len(input) == 0 {
			fmt.Println("Please enter a command!\nUse the 'help' command for the list of available commands.")
			continue
		}
		input = strings.TrimSpace(input)
		values := strings.SplitN(input, " ", 2)
		commandName := values[0]
		commandStruct, ok := getCommands()[commandName]
		if !ok {
			fmt.Print("Unknown command!\n")
			continue
		}
		commandCallback := commandStruct.callback
		var args = []string{}
		if len(values) == 2 {
			args = append(args, values[1])
		}
		err := commandCallback(args...)
		if err != nil {
			fmt.Print("An error occured: ", err)
		}
		fmt.Println()
	}
}
