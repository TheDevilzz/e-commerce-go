package middleware // ประกาศ package middleware

import ( // เริ่มส่วน import package ที่ไฟล์นี้ต้องใช้
	"net/http"          // นำเข้า package net/http
	"net/http/httptest" // นำเข้า package net/http/httptest
	"testing"           // นำเข้า package testing
	"time"              // นำเข้า package time

	"github.com/gin-gonic/gin"     // นำเข้า package github.com/gin-gonic/gin
	"github.com/golang-jwt/jwt/v5" // นำเข้า package github.com/golang-jwt/jwt/v5
) // ปิด block หรือกลุ่มรายการก่อนหน้า

func TestAuthMiddlewareRejectsMissingToken(t *testing.T) { // ประกาศฟังก์ชัน TestAuthMiddlewareRejectsMissingToken
	router := testRouter(AuthMiddleware())                       // กำหนดค่า router จาก testRouter(AuthMiddleware())
	response := httptest.NewRecorder()                           // กำหนดค่า response จาก httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/test", nil) // กำหนดค่า request จาก httptest.NewRequest(http.MethodGet, "/test", nil)

	router.ServeHTTP(response, request) // ทำงานคำสั่ง router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized { // ตรวจเงื่อนไข response.Code != http.StatusUnauthorized
		t.Fatalf("expected status 401, got %d", response.Code) // ทำงานคำสั่ง t.Fatalf("expected status 401, got %d", response.Code)
	} // ปิด block การทำงานปัจจุบัน
} // ปิด block การทำงานปัจจุบัน

func TestAuthMiddlewareRejectsInvalidToken(t *testing.T) { // ประกาศฟังก์ชัน TestAuthMiddlewareRejectsInvalidToken
	SetJWTSecret("test-secret")                                  // ทำงานคำสั่ง SetJWTSecret("test-secret")
	router := testRouter(AuthMiddleware())                       // กำหนดค่า router จาก testRouter(AuthMiddleware())
	response := httptest.NewRecorder()                           // กำหนดค่า response จาก httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/test", nil) // กำหนดค่า request จาก httptest.NewRequest(http.MethodGet, "/test", nil)
	request.Header.Set("Authorization", "Bearer invalid")        // ทำงานคำสั่ง request.Header.Set("Authorization", "Bearer invalid")

	router.ServeHTTP(response, request) // ทำงานคำสั่ง router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized { // ตรวจเงื่อนไข response.Code != http.StatusUnauthorized
		t.Fatalf("expected status 401, got %d", response.Code) // ทำงานคำสั่ง t.Fatalf("expected status 401, got %d", response.Code)
	} // ปิด block การทำงานปัจจุบัน
} // ปิด block การทำงานปัจจุบัน

func TestAuthMiddlewareAcceptsValidToken(t *testing.T) { // ประกาศฟังก์ชัน TestAuthMiddlewareAcceptsValidToken
	SetJWTSecret("test-secret")                                            // ทำงานคำสั่ง SetJWTSecret("test-secret")
	router := testRouter(AuthMiddleware())                                 // กำหนดค่า router จาก testRouter(AuthMiddleware())
	response := httptest.NewRecorder()                                     // กำหนดค่า response จาก httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/test", nil)           // กำหนดค่า request จาก httptest.NewRequest(http.MethodGet, "/test", nil)
	request.Header.Set("Authorization", "Bearer "+signedToken(t, "admin")) // ทำงานคำสั่ง request.Header.Set("Authorization", "Bearer "+signedToken(t, "admin"))

	router.ServeHTTP(response, request) // ทำงานคำสั่ง router.ServeHTTP(response, request)

	if response.Code != http.StatusOK { // ตรวจเงื่อนไข response.Code != http.StatusOK
		t.Fatalf("expected status 200, got %d", response.Code) // ทำงานคำสั่ง t.Fatalf("expected status 200, got %d", response.Code)
	} // ปิด block การทำงานปัจจุบัน
} // ปิด block การทำงานปัจจุบัน

func TestRequireRoleRejectsWrongRole(t *testing.T) { // ประกาศฟังก์ชัน TestRequireRoleRejectsWrongRole
	SetJWTSecret("test-secret")                                           // ทำงานคำสั่ง SetJWTSecret("test-secret")
	router := testRouter(AuthMiddleware(), RequireRole("admin"))          // กำหนดค่า router จาก testRouter(AuthMiddleware(), RequireRole("admin"))
	response := httptest.NewRecorder()                                    // กำหนดค่า response จาก httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/test", nil)          // กำหนดค่า request จาก httptest.NewRequest(http.MethodGet, "/test", nil)
	request.Header.Set("Authorization", "Bearer "+signedToken(t, "user")) // ทำงานคำสั่ง request.Header.Set("Authorization", "Bearer "+signedToken(t, "user"))

	router.ServeHTTP(response, request) // ทำงานคำสั่ง router.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden { // ตรวจเงื่อนไข response.Code != http.StatusForbidden
		t.Fatalf("expected status 403, got %d", response.Code) // ทำงานคำสั่ง t.Fatalf("expected status 403, got %d", response.Code)
	} // ปิด block การทำงานปัจจุบัน
} // ปิด block การทำงานปัจจุบัน

func testRouter(middleware ...gin.HandlerFunc) *gin.Engine { // ประกาศฟังก์ชัน testRouter
	gin.SetMode(gin.TestMode)                  // ทำงานคำสั่ง gin.SetMode(gin.TestMode)
	router := gin.New()                        // กำหนดค่า router จาก gin.New()
	router.Use(middleware...)                  // ติดตั้ง middleware ให้ router/group
	router.GET("/test", func(c *gin.Context) { // ลงทะเบียน endpoint GET
		c.JSON(http.StatusOK, gin.H{"message": "ok"}) // ส่ง JSON response กลับไปยัง client
	}) // ทำงานคำสั่ง })
	return router // คืนค่า router
} // ปิด block การทำงานปัจจุบัน

func signedToken(t *testing.T, role string) string { // ประกาศฟังก์ชัน signedToken
	t.Helper() // ทำงานคำสั่ง t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ // กำหนดค่า token จาก jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims
		"user_id": 1,                                // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"role":    role,                             // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"type":    "access",                         // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"exp":     time.Now().Add(time.Hour).Unix(), // ระบุค่ารายการหนึ่งในชุดข้อมูล
	}) // ทำงานคำสั่ง })

	tokenString, err := token.SignedString(JWTSecret()) // กำหนดค่า token String, err จาก token.SignedString(JWTSecret())
	if err != nil {                                     // ตรวจเงื่อนไข err != nil
		t.Fatalf("failed to sign token: %v", err) // ทำงานคำสั่ง t.Fatalf("failed to sign token: %v", err)
	} // ปิด block การทำงานปัจจุบัน
	return tokenString // คืนค่า tokenString
} // ปิด block การทำงานปัจจุบัน
