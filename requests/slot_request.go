package requests

type SlotStoreRequest struct {
	Name  string `form:"name" json:"name" binding:"required"`
	Total int    `form:"total" json:"total" binding:"required"`
}

type SlotUpdateRequest struct {
	Name  string `form:"name" json:"name" binding:"required"`
	Total int    `form:"total" json:"total" binding:"required"`
}
