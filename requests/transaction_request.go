package requests

type TransactionStartRequest struct {
	PlatNo string `form:"plat_no" json:"plat_no" binding:"required"`
}

type TransactionEndRequest struct {
	SlotName string `form:"slot_name" json:"slot_name" binding:"required"`
	PlatNo   string `form:"plat_no" json:"plat_no" binding:"required"`
}
