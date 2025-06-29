package gbot

import "time"

type Item struct {
	ID    string
	Name  string
	Value int64
}

type BotResult struct {
	ProjName string
	Date     time.Time
	Total    int64
	Night    int64

	BeforWeekDay int64

	Items []Item
}

type JobParam struct {
	ProjID string   `json:"proj_id"`
	Tags   []string `jaon:"tags"`
}
