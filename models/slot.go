package models

import "time"

type Slot struct {
	ID         uint      `json:"id,primary_key"`
	BuildingID uint      `json:"building_id"`
	Building   Building  `json:"building"`
	Name       string    `json:"name"`
	Total      int       `json:"total"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
