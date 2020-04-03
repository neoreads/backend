package models

import (
	"time"
)

type Note struct {
	ID       string    `json:"id"`
	Time     time.Time `json:"time"`
	NType    int       `json:"ntype"`
	PType    int       `json:"ptype"`
	PID      string    `json:"pid"`
	ColID    string    `json:"colid"`
	ArtID    string    `json:"artid"`
	ParaID   string    `json:"paraid"`
	SentID   string    `json:"sentid"`
	StartPos int       `json:"startpos"`
	EndPos   int       `json:"endpos"`
	Content  string    `json:"content"`
}
