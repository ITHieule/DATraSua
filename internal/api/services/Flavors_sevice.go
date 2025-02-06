package services

import (
	"fmt"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"
)

type FlavorsService struct {
	*BaseService
}

var Orderflavors = &FlavorsService{}

func (s *FlavorsService) FlavorsSevice() ([]types.Flavorstypes, error) {
	var orders []types.Flavorstypes

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
		SELECT * FROM OrderSystem.Flavors 

	`

	// Thực hiện truy vấn và ánh xạ kết quả vào struct
	err = db.Raw(query).Scan(&orders).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return orders, nil
}

func (s *FlavorsService) AddFlavorsSevice(requestParams *request.Flavorsrequest) ([]types.Flavorstypes, error) {
	var orders []types.Flavorstypes

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
		 INSERT INTO Flavors (
          name
        ) VALUES (?)

	`

	// Truyền tham số vào câu truy vấn
	err = db.Exec(query,
		requestParams.Name,
	).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return orders, nil
}

func (s *FlavorsService) UpdateFlavorsSevice(requestParams *request.Flavorsrequest) ([]types.Flavorstypes, error) {
	var orders []types.Flavorstypes

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
    UPDATE Flavors
    SET name = ?
    WHERE id = ?
`
	err = db.Exec(query,
		requestParams.Name,

		requestParams.Id,
	).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return orders, nil
}
func (s *FlavorsService) DeleteflavorsSevice(Id int) error {

	// Kết nối database
	db, err := database.DB1Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)

		return err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Thực thi lệnh DELETE
	result := db.Exec("DELETE FROM Flavors WHERE id = ?", Id)

	// Kiểm tra lỗi truy vấn
	if result.Error != nil {
		fmt.Println("Query execution error:", result.Error)
		return result.Error
	}

	// Kiểm tra số dòng bị ảnh hưởng (nếu ID không tồn tại, sẽ không xóa được)
	if result.RowsAffected == 0 {
		fmt.Println("No Flavors found with ID:", Id)
		return fmt.Errorf("không tìm thấy Flavors với ID %d", Id)
	}

	fmt.Println("Deleted book successfully!")
	return nil
}
