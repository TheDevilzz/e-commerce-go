package handler // ประกาศ package handler

import ( // เริ่มส่วน import package ที่ไฟล์นี้ต้องใช้
	"net/http" // นำเข้า package net/http
	"strings"  // นำเข้า package strings
	"time"     // นำเข้า package time

	"backend-go/internal/model" // นำเข้า package backend-go/internal/model

	"github.com/gin-gonic/gin" // นำเข้า package github.com/gin-gonic/gin
	"gorm.io/gorm"             // นำเข้า package gorm.io/gorm
) // ปิด block หรือกลุ่มรายการก่อนหน้า

type PromotionHandler struct { // ประกาศ struct PromotionHandler
	db *gorm.DB // ทำงานคำสั่ง db *gorm.DB
} // ปิด block การทำงานปัจจุบัน

func NewPromotionHandler(db *gorm.DB) *PromotionHandler { // ประกาศฟังก์ชัน NewPromotionHandler
	return &PromotionHandler{db: db} // คืนค่า &PromotionHandler{db: db}
} // ปิด block การทำงานปัจจุบัน

type PromotionRequest struct { // ประกาศ struct PromotionRequest
	Code        string  `json:"code" binding:"required,min=2,max=50"`                     // ประกาศ field Code พร้อม tag JSON/database
	Description string  `json:"description" binding:"required,min=2,max=255"`             // ประกาศ field Description พร้อม tag JSON/database
	Discount    float64 `json:"discount" binding:"required,gt=0"`                         // ประกาศ field Discount พร้อม tag JSON/database
	Type        string  `json:"type" binding:"required,oneof=percentage fixed"`           // ประกาศ field Type พร้อม tag JSON/database
	StartDate   string  `json:"start_date" binding:"required"`                            // ประกาศ field StartDate พร้อม tag JSON/database
	EndDate     string  `json:"end_date" binding:"required"`                              // ประกาศ field EndDate พร้อม tag JSON/database
	Status      string  `json:"status" binding:"required,oneof=active scheduled expired"` // ประกาศ field Status พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน

func (h *PromotionHandler) GetPromotions(c *gin.Context) { // ประกาศฟังก์ชัน GetPromotions
	var promotions []model.Promotion       // ประกาศตัวแปร promotions []model.Promotion
	query := h.db.Order("created_at desc") // กำหนดค่า query จาก h.db.Order("created_at desc")
	if c.Query("active") == "true" {       // ตรวจเงื่อนไข c.Query("active") == "true"
		query = query.Where("status = ?", "active") // กำหนดค่า query จาก query.Where("status = ?", "active")
	} // ปิด block การทำงานปัจจุบัน
	if err := query.Find(&promotions).Error; err != nil { // ตรวจเงื่อนไข err := query.Find(&promotions).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to get promotions") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to get promotions")
		return                                                               // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	Success(c, http.StatusOK, "success", promotions) // ทำงานคำสั่ง Success(c, http.StatusOK, "success", promotions)
} // ปิด block การทำงานปัจจุบัน

func (h *PromotionHandler) CreatePromotion(c *gin.Context) { // ประกาศฟังก์ชัน CreatePromotion
	h.savePromotion(c, 0) // ทำงานคำสั่ง h.savePromotion(c, 0)
} // ปิด block การทำงานปัจจุบัน

func (h *PromotionHandler) UpdatePromotion(c *gin.Context) { // ประกาศฟังก์ชัน UpdatePromotion
	id, err := parsePositiveID(c.Param("id")) // กำหนดค่า id, err จาก parsePositiveID(c.Param("id"))
	if err != nil {                           // ตรวจเงื่อนไข err != nil
		BadRequest(c, "invalid promotion id") // ทำงานคำสั่ง BadRequest(c, "invalid promotion id")
		return                                // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	h.savePromotion(c, id) // ทำงานคำสั่ง h.savePromotion(c, id)
} // ปิด block การทำงานปัจจุบัน

func (h *PromotionHandler) DeletePromotion(c *gin.Context) { // ประกาศฟังก์ชัน DeletePromotion
	id, err := parsePositiveID(c.Param("id")) // กำหนดค่า id, err จาก parsePositiveID(c.Param("id"))
	if err != nil {                           // ตรวจเงื่อนไข err != nil
		BadRequest(c, "invalid promotion id") // ทำงานคำสั่ง BadRequest(c, "invalid promotion id")
		return                                // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	if err := h.db.Delete(&model.Promotion{}, id).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Delete(&model.Promotion{}, id).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to delete promotion") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to delete promotion")
		return                                                                 // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	Success(c, http.StatusOK, "deleted", nil) // ทำงานคำสั่ง Success(c, http.StatusOK, "deleted", nil)
} // ปิด block การทำงานปัจจุบัน

func (h *PromotionHandler) savePromotion(c *gin.Context, id uint) { // ประกาศฟังก์ชัน savePromotion
	var req PromotionRequest                       // ประกาศตัวแปร req PromotionRequest
	if err := c.ShouldBindJSON(&req); err != nil { // ตรวจเงื่อนไข err := c.ShouldBindJSON(&req); err != nil
		BadRequest(c, "invalid request") // ทำงานคำสั่ง BadRequest(c, "invalid request")
		return                           // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	if _, err := time.Parse("2006-01-02", req.StartDate); err != nil { // ตรวจเงื่อนไข _, err := time.Parse("2006-01-02", req.StartDate); err != nil
		BadRequest(c, "start_date must use YYYY-MM-DD format") // ทำงานคำสั่ง BadRequest(c, "start_date must use YYYY-MM-DD format")
		return                                                 // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	if _, err := time.Parse("2006-01-02", req.EndDate); err != nil { // ตรวจเงื่อนไข _, err := time.Parse("2006-01-02", req.EndDate); err != nil
		BadRequest(c, "end_date must use YYYY-MM-DD format") // ทำงานคำสั่ง BadRequest(c, "end_date must use YYYY-MM-DD format")
		return                                               // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	promotion := model.Promotion{} // กำหนดค่า promotion จาก model.Promotion{}
	if id > 0 {                    // ตรวจเงื่อนไข id > 0
		if err := h.db.First(&promotion, id).Error; err != nil { // ตรวจเงื่อนไข err := h.db.First(&promotion, id).Error; err != nil
			Error(c, http.StatusNotFound, "promotion not found") // ทำงานคำสั่ง Error(c, http.StatusNotFound, "promotion not found")
			return                                               // หยุดการทำงานของฟังก์ชันนี้ทันที
		} // ปิด block การทำงานปัจจุบัน
	} // ปิด block การทำงานปัจจุบัน
	promotion.Code = strings.ToUpper(strings.TrimSpace(req.Code)) // กำหนดค่า Code ของ promotion จาก strings.ToUpper(strings.TrimSpace(req.Code))
	promotion.Description = strings.TrimSpace(req.Description)    // กำหนดค่า Description ของ promotion จาก strings.TrimSpace(req.Description)
	promotion.Discount = req.Discount                             // กำหนดค่า Discount ของ promotion จาก Discount ของ request
	promotion.Type = req.Type                                     // กำหนดค่า Type ของ promotion จาก Type ของ request
	promotion.StartDate = req.StartDate                           // กำหนดค่า Start Date ของ promotion จาก Start Date ของ request
	promotion.EndDate = req.EndDate                               // กำหนดค่า End Date ของ promotion จาก End Date ของ request
	promotion.Status = req.Status                                 // กำหนดค่า Status ของ promotion จาก Status ของ request

	if err := h.db.Save(&promotion).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Save(&promotion).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to save promotion") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to save promotion")
		return                                                               // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	Success(c, map[bool]int{true: http.StatusOK, false: http.StatusCreated}[id > 0], "saved", promotion) // ทำงานคำสั่ง Success(c, map[bool]int{true: http.StatusOK, false: http.StatusCreated}[id > 0], "saved", promotion)
} // ปิด block การทำงานปัจจุบัน

func (h *PromotionHandler) ApplyPromotion(c *gin.Context) { // ประกาศฟังก์ชัน ApplyPromotion
	var req struct { // ประกาศตัวแปร req struct
		Code     string  `json:"code" binding:"required"`           // ประกาศ field Code พร้อม tag JSON/database
		Subtotal float64 `json:"subtotal" binding:"required,gte=0"` // ประกาศ field Subtotal พร้อม tag JSON/database
	} // ปิด block การทำงานปัจจุบัน
	if err := c.ShouldBindJSON(&req); err != nil { // ตรวจเงื่อนไข err := c.ShouldBindJSON(&req); err != nil
		BadRequest(c, "invalid request") // ทำงานคำสั่ง BadRequest(c, "invalid request")
		return                           // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	var promotion model.Promotion                                                                                                                 // ประกาศตัวแปร promotion model.Promotion
	if err := h.db.Where("code = ? AND status = ?", strings.ToUpper(strings.TrimSpace(req.Code)), "active").First(&promotion).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Where("code = ? AND status = ?", strings.ToUpper(strings.TrimSpace(req.Code)), "active").First(&promotion).Error; err != nil
		Error(c, http.StatusNotFound, "promotion not found") // ทำงานคำสั่ง Error(c, http.StatusNotFound, "promotion not found")
		return                                               // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	discount := promotion.Discount      // กำหนดค่า discount จาก promotion.Discount
	if promotion.Type == "percentage" { // ตรวจเงื่อนไข promotion.Type == "percentage"
		discount = req.Subtotal * promotion.Discount / 100 // กำหนดค่า discount จาก req.Subtotal * promotion.Discount / 100
	} // ปิด block การทำงานปัจจุบัน
	Success(c, http.StatusOK, "applied", gin.H{"promotion": promotion, "discount": discount}) // ทำงานคำสั่ง Success(c, http.StatusOK, "applied", gin.H{"promotion": promotion, "discount": discount})
} // ปิด block การทำงานปัจจุบัน
