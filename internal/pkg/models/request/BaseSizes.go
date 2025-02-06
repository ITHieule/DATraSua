package request

type BaseSizesrequest struct {
	Id      int `json:"id"`
	Base_id int `json:"base_id"`
	Size_id int `json:"size_id"`
}

type Basesrequest struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
	Images string  `json:"images"`
}

type SizesRequest struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
