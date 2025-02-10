package services

import (
	"fmt"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"
)

type StatisticalService struct {
	*BaseService
}

var Statis = &StatisticalService{}

// GetOrderStatistics thống kê số lượng sách đã bán
func (s *StatisticalService) GetStatisticalSevice(requestParams *request.Statisticalrequest) ([]types.Statisticaltypes, error) {
	var Statistical []types.Statisticaltypes

	// Kết nối database
	db, err := database.DB1Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}

	// Lấy instance của database để đóng khi xong
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Truy vấn SQL lấy ngày đặt hàng và tổng số lượng sách đã bán
	query := `
		 SELECT Bases.name, 
       COUNT(*) AS Quantity, 
       SUM(OrderDetails.price) AS Price
FROM OrderSystem.OrderDetails
JOIN Orders ON OrderDetails.order_id = Orders.id 
JOIN Bases ON OrderDetails.base_id = Bases.id
WHERE Orders.status = 'Đã xác nhận'  and
 Bases.name = ?
      AND YEAR(Orders.order_date) = ?  -- Lọc theo năm 2024
   AND Month(Orders.order_date) = ?
GROUP BY Bases.name;
	`

	// Thực hiện truy vấn với tham số year và month
	err = db.Raw(query, requestParams.Name, requestParams.Year_Order, requestParams.Month_Order).Scan(&Statistical).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return Statistical, nil
}
