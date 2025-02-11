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

// OrderService xử lý logic đặt hàng
type OrderService struct {
	vnpay *VNPayService
}

// ✅ Hàm khởi tạo `OrderService`
func NewOrderService(vnpay *VNPayService) *OrderService {
	return &OrderService{vnpay: vnpay}
}	

func (s *OrderService) PlaceOrder(orderRequest request.OrderRequest) (*request.OrderRequest, string, error) {
	db, err := database.DB1Connection()
	if err != nil {
		return nil, "", err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// 🔹 Truy vấn giỏ hàng dựa trên `orderRequest.UserID`
	var cartItems []CartDB
	err = db.Where("user_id = ?", orderRequest.UserID).Find(&cartItems).Error
	if err != nil {
		return nil, "", err
	}
	if len(cartItems) == 0 {
		return nil, "", fmt.Errorf("Giỏ hàng trống")
	}

	// 🔹 Tính tổng tiền từ giỏ hàng
	totalAmount := 0
	for _, item := range cartItems {
		totalAmount += int(item.Price) * item.Quantity
	}

	// 🔹 Tạo đơn hàng mới
	order := request.OrderRequest{
		UserID:    orderRequest.UserID,
		OrderDate: time.Now(),
		Status:    "Đang xử lý",
	}
	err = db.Create(&order).Error
	if err != nil {
		return nil, "", err
	}

	// 🔹 Chuyển dữ liệu từ giỏ hàng → OrderDetails
	var orderDetails []request.OrderDetailsRequest
	for _, cart := range cartItems {
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
	}

	// 🔹 Lưu danh sách OrderDetails vào database
	err = db.Create(&orderDetails).Error
	if err != nil {
		return nil, "", err
	}

	// 🔹 Tạo URL thanh toán VNPay
	paymentURL, err := s.vnpay.GenerateVNPayURL(fmt.Sprintf("%d", order.ID), totalAmount)
	if err != nil {
		return nil, "", err
	}

	// 🔹 Xóa giỏ hàng sau khi tạo đơn hàng
	err = db.Where("user_id = ?", orderRequest.UserID).Delete(&CartDB{}).Error
	if err != nil {
		return nil, "", fmt.Errorf("Không thể xóa giỏ hàng sau khi đặt hàng")
	}

	// 🔹 Gán danh sách orderDetails vào order và trả về kết quả
	order.OrderDetails = orderDetails
	return &order, paymentURL, nil
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
