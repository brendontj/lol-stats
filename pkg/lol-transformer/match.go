package lol_transformer

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Match struct {
	ID uuid.UUID
	TeamAName string
	TeamBName string
	GameTime time.Time
}