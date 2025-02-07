package controllers

import (
	"net/http"
	"strconv"
	"web-api/internal/api/services"

	"github.com/gin-gonic/gin"
)

// OrderController xử lý API đặt hàng
type OrderController struct {
	orderService *services.OrderService
}

// NewOrderController khởi tạo controller mới
func NewOrderController() *OrderController {
	return &OrderController{
		orderService: services.NewOrderService(),
	}
}

// Đặt hàng từ giỏ hàng
func (c *OrderController) PlaceOrder(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("userID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	order, err := c.orderService.PlaceOrder(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, order)
}

// 🛒 API: Lấy danh sách OrderDetails theo OrderID
func (c *OrderController) GetOrderDetails(ctx *gin.Context) {
	// 📌 Lấy orderID từ URL
	orderID, err := strconv.Atoi(ctx.Param("orderID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// 📌 Gọi service để lấy danh sách OrderDetails
	orderDetails, err := c.orderService.GetOrderDetailsByOrderID(orderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ✅ Trả về JSON danh sách OrderDetails
	ctx.JSON(http.StatusOK, gin.H{"order_details": orderDetails})
}

// 📌 API: Lấy tất cả đơn hàng theo UserID
func (c *OrderController) GetOrdersByUserID(ctx *gin.Context) {
	// 🛒 Lấy userID từ URL
	userID, err := strconv.Atoi(ctx.Param("userID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 📌 Gọi service để lấy danh sách đơn hàng
	orders, err := c.orderService.GetOrdersByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ✅ Trả về danh sách đơn hàng
	ctx.JSON(http.StatusOK, gin.H{"orders": orders})
}

// 📌 API: Hủy đơn hàng theo OrderID
func (c *OrderController) CancelOrder(ctx *gin.Context) {
	// 🛒 Lấy orderID từ URL
	orderID, err := strconv.Atoi(ctx.Param("orderID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// 📌 Gọi service để hủy đơn hàng
	err = c.orderService.CancelOrder(orderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ✅ Trả về thông báo thành công
	ctx.JSON(http.StatusOK, gin.H{"message": "Order canceled successfully"})
}
