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

type Config struct {
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
}

var ApiConfig Config

func ApiLocationFetch(config *string) (PokeMap, error) {
	baseUrl := ""
	var data PokeMap

	if config == nil {
		baseUrl = "https://pokeapi.co/api/v2/location-area/"
	} else {
		baseUrl = *config
	}

	v, b := cache.Get(baseUrl)

	if b {
		if err := json.Unmarshal(v, &data); err != nil {
			return PokeMap{}, fmt.Errorf("could not unmarshal json from cache: %v", err)
		}
	} else {
		res, err := http.Get(baseUrl)
		if err != nil {
			return PokeMap{}, fmt.Errorf("error getting locations: %v", err)
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return PokeMap{}, fmt.Errorf("error reading body: %v", err)
		}

		if err := json.Unmarshal([]byte(body), &data); err != nil {
			return PokeMap{}, fmt.Errorf("could not unmarshal json: %v", err)
		}
		cache.Add(baseUrl, body)
	}

	ApiConfig.Next = data.Next
	ApiConfig.Previous = data.Previous

	return data, nil
}
