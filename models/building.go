package models

import "time"

type Building struct {
	ID        uint      `json:"id,primary_key"`
	Slots     []Slot    `json:"slots"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
