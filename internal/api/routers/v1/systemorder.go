package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterOrderSystemRouter(router *gin.RouterGroup) {

	router.GET("/get", controllers.Statistical.GetStatistical)

	router.POST("/register", controllers.User.Register)	
	router.PUT("/UpdateUsers", controllers.User.UpdateUsers)
	router.POST("/Login", controllers.User.Login)

	router.GET("/Getbases", controllers.Basese.Getbasese)
	router.POST("/Addbases", controllers.Basese.Addbases)
	router.PUT("/Updatebases", controllers.Basese.Updatebases)
	router.DELETE("/Deletebases", controllers.Basese.Deletebases)
	router.POST("/Searchbases", controllers.Basese.Searchbases)

	router.GET("/Getsize", controllers.Sizes.GetSizes)
	router.POST("/Addsize", controllers.Sizes.AddSizes)
	router.PUT("/Updatesize", controllers.Sizes.Updatesize)
	router.DELETE("/Deletesize", controllers.Sizes.Deletesize)
	router.POST("/Searchsize", controllers.Sizes.Searchsize)

	router.GET("/Getflavors", controllers.Flavors.GetFlavors)
	router.POST("/Addflavors", controllers.Flavors.Addflavors)
	router.PUT("/Updateflavors", controllers.Flavors.Updateflavors)
	router.DELETE("/Deleteflavors", controllers.Flavors.Deleteflavors)

	router.GET("/GeticeLevels", controllers.IceLevels.GetIceLevels)
	router.GET("/Getsweetness", controllers.Sweetness.GetSweetness)

	router.GET("/GetBaseSizes", controllers.BaseSizes.GetBaseSizes)
	router.POST("/BaseSizes", controllers.BaseSizes.AddBaseSizes)
	router.PUT("/UpdateBaseSizes", controllers.BaseSizes.UpdateBaseSizes)
	router.DELETE("/DeleteBaseSizes", controllers.BaseSizes.DeleteBaseSizes)
	router.POST("/SearchBaseSizes", controllers.BaseSizes.SearchBaseSizes)

	//Router order
	router.POST("/order/:userID", controllers.NewOrderController().PlaceOrder)               //🚀 API: checkout
	router.GET("/orders/:orderID/details", controllers.NewOrderController().GetOrderDetails) //🚀 API: lấy OrderDetails theo OrderID
	router.GET("/users/:userID/orders", controllers.NewOrderController().GetOrdersByUserID)  //🚀 API: lấy tất cả đơn hàng theo UserID
	router.PUT("/orders/:orderID/cancel", controllers.NewOrderController().CancelOrder)      // 🚀 API hủy đơn hàng

	//router cart
	router.GET("/cart/:userID", controllers.NewCartController().GetCart)                       //🚀 API: lấy giỏ hàng theo user Id
	router.POST("/cart/:userID", controllers.NewCartController().AddToCart)                    //🚀 API:  add to cart
	router.PUT("/cart/:userID/:cartItemID", controllers.NewCartController().UpdateCart)        //🚀 API:  update cart
	router.DELETE("/cart/:userID/:cartItemID", controllers.NewCartController().RemoveFromCart) //🚀 API:  xóa giỏ hàng theo userid và caarrt id

	//router admin orders
	router.GET("/admin/orders/status-list", controllers.NewAdminOrderController().GetOrderStatusList)    // 🚀 API: Lấy danh sách trạng thái đơn hàng
	router.PUT("/admin/orders/:orderID/status", controllers.NewAdminOrderController().UpdateOrderStatus) // 🚀 API: Admin cập nhật trạng thái đơn hàng

	//router tích hợp thanh toán vnpay
	// Khởi tạo cấu hình VNPay
	router.GET("/create_payment/:order_id", controllers.NewVNPayController().CreatePayment)
}
