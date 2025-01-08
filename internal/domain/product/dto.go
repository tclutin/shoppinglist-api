package product

import "time"

type ProductDTO struct {
	ProductID   uint64    `json:"product_id" db:"product_id"`
	ProductName string    `json:"product_name" db:"product_name"`
	Category    string    `json:"category" db:"category_name"`
	Price       *float64  `json:"price" db:"price"`
	Quantity    int       `json:"quantity" db:"quantity"`
	AddedBy     string    `json:"added_by" db:"added_by"`
	BoughtBy    *string   `json:"bought_by" db:"bought_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
