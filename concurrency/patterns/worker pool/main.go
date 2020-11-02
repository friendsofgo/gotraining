package main

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

const apiURL = "https://pokeapi.co/api/v2/pokemon/?limit=151"

func main() {
	// Creates a work pool with n gorutines.
	p := NewPool(10)

	resources, err := getResources()
	if err != nil {
		log.Fatal(err)
	}

	w, err := preparePokemonCSVFile()
	if err != nil {
		log.Fatal(err)
	}
	defer w.Flush()

	var wg sync.WaitGroup
	for _, r := range resources.Results {
		wg.Add(1)
		r := r
		go func() {
			pkmWorker := pokemonPrinter{
				url: r.Url,
				csv: w,
			}
			p.Add(pkmWorker)
			wg.Done()
		}()
	}
	wg.Wait()

	p.Shutdown()
}

func preparePokemonCSVFile() (*csv.Writer, error) {
	csvHeader := []string{"No", "Name", "Image", "Shiny Image"}
	file, err := os.Create("pokemons.csv")
	if err != nil {
		return nil, err
	}

	writer := csv.NewWriter(file)
	if err := writer.Write(csvHeader); err != nil {
		return nil, err
	}

	return writer, nil
}

func getResources() (apiResources, error) {
	res, err := http.Get(apiURL)
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return apiResources{}, err
	}
	defer res.Body.Close()
	var r apiResources
	if err := json.Unmarshal(b, &r); err != nil {
		return apiResources{}, err
	}

	return r, nil
}

type apiResources struct {
	Results []struct {
		Url string `json:"url"`
	} `json:"results"`
}

type pokemonPrinter struct {
	url string
	csv  *csv.Writer
}

type pokemon struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Sprite struct{
		DefaultImg string `json:"front_default"`
		ShinyImg string `json:"front_shiny"`
	} `json:"sprites"`
}

func (p pokemonPrinter) Task() {
	res, err := http.Get(p.url)
	if err != nil {
		log.Println(err)
		return
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	var pkm pokemon
	if err := json.Unmarshal(b, &pkm); err != nil {
		log.Println(err)
		return
	}

	p.csv.Write([]string{
		strconv.Itoa(pkm.ID),
		pkm.Name,
		pkm.Sprite.DefaultImg,
		pkm.Sprite.ShinyImg,
	})
}
