package models

import (
	"time"

	"github.com/lib/pq"
)

type Collection struct {
	ID      string         `json:"id"`
	PID     string         `json:"pid"`
	AddTime time.Time      `json:"addtime"`
	ModTime time.Time      `json:"modtime"`
	Title   string         `json:"title"`
	Intro   string         `json:"intro"`
	ArtIDs  pq.StringArray `json:"artids"`
}
