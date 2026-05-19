package handler // ประกาศ package handler

import ( // เริ่มส่วน import package ที่ไฟล์นี้ต้องใช้
	"net/http" // นำเข้า package net/http

	"backend-go/internal/model" // นำเข้า package backend-go/internal/model

	"github.com/gin-gonic/gin" // นำเข้า package github.com/gin-gonic/gin
	"gorm.io/gorm"             // นำเข้า package gorm.io/gorm
) // ปิด block หรือกลุ่มรายการก่อนหน้า

type DashboardHandler struct { // ประกาศ struct DashboardHandler
	db *gorm.DB // ทำงานคำสั่ง db *gorm.DB
} // ปิด block การทำงานปัจจุบัน

type categorySalesRow struct { // ประกาศ struct categorySalesRow
	Category string  `json:"category"` // ประกาศ field Category พร้อม tag JSON/database
	Sales    float64 `json:"sales"`    // ประกาศ field Sales พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน

func NewDashboardHandler(db *gorm.DB) *DashboardHandler { // ประกาศฟังก์ชัน NewDashboardHandler
	return &DashboardHandler{db: db} // คืนค่า &DashboardHandler{db: db}
} // ปิด block การทำงานปัจจุบัน

func (h *DashboardHandler) GetStats(c *gin.Context) { // ประกาศฟังก์ชัน GetStats
	var revenue float64     // ประกาศตัวแปร revenue float64
	var orderCount int64    // ประกาศตัวแปร orderCount int64
	var customerCount int64 // ประกาศตัวแปร customerCount int64
	var productCount int64  // ประกาศตัวแปร productCount int64

	if err := h.db.Model(&model.Order{}).Select("COALESCE(SUM(total_price), 0)").Scan(&revenue).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Model(&model.Order{}).Select("COALESCE(SUM(total_price), 0)").Scan(&revenue).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to calculate revenue") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to calculate revenue")
		return                                                                  // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	if err := h.db.Model(&model.Order{}).Count(&orderCount).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Model(&model.Order{}).Count(&orderCount).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to count orders") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to count orders")
		return                                                             // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	if err := h.db.Model(&model.User{}).Where("role = ?", "user").Count(&customerCount).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Model(&model.User{}).Where("role = ?", "user").Count(&customerCount).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to count customers") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to count customers")
		return                                                                // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	if err := h.db.Model(&model.Product{}).Count(&productCount).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Model(&model.Product{}).Count(&productCount).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to count products") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to count products")
		return                                                               // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	avgOrder := 0.0     // กำหนดค่า avg Order จาก 0.0
	if orderCount > 0 { // ตรวจเงื่อนไข orderCount > 0
		avgOrder = revenue / float64(orderCount) // กำหนดค่า avg Order จาก revenue / float64(orderCount)
	} // ปิด block การทำงานปัจจุบัน

	var categorySales []categorySalesRow // ประกาศตัวแปร categorySales []categorySalesRow
	if err := h.db.Table("order_items"). // ตรวจเงื่อนไข err := h.db.Table("order_items").
						Select("categories.name AS category, COALESCE(SUM(order_items.price * order_items.quantity), 0) AS sales"). // ทำงานคำสั่ง Select("categories.name AS category, COALESCE(SUM(order_items.price * order_items.quantity), 0) AS sales").
						Joins("JOIN products ON products.id = order_items.product_id").                                             // กำหนดค่า id ของ Joins("JOIN products ON products จาก order_items.product_id").
						Joins("JOIN categories ON categories.id = products.category_id").                                           // กำหนดค่า id ของ Joins("JOIN categories ON categories จาก products.category_id").
						Group("categories.name").                                                                                   // ทำงานคำสั่ง Group("categories.name").
						Scan(&categorySales).Error; err != nil {                                                                    // ทำงานคำสั่ง Scan(&categorySales).Error; err != nil {
		Error(c, http.StatusInternalServerError, "failed to calculate category sales") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to calculate category sales")
		return                                                                         // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	Success(c, http.StatusOK, "success", gin.H{ // ทำงานคำสั่ง Success(c, http.StatusOK, "success", gin.H{
		"total_revenue":  revenue,       // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"orders":         orderCount,    // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"customers":      customerCount, // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"products":       productCount,  // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"avg_order":      avgOrder,      // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"category_sales": categorySales, // ระบุค่ารายการหนึ่งในชุดข้อมูล
	}) // ทำงานคำสั่ง })
} // ปิด block การทำงานปัจจุบัน
