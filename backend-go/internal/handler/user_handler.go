package handler // ประกาศ package handler

import ( // เริ่มส่วน import package ที่ไฟล์นี้ต้องใช้
	"net/http" // นำเข้า package net/http
	"strings"  // นำเข้า package strings
	"time"     // นำเข้า package time

	"backend-go/internal/model" // นำเข้า package backend-go/internal/model

	"github.com/gin-gonic/gin" // นำเข้า package github.com/gin-gonic/gin
	"gorm.io/gorm"             // นำเข้า package gorm.io/gorm
) // ปิด block หรือกลุ่มรายการก่อนหน้า

type UserHandler struct { // ประกาศ struct UserHandler
	db *gorm.DB // ทำงานคำสั่ง db *gorm.DB
} // ปิด block การทำงานปัจจุบัน

func NewUserHandler(db *gorm.DB) *UserHandler { // ประกาศฟังก์ชัน NewUserHandler
	return &UserHandler{db: db} // คืนค่า &UserHandler{db: db}
} // ปิด block การทำงานปัจจุบัน

type UpdateProfileRequest struct { // ประกาศ struct UpdateProfileRequest
	Name        string `json:"name" binding:"required,min=2,max=100"` // ประกาศ field Name พร้อม tag JSON/database
	Phone       string `json:"phone" binding:"required,min=6,max=30"` // ประกาศ field Phone พร้อม tag JSON/database
	DateOfBirth string `json:"date_of_birth" binding:"required"`      // ประกาศ field DateOfBirth พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน

func (h *UserHandler) GetMe(c *gin.Context) { // ประกาศฟังก์ชัน GetMe
	userID := c.GetUint("userID") // กำหนดค่า user ID จาก c.GetUint("userID")

	var user model.User                                     // ประกาศตัวแปร user model.User
	if err := h.db.First(&user, userID).Error; err != nil { // ตรวจเงื่อนไข err := h.db.First(&user, userID).Error; err != nil
		Error(c, http.StatusNotFound, "user not found") // ทำงานคำสั่ง Error(c, http.StatusNotFound, "user not found")
		return                                          // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	Success(c, http.StatusOK, "success", userResponse(user)) // ทำงานคำสั่ง Success(c, http.StatusOK, "success", userResponse(user))
} // ปิด block การทำงานปัจจุบัน

func (h *UserHandler) UpdateMe(c *gin.Context) { // ประกาศฟังก์ชัน UpdateMe
	userID := c.GetUint("userID") // กำหนดค่า user ID จาก c.GetUint("userID")

	var req UpdateProfileRequest                   // ประกาศตัวแปร req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil { // ตรวจเงื่อนไข err := c.ShouldBindJSON(&req); err != nil
		BadRequest(c, "invalid request") // ทำงานคำสั่ง BadRequest(c, "invalid request")
		return                           // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	if _, err := time.Parse("2006-01-02", req.DateOfBirth); err != nil { // ตรวจเงื่อนไข _, err := time.Parse("2006-01-02", req.DateOfBirth); err != nil
		BadRequest(c, "date_of_birth must use YYYY-MM-DD format") // ทำงานคำสั่ง BadRequest(c, "date_of_birth must use YYYY-MM-DD format")
		return                                                    // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	var user model.User                                     // ประกาศตัวแปร user model.User
	if err := h.db.First(&user, userID).Error; err != nil { // ตรวจเงื่อนไข err := h.db.First(&user, userID).Error; err != nil
		Error(c, http.StatusNotFound, "user not found") // ทำงานคำสั่ง Error(c, http.StatusNotFound, "user not found")
		return                                          // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	user.Name = strings.TrimSpace(req.Name)               // กำหนดค่า Name ของ user จาก strings.TrimSpace(req.Name)
	user.Phone = strings.TrimSpace(req.Phone)             // กำหนดค่า Phone ของ user จาก strings.TrimSpace(req.Phone)
	user.DateOfBirth = strings.TrimSpace(req.DateOfBirth) // กำหนดค่า Date Of Birth ของ user จาก strings.TrimSpace(req.DateOfBirth)

	if err := h.db.Save(&user).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Save(&user).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to update user") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to update user")
		return                                                            // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	Success(c, http.StatusOK, "updated", userResponse(user)) // ทำงานคำสั่ง Success(c, http.StatusOK, "updated", userResponse(user))
} // ปิด block การทำงานปัจจุบัน

func (h *UserHandler) GetUsers(c *gin.Context) { // ประกาศฟังก์ชัน GetUsers
	var users []model.User                                                   // ประกาศตัวแปร users []model.User
	if err := h.db.Order("created_at desc").Find(&users).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Order("created_at desc").Find(&users).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to get users") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to get users")
		return                                                          // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	response := make([]gin.H, 0, len(users)) // กำหนดค่า response จาก make([]gin.H, 0, len(users))
	for _, user := range users {             // วนลูปตาม _, user := range users
		response = append(response, userResponse(user)) // กำหนดค่า response จาก append(response, userResponse(user))
	} // ปิด block การทำงานปัจจุบัน

	Success(c, http.StatusOK, "success", response) // ทำงานคำสั่ง Success(c, http.StatusOK, "success", response)
} // ปิด block การทำงานปัจจุบัน

func userResponse(user model.User) gin.H { // ประกาศฟังก์ชัน userResponse
	return gin.H{ // คืนค่า gin.H
		"id":            user.ID,          // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"username":      user.Username,    // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"email":         user.Email,       // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"name":          user.Name,        // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"phone":         user.Phone,       // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"date_of_birth": user.DateOfBirth, // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"role":          user.Role,        // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"created_at":    user.CreatedAt,   // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"updated_at":    user.UpdatedAt,   // ระบุค่ารายการหนึ่งในชุดข้อมูล
	} // ปิด block การทำงานปัจจุบัน
} // ปิด block การทำงานปัจจุบัน
