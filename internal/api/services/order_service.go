package services

import (
	"fmt"
	"time"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
)

// Struct trung gian KHÔNG chứa Extras
type CartDB struct {
	ID          int     `json:"id"`
	UserID      int     `json:"user_id"`
	BaseID      int     `json:"base_id"`
	SizeID      int     `json:"size_id"`
	FlavorID    int     `json:"flavor_id"`
	SweetnessID int     `json:"sweetness_id"`
	IceID       int     `json:"ice_id"`
	ExtraIDs    string  `json:"extra_ids"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

// 🚀 Chỉ định bảng thực sự là `carts`
func (CartDB) TableName() string {
	return "Cart" // Tên bảng thật trong database
}

type OrderService struct{}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) PlaceOrder(userID int) (*request.OrderRequest, error) {
	db, err := database.DB1Connection()
	if err != nil {
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// 🔹 Truy vấn giỏ hàng (DÙNG `CartDB` để tránh lỗi)
	var cartItems []CartDB
	err = db.Where("user_id = ?", userID).Find(&cartItems).Error
	if err != nil {
		return nil, err
	}
	if len(cartItems) == 0 {
		return nil, fmt.Errorf("Cart is empty")
	}

	// 🔹 Tạo đơn hàng
	order := request.OrderRequest{
		UserID:    userID,
		OrderDate: time.Now(),
		Status:    "Đang xử lý",
	}
	err = db.Create(&order).Error
	if err != nil {
		return nil, err
	}

	// 🔹 Chuyển từ Cart → OrderDetails
	var orderDetails []request.OrderDetailsRequest
	for _, cart := range cartItems {
		// 🔹 Lấy danh sách Extras từ ExtraIDs
		extras, err := GetExtrasFromIDs(db, cart.ExtraIDs)
		if err != nil {
			return nil, err
		}

		// 🔹 Thêm vào order_details
		orderDetails = append(orderDetails, request.OrderDetailsRequest{
			Order_id:     int(order.ID),
			Base_id:      cart.BaseID,
			Flavor_id:    cart.FlavorID,
			Sweetness_id: cart.SweetnessID,
			Ice_id:       cart.IceID,
			Size_id:      cart.SizeID,
			ExtraIDs:     cart.ExtraIDs,
			Price:        cart.Price,
		})

		// Debug danh sách Extras
		fmt.Printf("Cart ID: %d, Extras: %+v\n", cart.ID, extras)
	}

	// Lưu order_details vào DB
	for i := range orderDetails {
		orderDetails[i].Order_id = int(order.ID) // 🚀 Gán Order_id trước khi lưu
	}
	err = db.Create(&orderDetails).Error
	if err != nil {
		return nil, err
	}

	// err = db.Table("Cart").Where("user_id = ?", userID).Delete(nil).Error

	if err != nil {
		return nil, err
	}

	order.OrderDetails = orderDetails
	return &order, nil
}
