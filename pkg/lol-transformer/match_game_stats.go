package lol_transformer

import uuid "github.com/satori/go.uuid"

type MatchGameStats struct {
	GameID uuid.UUID
	MatchExternalReference string
	TeamAName string
	TeamBName string
	TeamBlueTotalBarons int
	TeamBlueDragons []string
	TeamBlueTotalInhibitors int
	TeamBlueTotalGold int
	TeamBlueTotalKills int
	TeamBlueTotalTowers int
	TeamRedTotalBarons int
	TeamRedDragons []string
	TeamRedTotalInhibitors int
	TeamRedTotalGold int
	TeamRedTotalKills int
	TeamRedTotalTowers int
}
