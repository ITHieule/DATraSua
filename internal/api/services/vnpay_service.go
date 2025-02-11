package services

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
	"web-api/internal/pkg/config"
)

type VNPayService struct {
	Config *config.VNPayConfig
}

func NewVNPayService(cfg *config.VNPayConfig) *VNPayService {
	return &VNPayService{Config: cfg}
}

// GenerateVNPayURL tạo URL thanh toán VNPay với chữ ký chính xác
func (v *VNPayService) GenerateVNPayURL(orderID string, amount int) (string, error) {
	log.Println("[VNPay] 🔹 Bắt đầu tạo URL thanh toán...")
	// 🔹 Kiểm tra nếu Config bị nil
	if v.Config == nil {
		log.Fatal(" VNPayService.Config chưa được khởi tạo!")
		return "", fmt.Errorf("VNPayService chưa được khởi tạo")
	}
	vnpParams := map[string]string{
		"vnp_Version":    "2.1.0",
		"vnp_Command":    "pay",
		"vnp_TmnCode":    v.Config.TmnCode,
		"vnp_Amount":     fmt.Sprintf("%d", amount*100),
		"vnp_CurrCode":   "VND",
		"vnp_TxnRef":     orderID,
		"vnp_OrderInfo":  fmt.Sprintf("Thanh toan don hang %s", orderID),
		"vnp_OrderType":  "billpayment",
		"vnp_Locale":     "vn",
		"vnp_ReturnUrl":  v.Config.ReturnURL,
		"vnp_IpAddr":     "127.0.0.1",
		"vnp_CreateDate": time.Now().Format("20060102150405"),
	}

	log.Printf("[VNPay] Thông tin đầu vào: %+v\n", vnpParams)

	// ✅ Encode từng tham số trước khi hash
	hashData := v.createHash(vnpParams)

	// ✅ Tạo query string đúng chuẩn
	queryString := v.createQueryString(vnpParams)

	// ✅ Tạo URL thanh toán chính xác
	paymentURL := fmt.Sprintf("%s?%s&vnp_SecureHash=%s", v.Config.BaseURL, queryString, hashData)

	log.Printf("[VNPay] ✅ URL thanh toán: %s\n", paymentURL)

	return paymentURL, nil
}

// createQueryString tạo query string đúng chuẩn
func (v *VNPayService) createQueryString(params map[string]string) string {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return values.Encode()
}

// createHash tạo mã checksum HMAC SHA512 chính xác
func (v *VNPayService) createHash(params map[string]string) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var rawData []string
	for _, k := range keys {
		if k != "vnp_SecureHash" && k != "vnp_SecureHashType" {
			// 🔥 Encode từng tham số trước khi hash
			rawData = append(rawData, fmt.Sprintf("%s=%s", k, url.QueryEscape(params[k])))
		}
	}

	hashString := strings.Join(rawData, "&")
	log.Printf("[VNPay] 🔹 Chuỗi dữ liệu trước khi hash: %s\n", hashString)

	h := hmac.New(sha512.New, []byte(v.Config.HashSecret))
	h.Write([]byte(hashString))
	hashResult := hex.EncodeToString(h.Sum(nil))

	log.Printf("[VNPay] ✅ Hash SHA512: %s\n", hashResult)
	return hashResult
}

func (v *VNPayService) HandleVNPayCallback(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	// Lấy order ID từ tham số callback (VNPay gửi `vnp_TxnRef`)
	orderID := params.Get("vnp_TxnRef")
	if orderID == "" {
		http.Error(w, "Thiếu order ID", http.StatusBadRequest)
		return
	}

	// Xác thực chữ ký VNPay
	if !v.VerifyVNPaySignature(params) {
		http.Error(w, "Chữ ký không hợp lệ", http.StatusForbidden)
		return
	}

	// Kiểm tra trạng thái giao dịch
	paymentStatus := params.Get("vnp_ResponseCode")
	var orderStatus string
	if paymentStatus == "00" {
		orderStatus = "Thành công"
	} else {
		orderStatus = "Thất bại"
	}

	log.Printf("[VNPay] ✅ Nhận callback cho đơn hàng %s, trạng thái: %s", orderID, orderStatus)

	// Phản hồi cho VNPay biết đã xử lý xong
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Callback xử lý thành công")
}

func (v *VNPayService) VerifyVNPaySignature(params url.Values) bool {
	// Lấy secret key từ config
	secretKey := v.Config.HashSecret
	if secretKey == "" {
		log.Println("[VNPay] ❌ Chưa cấu hình secret key")
		return false
	}

	// Lấy chữ ký từ request
	receivedSignature := params.Get("vnp_SecureHash")
	if receivedSignature == "" {
		log.Println("[VNPay] ❌ Thiếu vnp_SecureHash")
		return false
	}

	// Sắp xếp tham số theo thứ tự key (bỏ qua vnp_SecureHash)
	var keys []string
	for k := range params {
		if k != "vnp_SecureHash" && k != "vnp_SecureHashType" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// Tạo chuỗi dữ liệu để hash
	var rawData []string
	for _, k := range keys {
		// ✅ Sử dụng QueryEscape để đảm bảo encoding giống lúc tạo chữ ký ban đầu
		rawData = append(rawData, fmt.Sprintf("%s=%s", k, url.QueryEscape(params.Get(k))))
	}
	hashString := strings.Join(rawData, "&")

	// ✅ Log để kiểm tra chuỗi hash gửi đến VNPay
	log.Printf("[VNPay] 🔹 Chuỗi dữ liệu cần hash (từ callback): %s\n", hashString)

	// Tạo hash HMAC SHA512
	h := hmac.New(sha512.New, []byte(secretKey))
	h.Write([]byte(hashString))
	expectedHash := hex.EncodeToString(h.Sum(nil))

	// ✅ Log để kiểm tra hash server tạo ra
	log.Printf("[VNPay] 🔹 Hash SHA512 mong đợi: %s\n", expectedHash)

	// Kiểm tra chữ ký
	if strings.EqualFold(expectedHash, receivedSignature) {
		log.Println("[VNPay] ✅ Chữ ký hợp lệ")
		return true
	} else {
		log.Println("[VNPay] ❌ Chữ ký không hợp lệ")
		log.Printf("[VNPay] ❌ Hash nhận được từ VNPay: %s\n", receivedSignature)
		return false
	}
}
