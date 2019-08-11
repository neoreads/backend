package models

type ReviewNote struct {
	ID       string `json:"id"`
	NType    int    `json:"ntype"`
	ParaID   string `json:"paraid"`
	SentID   string `json:"sentid"`
	Content  string `json:"content"`
	Progress string `json:"progress"`
}
