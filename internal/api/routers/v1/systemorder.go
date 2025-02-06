package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterOrderSystemRouter(router *gin.RouterGroup) {

	router.POST("/register", controllers.User.Register)
	router.POST("/Login", controllers.User.Login)
	router.GET("/Getbases", controllers.Basese.Getbasese)

	router.GET("/Getsize", controllers.Sizes.GetSizes)
	router.POST("/Addsize", controllers.Sizes.AddSizes)
	router.PUT("/Updatesize", controllers.Sizes.Updatesize)
	router.DELETE("/Deletesize", controllers.Sizes.Deletesize)
	router.POST("/Searchsize", controllers.Sizes.Searchsize)

	router.GET("/Getflavors", controllers.Flavors.GetFlavors)
	router.GET("/GeticeLevels", controllers.IceLevels.GetIceLevels)
	router.GET("/Getsweetness", controllers.Sweetness.GetSweetness)
	router.GET("/GetBaseSizes", controllers.BaseSizes.GetBaseSizes)
	router.POST("/BaseSizes", controllers.BaseSizes.AddBaseSizes)

	router.POST("/createOrder", controllers.Order.CreateOrder)

}
