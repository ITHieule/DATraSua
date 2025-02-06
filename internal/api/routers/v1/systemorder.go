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

	//Router order
	router.POST("/createOrder", controllers.Order.CreateOrder)
	router.GET("/orders/:order_id/details", controllers.Order.GetOrderWithDetails)
	router.PUT("/orders/:order_id/cancel", controllers.Order.CancelOrder)

	//router cart
	router.GET("/cart/:userID", controllers.NewCartController().GetCart)
	router.POST("/cart/:userID", controllers.NewCartController().AddToCart)
	router.PUT("/cart/:userID/:cartItemID", controllers.NewCartController().UpdateCart)
	router.DELETE("/cart/:userID/:cartItemID", controllers.NewCartController().RemoveFromCart)
}
