package controllers

import (
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type IceLevelsController struct {
	*BaseController
}

var IceLevels = &IceLevelsController{}

func (c *IceLevelsController) GetIceLevels(ctx *gin.Context) {
	result, err := services.OrderIceLevels.IceLevelsSevice()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}
