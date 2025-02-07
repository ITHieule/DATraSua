package types

type BaseSizestypes struct {
	Id      int `json:"id"`
	Base_id int `json:"base_id"`
	Size_id int `json:"size_id"`
}

type Basestypes struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
	Images string  `json: "images"`
}

type SizesTypes struct {
	ID    int    `gorm:"primaryKey"`
	Name  string `gorm:"unique"`
	Price float64
}
