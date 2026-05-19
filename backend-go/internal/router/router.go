package router // ประกาศ package router

import ( // เริ่มส่วน import package ที่ไฟล์นี้ต้องใช้
	"backend-go/config"              // นำเข้า package backend-go/config
	"backend-go/internal/handler"    // นำเข้า package backend-go/internal/handler
	"backend-go/internal/middleware" // นำเข้า package backend-go/internal/middleware
	"net/http"                       // นำเข้า package net/http

	"github.com/gin-gonic/gin" // นำเข้า package github.com/gin-gonic/gin
	"gorm.io/gorm"             // นำเข้า package gorm.io/gorm
) // ปิด block หรือกลุ่มรายการก่อนหน้า

func SetupRouter( // ประกาศฟังก์ชัน SetupRouter
	cfg config.Config, // ระบุค่ารายการหนึ่งในชุดข้อมูล
	db *gorm.DB, // ระบุค่ารายการหนึ่งในชุดข้อมูล
	productHandler *handler.ProductHandler, // ระบุค่ารายการหนึ่งในชุดข้อมูล
	authHandler *handler.AuthHandler, // ระบุค่ารายการหนึ่งในชุดข้อมูล
	userHandler *handler.UserHandler, // ระบุค่ารายการหนึ่งในชุดข้อมูล
	orderHandler *handler.OrderHandler, // ระบุค่ารายการหนึ่งในชุดข้อมูล
	paymentHandler *handler.PaymentHandler, // ระบุค่ารายการหนึ่งในชุดข้อมูล
	customerHandler *handler.CustomerHandler, // ระบุค่ารายการหนึ่งในชุดข้อมูล
	promotionHandler *handler.PromotionHandler, // ระบุค่ารายการหนึ่งในชุดข้อมูล
	dashboardHandler *handler.DashboardHandler, // ระบุค่ารายการหนึ่งในชุดข้อมูล
) *gin.Engine { // ทำงานคำสั่ง ) *gin.Engine {
	r := gin.Default() // กำหนดค่า r จาก gin.Default()

	r.Use(middleware.ErrorHandler())              // ติดตั้ง middleware ให้ router/group
	r.Use(middleware.HTTPSOnly(cfg.EnforceHTTPS)) // ติดตั้ง middleware ให้ router/group
	r.Use(middleware.SecurityHeaders())           // ติดตั้ง middleware ให้ router/group
	r.Use(corsMiddleware(cfg.AllowedOrigins))     // ติดตั้ง middleware ให้ router/group

	r.Static("/uploads", "./uploads") // ทำงานคำสั่ง r.Static("/uploads", "./uploads")
	r.Static("/slips", "./slips")     // ทำงานคำสั่ง r.Static("/slips", "./slips")

	api := r.Group("/api") // กำหนดค่า api จาก r.Group("/api")
	{                      // เปิด block การทำงานใหม่
		api.POST("/register", authHandler.Register)                               // ลงทะเบียน endpoint POST
		api.POST("/login", authHandler.Login)                                     // ลงทะเบียน endpoint POST
		api.POST("/refresh", authHandler.Refresh)                                 // ลงทะเบียน endpoint POST
		api.POST("/logout", authHandler.Logout)                                   // ลงทะเบียน endpoint POST
		api.GET("/products", productHandler.GetProducts)                          // ลงทะเบียน endpoint GET
		api.GET("/products/:id", productHandler.GetProductByID)                   // ลงทะเบียน endpoint GET
		api.GET("/products/search", productHandler.SearchProducts)                // ลงทะเบียน endpoint GET
		api.GET("/categories", productHandler.GetCategories)                      // ลงทะเบียน endpoint GET
		api.GET("/categories/:id/products", productHandler.GetProductsByCategory) // ลงทะเบียน endpoint GET
		api.GET("/promotions", promotionHandler.GetPromotions)                    // ลงทะเบียน endpoint GET

		protected := api.Group("")                 // กำหนดค่า protected จาก api.Group("")
		protected.Use(middleware.AuthMiddleware()) // ติดตั้ง middleware ให้ router/group
		{                                          // เปิด block การทำงานใหม่
			admin := protected.Group("")               // กำหนดค่า admin จาก protected.Group("")
			admin.Use(middleware.RequireRole("admin")) // ติดตั้ง middleware ให้ router/group
			{                                          // เปิด block การทำงานใหม่
				admin.POST("/products", productHandler.CreateProduct)             // ลงทะเบียน endpoint POST
				admin.PUT("/products/:id", productHandler.UpdateProduct)          // ลงทะเบียน endpoint PUT
				admin.DELETE("/products/:id", productHandler.DeleteProduct)       // ลงทะเบียน endpoint DELETE
				admin.POST("/upload", handler.UploadImage)                        // ลงทะเบียน endpoint POST
				admin.POST("/categories", productHandler.CreateCategory)          // ลงทะเบียน endpoint POST
				admin.PUT("/categories/:id", productHandler.UpdateCategory)       // ลงทะเบียน endpoint PUT
				admin.DELETE("/categories/:id", productHandler.DeleteCategory)    // ลงทะเบียน endpoint DELETE
				admin.GET("/users", userHandler.GetUsers)                         // ลงทะเบียน endpoint GET
				admin.GET("/orders", orderHandler.GetOrders)                      // ลงทะเบียน endpoint GET
				admin.GET("/dashboard/stats", dashboardHandler.GetStats)          // ลงทะเบียน endpoint GET
				admin.POST("/promotions", promotionHandler.CreatePromotion)       // ลงทะเบียน endpoint POST
				admin.PUT("/promotions/:id", promotionHandler.UpdatePromotion)    // ลงทะเบียน endpoint PUT
				admin.DELETE("/promotions/:id", promotionHandler.DeletePromotion) // ลงทะเบียน endpoint DELETE
			} // ปิด block การทำงานปัจจุบัน

			protected.GET("/me", userHandler.GetMe)                                       // ลงทะเบียน endpoint GET
			protected.PUT("/me", userHandler.UpdateMe)                                    // ลงทะเบียน endpoint PUT
			protected.POST("/orders", orderHandler.CreateOrder)                           // ลงทะเบียน endpoint POST
			protected.GET("/orders/my", orderHandler.GetMyOrders)                         // ลงทะเบียน endpoint GET
			protected.POST("/payment", paymentHandler.UploadSlip)                         // ลงทะเบียน endpoint POST
			protected.POST("/promotions/apply", promotionHandler.ApplyPromotion)          // ลงทะเบียน endpoint POST
			protected.GET("/cart", customerHandler.GetCart)                               // ลงทะเบียน endpoint GET
			protected.POST("/cart", customerHandler.AddCartItem)                          // ลงทะเบียน endpoint POST
			protected.PUT("/cart/:id", customerHandler.UpdateCartItem)                    // ลงทะเบียน endpoint PUT
			protected.DELETE("/cart/:id", customerHandler.DeleteCartItem)                 // ลงทะเบียน endpoint DELETE
			protected.GET("/wishlist", customerHandler.GetWishlist)                       // ลงทะเบียน endpoint GET
			protected.POST("/wishlist", customerHandler.AddWishlistItem)                  // ลงทะเบียน endpoint POST
			protected.DELETE("/wishlist/:id", customerHandler.DeleteWishlistItem)         // ลงทะเบียน endpoint DELETE
			protected.GET("/addresses", customerHandler.GetAddresses)                     // ลงทะเบียน endpoint GET
			protected.POST("/addresses", customerHandler.CreateAddress)                   // ลงทะเบียน endpoint POST
			protected.PUT("/addresses/:id", customerHandler.UpdateAddress)                // ลงทะเบียน endpoint PUT
			protected.DELETE("/addresses/:id", customerHandler.DeleteAddress)             // ลงทะเบียน endpoint DELETE
			protected.GET("/payment-methods", customerHandler.GetPaymentMethods)          // ลงทะเบียน endpoint GET
			protected.POST("/payment-methods", customerHandler.CreatePaymentMethod)       // ลงทะเบียน endpoint POST
			protected.PUT("/payment-methods/:id", customerHandler.UpdatePaymentMethod)    // ลงทะเบียน endpoint PUT
			protected.DELETE("/payment-methods/:id", customerHandler.DeletePaymentMethod) // ลงทะเบียน endpoint DELETE
		} // ปิด block การทำงานปัจจุบัน
	} // ปิด block การทำงานปัจจุบัน

	return r // คืนค่า r
} // ปิด block การทำงานปัจจุบัน

func corsMiddleware(allowedOrigins []string) gin.HandlerFunc { // ประกาศฟังก์ชัน corsMiddleware
	allowed := map[string]bool{}            // กำหนดค่า allowed จาก map[string]bool{}
	for _, origin := range allowedOrigins { // วนลูปตาม _, origin := range allowedOrigins
		allowed[origin] = true // กำหนดค่า allowed[origin] จาก true
	} // ปิด block การทำงานปัจจุบัน

	return func(c *gin.Context) { // คืนค่า func(c *gin.Context)
		origin := c.GetHeader("Origin") // กำหนดค่า origin จาก c.GetHeader("Origin")

		if allowed[origin] { // ตรวจเงื่อนไข allowed[origin]
			c.Header("Access-Control-Allow-Origin", origin)                             // ตั้งค่า HTTP header ให้ response
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")     // ตั้งค่า HTTP header ให้ response
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS") // ตั้งค่า HTTP header ให้ response
			c.Header("Vary", "Origin")                                                  // ตั้งค่า HTTP header ให้ response
		} // ปิด block การทำงานปัจจุบัน

		if c.Request.Method == http.MethodOptions { // ตรวจเงื่อนไข c.Request.Method == http.MethodOptions
			c.AbortWithStatus(http.StatusNoContent) // หยุด request ปัจจุบันไม่ให้ไป middleware ถัดไป
			return                                  // หยุดการทำงานของฟังก์ชันนี้ทันที
		} // ปิด block การทำงานปัจจุบัน

		c.Next() // ส่ง request ไป middleware หรือ handler ถัดไป
	} // ปิด block การทำงานปัจจุบัน
} // ปิด block การทำงานปัจจุบัน
