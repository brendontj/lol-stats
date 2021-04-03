package lolsports

type EventDetailData struct {
	Data Data `json:"data"`
}

type Tournament struct {
	TournamentID string `json:"id"`
}

type Games struct {
	Number int     `json:"number"`
	ID     string  `json:"id"`
	State  string  `json:"state"`
	Teams  []Teams `json:"teams"`
}

type Event struct {
	ID         string        `json:"id"`
	Type       string        `json:"type"`
	Tournament Tournament    `json:"tournament"`
	League     League        `json:"league"`
	Match      Match         `json:"match"`
}

type Data struct {
	Event Event `json:"event"`
}