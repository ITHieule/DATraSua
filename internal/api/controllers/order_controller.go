package controllers

import (
	"log"
	"net/http"
	"strconv"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"

	"github.com/gin-gonic/gin"
)

// OrderController xử lý API đặt hàng
type OrderController struct {
	orderService *services.OrderService
	vnpayService *services.VNPayService
}

// 🔹 Biến toàn cục `Orders`
var Orders *OrderController

// Getter cho orderService
func (c *OrderController) GetOrderService() *services.OrderService {
	return c.orderService
}

// Hàm khởi tạo `Orders`
func InitOrderController(orderService *services.OrderService, vnpayService *services.VNPayService) {
	Orders = &OrderController{
		orderService: orderService,
		vnpayService: vnpayService,
	}
	log.Println(" Orders đã được khởi tạo thành công!")
}

// API Đặt hàng
func (c *OrderController) PlaceOrder(ctx *gin.Context) {
	// 🔹 Kiểm tra xem orderService đã được khởi tạo chưa
	if c.orderService == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "OrderService chưa được khởi tạo"})
		return
	}

	var orderRequest request.OrderRequest
	if err := ctx.ShouldBindJSON(&orderRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	// 🔹 Gọi service để xử lý đặt hàng
	order, paymentURL, err := c.orderService.PlaceOrder(orderRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"order":       order,
		"payment_url": paymentURL,
	})
}

// VNPayCallback xử lý phản hồi từ VNPay
func (c *OrderController) VNPayCallback(ctx *gin.Context) {
	params := ctx.Request.URL.Query()

	// ✅ Kiểm tra chữ ký VNPay trước khi xử lý
	if !c.vnpayService.VerifyVNPaySignature(params) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Chữ ký không hợp lệ"})
		return
	}

	orderID := params.Get("vnp_TxnRef")
	transactionStatus := params.Get("vnp_ResponseCode")

	if orderID == "" || transactionStatus == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Thiếu thông tin phản hồi từ VNPay"})
		return
	}

	var statusMessage string
	if transactionStatus == "00" {
		statusMessage = "Giao dịch thành công"
	} else {
		statusMessage = "Giao dịch thất bại"
	}

	log.Printf("[VNPay] ✅ Nhận callback cho đơn hàng %s, trạng thái: %s", orderID, statusMessage)

	// Phản hồi kết quả mà không cập nhật database
	ctx.JSON(http.StatusOK, gin.H{"message": "Xử lý callback thành công", "transaction_status": statusMessage})
}

// API: Lấy danh sách OrderDetails theo OrderID
func (c *OrderController) GetOrderDetails(ctx *gin.Context) {
	//  Lấy orderID từ URL
	orderID, err := strconv.Atoi(ctx.Param("orderID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Order ID không hợp lệ"})
		return
	}

	//  Gọi service để lấy danh sách OrderDetails
	orderDetails, err := c.orderService.GetOrderDetailsByOrderID(orderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//  Trả về JSON danh sách OrderDetails
	ctx.JSON(http.StatusOK, gin.H{"order_details": orderDetails})
}

// API: Lấy tất cả đơn hàng theo UserID
func (c *OrderController) GetOrdersByUserID(ctx *gin.Context) {
	//  Lấy userID từ URL
	userID, err := strconv.Atoi(ctx.Param("userID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID không hợp lệ"})
		return
	}

	//  Gọi service để lấy danh sách đơn hàng
	orders, err := c.orderService.GetOrdersByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//  Trả về danh sách đơn hàng
	ctx.JSON(http.StatusOK, gin.H{"orders": orders})
}

// API: Hủy đơn hàng theo OrderID
func (c *OrderController) CancelOrder(ctx *gin.Context) {
	//  Lấy orderID từ URL
	orderID, err := strconv.Atoi(ctx.Param("orderID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Order ID không hợp lệ"})
		return
	}

	//  Gọi service để hủy đơn hàng
	err = c.orderService.CancelOrder(orderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//  Trả về thông báo thành công
	ctx.JSON(http.StatusOK, gin.H{"message": "Đơn hàng đã được hủy thành công"})
}
