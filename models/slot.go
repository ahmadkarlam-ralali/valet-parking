package models

import "time"

type Slot struct {
	ID        int       `json:"id,primary_key"`
	Name      string    `gorm:"column:name" json:"name"`
	Status    string    `gorm:"column:status" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
