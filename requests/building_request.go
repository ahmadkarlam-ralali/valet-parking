package requests

type BuildingStoreRequest struct {
	Name string `form:"name" json:"name" binding:"required"`
}

type BuildingUpdateRequest struct {
	Name string `form:"name" json:"name" binding:"required"`
}
