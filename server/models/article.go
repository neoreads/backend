package models

import "time"

type Article struct {
	ID      string    `json:"id"`
	PID     string    `json:"pid"`
	AddTime time.Time `json:"addtime"`
	ModTime time.Time `json:"modtime"`
	Kind    int       `json:"kind"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Author  string    `json:"author"`
}

type ArticleKind int

const (
	ChapterKind ArticleKind = iota
	BlogKind
	PeomKind
	EmarkKind
)
