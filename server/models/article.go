package models

import "time"

type Article struct {
	ID      string    `json:"id"`
	PID     string    `json:"pid"`
	AddTime time.Time `json:"addtime"`
	ModTime time.Time `json:"modtime"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
}
