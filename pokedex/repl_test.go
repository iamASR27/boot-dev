package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iamASR27/pokedex/internal/pokeapi"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{"  hello world  ", []string{"hello", "world"}},
		{"foo bar baz", []string{"foo", "bar", "baz"}},
		{"   singleword   ", []string{"singleword"}},
		{"", []string{}},
	}

	for _, c := range cases {
		result := cleanInput(c.input)
		if len(result) != len(c.expected) {
			t.Errorf("cleanInput(%q) = %v; expected %v", c.input, result, c.expected)
			continue
		}
		for i := range result {
			if result[i] != c.expected[i] {
				t.Errorf("cleanInput(%q) = %v; expected %v", c.input, result, c.expected)
				break
			}
		}
	}
}

func TestFetchLocationAreaPageReturnsNextAndPrevious(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/location-area" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("offset"); got != "20" {
			t.Fatalf("unexpected offset query: %s", got)
		}
		if got := r.URL.Query().Get("limit"); got != "20" {
			t.Fatalf("unexpected limit query: %s", got)
		}

		_, _ = w.Write([]byte(`{"next":"/location-area?offset=40&limit=20","previous":"/location-area?offset=0&limit=20","results":[{"name":"area-21"},{"name":"area-22"}]}`))
	}))
	defer server.Close()

	page, err := pokeapi.FetchLocationAreaPage(server.URL + "/location-area?offset=20&limit=20")
	if err != nil {
		t.Fatalf("fetchLocationAreaPage returned error: %v", err)
	}

	if len(page.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(page.Results))
	}
	if page.Results[0] != "area-21" || page.Results[1] != "area-22" {
		t.Fatalf("unexpected names: %v", page.Results)
	}
	if page.Next != server.URL+"/location-area?offset=40&limit=20" {
		t.Fatalf("unexpected next URL: %s", page.Next)
	}
	if page.Previous != server.URL+"/location-area?offset=0&limit=20" {
		t.Fatalf("unexpected previous URL: %s", page.Previous)
	}
}
