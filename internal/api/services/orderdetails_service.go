package services

import (
	"fmt"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"
)

type OrderDetailsService struct {
	*BaseService
}

var OrderDetails = &OrderDetailsService{}

// Hàm tính giá dựa trên base_id và size_id
func (s *OrderDetailsService) CalculatePrice(baseID, sizeID int) (float64, error) {
	var basePrice, sizePrice float64

	// Kết nối cơ sở dữ liệu
	db, err := database.DB1Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return 0, err
	}
	// Đảm bảo đóng kết nối sau khi hoàn thành
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Truy vấn giá từ bảng Bases
	queryBase := "SELECT price FROM Bases WHERE id = ?"
	err = db.Raw(queryBase, baseID).Scan(&basePrice).Error
	if err != nil {
		fmt.Println("Error fetching price from Bases:", err)
		return 0, err
	}

	// Truy vấn giá từ bảng Sizes
	querySize := "SELECT price FROM Sizes WHERE id = ?"
	err = db.Raw(querySize, sizeID).Scan(&sizePrice).Error
	if err != nil {
		fmt.Println("Error fetching price from Sizes:", err)
		return 0, err
	}

	// Cộng giá từ Base và Size để tính tổng
	totalPrice := basePrice + sizePrice

	// Trả về tổng giá
	return totalPrice, nil
}

// Hàm tạo chi tiết đơn hàng
func (s *OrderDetailsService) Create(requestParams *request.OrderDetailsRequest) ([]types.OrderDetailsTypes, error) {
	var orderDetails []types.OrderDetailsTypes

	// Tính giá từ base_id và size_id
	price, err := s.CalculatePrice(requestParams.Base_id, requestParams.Size_id)
	if err != nil {
		return nil, fmt.Errorf("Error calculating price: %v", err)
	}

	// Gán giá tính được vào requestParams.Price
	requestParams.Price = price

	// Kết nối cơ sở dữ liệu
	db, err := database.DB1Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	// Đảm bảo đóng kết nối sau khi hoàn thành
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Câu lệnh INSERT vào bảng OrderDetails
	query := "INSERT INTO OrderDetails (order_id, base_id, flavor_id, sweetness_id, ice_id, size_id, price) VALUES (?, ?, ?, ?, ?, ?, ?)"
	err = db.Raw(query,
		requestParams.Order_id,
		requestParams.Base_id,
		requestParams.Flavor_id,
		requestParams.Sweetness_id,
		requestParams.Ice_id,
		requestParams.Size_id,
		requestParams.Price,
	).Scan(&orderDetails).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}

	// Truy vấn lại để lấy thông tin chi tiết đơn hàng đã được thêm vào
	err = db.Raw("SELECT * FROM OrderDetails WHERE order_id = ?", requestParams.Order_id).Scan(&orderDetails).Error
	if err != nil {
		fmt.Println("Error fetching created order details:", err)
		return nil, err
	}

	// Trả về chi tiết đơn hàng vừa tạo
	return orderDetails, nil
}
