package requests

type TransactionStartRequest struct {
	PlatNo     string `form:"plat_no" json:"plat_no" binding:"required"`
	BuildingID uint   `form:"building_id" json:"building_id" binding:"required"`
}

type TransactionEndRequest struct {
	Code string `form:"code" json:"code" binding:"required"`
}
