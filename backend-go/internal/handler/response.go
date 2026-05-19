package handler // ประกาศ package handler

import ( // เริ่มส่วน import package ที่ไฟล์นี้ต้องใช้
	"net/http" // นำเข้า package net/http

	"github.com/gin-gonic/gin" // นำเข้า package github.com/gin-gonic/gin
) // ปิด block หรือกลุ่มรายการก่อนหน้า

func Success(c *gin.Context, status int, message string, data any) { // ประกาศฟังก์ชัน Success
	response := gin.H{"message": message} // กำหนดค่า response จาก gin.H{"message": message}
	if data != nil {                      // ตรวจเงื่อนไข data != nil
		response["data"] = data // กำหนดค่า response["data"] จาก data
	} // ปิด block การทำงานปัจจุบัน
	c.JSON(status, response) // ส่ง JSON response กลับไปยัง client
} // ปิด block การทำงานปัจจุบัน

func Error(c *gin.Context, status int, message string) { // ประกาศฟังก์ชัน Error
	c.JSON(status, gin.H{"message": message}) // ส่ง JSON response กลับไปยัง client
} // ปิด block การทำงานปัจจุบัน

func BadRequest(c *gin.Context, message string) { // ประกาศฟังก์ชัน BadRequest
	Error(c, http.StatusBadRequest, message) // ทำงานคำสั่ง Error(c, http.StatusBadRequest, message)
} // ปิด block การทำงานปัจจุบัน
