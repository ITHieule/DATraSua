package config

import (
	"os"
)

type VNPayConfig struct {
	TmnCode    string
	HashSecret string
	BaseURL    string
	ReturnURL  string
}

func LoadVNPayConfig() *VNPayConfig {
	return &VNPayConfig{
		TmnCode:    getEnv("VNPAY_TMN_CODE", "CE7KSU2X"), // Giá trị mặc định
		HashSecret: getEnv("VNPAY_HASH_SECRET", "46WE9BSARW85G3D705FP5TYXK270P7TX"),
		BaseURL:    "https://sandbox.vnpayment.vn/paymentv2/vpcpay.html",
		ReturnURL:  "http://localhost:8081/api/v1/System/callback",
	}
}

// getEnv lấy giá trị biến môi trường hoặc dùng giá trị mặc định
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
