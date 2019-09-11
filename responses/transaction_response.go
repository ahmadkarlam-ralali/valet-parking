package responses

type TransactionStartResponse struct {
	PlatNo   string `json:"plat_no"`
	SlotName string `json:"slot_name"`
	SlotID   uint   `json:"slot_id"`
}

type TransactionEndResponse struct {
	Total uint `json:"total"`
}
