package services

import (
	"github.com/go-gomail/gomail"
	"github.com/spf13/viper"
	"log"
	"strconv"
)

// Hàm gửi email
func SendEmail(to string, subject string, body string) error {

	viper.SetConfigName("config")   // Tên file (không bao gồm đuôi .yaml)
	viper.SetConfigType("yaml")     // Loại file
	viper.AddConfigPath("././data") // Đường dẫn file (thư mục hiện tại)

	// Đọc file cấu hình
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Lấy các thông tin cấu hình từ biến môi trường
	emailHost := viper.GetString("email.host")
	emailPort := viper.GetString("email.port")
	emailUsername := viper.GetString("email.username")
	emailPassword := viper.GetString("email.password")

	// Cấu hình email
	mail := gomail.NewMessage()
	mail.SetHeader("From", emailUsername)
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/plain", body)

	// Cấu hình SMTP server
	portNumber, _ := strconv.Atoi(emailPort) // Chuyển port sang kiểu số
	dialer := gomail.NewDialer(emailHost, portNumber, emailUsername, emailPassword)

	// Gửi email
	err := dialer.DialAndSend(mail)
	if err != nil {
		return err
	}
	return nil
}
