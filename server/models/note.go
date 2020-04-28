package models

import (
	"time"
)

type Note struct {
	ID       string    `json:"id"`
	Time     time.Time `json:"time"`
	NType    int       `json:"ntype"` // Note Type
	PType    int       `json:"ptype"` // Position Type
	PID      string    `json:"pid"`
	ColID    string    `json:"colid"`
	ArtID    string    `json:"artid"`
	ParaID   string    `json:"paraid"`
	SentID   string    `json:"sentid"`
	StartPos int       `json:"startpos"`
	EndPos   int       `json:"endpos"`
	Content  string    `json:"content"`
	Value    int       `json:"value"`
}

type NoteCard struct {
	Note

	PName string `json:"pname"` // Person Name
}
