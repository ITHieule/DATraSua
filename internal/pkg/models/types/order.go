package types

type OrderTypes struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	UserID    int    `json:"user_id"`
	OrderDate string `json:"order_date"`
	Status    string `json:"status"`
}

type OrdersTypes struct {
	Order_id      uint   `json:"id" gorm:"primaryKey"`
	BasesName     string `json:"basesname"`
	FlavorsName   string `json:"flavorsName"`
	SweetnessName string `json:"sweetnessName"`
	IceLevelsName string `json:"iceLevelsName"`
	SizeName      string `json:"sizeName"`
}
