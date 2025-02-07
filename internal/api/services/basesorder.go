package services

import (
	"fmt"
	"mime/multipart"
	"os"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"

	"github.com/gin-gonic/gin"
)

type BasesService struct {
	*BaseService
}

var Order = &BasesService{}

func (s *BasesService) BasesSevice() ([]types.Basestypes, error) {
	var orders []types.Basestypes

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
		SELECT * FROM OrderSystem.Bases 

	`

	// Thực hiện truy vấn và ánh xạ kết quả vào struct
	err = db.Raw(query).Scan(&orders).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return orders, nil
}

func (s *BasesService) AddbasesSevice(requestParams *request.Basesrequest) ([]types.Basestypes, error) {
	var Sizes []types.Basestypes

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
		 INSERT INTO Bases (
          name,price,images
        ) VALUES (?, ?, ?)

	`

	// Truyền tham số vào câu truy vấn
	err = db.Exec(query,
		requestParams.Name,
		requestParams.Price,
		requestParams.Images,
	).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}

	return Sizes, nil
}
func (s *BasesService) UpdatebasesSevice(requestParams *request.Basesrequest) ([]types.Basestypes, error) {
	var Sizes []types.Basestypes

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
    UPDATE Bases
    SET name = ?, price = ?
    WHERE id = ?
`
	err = db.Exec(query,
		requestParams.Name,
		requestParams.Price,
		requestParams.Images,
		requestParams.Id,
	).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return Sizes, nil
}

func (s *BasesService) DeletebasesSevice(Id int) error {

	// Kết nối database
	db, err := database.DB1Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)

		return err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Thực thi lệnh DELETE
	result := db.Exec("DELETE FROM Bases WHERE id = ?", Id)

	// Kiểm tra lỗi truy vấn
	if result.Error != nil {
		fmt.Println("Query execution error:", result.Error)
		return result.Error
	}

	// Kiểm tra số dòng bị ảnh hưởng (nếu ID không tồn tại, sẽ không xóa được)
	if result.RowsAffected == 0 {
		fmt.Println("No Bases found with ID:", Id)
		return fmt.Errorf("không tìm thấy Bases với ID %d", Id)
	}

	fmt.Println("Deleted Bases successfully!")
	return nil
}

func (s *BasesService) SearchbasesSevice(requestParams *request.Basesrequest) ([]types.Basestypes, error) {
	var Sizes []types.Basestypes

	// Kết nối database
	db, err := database.DB1Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)

		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Truy vấn SQL lấy ngày đặt hàng và tổng số lượng sách đã bán

	err = db.Raw("SELECT * FROM Bases WHERE name = ? OR price = ? OR images = ? OR id = ?",
		requestParams.Name,
		requestParams.Price,
		requestParams.Images,
		requestParams.Id,
	).Scan(&Sizes).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return Sizes, nil
}

func (s *BasesService) SaveImage(path string, file *multipart.FileHeader, ctx *gin.Context) (string, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := ctx.SaveUploadedFile(file, path); err != nil {
			return "Không lưu được hình ảnh", err
		}
	} else if err != nil {
		return err.Error(), err
	} else {
		return "File đã tồn tại", nil
	}
	return "OK", nil
}
