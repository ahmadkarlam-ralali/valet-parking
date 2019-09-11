package models

import "time"

type Transaction struct {
	ID         uint `json:"id,primary_key"`
	Slot       Slot
	SlotId     uint `json:"slot_id"`
	Employee   User
	EmployeeId uint      `json:"employee_id"`
	PlatNo     string    `json:"name"`
	Total      uint      `json:"total"`
	StartAt    string    `json:"start_at"`
	EndAt      string    `json:"end_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
