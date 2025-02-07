package controllers

import (
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/config"

	"github.com/gin-gonic/gin"
)

type VNPayController struct {
	Service *services.VNPayService
}

func NewVNPayController() *VNPayController {
	cfg := config.LoadVNPayConfig()          // ✅ Load cấu hình
	service := services.NewVNPayService(cfg) // ✅ Tạo service
	return &VNPayController{Service: service}
}

// CreatePayment - API tạo thanh toán VNPay
func (v *VNPayController) CreatePayment(c *gin.Context) {
	orderID := c.Param("order_id")
	amount := 100000 // 🚀 Số tiền thanh toán (lấy từ DB nếu cần)

	paymentURL, err := v.Service.GenerateVNPayURL(orderID, amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi tạo URL thanh toán"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment_url": paymentURL})
}
