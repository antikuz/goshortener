package models

import "time"

type ShortenURL struct {
	Key        string    `json:"key"`
	Secret_key string    `json:"secret_key"`
	Target_url string    `json:"target_url"`
	Is_active  bool      `json:"is_active"`
	Clicks     int       `json:"clicks"`
	Expires    time.Time `json:"expires"`
}
