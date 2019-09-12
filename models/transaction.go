package models

import "time"

type Transaction struct {
	ID        uint   `json:"id,primary_key"`
	Code      string `json:"code"`
	Slot      Slot
	SlotId    uint      `json:"slot_id"`
	PlatNo    string    `json:"name"`
	Total     uint      `json:"total"`
	StartAt   time.Time `json:"start_at"`
	EndAt     time.Time `json:"end_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
