package types

// An item is a purchased item on a receipt
type Item struct {
	ShortDescription string `json:"shortDescription" binding:"required"`
	Price            string `json:"price" binding:"required"`
}

// A receipt is a listing of purchased items, and other metadata
type Receipt struct {
	ID           string `json:"id"`
	Retailer     string `json:"retailer" binding:"required"`
	PurchaseDate string `json:"purchaseDate" binding:"required" time_format:"2006-01-02"`
	PurchaseTime string `json:"purchaseTime" binding:"required" time_format:"hh:mm"`
	Total        string `json:"total" binding:"required"`
	Items        []Item `json:"items" binding:"required"`
}

// Stores the result of the point calculation of every item in a receipt
type PointRuleItem struct {
	Description       string
	Price             string
	DescriptionLength int
	Value             float64
}

// Stores the calculated results for a point calculation for a receipt
type PointRules struct {
	AlphanumericPoints int
	RoundDollarPoints  int
	MultiplierPoints   int
	GroupingPoints     int
	GroupingAmount     int
	PurchaseDatePoints int
	PurchaseTimePoints int
	RuleItems          []PointRuleItem
}
