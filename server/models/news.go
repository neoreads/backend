package models

import (
	"time"
)

type News struct {
	ID       string    `json:"id"`
	PID      string    `json:"pid"`
	Kind     int       `json:"kind"`
	AddTime  time.Time `json:"addtime"`
	ModTime  time.Time `json:"modtime"`
	Link     string    `json:"link"`
	Source   string    `json:"source"`
	Title    string    `json:"title"`
	Summary  string    `json:"summary"`
	Content  string    `json:"content"`
	Tags     []Tag     `json:"tags"`
	TagsJSON string    `json:"tagsjson"`
}
