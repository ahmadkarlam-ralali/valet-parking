package models

import "time"

type Building struct {
	ID        uint      `json:"id,primary_key"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
