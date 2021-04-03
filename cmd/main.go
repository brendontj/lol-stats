package main

import (
	"encoding/json"
	"github.com/brendontj/lol-stats/pkg/lolsports"
	"log"
	"net/http"
)

func main() {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://esports-api.lolesports.com/persisted/gw/getLeagues?hl=pt-BR",
		nil,
	)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-api-key", "0TvQnueqKa5mxJntVWt0w4LpLfEkrV1Ta8rQBb9Z")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error sending HTTP request: %v", err)
	}

	var data lolsports.DataLeague
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		log.Fatalf("error deserializing weather data")
	}

	log.Println("We got the response:", data.DataLeagues.Leagues[0])
}
