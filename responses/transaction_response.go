package responses

type TransactionStartResponse struct {
	PlatNo   string `json:"plat_no"`
	SlotName string `json:"slot_name"`
	Code     string `json:"code"`
}

type TransactionEndResponse struct {
	Total uint `json:"total"`
}
