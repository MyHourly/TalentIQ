package model

type Technology struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	DemandCount int    `json:"demandCount"`
}
