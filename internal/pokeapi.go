package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PokeMap struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func ApiLocationFetch() error {
	fullURL := "https://pokeapi.co/api/v2/location-area"

	res, err := http.Get(fullURL)
	if err != nil {
		return fmt.Errorf("error getting locations: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading body: %v", err)
	}

	var data PokeMap
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return fmt.Errorf("could not unmarshal json: %v", err)
	}
	return nil
}
