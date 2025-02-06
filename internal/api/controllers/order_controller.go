package controllers

import (
	"net/http"
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
