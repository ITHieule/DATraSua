package request

type OrderDetailsRequest struct {
	ID           int     `json:"id"`
	Order_id     int     `json:"order_id"`
	Base_id      int     `json:"base_id"`
	Flavor_id    int     `json:"flavor_id"`
	Sweetness_id int     `json:"sweetness_id"`
	Ice_id       int     `json:"ice_id"`
	Size_id      int     `json:"size_id"`
	Price        float64 `json:"price"`
}
