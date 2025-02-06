package controllers

import (
	"net/http"
	"strconv"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	*BaseController
}

var Order = &OrderController{}

func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var requestParams request.OrderRequest

	// Xác thực tham số yêu cầu
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	// Tạo đơn hàng và chi tiết đơn hàng
	orderDetails, err := services.OrderDetails.Create(&requestParams) // Gọi service để tạo đơn hàng và chi tiết
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	// Trả về kết quả thành công với thông tin đơn hàng và chi tiết đơn hàng
	response.OkWithData(ctx, gin.H{"order_id": requestParams.ID, "order_details": orderDetails})
}

func (c *OrderController) GetOrderWithDetails(ctx *gin.Context) {
	// Lấy order_id từ URL (hoặc query parameter)
	orderIDStr := ctx.Param("order_id") // /orders/:order_id/details

	// Chuyển orderID từ string sang int
	orderID, err := strconv.Atoi(orderIDStr) // strconv.Atoi chuyển string sang int
	if err != nil {
		// Trả về lỗi nếu không thể chuyển đổi
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Invalid order_id")
		return
	}

	// Gọi hàm GetOrderWithDetails từ OrderDetailsService
	orderWithDetails, err := services.OrderDetails.GetOrderWithDetails(orderID)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusNotFound, nil, err.Error())
		return
	}

	// Trả về kết quả
	response.OkWithData(ctx, orderWithDetails)
}

// Hàm xử lý API hủy đơn hàng
func (c *OrderController) CancelOrder(ctx *gin.Context) {
	// Lấy order_id từ URL
	orderIDStr := ctx.Param("order_id")

	// Chuyển orderID từ string sang int
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		// Trả về lỗi nếu không thể chuyển đổi
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Invalid order_id")
		return
	}

	// Gọi hàm CancelOrder từ dịch vụ để hủy đơn hàng
	err = services.Orders.CancelOrder(orderID)
	if err != nil {
		// Trả về lỗi nếu không thể hủy đơn hàng
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	// Trả về phản hồi thành công
	response.OkWithMessage(ctx, "Order successfully canceled")
}
