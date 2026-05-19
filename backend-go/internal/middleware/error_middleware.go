package middleware // ประกาศ package middleware

import ( // เริ่มส่วน import package ที่ไฟล์นี้ต้องใช้
	"net/http" // นำเข้า package net/http

	"github.com/gin-gonic/gin" // นำเข้า package github.com/gin-gonic/gin
) // ปิด block หรือกลุ่มรายการก่อนหน้า

func ErrorHandler() gin.HandlerFunc { // ประกาศฟังก์ชัน ErrorHandler
	return func(c *gin.Context) { // คืนค่า func(c *gin.Context)
		defer func() { // เลื่อนคำสั่งนี้ไปทำตอนฟังก์ชันจบ: func() {
			if recovered := recover(); recovered != nil { // ตรวจเงื่อนไข recovered := recover(); recovered != nil
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{ // หยุด request ปัจจุบันไม่ให้ไป middleware ถัดไป
					"message": "internal server error", // ระบุค่ารายการหนึ่งในชุดข้อมูล
				}) // ทำงานคำสั่ง })
			} // ปิด block การทำงานปัจจุบัน
		}() // ทำงานคำสั่ง }()

		c.Next() // ส่ง request ไป middleware หรือ handler ถัดไป
	} // ปิด block การทำงานปัจจุบัน
} // ปิด block การทำงานปัจจุบัน
