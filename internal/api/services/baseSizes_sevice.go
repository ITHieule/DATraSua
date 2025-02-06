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

func (s *BaseSizesService) GetBaseSizesSevice() ([]types.BaseSizestypes, error) {
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
	SELECT * FROM OrderSystem.BaseSizes 

`

	// Thực hiện truy vấn và ánh xạ kết quả vào struct
	err = db.Raw(query).Scan(&BaseSizes).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return BaseSizes, nil
}
func (s *BaseSizesService) UpdateBaseSizesSevice(requestParams *request.BaseSizesrequest) ([]types.BaseSizestypes, error) {
	var Sizes []types.BaseSizestypes

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
    UPDATE BaseSizes
    SET base_id = ?, size_id = ?
    WHERE id = ?
`
	err = db.Exec(query,
		requestParams.Base_id,
		requestParams.Size_id,
		requestParams.Id,
	).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return Sizes, nil
}

func (s *BaseSizesService) DeleteBaseSizesSevice(Id int) error {

	// Kết nối database
	db, err := database.DB1Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)

		return err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Thực thi lệnh DELETE
	result := db.Exec("DELETE FROM BaseSizes WHERE id = ?", Id)

	// Kiểm tra lỗi truy vấn
	if result.Error != nil {
		fmt.Println("Query execution error:", result.Error)
		return result.Error
	}

	// Kiểm tra số dòng bị ảnh hưởng (nếu ID không tồn tại, sẽ không xóa được)
	if result.RowsAffected == 0 {
		fmt.Println("No BaseSizes found with ID:", Id)
		return fmt.Errorf("không tìm thấy BaseSizes với ID %d", Id)
	}

	fmt.Println("Deleted BaseSizes successfully!")
	return nil
}

func (s *BaseSizesService) SearchBaseSizes(requestParams *request.BaseSizesrequest) ([]types.BaseSizestypes, error) {
	var Sizes []types.BaseSizestypes

	// Kết nối database
	db, err := database.DB1Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)

		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Truy vấn SQL lấy ngày đặt hàng và tổng số lượng sách đã bán

	err = db.Raw("SELECT * FROM BaseSizes WHERE base_id = ? OR size_id = ? OR id = ?",
		requestParams.Base_id,
		requestParams.Size_id,
		requestParams.Id,
	).Scan(&Sizes).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return Sizes, nil
}
