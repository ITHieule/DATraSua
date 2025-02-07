package services

import (
	"fmt"
	"web-api/internal/api/until"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	*BaseService
}

var User = &UserService{}

func (s *UserService) Register(requestParams *request.User) ([]types.Usertypes, error) {
	var User []types.Usertypes
	// kết nối cơ sở dữ liệu
	db, err := database.DB1Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	// Mã hóa mật khẩu
	hashedPassword, err := HashPassword(requestParams.Password_hash)
	if err != nil {

		return nil, err
	}
	requestParams.Password_hash = hashedPassword
	// Mặc định Is_verified = 1 nếu không có giá trị
	if requestParams.Is_verified == false {
		requestParams.Is_verified = true
	}
	// Kiểm tra dữ liệu đầu vào
	dbInstance, _ := db.DB()
	defer dbInstance.Close()
	query := "INSERT INTO Users (Username, Password_hash, Email,Phone,Role ,Is_verified) VALUES (?, ?, ?, ?, ?, ?)"
	err = db.Raw(query,
		requestParams.Username,
		requestParams.Password_hash,
		requestParams.Email,
		requestParams.Phone,
		requestParams.Role,
		requestParams.Is_verified,
	).Scan(&User).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
	}

	return User, nil
}

func (s *UserService) Login(requestParams *request.User) (string, error) {
	var user types.Usertypes

	// Kết nối cơ sở dữ liệu
	db, err := database.DB1Connection()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return "", err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Truy vấn thông tin người dùng dựa trên email
	query := "SELECT * FROM Users WHERE Username = ?"
	err = db.Raw(query, requestParams.Username).Scan(&user).Error

	if err != nil {
		fmt.Println("Query execution error:", err)
		return "", err
	}

	// So sánh mật khẩu đã mã hóa với mật khẩu người dùng nhập vào
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password_hash), []byte(requestParams.Password_hash)); err != nil {
		return "", nil
	}

	// Tạo JWT token
	token, err := until.GenerateJWT(user.Id, user.Role)
	if err != nil {
		fmt.Println("Error generating token:", err)
		return token, nil
	}

	// Trả về thông tin người dùng và token
	return token, nil

}

// Hàm mã hóa mật khẩu
func HashPassword(Password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
