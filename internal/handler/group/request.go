package group

type CreateGroupRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"required,max=255"`
}

type JoinToGroupRequest struct {
	Code string `json:"code" binding:"required"`
}

type CreateProductRequest struct {
	ProductNameID uint64 `json:"product_name_id" binding:"required"`
	Quantity      int    `json:"quantity" binding:"required,min=1,max=1000"`
}

type UpdateProductRequest struct {
	Price    *float64 `json:"price"`
	Quantity int      `json:"quantity" binding:"required,min=1,max=1000"`
	Status   string   `json:"status" binding:"required,oneof=open closed"`
}
