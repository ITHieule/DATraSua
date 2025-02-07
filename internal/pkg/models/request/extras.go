package request

type Extrasrequest struct {
	Id    int     `json:"id" gorm:"primaryKey"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// 🚀 Gán tên bảng đúng
func (Extrasrequest) TableName() string {
	return "Extras" // ⚠️tên đúng trong DB
}
