package controllers

import (
	"net/http"
	"strconv"
	"web-api/internal/api/services"

	"github.com/gin-gonic/gin"
)

type AdminOrderController struct {
	orderService *services.OrderService
}

func NewAdminOrderController() *AdminOrderController {
	return &AdminOrderController{
		orderService: services.NewOrderService(),
	}
}

// 📌 API: Lấy danh sách trạng thái đơn hàng
func (c *AdminOrderController) GetOrderStatusList(ctx *gin.Context) {
	statusList := []string{"Đang xử lý", "Đã xác nhận", "Đang giao", "Hoàn thành", "Đã hủy"}
	ctx.JSON(http.StatusOK, gin.H{"status_list": statusList})
}

// 📌 API: Admin cập nhật trạng thái đơn hàng
func (c *AdminOrderController) UpdateOrderStatus(ctx *gin.Context) {
	// 🛒 Lấy orderID từ URL
	orderID, err := strconv.Atoi(ctx.Param("orderID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// 🛒 Lấy status mới từ body request
	var requestBody struct {
		Status string `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 📌 Gọi service để cập nhật trạng thái đơn hàng
	err = c.orderService.UpdateOrderStatus(orderID, requestBody.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ✅ Trả về thông báo thành công
	ctx.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}
