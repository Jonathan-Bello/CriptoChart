package model

import "time"

type Response struct {
	Currency   string      `json:"currency"`
	Timestamps []time.Time `json:"timestamps"`
	Prices     []string    `json:"prices"`
}
