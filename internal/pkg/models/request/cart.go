package request

type Cartrequest struct {
	ID          int              `json:"id"`
	UserID      int              `json:"user_id"`
	BaseID      int              `json:"base_id"`
	SizeID      int              `json:"size_id"`
	FlavorID    int              `json:"flavor_id"`
	SweetnessID int              `json:"sweetness_id"`
	IceID       int              `json:"ice_id"`
	ExtraIDs    string           `json:"extra_ids"` // Dữ liệu phụ kiện dưới dạng chuỗi
	Quantity    int              `json:"quantity"`
	Price       float64          `json:"price"`
	Base        Basesrequest     `json:"base"`
	Size        SizesRequest     `json:"size"`
	Flavor      Flavorsrequest   `json:"flavor"`
	Sweetness   Sweetnessrequest `json:"sweetness"`
	Ice         IceLevelsrequest `json:"ice"`
	Extras      []Extrasrequest  `json:"extras"` // Danh sách phụ kiện liên kết
}

type CartDetails struct {
	ID            int     `json:"id"`
	UserID        int     `json:"user_id"`
	BaseID        int     `json:"base_id"`
	SizeID        int     `json:"size_id"`
	FlavorID      int     `json:"flavor_id"`
	SweetnessID   int     `json:"sweetness_id"`
	IceID         int     `json:"ice_id"`
	ExtraIDs      string  `json:"extra_ids"`
	Quantity      int     `json:"quantity"`
	Price         float64 `json:"price"`
	BaseName      string  `json:"base_name"`
	BasePrice     float64 `json:"base_price"`
	SizeName      string  `json:"size_name"`
	SizePrice     float64 `json:"size_price"`
	FlavorName    string  `json:"flavor_name"`
	SweetnessName string  `json:"sweetness_name"`
	IceName       string  `json:"ice_name"`
	ExtraID       int     `json:"extra_id"`
	ExtraName     string  `json:"extra_name"`
	ExtraPrice    float64 `json:"extra_price"`
}
