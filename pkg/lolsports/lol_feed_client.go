package lolsports

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type FeedAPIScrapper interface {
	GetDataFromLiveMatch(gameID string, startingTime time.Time) (*LiveMatchData, error)
	GetDetailsFromLiveMatch(gameID string, startingTime time.Time) (*LiveMatchDetailData, error)
	Close()
}

type LolFeedClient struct {
	baseURI string
	token string
	httpClient *http.Client
}

func (l LolFeedClient) GetDataFromLiveMatch(gameID string, startingTime time.Time) (*LiveMatchData, error) {
	timeOfBegin := time.Date(1950,1,1,0,0,0,0,time.UTC)
	currentTimestamp := ""
	if startingTime.After(timeOfBegin) {
		currentTimestamp = fmt.Sprint("?startingTime=", startingTime.String())
	}
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprint(l.baseURI,"details/", gameID, currentTimestamp),
		nil,
	)
	if err != nil {
		log.Printf("error creating HTTP request for get data from live match: gameID: %s, currentTime: %s, err: %v\n",gameID, startingTime.String(), err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-api-key", l.token)

	res, err := l.httpClient.Do(req)
	if err != nil {
		log.Printf("error sending HTTP request for get data from live match: gameID: %s, currentTime: %s, err: %v\n",gameID, startingTime.String(), err)
		return nil, err
	}

	liveMatch := new(LiveMatchData)

	if err := json.NewDecoder(res.Body).Decode(&liveMatch); err != nil {
		log.Printf("error deserializing weather data\n")
		return nil, err
	}
	return liveMatch, nil
}

func (l LolFeedClient) GetDetailsFromLiveMatch(gameID string, startingTime time.Time) (*LiveMatchDetailData, error) {
	timeOfBegin := time.Date(1950,1,1,0,0,0,0,time.UTC)
	currentTimestamp := ""
	if startingTime.After(timeOfBegin) {
		currentTimestamp = fmt.Sprint("?startingTime=", startingTime.String())
	}
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprint(l.baseURI,"window/", gameID, currentTimestamp),
		nil,
	)
	if err != nil {
		log.Printf("error creating HTTP request for get data from live match detail: gameID: %s, currentTime: %s, err: %v\n",gameID, startingTime.String(), err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-api-key", l.token)

	res, err := l.httpClient.Do(req)
	if err != nil {
		log.Printf("error sending HTTP request for get data from live match detail: gameID: %s, currentTime: %s, err: %v\n",gameID, startingTime.String(), err)
		return nil, err
	}

	liveMatchDetail := new(LiveMatchDetailData)

	if err := json.NewDecoder(res.Body).Decode(&liveMatchDetail); err != nil {
		log.Printf("error deserializing weather data\n")
		return nil, err
	}
	return liveMatchDetail, nil
}

func (l *LolFeedClient) Close() {
	l.httpClient = nil
}

func NewLolFeedClient(baseURI, token string) FeedAPIScrapper {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	return &LolFeedClient{
		baseURI: baseURI,
		token:   token,
		httpClient: &http.Client{
			Timeout:   10 * time.Second,
			Transport: t,
		},
	}
}