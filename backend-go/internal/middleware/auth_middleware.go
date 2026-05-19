package middleware // ประกาศ package middleware

import ( // เริ่มส่วน import package ที่ไฟล์นี้ต้องใช้
	"errors"   // นำเข้า package errors
	"net/http" // นำเข้า package net/http
	"strings"  // นำเข้า package strings
	"time"     // นำเข้า package time

	"github.com/gin-gonic/gin"     // นำเข้า package github.com/gin-gonic/gin
	"github.com/golang-jwt/jwt/v5" // นำเข้า package github.com/golang-jwt/jwt/v5
) // ปิด block หรือกลุ่มรายการก่อนหน้า

var jwtSecret []byte // ประกาศตัวแปร jwtSecret []byte

func SetJWTSecret(secret string) { // ประกาศฟังก์ชัน SetJWTSecret
	jwtSecret = []byte(strings.TrimSpace(secret)) // กำหนดค่า jwt Secret จาก []byte(strings.TrimSpace(secret))
} // ปิด block การทำงานปัจจุบัน

func JWTSecret() []byte { // ประกาศฟังก์ชัน JWTSecret
	return jwtSecret // คืนค่า jwtSecret
} // ปิด block การทำงานปัจจุบัน

func AuthMiddleware() gin.HandlerFunc { // ประกาศฟังก์ชัน AuthMiddleware
	return func(c *gin.Context) { // คืนค่า func(c *gin.Context)
		authHeader := c.GetHeader("Authorization") // กำหนดค่า auth Header จาก c.GetHeader("Authorization")
		if authHeader == "" {                      // ตรวจเงื่อนไข authHeader == ""
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is required"}) // ส่ง JSON response กลับไปยัง client
			c.Abort()                                                                             // หยุด request ปัจจุบันไม่ให้ไป middleware ถัดไป
			return                                                                                // หยุดการทำงานของฟังก์ชันนี้ทันที
		} // ปิด block การทำงานปัจจุบัน

		if !strings.HasPrefix(authHeader, "Bearer ") { // ตรวจเงื่อนไข !strings.HasPrefix(authHeader, "Bearer ")
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization scheme"}) // ส่ง JSON response กลับไปยัง client
			c.Abort()                                                                         // หยุด request ปัจจุบันไม่ให้ไป middleware ถัดไป
			return                                                                            // หยุดการทำงานของฟังก์ชันนี้ทันที
		} // ปิด block การทำงานปัจจุบัน

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer ")) // กำหนดค่า token String จาก strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if tokenString == "" || len(jwtSecret) == 0 {                               // ตรวจเงื่อนไข tokenString == "" || len(jwtSecret) == 0
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"}) // ส่ง JSON response กลับไปยัง client
			c.Abort()                                                          // หยุด request ปัจจุบันไม่ให้ไป middleware ถัดไป
			return                                                             // หยุดการทำงานของฟังก์ชันนี้ทันที
		} // ปิด block การทำงานปัจจุบัน

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { // กำหนดค่า token, err จาก jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // ตรวจเงื่อนไข _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok
				return nil, jwt.ErrSignatureInvalid // คืนค่า nil, jwt.ErrSignatureInvalid
			} // ปิด block การทำงานปัจจุบัน
			return jwtSecret, nil // คืนค่า jwtSecret, nil
		}) // ทำงานคำสั่ง })
		if err != nil || !token.Valid { // ตรวจเงื่อนไข err != nil || !token.Valid
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"}) // ส่ง JSON response กลับไปยัง client
			c.Abort()                                                          // หยุด request ปัจจุบันไม่ให้ไป middleware ถัดไป
			return                                                             // หยุดการทำงานของฟังก์ชันนี้ทันที
		} // ปิด block การทำงานปัจจุบัน

		claims, ok := token.Claims.(jwt.MapClaims)                        // กำหนดค่า claims, ok จาก token.Claims.(jwt.MapClaims)
		if !ok || !isAccessToken(claims) || !hasValidExpiration(claims) { // ตรวจเงื่อนไข !ok || !isAccessToken(claims) || !hasValidExpiration(claims)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"}) // ส่ง JSON response กลับไปยัง client
			c.Abort()                                                          // หยุด request ปัจจุบันไม่ให้ไป middleware ถัดไป
			return                                                             // หยุดการทำงานของฟังก์ชันนี้ทันที
		} // ปิด block การทำงานปัจจุบัน

		userID, err := getUintClaim(claims, "user_id") // กำหนดค่า user ID, err จาก getUintClaim(claims, "user_id")
		if err != nil {                                // ตรวจเงื่อนไข err != nil
			userID, err = getUintClaim(claims, "userID") // กำหนดค่า user ID, err จาก getUintClaim(claims, "userID")
		} // ปิด block การทำงานปัจจุบัน
		if err != nil { // ตรวจเงื่อนไข err != nil
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"}) // ส่ง JSON response กลับไปยัง client
			c.Abort()                                                          // หยุด request ปัจจุบันไม่ให้ไป middleware ถัดไป
			return                                                             // หยุดการทำงานของฟังก์ชันนี้ทันที
		} // ปิด block การทำงานปัจจุบัน

		role, _ := claims["role"].(string) // กำหนดค่า role,   จาก claims["role"].(string)
		c.Set("userID", userID)            // ทำงานคำสั่ง c.Set("userID", userID)
		c.Set("userId", userID)            // ทำงานคำสั่ง c.Set("userId", userID)
		c.Set("role", role)                // ทำงานคำสั่ง c.Set("role", role)
		c.Next()                           // ส่ง request ไป middleware หรือ handler ถัดไป
	} // ปิด block การทำงานปัจจุบัน
} // ปิด block การทำงานปัจจุบัน

func RequireRole(roles ...string) gin.HandlerFunc { // ประกาศฟังก์ชัน RequireRole
	allowed := map[string]bool{} // กำหนดค่า allowed จาก map[string]bool{}
	for _, role := range roles { // วนลูปตาม _, role := range roles
		allowed[role] = true // กำหนดค่า allowed[role] จาก true
	} // ปิด block การทำงานปัจจุบัน

	return func(c *gin.Context) { // คืนค่า func(c *gin.Context)
		role, _ := c.Get("role")      // กำหนดค่า role,   จาก c.Get("role")
		roleValue, _ := role.(string) // กำหนดค่า role Value,   จาก role.(string)
		if !allowed[roleValue] {      // ตรวจเงื่อนไข !allowed[roleValue]
			c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden"}) // ส่ง JSON response กลับไปยัง client
			c.Abort()                                                   // หยุด request ปัจจุบันไม่ให้ไป middleware ถัดไป
			return                                                      // หยุดการทำงานของฟังก์ชันนี้ทันที
		} // ปิด block การทำงานปัจจุบัน
		c.Next() // ส่ง request ไป middleware หรือ handler ถัดไป
	} // ปิด block การทำงานปัจจุบัน
} // ปิด block การทำงานปัจจุบัน

func isAccessToken(claims jwt.MapClaims) bool { // ประกาศฟังก์ชัน isAccessToken
	tokenType, _ := claims["type"].(string) // กำหนดค่า token Type,   จาก claims["type"].(string)
	return tokenType == "access"            // คืนค่า tokenType == "access"
} // ปิด block การทำงานปัจจุบัน

func hasValidExpiration(claims jwt.MapClaims) bool { // ประกาศฟังก์ชัน hasValidExpiration
	exp, err := getUnixClaim(claims, "exp") // กำหนดค่า exp, err จาก getUnixClaim(claims, "exp")
	if err != nil {                         // ตรวจเงื่อนไข err != nil
		return false // คืนค่า false
	} // ปิด block การทำงานปัจจุบัน
	return time.Now().Unix() < exp // คืนค่า time.Now().Unix() < exp
} // ปิด block การทำงานปัจจุบัน

func getUnixClaim(claims jwt.MapClaims, key string) (int64, error) { // ประกาศฟังก์ชัน getUnixClaim
	value, ok := claims[key] // กำหนดค่า value, ok จาก claims[key]
	if !ok {                 // ตรวจเงื่อนไข !ok
		return 0, errors.New("missing claim") // คืนค่า 0, errors.New("missing claim")
	} // ปิด block การทำงานปัจจุบัน

	switch v := value.(type) { // เลือก case จาก v := value.(type)
	case float64: // กรณี float64
		return int64(v), nil // คืนค่า int64(v), nil
	case int64: // กรณี int64
		return v, nil // คืนค่า v, nil
	case int: // กรณี int
		return int64(v), nil // คืนค่า int64(v), nil
	default: // กรณี default เมื่อไม่ตรง case อื่น
		return 0, errors.New("invalid claim") // คืนค่า 0, errors.New("invalid claim")
	} // ปิด block การทำงานปัจจุบัน
} // ปิด block การทำงานปัจจุบัน

func getUintClaim(claims jwt.MapClaims, key string) (uint, error) { // ประกาศฟังก์ชัน getUintClaim
	value, ok := claims[key] // กำหนดค่า value, ok จาก claims[key]
	if !ok {                 // ตรวจเงื่อนไข !ok
		return 0, errors.New("missing claim") // คืนค่า 0, errors.New("missing claim")
	} // ปิด block การทำงานปัจจุบัน

	switch v := value.(type) { // เลือก case จาก v := value.(type)
	case float64: // กรณี float64
		if v <= 0 { // ตรวจเงื่อนไข v <= 0
			return 0, errors.New("invalid claim") // คืนค่า 0, errors.New("invalid claim")
		} // ปิด block การทำงานปัจจุบัน
		return uint(v), nil // คืนค่า uint(v), nil
	case int: // กรณี int
		if v <= 0 { // ตรวจเงื่อนไข v <= 0
			return 0, errors.New("invalid claim") // คืนค่า 0, errors.New("invalid claim")
		} // ปิด block การทำงานปัจจุบัน
		return uint(v), nil // คืนค่า uint(v), nil
	default: // กรณี default เมื่อไม่ตรง case อื่น
		return 0, errors.New("invalid claim") // คืนค่า 0, errors.New("invalid claim")
	} // ปิด block การทำงานปัจจุบัน
} // ปิด block การทำงานปัจจุบัน
