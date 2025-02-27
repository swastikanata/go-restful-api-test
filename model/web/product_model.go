package web

type ProductCreateRequest struct {
	Name        string  `json:"name" validate:"required,max=32,min=10"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gte=0"`
	StockQty    int     `json:"stock_qty" validate:"required,gte=0"`
	CategoryID  int     `json:"category" validate:"required"`
	SKU         string  `json:"sku" validate:"required"`
	TaxRate     float64 `json:"tax_rate" validate:"required,gte=0"`
}

type ProductUpdateRequest struct {
	Id          uint64  `json:"id" validate:"required,gte=0"`
	Name        string  `json:"name" validate:"required,max=32,min=10"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gte=0"`
	StockQty    int     `json:"stock_qty" validate:"required,gte=0"`
	CategoryID  int     `json:"category_id" validate:"required"`
	SKU         string  `json:"sku" validate:"required"`
	TaxRate     float64 `json:"tax_rate" validate:"required,gte=0"`
}

type ProductResponse struct {
	Id          uint64  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	StockQty    int     `json:"stock_qty"`
	CategoryID  uint64  `json:"category_id"`
	SKU         string  `json:"sku"`
	TaxRate     float64 `json:"tax_rate"`
}
