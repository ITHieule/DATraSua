package services

import (
	"fmt"
	"strings"
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

// Hàm tạo đơn hàng và chi tiết đơn hàng
func (s *OrderDetailsService) Create(requestParams *request.OrderRequest) ([]types.OrderDetailsTypes, error) {
	var orderDetails []types.OrderDetailsTypes

	// Kết nối cơ sở dữ liệu
	db, err := database.DB1Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	// Đảm bảo đóng kết nối sau khi hoàn thành
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Câu lệnh INSERT để tạo đơn hàng
	orderQuery := "INSERT INTO Orders (user_id, order_date, status) VALUES (?, ?, ?)"
	res := db.Exec(orderQuery, requestParams.UserID, requestParams.OrderDate, requestParams.Status)
	if res.Error != nil {
		fmt.Println("Error creating order:", res.Error)
		return nil, res.Error
	}

	// Lấy giá trị order_id tự sinh sau khi thực hiện INSERT
	var orderID int
	err = db.Raw("SELECT LAST_INSERT_ID()").Scan(&orderID).Error
	if err != nil {
		fmt.Println("Error fetching last insert ID:", err)
		return nil, err
	}

	// Gán order_id cho requestParams
	requestParams.ID = uint(orderID)

	// Kiểm tra lại xem order_id đã được sinh chưa
	if requestParams.ID == 0 {
		return nil, fmt.Errorf("Failed to create order, order_id is missing")
	}

	// Duyệt qua từng OrderDetails trong mảng và tạo từng chi tiết
	for _, detail := range requestParams.OrderDetails {
		// Tính giá cho từng chi tiết đơn hàng
		price, err := s.CalculatePrice(detail.Base_id, detail.Size_id)
		if err != nil {
			return nil, fmt.Errorf("Error calculating price for detail: %v", err)
		}

		// Gán giá tính được vào detail
		detail.Price = price
		detail.Order_id = int(requestParams.ID) // Gán order_id đã tự sinh

		// Câu lệnh INSERT vào bảng OrderDetails
		query := "INSERT INTO OrderDetails (order_id, base_id, flavor_id, sweetness_id, ice_id, size_id, price) VALUES (?, ?, ?, ?, ?, ?, ?)"
		err = db.Raw(query,
			detail.Order_id,
			detail.Base_id,
			detail.Flavor_id,
			detail.Sweetness_id,
			detail.Ice_id,
			detail.Size_id,
			detail.Price,
		).Scan(&orderDetails).Error
		if err != nil {
			fmt.Println("Error inserting order details:", err)
			return nil, err
		}
	}

	// Truy vấn lại để lấy thông tin chi tiết đơn hàng đã được thêm vào
	err = db.Raw("SELECT * FROM OrderDetails WHERE order_id = ?", requestParams.ID).Scan(&orderDetails).Error
	if err != nil {
		fmt.Println("Error fetching created order details:", err)
		return nil, err
	}

	var user types.Usertypes
	err = db.Raw("SELECT * FROM Users WHERE id = ?", requestParams.UserID).Scan(&user).Error
	if err != nil {
		fmt.Println("Error fetching created order details:", err)
		return nil, err
	}

	var order []types.OrdersTypes
	err = db.Raw(`
	
	SELECT OrderDetails.order_id,Bases.name as 'BasesName', Flavors.name  as 'FlavorsName',Sweetness.name  as 'SweetnessName',IceLevels.name  as 'IceLevelsName',Sizes.name  as 'SizeName'  FROM OrderSystem.OrderDetails
	join Bases on OrderDetails.base_id = Bases.id
	join Flavors on OrderDetails.flavor_id = Flavors.id
	join Sweetness on OrderDetails.sweetness_id = Sweetness.id
	join IceLevels on OrderDetails.ice_id = IceLevels.id
	join Sizes on OrderDetails.size_id = Sizes.id
	where OrderDetails.order_id = ?
	`, requestParams.ID).Scan(&order).Error
	if err != nil {
		fmt.Println("Error fetching created order details:", err)
		return nil, err
	}

	body := formatOrderEmailContent(requestParams.ID, order)

	SendEmail(user.Email, "Đơn hàng", body)

	// Trả về chi tiết đơn hàng vừa tạo
	return orderDetails, nil
}

func (s *OrderDetailsService) GetOrderWithDetails(orderID int) (types.OrderWithDetails, error) {
	var order types.OrderTypes
	var orderDetails []types.OrderDetailsTypes

	// Kết nối cơ sở dữ liệu
	db, err := database.DB1Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return types.OrderWithDetails{}, err
	}
	// Đảm bảo đóng kết nối sau khi hoàn thành
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Truy vấn thông tin đơn hàng từ bảng Orders
	orderQuery := "SELECT * FROM Orders WHERE id = ?"
	err = db.Raw(orderQuery, orderID).Scan(&order).Error
	if err != nil {
		fmt.Println("Error fetching order:", err)
		return types.OrderWithDetails{}, err
	}

	// Nếu không tìm thấy đơn hàng
	if (order == types.OrderTypes{}) {
		return types.OrderWithDetails{}, fmt.Errorf("Order with id %d not found", orderID)
	}

	// Truy vấn chi tiết đơn hàng từ bảng OrderDetails theo order_id
	orderDetailsQuery := "SELECT * FROM OrderDetails WHERE order_id = ?"
	err = db.Raw(orderDetailsQuery, orderID).Scan(&orderDetails).Error
	if err != nil {
		fmt.Println("Error fetching order details:", err)
		return types.OrderWithDetails{}, err
	}

	// Trả về một object chứa cả thông tin đơn hàng và chi tiết đơn hàng
	return types.OrderWithDetails{
		Order:        order,
		OrderDetails: orderDetails,
	}, nil
}

// Hàm hủy đơn hàng (set trạng thái thành "Đã hủy")
func (s *OrderService) CancelOrder(orderID int) error {
	// Kết nối cơ sở dữ liệu
	db, err := database.DB1Connection()
	if err != nil {
		return fmt.Errorf("Database connection error: %v", err)
	}
	// Đảm bảo đóng kết nối sau khi hoàn thành
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Cập nhật trạng thái đơn hàng thành "Đã hủy"
	query := "UPDATE Orders SET status = ? WHERE id = ?"
	res := db.Exec(query, "Đã hủy", orderID)
	if res.Error != nil {
		return fmt.Errorf("Error updating order status: %v", res.Error)
	}

	// Kiểm tra nếu không có đơn hàng nào bị ảnh hưởng
	if res.RowsAffected == 0 {
		return fmt.Errorf("Order with id %d not found", orderID)
	}

	return nil
}
func formatOrderEmailContent(orderID uint, details []types.OrdersTypes) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Chào bạn,\n\nĐơn hàng #%d của bạn đã được đặt thành công.\n\n", orderID))
	sb.WriteString("📌 *Chi tiết đơn hàng:*\n\n")

	for _, item := range details {
		sb.WriteString(fmt.Sprintf(
			"- 🍵 Tên món: %s\n  🌿 Hương vị: %s\n  🍯 Độ ngọt: %s\n  ❄️ Mức đá: %s\n  📏 Kích thước: %s\n\n",
			item.BasesName,     // Ví dụ: "Trà Đào"
			item.FlavorsName,   // Ví dụ: "Vani"
			item.SweetnessName, // Ví dụ: "Ít ngọt"
			item.IceLevelsName, // Ví dụ: "Ít đá"
			item.SizeName,      // Ví dụ: "Vừa"
		))
	}

	sb.WriteString("🎉 Cảm ơn bạn đã mua hàng!\n")
	return sb.String()
}
