package typesjson

type ProductGroupForBuyer struct {
	Data  ProductGroupForBuyerData `json:"data"`
	Error string                   `json:"error,omitempty"`
}

type ProductGroupForBuyerData struct {
	CalculatedTotal int    `json:"calculatedTotal"`
	InProgress      bool   `json:"inProgress"`
	OrderHashId     string `json:"orderHashId"`
	Originator      string `json:"originator"`
	ShopName        string `json:"shopName"`
	ShopTelNo       string `json:"shopTelNo"`
	Size            int    `json:"size"`
	Total           int    `json:"total"`
	TotalBuyer      int    `json:"totalBuyer"`

	Rows []ProductGroupForBuyerRow `json:"rows"`
}

type ProductGroupForBuyerRow struct {
	CalculatedTotal      int    `json:"calculatedTotal"`
	FirstCalculatedPrice int    `json:"firstCalculatedPrice"`
	FirstPrice           int    `json:"firstPrice"`
	MergedKey            string `json:"mergedKey"`
	Name                 string `json:"name"`
	Size                 int    `json:"size"`
	Total                int    `json:"total"`

	Items []ProductGroupForBuyerItem `json:"items"`
}

type ProductGroupForBuyerItem struct {
	Cancelable   bool   `json:"cancelable"`
	Comment      string `json:"comment"`
	FullName     string `json:"fullName"`
	MergedKey    string `json:"mergedKey"`
	MergedName   string `json:"mergedName"`
	OrderItemIds []int  `json:"orderItemIds"`
	Paid         bool   `json:"paid"`
	Shipped      bool   `json:"shipped"`
	Size         int    `json:"size"`
	Total        int    `json:"total"`
}
