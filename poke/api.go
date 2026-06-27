package poke

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PokeApiResp struct {
	Results  []PokeLocation
	Next     string
	Previous string
	Count    int
}

type PokeLocation struct {
	Name string
	Url  string
}

type LocationDetail struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Stats          []PokemonStats `josn:"stats"`
	Types          []PokemonType  `json:"types"`
	Name           string         `json:"name"`
	URL            string         `json:"url"`
	BaseExperience int            `json:"base_experience"`
	Height         int            `json:"height"`
	Weight         int            `json:"weight"`
}

type PokemonStats struct {
	Stat struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"stat"`
	Effort   int `json:"effor"`
	BaseStat int `json:"base_stat"`
}

type PokemonType struct {
	Type struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"type"`
}

type CacheInterface interface {
	Add(key string, val []byte)
	Get(key string) ([]byte, bool)
}

const (
	pokeBaseUrl              = "https://pokeapi.co/api/v2/"
	PokeLocationUrlFirstPage = pokeBaseUrl + "location-area/"
)

func FetchPokeLocation(url string, cache CacheInterface) (PokeApiResp, bool, error) {
	cacheVal, ok := cache.Get(url)
	var body []byte
	var err error
	cachedResp := false
	if ok {
		body = cacheVal
		cachedResp = true
	} else {
		body, err = get(url)
		if err != nil {
			return PokeApiResp{}, cachedResp, err
		}
		cache.Add(url, body)
	}
	decoder := json.NewDecoder(bytes.NewReader(body))
	apiResponse := PokeApiResp{}
	if err := decoder.Decode(&apiResponse); err != nil {
		return apiResponse, cachedResp, err
	}
	return apiResponse, cachedResp, nil
}

func FetchPokemonsInLocation(location string, cache CacheInterface) ([]PokemonEncounter, bool, error) {
	url := fmt.Sprintf("%s%s/", PokeLocationUrlFirstPage, location)
	fmt.Println(url)
	cacheVal, ok := cache.Get(url)
	var body []byte
	var err error
	cachedResp := false
	if ok {
		body = cacheVal
		cachedResp = true
	} else {
		body, err = get(url)
		if err != nil {
			return []PokemonEncounter{}, cachedResp, err
		}
		cache.Add(url, body)
	}
	decoder := json.NewDecoder(bytes.NewReader(body))
	apiResponse := LocationDetail{}
	if err := decoder.Decode(&apiResponse); err != nil {
		return apiResponse.PokemonEncounters, cachedResp, err
	}
	return apiResponse.PokemonEncounters, cachedResp, nil
}

func FetchPokemon(name string, cache CacheInterface) (Pokemon, error) {
	url := fmt.Sprintf("%spokemon/%s", pokeBaseUrl, name)
	cacheValue, ok := cache.Get(url)
	var body []byte
	var err error
	if ok {
		body = cacheValue
	} else {
		body, err = get(url)
		if err != nil {
			return Pokemon{}, err
		}
		cache.Add(url, body)
	}
	decoder := json.NewDecoder(bytes.NewReader(body))
	result := Pokemon{}
	if err := decoder.Decode(&result); err != nil {
		return result, err
	}
	return result, nil

}

func get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	if res.StatusCode > 299 {
		return []byte{}, fmt.Errorf("API issue code: %v", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()
	return body, nil
}
