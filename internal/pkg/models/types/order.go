package types

type OrderTypes struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	UserID    int    `json:"user_id"`
	OrderDate string `json:"order_date"`
	Status    string `json:"status"`
}
