package poke

import (
	"bytes"
	"encoding/json"
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

type CacheInterface interface {
	Add(key string, val []byte)
	Get(key string) ([]byte, bool)
}

const PokeLocationUrlFirstPage = "https://pokeapi.co/api/v2/location-area/"

func FetchPokeLocation(url string, cache CacheInterface) (PokeApiResp, bool, error) {
	cacheVal, ok := cache.Get(url)
	var body []byte
	cachedResp := false
	if ok {
		body = cacheVal
		cachedResp = true
	} else {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return PokeApiResp{}, cachedResp, err
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return PokeApiResp{}, cachedResp, err
		}
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return PokeApiResp{}, cachedResp, err
		}
		defer resp.Body.Close()
	}
	decoder := json.NewDecoder(bytes.NewReader(body))
	apiResponse := PokeApiResp{}
	if err := decoder.Decode(&apiResponse); err != nil {
		return apiResponse, cachedResp, err
	}
	return apiResponse, cachedResp, nil
}
