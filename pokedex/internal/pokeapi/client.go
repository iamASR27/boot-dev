package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const BaseURL = "https://pokeapi.co/api/v2"

type LocationAreaPage struct {
	Results  []string
	Next     string
	Previous string
}

func FetchLocationAreaPage(endpoint string) (LocationAreaPage, error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return LocationAreaPage{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return LocationAreaPage{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var payload struct {
		Next     string `json:"next"`
		Previous string `json:"previous"`
		Results  []struct {
			Name string `json:"name"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
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
