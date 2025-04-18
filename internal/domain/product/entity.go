package product

import "time"

type Product struct {
	ProductID     uint64
	GroupID       uint64
	ProductNameID uint64
	Price         *float64
	Status        string
	Quantity      int
	AddedBy       uint64
	BoughtBy      *uint64
	CreatedAt     time.Time
}

type Category struct {
	CategoryID uint64 `json:"category-id" db:"category_id"`
	Name       string `json:"name" db:"name"`
}

type ProductName struct {
	ProductNameID uint64 `json:"product_name_id" db:"product_name_id"`
	CategoryID    uint64 `json:"category_id" db:"category_id"`
	Name          string `json:"name" db:"name"`
}
