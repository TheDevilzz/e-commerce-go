package main // ประกาศ package main

import ( // เริ่มส่วน import package ที่ไฟล์นี้ต้องใช้
	"backend-go/config"              // นำเข้า package backend-go/config
	"backend-go/database"            // นำเข้า package backend-go/database
	"backend-go/internal/handler"    // นำเข้า package backend-go/internal/handler
	"backend-go/internal/middleware" // นำเข้า package backend-go/internal/middleware
	"backend-go/internal/model"      // นำเข้า package backend-go/internal/model
	"backend-go/internal/repository" // นำเข้า package backend-go/internal/repository
	"backend-go/internal/router"     // นำเข้า package backend-go/internal/router
	"backend-go/internal/service"    // นำเข้า package backend-go/internal/service
) // ปิด block หรือกลุ่มรายการก่อนหน้า

func main() { // ประกาศฟังก์ชัน main
	cfg := config.LoadConfig()             // กำหนดค่า config จาก config.LoadConfig()
	middleware.SetJWTSecret(cfg.JWTSecret) // ทำงานคำสั่ง middleware.SetJWTSecret(cfg.JWTSecret)

	db := database.ConnectDB(cfg) // กำหนดค่า database จาก database.ConnectDB(cfg)

	err := db.AutoMigrate( // กำหนดค่า error จาก db.AutoMigrate(
		&model.User{},          // ระบุค่ารายการหนึ่งในชุดข้อมูล
		&model.RefreshToken{},  // ระบุค่ารายการหนึ่งในชุดข้อมูล
		&model.Category{},      // ระบุค่ารายการหนึ่งในชุดข้อมูล
		&model.Product{},       // ระบุค่ารายการหนึ่งในชุดข้อมูล
		&model.Order{},         // ระบุค่ารายการหนึ่งในชุดข้อมูล
		&model.OrderItem{},     // ระบุค่ารายการหนึ่งในชุดข้อมูล
		&model.Payment{},       // ระบุค่ารายการหนึ่งในชุดข้อมูล
		&model.CartItem{},      // ระบุค่ารายการหนึ่งในชุดข้อมูล
		&model.WishlistItem{},  // ระบุค่ารายการหนึ่งในชุดข้อมูล
		&model.Address{},       // ระบุค่ารายการหนึ่งในชุดข้อมูล
		&model.PaymentMethod{}, // ระบุค่ารายการหนึ่งในชุดข้อมูล
		&model.Promotion{},     // ระบุค่ารายการหนึ่งในชุดข้อมูล
	) // ปิด block หรือกลุ่มรายการก่อนหน้า
	if err != nil { // ตรวจเงื่อนไข err != nil
		panic("Failed to migrate database: " + err.Error()) // ทำงานคำสั่ง panic("Failed to migrate database: " + err.Error())
	} // ปิด block การทำงานปัจจุบัน

	productRepo := repository.NewProductRepository(db)          // กำหนดค่า product Repo จาก repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)    // กำหนดค่า product Service จาก service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService) // กำหนดค่า product Handler จาก handler.NewProductHandler(productService)

	authHandler := handler.NewAuthHandler(db, cfg.RefreshTokenDays) // กำหนดค่า auth Handler จาก handler.NewAuthHandler(db, cfg.RefreshTokenDays)
	userHandler := handler.NewUserHandler(db)                       // กำหนดค่า user Handler จาก handler.NewUserHandler(db)
	orderHandler := handler.NewOrderHandler(db)                     // กำหนดค่า order Handler จาก handler.NewOrderHandler(db)
	paymentHandler := handler.NewPaymentHandler(db)                 // กำหนดค่า payment Handler จาก handler.NewPaymentHandler(db)
	customerHandler := handler.NewCustomerHandler(db)               // กำหนดค่า customer Handler จาก handler.NewCustomerHandler(db)
	promotionHandler := handler.NewPromotionHandler(db)             // กำหนดค่า promotion Handler จาก handler.NewPromotionHandler(db)
	dashboardHandler := handler.NewDashboardHandler(db)             // กำหนดค่า dashboard Handler จาก handler.NewDashboardHandler(db)

	r := router.SetupRouter(cfg, db, productHandler, authHandler, userHandler, orderHandler, paymentHandler, customerHandler, promotionHandler, dashboardHandler) // กำหนดค่า r จาก router.SetupRouter(cfg, db, productHandler, authHandler, userHandler, orderHandler, paymentHandler, customerHandler, promotionHandler, dashboardHandler)

	r.Run(":" + cfg.Port) // ทำงานคำสั่ง r.Run(":" + cfg.Port)
} // ปิด block การทำงานปัจจุบัน
