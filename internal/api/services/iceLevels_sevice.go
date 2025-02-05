package services

import (
	"fmt"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/types"
)

type IceLevelsService struct {
	*BaseService
}

var OrderIceLevels = &IceLevelsService{}

func (s *IceLevelsService) IceLevelsSevice() ([]types.IceLevels, error) {
	var orders []types.IceLevels

	// Kết nối database
	db, err := database.DB1Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)

		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Truy vấn SQL lấy ngày đặt hàng và tổng số lượng sách đã bán
	query := `
		SELECT * FROM OrderSystem.IceLevels 

	`

	// Thực hiện truy vấn và ánh xạ kết quả vào struct
	err = db.Raw(query).Scan(&orders).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return orders, nil
}
