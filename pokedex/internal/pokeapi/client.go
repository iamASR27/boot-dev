package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/iamASR27/pokedex/internal/pokecache"
)

const BaseURL = "https://pokeapi.co/api/v2"

type LocationAreaPage struct {
	Results  []string
	Next     string
	Previous string
}

type LocationAreaDetails struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type PokemonData struct {
	Name           string `json:"name"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	BaseExperience int    `json:"base_experience"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

type Client struct {
	cache *pokecache.Cache
}

// NewClient creates a new PokeAPI client with caching
func NewClient(cacheInterval time.Duration) *Client {
	return &Client{
		cache: pokecache.NewCache(cacheInterval),
	}
}

// FetchLocationAreaPage fetches a page of location areas from the PokeAPI
func (c *Client) FetchLocationAreaPage(endpoint string) (LocationAreaPage, error) {
	// Check if data is in cache
	if data, ok := c.cache.Get(endpoint); ok {
		fmt.Println("(from cache)")
		var payload struct {
			Next     string `json:"next"`
			Previous string `json:"previous"`
			Results  []struct {
				Name string `json:"name"`
			} `json:"results"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return LocationAreaPage{}, err
		}

		results := make([]string, 0, len(payload.Results))
		for _, result := range payload.Results {
			results = append(results, result.Name)
		}

		baseURL, err := url.Parse(endpoint)
		if err != nil {
			return LocationAreaPage{}, err
		}

		nextURL := ""
		if payload.Next != "" {
			nextRef, err := url.Parse(payload.Next)
			if err != nil {
				return LocationAreaPage{}, err
			}
			nextURL = baseURL.ResolveReference(nextRef).String()
		}

		previousURL := ""
		if payload.Previous != "" {
			previousRef, err := url.Parse(payload.Previous)
			if err != nil {
				return LocationAreaPage{}, err
			}
			previousURL = baseURL.ResolveReference(previousRef).String()
		}

		return LocationAreaPage{
			Results:  results,
			Next:     nextURL,
			Previous: previousURL,
		}, nil
	}

	// Make HTTP request
	resp, err := http.Get(endpoint)
	if err != nil {
		return LocationAreaPage{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return LocationAreaPage{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the body into memory so we can cache it
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaPage{}, err
	}

	// Cache the response
	c.cache.Add(endpoint, body)

	var payload struct {
		Next     string `json:"next"`
		Previous string `json:"previous"`
		Results  []struct {
			Name string `json:"name"`
		} `json:"results"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return LocationAreaPage{}, err
	}

	results := make([]string, 0, len(payload.Results))
	for _, result := range payload.Results {
		results = append(results, result.Name)
	}

	baseURL, err := url.Parse(endpoint)
	if err != nil {
		return LocationAreaPage{}, err
	}

	nextURL := ""
	if payload.Next != "" {
		nextRef, err := url.Parse(payload.Next)
		if err != nil {
			return LocationAreaPage{}, err
		}
		nextURL = baseURL.ResolveReference(nextRef).String()
	}

	previousURL := ""
	if payload.Previous != "" {
		previousRef, err := url.Parse(payload.Previous)
		if err != nil {
			return LocationAreaPage{}, err
		}
		previousURL = baseURL.ResolveReference(previousRef).String()
	}

	return LocationAreaPage{
		Results:  results,
		Next:     nextURL,
		Previous: previousURL,
	}, nil
}

// FetchLocationArea fetches detailed information about a location area including pokemon
func (c *Client) FetchLocationArea(locationAreaName string) (LocationAreaDetails, error) {
	endpoint := BaseURL + "/location-area/" + locationAreaName

	// Check if data is in cache
	if data, ok := c.cache.Get(endpoint); ok {
		fmt.Println("(from cache)")
		var locationArea LocationAreaDetails
		if err := json.Unmarshal(data, &locationArea); err != nil {
			return LocationAreaDetails{}, err
		}
		return locationArea, nil
	}

	// Make HTTP request
	resp, err := http.Get(endpoint)
	if err != nil {
		return LocationAreaDetails{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return LocationAreaDetails{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the body into memory so we can cache it
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaDetails{}, err
	}

	// Cache the response
	c.cache.Add(endpoint, body)

	var locationArea LocationAreaDetails
	if err := json.Unmarshal(body, &locationArea); err != nil {
		return LocationAreaDetails{}, err
	}

	return locationArea, nil
}

// FetchPokemon fetches information about a specific Pokemon
func (c *Client) FetchPokemon(pokemonName string) (PokemonData, error) {
	endpoint := BaseURL + "/pokemon/" + pokemonName

	// Check if data is in cache
	if data, ok := c.cache.Get(endpoint); ok {
		fmt.Println("(from cache)")
		var pokemon PokemonData
		if err := json.Unmarshal(data, &pokemon); err != nil {
			return PokemonData{}, err
		}
		return pokemon, nil
	}

	// Make HTTP request
	resp, err := http.Get(endpoint)
	if err != nil {
		return PokemonData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PokemonData{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the body into memory so we can cache it
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokemonData{}, err
	}

	// Cache the response
	c.cache.Add(endpoint, body)

	var pokemon PokemonData
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return PokemonData{}, err
	}

	return pokemon, nil
}
