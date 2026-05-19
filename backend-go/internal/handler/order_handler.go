package handler // ประกาศ package handler

import ( // เริ่มส่วน import package ที่ไฟล์นี้ต้องใช้
	"net/http" // นำเข้า package net/http

	"backend-go/internal/model" // นำเข้า package backend-go/internal/model

	"github.com/gin-gonic/gin" // นำเข้า package github.com/gin-gonic/gin
	"gorm.io/gorm"             // นำเข้า package gorm.io/gorm
) // ปิด block หรือกลุ่มรายการก่อนหน้า

type OrderHandler struct { // ประกาศ struct OrderHandler
	db *gorm.DB // ทำงานคำสั่ง db *gorm.DB
} // ปิด block การทำงานปัจจุบัน

func NewOrderHandler(db *gorm.DB) *OrderHandler { // ประกาศฟังก์ชัน NewOrderHandler
	return &OrderHandler{db: db} // คืนค่า &OrderHandler{db: db}
} // ปิด block การทำงานปัจจุบัน

type CreateOrderRequest struct { // ประกาศ struct CreateOrderRequest
	Items []CreateOrderItem `json:"items" binding:"required,min=1,dive"` // ประกาศ field Items พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน

type CreateOrderItem struct { // ประกาศ struct CreateOrderItem
	ProductID uint `json:"product_id" binding:"required,gt=0"`      // ประกาศ field ProductID พร้อม tag JSON/database
	Quantity  int  `json:"quantity" binding:"required,gt=0,lte=99"` // ประกาศ field Quantity พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน

func (h *OrderHandler) CreateOrder(c *gin.Context) { // ประกาศฟังก์ชัน CreateOrder
	userId := c.GetUint("userId") // กำหนดค่า user Id จาก c.GetUint("userId")

	var req CreateOrderRequest                     // ประกาศตัวแปร req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil { // ตรวจเงื่อนไข err := c.ShouldBindJSON(&req); err != nil
		BadRequest(c, "invalid request") // ทำงานคำสั่ง BadRequest(c, "invalid request")
		return                           // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	tx := h.db.Begin()   // กำหนดค่า tx จาก h.db.Begin()
	if tx.Error != nil { // ตรวจเงื่อนไข tx.Error != nil
		Error(c, http.StatusInternalServerError, "failed to start transaction") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to start transaction")
		return                                                                  // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	order := model.Order{ // กำหนดค่า order จาก model.Order
		UserID: userId,    // ระบุค่ารายการหนึ่งในชุดข้อมูล
		Status: "pending", // ระบุค่ารายการหนึ่งในชุดข้อมูล
	} // ปิด block การทำงานปัจจุบัน

	if err := tx.Create(&order).Error; err != nil { // ตรวจเงื่อนไข err := tx.Create(&order).Error; err != nil
		tx.Rollback()                                                      // ทำงานคำสั่ง tx.Rollback()
		Error(c, http.StatusInternalServerError, "failed to create order") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to create order")
		return                                                             // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	var total float64 // ประกาศตัวแปร total float64

	for _, item := range req.Items { // วนลูปตาม _, item := range req.Items
		var product model.Product // ประกาศตัวแปร product model.Product

		if err := tx.First(&product, item.ProductID).Error; err != nil { // ตรวจเงื่อนไข err := tx.First(&product, item.ProductID).Error; err != nil
			tx.Rollback()                                      // ทำงานคำสั่ง tx.Rollback()
			Error(c, http.StatusNotFound, "product not found") // ทำงานคำสั่ง Error(c, http.StatusNotFound, "product not found")
			return                                             // หยุดการทำงานของฟังก์ชันนี้ทันที
		} // ปิด block การทำงานปัจจุบัน

		if product.Stock < item.Quantity { // ตรวจเงื่อนไข product.Stock < item.Quantity
			tx.Rollback()                     // ทำงานคำสั่ง tx.Rollback()
			BadRequest(c, "stock not enough") // ทำงานคำสั่ง BadRequest(c, "stock not enough")
			return                            // หยุดการทำงานของฟังก์ชันนี้ทันที
		} // ปิด block การทำงานปัจจุบัน

		product.Stock -= item.Quantity // ทำงานคำสั่ง product.Stock -= item.Quantity

		if err := tx.Save(&product).Error; err != nil { // ตรวจเงื่อนไข err := tx.Save(&product).Error; err != nil
			tx.Rollback()                                                      // ทำงานคำสั่ง tx.Rollback()
			Error(c, http.StatusInternalServerError, "failed to update stock") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to update stock")
			return                                                             // หยุดการทำงานของฟังก์ชันนี้ทันที
		} // ปิด block การทำงานปัจจุบัน

		orderItem := model.OrderItem{ // กำหนดค่า order Item จาก model.OrderItem
			OrderID:     order.ID,         // ระบุค่ารายการหนึ่งในชุดข้อมูล
			ProductID:   uint(product.ID), // ระบุค่ารายการหนึ่งในชุดข้อมูล
			ProductName: product.Name,     // ระบุค่ารายการหนึ่งในชุดข้อมูล
			Quantity:    item.Quantity,    // ระบุค่ารายการหนึ่งในชุดข้อมูล
			Price:       product.Price,    // ระบุค่ารายการหนึ่งในชุดข้อมูล
		} // ปิด block การทำงานปัจจุบัน

		if err := tx.Create(&orderItem).Error; err != nil { // ตรวจเงื่อนไข err := tx.Create(&orderItem).Error; err != nil
			tx.Rollback()                                                           // ทำงานคำสั่ง tx.Rollback()
			Error(c, http.StatusInternalServerError, "failed to create order item") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to create order item")
			return                                                                  // หยุดการทำงานของฟังก์ชันนี้ทันที
		} // ปิด block การทำงานปัจจุบัน

		total += product.Price * float64(item.Quantity) // ทำงานคำสั่ง total += product.Price * float64(item.Quantity)
	} // ปิด block การทำงานปัจจุบัน

	order.TotalPrice = total // กำหนดค่า Total Price ของ order จาก total

	if err := tx.Save(&order).Error; err != nil { // ตรวจเงื่อนไข err := tx.Save(&order).Error; err != nil
		tx.Rollback()                                                            // ทำงานคำสั่ง tx.Rollback()
		Error(c, http.StatusInternalServerError, "failed to update order total") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to update order total")
		return                                                                   // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	if err := tx.Commit().Error; err != nil { // ตรวจเงื่อนไข err := tx.Commit().Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to commit order") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to commit order")
		return                                                             // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	Success(c, http.StatusCreated, "order created", order) // ทำงานคำสั่ง Success(c, http.StatusCreated, "order created", order)
} // ปิด block การทำงานปัจจุบัน

func (h *OrderHandler) GetMyOrders(c *gin.Context) { // ประกาศฟังก์ชัน GetMyOrders
	userId := c.GetUint("userId") // กำหนดค่า user Id จาก c.GetUint("userId")

	var orders []model.Order // ประกาศตัวแปร orders []model.Order

	if err := h.db.Preload("Items").Where("user_id = ?", userId).Find(&orders).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Preload("Items").Where("user_id = ?", userId).Find(&orders).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to get orders") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to get orders")
		return                                                           // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	Success(c, http.StatusOK, "success", orders) // ทำงานคำสั่ง Success(c, http.StatusOK, "success", orders)
} // ปิด block การทำงานปัจจุบัน

func (h *OrderHandler) GetOrders(c *gin.Context) { // ประกาศฟังก์ชัน GetOrders
	var orders []model.Order // ประกาศตัวแปร orders []model.Order

	if err := h.db.Preload("Items").Order("created_at desc").Find(&orders).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Preload("Items").Order("created_at desc").Find(&orders).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to get orders") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to get orders")
		return                                                           // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	Success(c, http.StatusOK, "success", orders) // ทำงานคำสั่ง Success(c, http.StatusOK, "success", orders)
} // ปิด block การทำงานปัจจุบัน
