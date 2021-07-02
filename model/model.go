package model

import "time"

// Response is a model of response of API Nomics
type Response struct {
	Currency   string      `json:"currency"`
	Timestamps []time.Time `json:"timestamps"`
	Prices     []string    `json:"prices"`
}
