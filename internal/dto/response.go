package dto

type PriceBreak struct {
	Quantity              int     `json:"quantity"`
	Price                 float64 `json:"price"`
	CostWithDelivery      float64 `json:"cost_with_delivery"`
	TargetPricePurchasing float64 `json:"target_price_purchasing"`
	TargetPriceSales      float64 `json:"target_price_sales"`
	Currency              string  `json:"currency"`
}

type SearchItem struct {
	MPN          string       `json:"mpn"`
	RequestedMPN string       `json:"requested_mpn"`
	RequestedQty int          `json:"requested_quantity"`
	Manufacturer string       `json:"manufacturer"`
	Stock        int          `json:"stock"`
	Status       string       `json:"status"`
	Price        float64      `json:"price"`
	Currency     string       `json:"currency"`
	PriceBreaks  []PriceBreak `json:"priceBreaks"`
	Supplier     string       `json:"supplier"`
}
