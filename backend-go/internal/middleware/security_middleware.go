package middleware // ประกาศ package middleware

import ( // เริ่มส่วน import package ที่ไฟล์นี้ต้องใช้
	"net/http" // นำเข้า package net/http
	"strings"  // นำเข้า package strings

	"github.com/gin-gonic/gin" // นำเข้า package github.com/gin-gonic/gin
) // ปิด block หรือกลุ่มรายการก่อนหน้า

func HTTPSOnly(enforce bool) gin.HandlerFunc { // ประกาศฟังก์ชัน HTTPSOnly
	return func(c *gin.Context) { // คืนค่า func(c *gin.Context)
		if !enforce || isLocalRequest(c.Request.Host) { // ตรวจเงื่อนไข !enforce || isLocalRequest(c.Request.Host)
			c.Next() // ส่ง request ไป middleware หรือ handler ถัดไป
			return   // หยุดการทำงานของฟังก์ชันนี้ทันที
		} // ปิด block การทำงานปัจจุบัน

		forwardedProto := c.GetHeader("X-Forwarded-Proto")                       // กำหนดค่า forwarded Proto จาก c.GetHeader("X-Forwarded-Proto")
		if c.Request.TLS == nil && !strings.EqualFold(forwardedProto, "https") { // ตรวจเงื่อนไข c.Request.TLS == nil && !strings.EqualFold(forwardedProto, "https")
			c.AbortWithStatusJSON(http.StatusUpgradeRequired, gin.H{"message": "HTTPS is required"}) // หยุด request ปัจจุบันไม่ให้ไป middleware ถัดไป
			return                                                                                   // หยุดการทำงานของฟังก์ชันนี้ทันที
		} // ปิด block การทำงานปัจจุบัน

		c.Next() // ส่ง request ไป middleware หรือ handler ถัดไป
	} // ปิด block การทำงานปัจจุบัน
} // ปิด block การทำงานปัจจุบัน

func SecurityHeaders() gin.HandlerFunc { // ประกาศฟังก์ชัน SecurityHeaders
	return func(c *gin.Context) { // คืนค่า func(c *gin.Context)
		c.Header("X-Content-Type-Options", "nosniff")                                             // ตั้งค่า HTTP header ให้ response
		c.Header("X-Frame-Options", "DENY")                                                       // ตั้งค่า HTTP header ให้ response
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")                            // ตั้งค่า HTTP header ให้ response
		if c.Request.TLS != nil || strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https") { // ตรวจเงื่อนไข c.Request.TLS != nil || strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https")
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains") // ตั้งค่า HTTP header ให้ response
		} // ปิด block การทำงานปัจจุบัน
		c.Next() // ส่ง request ไป middleware หรือ handler ถัดไป
	} // ปิด block การทำงานปัจจุบัน
} // ปิด block การทำงานปัจจุบัน

func isLocalRequest(host string) bool { // ประกาศฟังก์ชัน isLocalRequest
	return strings.HasPrefix(host, "localhost") || // คืนค่า strings.HasPrefix(host, "localhost") ||
		strings.HasPrefix(host, "127.0.0.1") || // ทำงานคำสั่ง strings.HasPrefix(host, "127.0.0.1") ||
		strings.HasPrefix(host, "[::1]") || // ทำงานคำสั่ง strings.HasPrefix(host, "[::1]") ||
		strings.HasPrefix(host, "::1") // ทำงานคำสั่ง strings.HasPrefix(host, "::1")
} // ปิด block การทำงานปัจจุบัน
