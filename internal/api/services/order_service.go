package services

import (
	"fmt"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"
)

type OrderService struct {
	*BaseService
}

var Orders = &OrderService{}

func (s *OrderService) Create(requestParams *request.OrderRequest) ([]types.OrderTypes, error) {
	var Order []types.OrderTypes
	// kết nối cơ sở dữ liệu
	db, err := database.DB1Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	// Kiểm tra dữ liệu đầu vào
	dbInstance, _ := db.DB()
	defer dbInstance.Close()
	query := "INSERT INTO Orders (user_id,  order_date, status) VALUES (?, ?, ?)"
	err = db.Raw(query,
		requestParams.UserID,
		requestParams.OrderDate,
		requestParams.Status,
	).Scan(&Order).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
	}

	return Order, nil
}
