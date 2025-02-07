package controllers

import (
	"net/http"
	"strconv"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type BaseseController struct {
	*BaseController
}

var Basese = &BaseseController{}

func (c *BaseseController) Getbasese(ctx *gin.Context) {
	result, err := services.Order.BasesSevice()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}

func (c *BaseseController) Addbases(ctx *gin.Context) {
	var requestParams request.Basesrequest

	// Nhận text từ form-data
	requestParams.Name = ctx.PostForm("name")
	// Ép kiểu price từ string -> float64
	priceStr := ctx.PostForm("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Giá trị price không hợp lệ")
		return
	}
	requestParams.Price = price

	// Nhận file từ form-data
	file, err := ctx.FormFile("images")
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, "Không tìm thấy file")
		return
	}

	// Lưu file vào thư mục
	filePath := "D:/Image/" + file.Filename
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, "Lỗi lưu file")
		return
	}

	requestParams.Images = file.Filename // Lưu path file vào struct

	// Gọi service xử lý
	result, err := services.Order.AddbasesSevice(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.OkWithData(ctx, result)
}

func (c *BaseseController) Updatebases(ctx *gin.Context) {
	var requestParams request.Basesrequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	result, err := services.Order.UpdatebasesSevice(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}

func (c *BaseseController) Deletebases(ctx *gin.Context) {
	var requestParams request.Basesrequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	err := services.Order.DeletebasesSevice(requestParams.Id)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, nil)
}
func (c *BaseseController) Searchbases(ctx *gin.Context) {
	var requestParams request.Basesrequest

	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	result, err := services.Order.SearchbasesSevice(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}
