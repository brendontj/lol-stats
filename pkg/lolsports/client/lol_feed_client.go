package httpClient

import (
	"github.com/brendontj/lol-stats/pkg/lolsports"
	"net/http"
	"time"
)

type FeedAPIScrapper interface {
	GetDataFromLiveMatch(gameID, startingTime string) (*lolsports.LiveMatchData, error)
	GetDetailsFromLiveMatch(gameID string) (*lolsports.LiveMatchDetailData, error)
}

type LolFeedClient struct {
	baseURI string
	token string
	httpClient *http.Client
}

func (l LolFeedClient) GetDataFromLiveMatch(gameID, startingTime string) (*lolsports.LiveMatchData, error) {
	panic("implement me")
}

func (l LolFeedClient) GetDetailsFromLiveMatch(gameID string) (*lolsports.LiveMatchDetailData, error) {
	panic("implement me")
}

func NewLolFeedClient(baseURI, token string) EsportsAPIScrapper {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	return &PersistedDataClient{
		baseURI: baseURI,
		token:   token,
		httpClient: &http.Client{
			Timeout:   10 * time.Second,
			Transport: t,
		},
	}
}