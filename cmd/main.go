package main

import (
	"github.com/brendontj/lol-stats/app"
)

func main() {
	baseURI := "https://esports-api.lolesports.com/persisted/gw/" //Todo Add env vars
	token := "0TvQnueqKa5mxJntVWt0w4LpLfEkrV1Ta8rQBb9Z"           //Todo Add env vars
	baseURIFeed := "https://feed.lolesports.com/livestats/v1/"    //Todo Add env vars

	application := app.NewApplication(baseURI, token, baseURIFeed)
	application.Start()
	application.PopulateHistoricalData()
	application.Close()
}
