package httpClient

import "github.com/brendontj/lol-stats/pkg/lolsports"

type LolEsportsScrapper interface {
	GetLeagues(region string) ([]*lolsports.ScheduleLeague, error)
	GetSchedule(region, leagueID, pageToken string)

}

type httpClient struct {
	baseURI string
	token string
}

func NewLolStatsClient(baseURI, token string) LolEsportsScrapper {
	return &httpClient{
		baseURI: baseURI,
		token:   token,
	}
}