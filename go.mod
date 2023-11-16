module github.com/OkoliDaniel/pokedexcli

go 1.21.2

replace internal/pokeapi v0.0.0 => ./internal/pokeapi
replace internal/pokecache v0.0.0 => ./internal/pokecache

require(
    internal/pokeapi v0.0.0
    internal/pokecache v0.0.0
)
