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

func (s *OrderService) GetOrderDetailsByOrderID(orderID int) ([]request.OrderDetailsRequest, error) {
	db, err := database.DB1Connection()
	if err != nil {
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// 🔹 Truy vấn danh sách OrderDetails theo orderID
	var orderDetails []request.OrderDetailsRequest
	err = db.Where("order_id = ?", orderID).Find(&orderDetails).Error
	if err != nil {
		return nil, err
	}

	// 🔹 Debug danh sách trả về
	fmt.Printf("Order ID: %d, Details: %+v\n", orderID, orderDetails)

	return orderDetails, nil
}

func (s *OrderService) GetOrdersByUserID(userID int) ([]request.OrderRequest, error) {
	db, err := database.DB1Connection()
	if err != nil {
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// 📌 Truy vấn danh sách đơn hàng theo UserID
	var orders []request.OrderRequest
	err = db.Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	// 📌 Lặp qua từng đơn hàng để lấy danh sách OrderDetails
	for i := range orders {
		var orderDetails []request.OrderDetailsRequest
		err := db.Where("order_id = ?", orders[i].ID).Find(&orderDetails).Error
		if err != nil {
			return nil, err
		}
		orders[i].OrderDetails = orderDetails
	}

	return orders, nil
}

func (s *OrderService) CancelOrder(orderID int) error {
	db, err := database.DB1Connection()
	if err != nil {
		return err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// 📌 Kiểm tra đơn hàng có tồn tại không
	var order request.OrderRequest
	err = db.Where("id = ?", orderID).First(&order).Error
	if err != nil {
		return err
	}

	// 📌 Cập nhật trạng thái đơn hàng thành "Đã hủy"
	order.Status = "Đã hủy"
	err = db.Save(&order).Error
	if err != nil {
		return err
	}

	return nil
}

// 📌 Cập nhật trạng thái đơn hàng
func (s *OrderService) UpdateOrderStatus(orderID int, status string) error {
	db, err := database.DB1Connection()
	if err != nil {
		return err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// 📌 Kiểm tra đơn hàng có tồn tại không
	var order request.OrderRequest
	err = db.Where("id = ?", orderID).First(&order).Error
	if err != nil {
		return err
	}

	// 📌 Cập nhật trạng thái đơn hàng
	order.Status = status
	err = db.Save(&order).Error
	if err != nil {
		return err
	}

	return nil
}
