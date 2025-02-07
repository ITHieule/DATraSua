package request

import "time"

type OrderRequest struct {
	ID           uint                  `json:"id" gorm:"primaryKey"`
	UserID       int                   `json:"user_id"`
	OrderDate    time.Time             `json:"order_date"`
	Status       string                `json:"status"`
	OrderDetails []OrderDetailsRequest `json:"order_details" gorm:"-"` // Sửa thành mảng
}

type OrderDetailsRequest struct {
	ID           int     `json:"id" gorm:"primaryKey"`
	Order_id     int     `json:"order_id"`
	Base_id      int     `json:"base_id"`
	Flavor_id    int     `json:"flavor_id"`
	Sweetness_id int     `json:"sweetness_id"`
	Ice_id       int     `json:"ice_id"`
	Size_id      int     `json:"size_id"`
	ExtraIDs     string  `json:"extra_ids"`
	Price        float64 `json:"price"`
}

// 🚀 Gán tên bảng đúng
func (OrderRequest) TableName() string {
	return "Orders" // ⚠️tên đúng trong DB
}

// 🚀 Gán tên bảng đúng
func (OrderDetailsRequest) TableName() string {
	return "OrderDetails" // ⚠️tên đúng trong DB
}
