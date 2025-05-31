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

type Area struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Cameras []string `json:"cameras"`
}

type JobParam struct {
	Proj  string   `json:"proj"`
	Entry []string `json:"entry"`

	Areas []Area `json:"areas,omitempty"`

	Date string `json:"-"`
}
