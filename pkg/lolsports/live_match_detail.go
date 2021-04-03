package lolsports

type LiveMatchDetailData struct {
	EsportsGameID  string       `json:"esportsGameId"`
	EsportsMatchID string       `json:"esportsMatchId"`
	GameMetadata   GameMetaData `json:"gameMetadata"`
	Frames         []Frames     `json:"frames"`
}

type ParticipantMetadata struct {
	ParticipantID   int    `json:"participantId"`
	EsportsPlayerID string `json:"esportsPlayerId"`
	SummonerName    string `json:"summonerName"`
	ChampionID      string `json:"championId"`
	Role            string `json:"role"`
}

type BlueTeamMetadata struct {
	EsportsTeamID       string                `json:"esportsTeamId"`
	ParticipantMetadata []ParticipantMetadata `json:"participantMetadata"`
}

type RedTeamMetadata struct {
	EsportsTeamID       string                `json:"esportsTeamId"`
	ParticipantMetadata []ParticipantMetadata `json:"participantMetadata"`
}

type GameMetaData struct {
	PatchVersion     string           `json:"patchVersion"`
	BlueTeamMetadata BlueTeamMetadata `json:"blueTeamMetadata"`
	RedTeamMetadata  RedTeamMetadata  `json:"redTeamMetadata"`
}
