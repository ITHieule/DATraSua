package controllers

import (
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type SweetnessController struct {
	*BaseController
}

var Sweetness = &SweetnessController{}

func (c *SweetnessController) GetSweetness(ctx *gin.Context) {
	result, err := services.OrderSweetness.SweetnessSevice()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}
