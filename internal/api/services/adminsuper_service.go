package services

import (
	"fmt"
	"os"
	"time"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// FILE - ADMIN SERVICE
type AdminSuperService struct {
	*BaseService
}

var AdminSuper = &AdminSuperService{}

// func - LoginService

func (s *AdminSuperService) LoginAdminSuperService(requestParams *request.AdminSuper) (string, error) {
	var admin types.AdminSuper

	// Kết nối cơ sở dữ liệu
	db, err := database.DB1Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return "", err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Truy vấn lấy thông tin admin
	query := "SELECT id, username, password_hash, role FROM Users WHERE username = ? AND role = 'ADMIN'"
	err = db.Raw(query, requestParams.Username).Scan(&admin).Error
	if err != nil {
		fmt.Println("Query error:", err)
		return "", err
	}

	// So sánh mật khẩu
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password_hash), []byte(requestParams.Password_hash)); err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	// Tạo JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"adminid":  admin.Adminid,
		"username": admin.Username,
		"role":     admin.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token hết hạn sau 24 giờ
	})

	// Ký token
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		fmt.Println("Token signing error:", err)
		return "", err
	}

	return tokenString, nil
}

func (s *AdminSuperService) GetUsersSevice() ([]types.Usertypes, error) {
	var orders []types.Usertypes

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
		SELECT * FROM OrderSystem.Users 

	`

	// Thực hiện truy vấn và ánh xạ kết quả vào struct
	err = db.Raw(query).Scan(&orders).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return orders, nil
}

func (s *AdminSuperService) UpdateAdmidsuperSevice(requestParams *request.AdminSuper) ([]types.AdminSuper, error) {
	var Admin []types.AdminSuper

	// Kết nối database
	db, err := database.DB1Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Mã hóa mật khẩu nếu có thay đổi
	var hashedPassword string
	if requestParams.Password_hash != "" { // Kiểm tra nếu mật khẩu được cung cấp
		hashedPassword, err = HashPassword(requestParams.Password_hash)
		if err != nil {
			fmt.Println("Password hashing error:", err)
			return nil, err
		}
	}

	// Truy vấn SQL cập nhật thông tin người dùng
	query := `
        UPDATE Users
        SET password_hash = COALESCE(?, password_hash),
            email = ?,
            phone = ?,
			role = ?,
			is_verified = ?
        WHERE id = ?
    `

	err = db.Exec(query,
		hashedPassword, // Mật khẩu đã mã hóa hoặc NULL để giữ nguyên
		requestParams.Email,
		requestParams.Phone,
		requestParams.Role,
		requestParams.Is_verified,
		requestParams.Id,
	).Error

	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}

	return Admin, nil
}

func (s *AdminSuperService) DeleteAdmidsuperSevice(Id int) error {
	

		// Kết nối database
		db, err := database.DB1Connection()
		if err != nil {
			fmt.Println("Database connection error:", err)
	
			return err
		}
		dbInstance, _ := db.DB()
		defer dbInstance.Close()
	
		// Thực thi lệnh DELETE
		result := db.Exec("DELETE FROM Users WHERE id = ?", Id)
	
		// Kiểm tra lỗi truy vấn
		if result.Error != nil {
			fmt.Println("Query execution error:", result.Error)
			return result.Error
		}
	
		// Kiểm tra số dòng bị ảnh hưởng (nếu ID không tồn tại, sẽ không xóa được)
		if result.RowsAffected == 0 {
			fmt.Println("No Users found with ID:", Id)
			return fmt.Errorf("không tìm thấy Users với ID %d", Id)
		}
	
		fmt.Println("Deleted book successfully!")
		return nil
}
