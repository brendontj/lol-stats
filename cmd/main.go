package main

import (
	"fmt"
	httpClient "github.com/brendontj/lol-stats/pkg/lolsports/client"
)

func main() {
	baseURI := "https://esports-api.lolesports.com/persisted/gw/"
	token := "0TvQnueqKa5mxJntVWt0w4LpLfEkrV1Ta8rQBb9Z"
	scrapper := httpClient.NewLolStatsClient(baseURI, token)
	leagues, err := scrapper.GetLeagues("pt-BR")
	if err != nil {
		panic("Error getting leagues")
	}
	fmt.Println(leagues)
}
