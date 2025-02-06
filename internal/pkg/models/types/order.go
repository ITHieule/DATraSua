package types

type OrderTypes struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	UserID    int    `json:"user_id"`
	OrderDate string `json:"order_date"`
	Status    string `json:"status"`
}

// Cấu trúc chứa thông tin đơn hàng và chi tiết đơn hàng
type OrderWithDetails struct {
	Order        OrderTypes          `json:"order"`
	OrderDetails []OrderDetailsTypes `json:"order_details"`
}
