package main

import (
	"github.com/brendontj/lol-stats/app"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	initRoutes()
}

func initRoutes() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		 c.Data(http.StatusOK, "text/html", []byte("pong"))
	})

	router.GET("/sync_data", func(c *gin.Context) {
		baseURI := "https://esports-api.lolesports.com/persisted/gw/" //Todo Add env vars
		token := "0TvQnueqKa5mxJntVWt0w4LpLfEkrV1Ta8rQBb9Z"           //Todo Add env vars
		baseURIFeed := "https://feed.lolesports.com/livestats/v1/"    //Todo Add env vars

		application := app.NewApplication(baseURI, token, baseURIFeed)
		application.Start()
		application.PopulateHistoricalData()
		application.Close()

		c.Data(http.StatusOK, "text/html", []byte("synced"))
	})

	if err := router.Run(); err != nil {
		panic(err)
	}
}
