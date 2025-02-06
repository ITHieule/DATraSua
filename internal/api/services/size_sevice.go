package services

import (
	"fmt"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"
)

type SizeService struct {
	*BaseService
}

var OrderSize = &SizeService{}

func (s *SizeService) SizesSevice() ([]types.Sizestypes, error) {
	var orders []types.Sizestypes

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
		SELECT * FROM OrderSystem.Sizes 

	`

	// Thực hiện truy vấn và ánh xạ kết quả vào struct
	err = db.Raw(query).Scan(&orders).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return orders, nil
}

func (s *SizeService) SizesesSevice(requestParams *request.Sizesrequest) ([]types.Sizestypes, error) {
	var Sizes []types.Sizestypes

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
		 INSERT INTO Sizes (
          name,price
        ) VALUES (?, ?)

	`

	// Truyền tham số vào câu truy vấn
	err = db.Exec(query,
		requestParams.Name,
		requestParams.Price,
	).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return Sizes, nil
}
func (s *SizeService) UpdatSizesesSevice(requestParams *request.Sizesrequest) ([]types.Sizestypes, error) {
	var Sizes []types.Sizestypes

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
    UPDATE Sizes
    SET name = ?, price = ?
    WHERE id = ?
`
	err = db.Exec(query,
		requestParams.Name,
		requestParams.Price,
		requestParams.Id,
	).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return Sizes, nil
}

func (s *SizeService) DeleteSizesesSevice(Id int) error {

	// Kết nối database
	db, err := database.DB1Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)

		return err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Thực thi lệnh DELETE
	result := db.Exec("DELETE FROM Sizes WHERE id = ?", Id)

	// Kiểm tra lỗi truy vấn
	if result.Error != nil {
		fmt.Println("Query execution error:", result.Error)
		return result.Error
	}

	// Kiểm tra số dòng bị ảnh hưởng (nếu ID không tồn tại, sẽ không xóa được)
	if result.RowsAffected == 0 {
		fmt.Println("No book found with ID:", Id)
		return fmt.Errorf("không tìm thấy sách với ID %d", Id)
	}

	fmt.Println("Deleted book successfully!")
	return nil
}

func (s *SizeService) SearchSizesesSevice(requestParams *request.Sizesrequest) ([]types.Sizestypes, error) {
	var Sizes []types.Sizestypes

	// Kết nối database
	db, err := database.DB1Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)

		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Truy vấn SQL lấy ngày đặt hàng và tổng số lượng sách đã bán

	err = db.Raw("SELECT * FROM Sizes WHERE name = ? OR price = ? OR id = ?",
		requestParams.Name,
		requestParams.Price,
		requestParams.Id,
	).Scan(&Sizes).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return Sizes, nil
}
