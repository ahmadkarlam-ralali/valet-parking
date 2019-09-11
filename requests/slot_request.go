package requests

type SlotStoreRequest struct {
	Name string `form:"name" json:"name" binding:"required"`
}

type SlotUpdateRequest struct {
	Name   string `form:"name" json:"name"`
	Status string `form:"status" json:"status"`
}
