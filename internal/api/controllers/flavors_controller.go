package controllers

import (
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type FlavorsController struct {
	*BaseController
}

var Flavors = &FlavorsController{}

func (c *FlavorsController) GetFlavors(ctx *gin.Context) {
	result, err := services.Orderflavors.FlavorsSevice()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}
