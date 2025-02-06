package services

import (
	"fmt"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"
)

type BaseSizesService struct {
	*BaseService
}

var OrderBaseSizes = &BaseSizesService{}

func (s *BaseSizesService) BaseSizesSevice(requestParams *request.BaseSizesrequest) ([]types.BaseSizestypes, error) {
	var BaseSizes []types.BaseSizestypes

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
		 INSERT INTO BaseSizes (
          base_id,size_id
        ) VALUES (?, ?)

	`

	// Truyền tham số vào câu truy vấn
	err = db.Exec(query,
		requestParams.Base_id,
		requestParams.Size_id,
	).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return BaseSizes, nil
}
