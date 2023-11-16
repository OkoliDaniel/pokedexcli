package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(args ...string) error
}

type config struct {
	previous string
	next     string
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    helpCommand,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    exitCommand,
		},
		"map": {
			name:        "map",
			description: "Displays loaction areas in batches of 20",
			callback:    mapCommand,
		},
		"mapb": {
			name:        "map back",
			description: "Navigate backwards between location area batches",
			callback:    mapBackCommand,
		},
		"explore": {
			name:        "explore",
			description: "Takes a location area name as input and returns the names of pokemons encountered in that area",
			callback:    exploreCommand,
		},
		"catch": {
			name:        "catch",
			description: "Takes a pokemon name as input and tries to catch that pokemon, adding it to your pokedex upon success",
			callback:    catchCommand,
		},
		"inspect": {
			name:        "inspect",
			description: "Takes a pokemon name as input. If you've caught it, inspect command displays its information",
			callback:    inspectCommand,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays the names of the pokemons in your pokedex",
			callback:    pokedexCommand,
		},
	}
}

func helpCommand(args ...string) error {
	if len(args) != 0 {
		return fmt.Errorf("help command does not accept arguments")
	}
	fmt.Print("\nWelcome to the Pokedex:\nUsage:\n\n")
	commandsMap := getCommands()
	for _, v := range commandsMap {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	return nil
}

func exitCommand(args ...string) error {
	if len(args) != 0 {
		return fmt.Errorf("exit command does not accept arguments")
	}
	os.Exit(0)
	return nil
}

func mapCommand(args ...string) error {
	if len(args) != 0 {
		return fmt.Errorf("map command does not accept arguments")
	}
	if conf.next == "" {
		return fmt.Errorf("end of location areas reached")
	}

	data, ok := cache.Get(conf.next)
	if !ok {
		res, err := http.Get(conf.next)
		if err != nil {
			return fmt.Errorf("something went wrong while fetching location areas: %s", err)
		}
		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			return fmt.Errorf("an error occured: %s", res.Status)
		}
		if err != nil {
			return fmt.Errorf("an error occured while reading response: %s", err)
		}

		data = []byte(body)
		cache.Add(conf.next, data)
	}

	resourceList := pokeapi.ResourceList{}
	err := json.Unmarshal(data, &resourceList)
	if err != nil {
		return err
	}
	if resourceList.Next == nil {
		conf.next = ""
	} else {
		conf.next = *resourceList.Next
	}
	if resourceList.Previous == nil {
		conf.previous = ""
	} else {
		conf.previous = *resourceList.Previous
	}
	locationAreas := resourceList.Results
	for _, locationArea := range locationAreas {
		fmt.Println(locationArea.Name)
	}

	return nil
}

func mapBackCommand(args ...string) error {
	if len(args) != 0 {
		return fmt.Errorf("mapb command does not accept arguments")
	}
	if conf.previous == "" {
		return fmt.Errorf("end of location areas reached")
	}

	data, ok := cache.Get(conf.previous)
	if !ok {
		res, err := http.Get(conf.previous)
		if err != nil {
			return fmt.Errorf("something went wrong while fetching location areas: %s", err)
		}
		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			return fmt.Errorf("an error occured: %s", res.Status)
		}
		if err != nil {
			return fmt.Errorf("an error occured while reading response: %s", err)
		}

		data = []byte(body)
		cache.Add(conf.previous, data)
	}

	resourceList := pokeapi.ResourceList{}
	err := json.Unmarshal(data, &resourceList)
	if err != nil {
		return err
	}
	if resourceList.Next == nil {
		conf.next = ""
	} else {
		conf.next = *resourceList.Next
	}
	if resourceList.Previous == nil {
		conf.previous = ""
	} else {
		conf.previous = *resourceList.Previous
	}
	locationAreas := resourceList.Results
	for _, locationArea := range locationAreas {
		fmt.Println(locationArea.Name)
	}

	return nil
}

func exploreCommand(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("explore command takes exactly one argument: name of location area to explore")
	}
	trimmed := strings.TrimSpace(args[0])
	locationName := strings.ReplaceAll(trimmed, " ", "-")
	fmt.Printf("Exploring %s area...\n", locationName)
	locationAreaURL := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", locationName)
	data, ok := cache.Get(locationAreaURL)
	if !ok {
		res, err := http.Get(locationAreaURL)
		if err != nil {
			return fmt.Errorf("something went wrong while fetching location info: %s", err)
		}
		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			return fmt.Errorf("an error occured: %s", res.Status)
		}
		if err != nil {
			return fmt.Errorf("an error occured while reading response: %s", err)
		}

		data = []byte(body)
		cache.Add(locationAreaURL, data)
	}
	locationArea := pokeapi.LocationArea{}
	err := json.Unmarshal(data, &locationArea)
	if err != nil {
		return err
	}
	pokemonEncounters := locationArea.PokemonEncounters
	fmt.Println("Found the following Pokemon:")
	for _, pokemonEncounter := range pokemonEncounters {
		fmt.Println(pokemonEncounter.Pokemon.Name)
	}
	return nil
}

func catchCommand(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("catch command takes exactly one argument: name of pokemon to catch")
	}
	trimmed := strings.TrimSpace(args[0])
	pokemonName := strings.ReplaceAll(trimmed, " ", "-")
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	pokemonURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemonName)
	data, ok := cache.Get(pokemonName)
	if !ok {
		res, err := http.Get(pokemonURL)
		if err != nil {
			return fmt.Errorf("something went wrong while fetching pokemon info: %s", err)
		}
		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			return fmt.Errorf("an error occured: %s", res.Status)
		}
		if err != nil {
			return fmt.Errorf("an error occured while reading response: %s", err)
		}

		data = []byte(body)
		cache.Add(pokemonName, data)
	}
	pokemon := pokeapi.Pokemon{}
	err := json.Unmarshal(data, &pokemon)
	if err != nil {
		return err
	}
	chance := pokemon.BaseExperience
	fmt.Printf("Probability of catching this pokemon is 1 in %d\n", chance)
	randNum := rand.Intn(chance) + 1
	if randNum == chance {
		fmt.Printf("%s was caught!\nYou may now inspect it with the inspect command.", pokemonName)
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}
	pokedex[pokemonName] = pokemon

	return nil
}

func inspectCommand(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("inspect command takes exactly one argument: name of pokemon to inspect")
	}
	trimmed := strings.TrimSpace(args[0])
	pokemonName := strings.ReplaceAll(trimmed, " ", "-")
	pokemon, ok := pokedex[pokemonName]
	if !ok {
		fmt.Printf("You haven't caught %s yet!\n", pokemonName)
		return nil
	}
	fmt.Println("Name: ", pokemon.Name)
	fmt.Println("Height: ", pokemon.Height)
	fmt.Println("Weight: ", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stats := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stats.Stat.Name, stats.BaseStat)
	}
	fmt.Println("Types:")
	for _, types := range pokemon.Types {
		fmt.Printf("  - %s\n", types.Type.Name)
	}

	return nil
}

func pokedexCommand(args ...string) error {
	if len(args) != 0 {
		return fmt.Errorf("pokedex command does not accept arguments")
	}
	if len(pokedex) == 0 {
		fmt.Println("Your pokedex is empty!. You need to catch some pokemons!.")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for pokemonName := range pokedex {
		fmt.Printf("  - %s\n", pokemonName)
	}

	return nil
}
