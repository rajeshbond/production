package products

type CreateProductRequest struct {
	ProductName string `json:"product_name" validate:"required"`
	ProductNo   string `json:"product_no" validate:"required"`
}

type ProductResponse struct {
	ID          int64  `json:"id"`
	ProductName string `json:"product_name"`
	ProductNo   string `json:"product_no"`
}
