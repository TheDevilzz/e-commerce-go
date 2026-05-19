package handler // ประกาศ package handler

import ( // เริ่มส่วน import package ที่ไฟล์นี้ต้องใช้
	"crypto/rand"     // นำเข้า package crypto/rand
	"crypto/sha256"   // นำเข้า package crypto/sha256
	"encoding/base64" // นำเข้า package encoding/base64
	"encoding/hex"    // นำเข้า package encoding/hex
	"net/http"        // นำเข้า package net/http
	"strings"         // นำเข้า package strings
	"sync"            // นำเข้า package sync
	"time"            // นำเข้า package time

	"backend-go/internal/middleware" // นำเข้า package backend-go/internal/middleware
	"backend-go/internal/model"      // นำเข้า package backend-go/internal/model

	"github.com/gin-gonic/gin"     // นำเข้า package github.com/gin-gonic/gin
	"github.com/golang-jwt/jwt/v5" // นำเข้า package github.com/golang-jwt/jwt/v5
	"golang.org/x/crypto/bcrypt"   // นำเข้า package golang.org/x/crypto/bcrypt
	"gorm.io/gorm"                 // นำเข้า package gorm.io/gorm
) // ปิด block หรือกลุ่มรายการก่อนหน้า

const accessTokenTTL = 15 * time.Minute // ประกาศค่าคงที่ accessTokenTTL = 15 * time.Minute

type AuthHandler struct { // ประกาศ struct AuthHandler
	DB              *gorm.DB      // ทำงานคำสั่ง DB *gorm.DB
	refreshTokenTTL time.Duration // ทำงานคำสั่ง refreshTokenTTL time.Duration
} // ปิด block การทำงานปัจจุบัน

func NewAuthHandler(db *gorm.DB, refreshTokenDays int) *AuthHandler { // ประกาศฟังก์ชัน NewAuthHandler
	if refreshTokenDays < 7 { // ตรวจเงื่อนไข refreshTokenDays < 7
		refreshTokenDays = 7 // กำหนดค่า refresh Token Days จาก 7
	} // ปิด block การทำงานปัจจุบัน
	if refreshTokenDays > 30 { // ตรวจเงื่อนไข refreshTokenDays > 30
		refreshTokenDays = 30 // กำหนดค่า refresh Token Days จาก 30
	} // ปิด block การทำงานปัจจุบัน
	return &AuthHandler{ // คืนค่า &AuthHandler
		DB:              db,                                               // ระบุค่ารายการหนึ่งในชุดข้อมูล
		refreshTokenTTL: time.Duration(refreshTokenDays) * 24 * time.Hour, // ระบุค่ารายการหนึ่งในชุดข้อมูล
	} // ปิด block การทำงานปัจจุบัน
} // ปิด block การทำงานปัจจุบัน

type RegisterRequest struct { // ประกาศ struct RegisterRequest
	Username    string `json:"username" binding:"required,min=3,max=50"` // ประกาศ field Username พร้อม tag JSON/database
	Name        string `json:"name" binding:"required,min=2,max=100"`    // ประกาศ field Name พร้อม tag JSON/database
	Email       string `json:"email" binding:"required,email"`           // ประกาศ field Email พร้อม tag JSON/database
	Phone       string `json:"phone" binding:"required,min=6,max=30"`    // ประกาศ field Phone พร้อม tag JSON/database
	DateOfBirth string `json:"date_of_birth" binding:"required"`         // ประกาศ field DateOfBirth พร้อม tag JSON/database
	Password    string `json:"password" binding:"required,min=8,max=72"` // ประกาศ field Password พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน

type LoginRequest struct { // ประกาศ struct LoginRequest
	Username string `json:"username"`                                 // ประกาศ field Username พร้อม tag JSON/database
	Email    string `json:"email"`                                    // ประกาศ field Email พร้อม tag JSON/database
	Password string `json:"password" binding:"required,min=8,max=72"` // ประกาศ field Password พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน

type RefreshRequest struct { // ประกาศ struct RefreshRequest
	RefreshToken string `json:"refresh_token" binding:"required"` // ประกาศ field RefreshToken พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน

type LogoutRequest struct { // ประกาศ struct LogoutRequest
	RefreshToken string `json:"refresh_token" binding:"required"` // ประกาศ field RefreshToken พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน

type loginAttempt struct { // ประกาศ struct loginAttempt
	Count     int       // ทำงานคำสั่ง Count int
	ResetTime time.Time // ทำงานคำสั่ง ResetTime time.Time
} // ปิด block การทำงานปัจจุบัน

var ( // ประกาศตัวแปร (
	loginMu       sync.Mutex                  // ทำงานคำสั่ง loginMu sync.Mutex
	loginAttempts = map[string]loginAttempt{} // กำหนดค่า login Attempts จาก map[string]loginAttempt{}
) // ปิด block หรือกลุ่มรายการก่อนหน้า

func (h *AuthHandler) Register(c *gin.Context) { // ประกาศฟังก์ชัน Register
	var req RegisterRequest                        // ประกาศตัวแปร req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil { // ตรวจเงื่อนไข err := c.ShouldBindJSON(&req); err != nil
		BadRequest(c, "invalid request") // ทำงานคำสั่ง BadRequest(c, "invalid request")
		return                           // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	req.Username = strings.TrimSpace(req.Username)            // กำหนดค่า Username ของ request จาก strings.TrimSpace(req.Username)
	req.Name = strings.TrimSpace(req.Name)                    // กำหนดค่า Name ของ request จาก strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email)) // กำหนดค่า Email ของ request จาก strings.ToLower(strings.TrimSpace(req.Email))
	req.Phone = strings.TrimSpace(req.Phone)                  // กำหนดค่า Phone ของ request จาก strings.TrimSpace(req.Phone)
	req.DateOfBirth = strings.TrimSpace(req.DateOfBirth)      // กำหนดค่า Date Of Birth ของ request จาก strings.TrimSpace(req.DateOfBirth)

	if _, err := time.Parse("2006-01-02", req.DateOfBirth); err != nil { // ตรวจเงื่อนไข _, err := time.Parse("2006-01-02", req.DateOfBirth); err != nil
		BadRequest(c, "date_of_birth must use YYYY-MM-DD format") // ทำงานคำสั่ง BadRequest(c, "date_of_birth must use YYYY-MM-DD format")
		return                                                    // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost) // กำหนดค่า hash Password, err จาก bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {                                                                            // ตรวจเงื่อนไข err != nil
		Error(c, http.StatusInternalServerError, "failed to process password") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to process password")
		return                                                                 // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	if err := h.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&model.User{}).Error; err == nil { // ตรวจเงื่อนไข err := h.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&model.User{}).Error; err == nil
		BadRequest(c, "username or email already exists") // ทำงานคำสั่ง BadRequest(c, "username or email already exists")
		return                                            // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	user := model.User{ // กำหนดค่า user จาก model.User
		Username:    req.Username,         // ระบุค่ารายการหนึ่งในชุดข้อมูล
		Name:        req.Name,             // ระบุค่ารายการหนึ่งในชุดข้อมูล
		Email:       req.Email,            // ระบุค่ารายการหนึ่งในชุดข้อมูล
		Phone:       req.Phone,            // ระบุค่ารายการหนึ่งในชุดข้อมูล
		DateOfBirth: req.DateOfBirth,      // ระบุค่ารายการหนึ่งในชุดข้อมูล
		Password:    string(hashPassword), // ระบุค่ารายการหนึ่งในชุดข้อมูล
		Role:        "user",               // ระบุค่ารายการหนึ่งในชุดข้อมูล
	} // ปิด block การทำงานปัจจุบัน

	if err := h.DB.Create(&user).Error; err != nil { // ตรวจเงื่อนไข err := h.DB.Create(&user).Error; err != nil
		BadRequest(c, "failed to create user") // ทำงานคำสั่ง BadRequest(c, "failed to create user")
		return                                 // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	Success(c, http.StatusCreated, "registration successful", userPayload(user)) // ทำงานคำสั่ง Success(c, http.StatusCreated, "registration successful", userPayload(user))
} // ปิด block การทำงานปัจจุบัน

func (h *AuthHandler) Login(c *gin.Context) { // ประกาศฟังก์ชัน Login
	var req LoginRequest                           // ประกาศตัวแปร req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil { // ตรวจเงื่อนไข err := c.ShouldBindJSON(&req); err != nil
		BadRequest(c, "invalid request") // ทำงานคำสั่ง BadRequest(c, "invalid request")
		return                           // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	key := c.ClientIP() + ":" + strings.ToLower(strings.TrimSpace(req.Email)) + ":" + strings.ToLower(strings.TrimSpace(req.Username)) // กำหนดค่า key จาก c.ClientIP() + ":" + strings.ToLower(strings.TrimSpace(req.Email)) + ":" + strings.ToLower(strings.TrimSpace(req.Username))
	if isRateLimited(key) {                                                                                                            // ตรวจเงื่อนไข isRateLimited(key)
		Error(c, http.StatusTooManyRequests, "too many login attempts, please try again later") // ทำงานคำสั่ง Error(c, http.StatusTooManyRequests, "too many login attempts, please try again later")
		return                                                                                  // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	req.Email = strings.ToLower(strings.TrimSpace(req.Email)) // กำหนดค่า Email ของ request จาก strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.TrimSpace(req.Username)            // กำหนดค่า Username ของ request จาก strings.TrimSpace(req.Username)
	if req.Email == "" && req.Username == "" {                // ตรวจเงื่อนไข req.Email == "" && req.Username == ""
		BadRequest(c, "email or username is required") // ทำงานคำสั่ง BadRequest(c, "email or username is required")
		return                                         // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	var user model.User                                                                                         // ประกาศตัวแปร user model.User
	if err := h.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&user).Error; err != nil { // ตรวจเงื่อนไข err := h.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&user).Error; err != nil
		recordFailedLogin(key)                                                    // ทำงานคำสั่ง recordFailedLogin(key)
		Error(c, http.StatusUnauthorized, "invalid username, email, or password") // ทำงานคำสั่ง Error(c, http.StatusUnauthorized, "invalid username, email, or password")
		return                                                                    // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil { // ตรวจเงื่อนไข err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil
		recordFailedLogin(key)                                                    // ทำงานคำสั่ง recordFailedLogin(key)
		Error(c, http.StatusUnauthorized, "invalid username, email, or password") // ทำงานคำสั่ง Error(c, http.StatusUnauthorized, "invalid username, email, or password")
		return                                                                    // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	clearLoginAttempts(key) // ทำงานคำสั่ง clearLoginAttempts(key)

	accessToken, refreshToken, expiresAt, err := h.issueTokenPair(user) // กำหนดค่า access Token, refresh Token, expires At, err จาก h.issueTokenPair(user)
	if err != nil {                                                     // ตรวจเงื่อนไข err != nil
		Error(c, http.StatusInternalServerError, "failed to create token") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to create token")
		return                                                             // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	Success(c, http.StatusOK, "login successful", gin.H{ // ทำงานคำสั่ง Success(c, http.StatusOK, "login successful", gin.H{
		"token":         accessToken,                   // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"access_token":  accessToken,                   // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"refresh_token": refreshToken,                  // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"expires_in":    int(accessTokenTTL.Seconds()), // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"expires_at":    expiresAt,                     // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"user":          userPayload(user),             // ระบุค่ารายการหนึ่งในชุดข้อมูล
	}) // ทำงานคำสั่ง })
} // ปิด block การทำงานปัจจุบัน

func (h *AuthHandler) Refresh(c *gin.Context) { // ประกาศฟังก์ชัน Refresh
	var req RefreshRequest                         // ประกาศตัวแปร req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil { // ตรวจเงื่อนไข err := c.ShouldBindJSON(&req); err != nil
		BadRequest(c, "invalid request") // ทำงานคำสั่ง BadRequest(c, "invalid request")
		return                           // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	tokenHash := hashRefreshToken(strings.TrimSpace(req.RefreshToken))                                                                              // กำหนดค่า token Hash จาก hashRefreshToken(strings.TrimSpace(req.RefreshToken))
	var storedToken model.RefreshToken                                                                                                              // ประกาศตัวแปร storedToken model.RefreshToken
	if err := h.DB.Where("token_hash = ? AND revoked_at IS NULL AND expires_at > ?", tokenHash, time.Now()).First(&storedToken).Error; err != nil { // ตรวจเงื่อนไข err := h.DB.Where("token_hash = ? AND revoked_at IS NULL AND expires_at > ?", tokenHash, time.Now()).First(&storedToken).Error; err != nil
		Error(c, http.StatusUnauthorized, "invalid refresh token") // ทำงานคำสั่ง Error(c, http.StatusUnauthorized, "invalid refresh token")
		return                                                     // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	var user model.User                                                 // ประกาศตัวแปร user model.User
	if err := h.DB.First(&user, storedToken.UserID).Error; err != nil { // ตรวจเงื่อนไข err := h.DB.First(&user, storedToken.UserID).Error; err != nil
		Error(c, http.StatusUnauthorized, "invalid refresh token") // ทำงานคำสั่ง Error(c, http.StatusUnauthorized, "invalid refresh token")
		return                                                     // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	now := time.Now()                                                                 // กำหนดค่า now จาก time.Now()
	if err := h.DB.Model(&storedToken).Update("revoked_at", &now).Error; err != nil { // ตรวจเงื่อนไข err := h.DB.Model(&storedToken).Update("revoked_at", &now).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to rotate refresh token") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to rotate refresh token")
		return                                                                     // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	accessToken, refreshToken, expiresAt, err := h.issueTokenPair(user) // กำหนดค่า access Token, refresh Token, expires At, err จาก h.issueTokenPair(user)
	if err != nil {                                                     // ตรวจเงื่อนไข err != nil
		Error(c, http.StatusInternalServerError, "failed to create token") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to create token")
		return                                                             // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	Success(c, http.StatusOK, "token refreshed", gin.H{ // ทำงานคำสั่ง Success(c, http.StatusOK, "token refreshed", gin.H{
		"token":         accessToken,                   // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"access_token":  accessToken,                   // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"refresh_token": refreshToken,                  // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"expires_in":    int(accessTokenTTL.Seconds()), // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"expires_at":    expiresAt,                     // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"user":          userPayload(user),             // ระบุค่ารายการหนึ่งในชุดข้อมูล
	}) // ทำงานคำสั่ง })
} // ปิด block การทำงานปัจจุบัน

func (h *AuthHandler) Logout(c *gin.Context) { // ประกาศฟังก์ชัน Logout
	var req LogoutRequest                          // ประกาศตัวแปร req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil { // ตรวจเงื่อนไข err := c.ShouldBindJSON(&req); err != nil
		BadRequest(c, "invalid request") // ทำงานคำสั่ง BadRequest(c, "invalid request")
		return                           // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	now := time.Now()                                                  // กำหนดค่า now จาก time.Now()
	tokenHash := hashRefreshToken(strings.TrimSpace(req.RefreshToken)) // กำหนดค่า token Hash จาก hashRefreshToken(strings.TrimSpace(req.RefreshToken))
	h.DB.Model(&model.RefreshToken{}).                                 // ทำงานคำสั่ง h.DB.Model(&model.RefreshToken{}).
										Where("token_hash = ? AND revoked_at IS NULL", tokenHash). // กำหนดค่า Where("token hash จาก ? AND revoked_at IS NULL", tokenHash).
										Update("revoked_at", &now)                                 // ทำงานคำสั่ง Update("revoked_at", &now)

	Success(c, http.StatusOK, "logged out", nil) // ทำงานคำสั่ง Success(c, http.StatusOK, "logged out", nil)
} // ปิด block การทำงานปัจจุบัน

func (h *AuthHandler) issueTokenPair(user model.User) (string, string, time.Time, error) { // ประกาศฟังก์ชัน issueTokenPair
	accessToken, expiresAt, err := createAccessToken(user) // กำหนดค่า access Token, expires At, err จาก createAccessToken(user)
	if err != nil {                                        // ตรวจเงื่อนไข err != nil
		return "", "", time.Time{}, err // คืนค่า "", "", time.Time{}, err
	} // ปิด block การทำงานปัจจุบัน

	refreshToken, err := randomRefreshToken() // กำหนดค่า refresh Token, err จาก randomRefreshToken()
	if err != nil {                           // ตรวจเงื่อนไข err != nil
		return "", "", time.Time{}, err // คืนค่า "", "", time.Time{}, err
	} // ปิด block การทำงานปัจจุบัน

	refreshExpiresAt := time.Now().Add(h.refreshTokenTTL) // กำหนดค่า refresh Expires At จาก time.Now().Add(h.refreshTokenTTL)
	if err := h.DB.Create(&model.RefreshToken{            // ตรวจเงื่อนไข err := h.DB.Create(&model.RefreshToken
		UserID:    user.ID,                        // ระบุค่ารายการหนึ่งในชุดข้อมูล
		TokenHash: hashRefreshToken(refreshToken), // ระบุค่ารายการหนึ่งในชุดข้อมูล
		ExpiresAt: refreshExpiresAt,               // ระบุค่ารายการหนึ่งในชุดข้อมูล
	}).Error; err != nil { // ทำงานคำสั่ง }).Error; err != nil {
		return "", "", time.Time{}, err // คืนค่า "", "", time.Time{}, err
	} // ปิด block การทำงานปัจจุบัน

	return accessToken, refreshToken, expiresAt, nil // คืนค่า accessToken, refreshToken, expiresAt, nil
} // ปิด block การทำงานปัจจุบัน

func createAccessToken(user model.User) (string, time.Time, error) { // ประกาศฟังก์ชัน createAccessToken
	expiresAt := time.Now().Add(accessTokenTTL)                       // กำหนดค่า expires At จาก time.Now().Add(accessTokenTTL)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ // กำหนดค่า token จาก jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims
		"user_id": user.ID,           // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"role":    user.Role,         // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"type":    "access",          // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"exp":     expiresAt.Unix(),  // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"iat":     time.Now().Unix(), // ระบุค่ารายการหนึ่งในชุดข้อมูล
	}) // ทำงานคำสั่ง })

	tokenString, err := token.SignedString(middleware.JWTSecret()) // กำหนดค่า token String, err จาก token.SignedString(middleware.JWTSecret())
	if err != nil {                                                // ตรวจเงื่อนไข err != nil
		return "", time.Time{}, err // คืนค่า "", time.Time{}, err
	} // ปิด block การทำงานปัจจุบัน
	return tokenString, expiresAt, nil // คืนค่า tokenString, expiresAt, nil
} // ปิด block การทำงานปัจจุบัน

func randomRefreshToken() (string, error) { // ประกาศฟังก์ชัน randomRefreshToken
	bytes := make([]byte, 32)                   // กำหนดค่า bytes จาก make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil { // ตรวจเงื่อนไข _, err := rand.Read(bytes); err != nil
		return "", err // คืนค่า "", err
	} // ปิด block การทำงานปัจจุบัน
	return base64.RawURLEncoding.EncodeToString(bytes), nil // คืนค่า base64.RawURLEncoding.EncodeToString(bytes), nil
} // ปิด block การทำงานปัจจุบัน

func hashRefreshToken(token string) string { // ประกาศฟังก์ชัน hashRefreshToken
	sum := sha256.Sum256([]byte(token)) // กำหนดค่า sum จาก sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])   // คืนค่า hex.EncodeToString(sum[:])
} // ปิด block การทำงานปัจจุบัน

func userPayload(user model.User) gin.H { // ประกาศฟังก์ชัน userPayload
	return gin.H{ // คืนค่า gin.H
		"id":            user.ID,          // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"username":      user.Username,    // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"email":         user.Email,       // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"name":          user.Name,        // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"phone":         user.Phone,       // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"date_of_birth": user.DateOfBirth, // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"role":          user.Role,        // ระบุค่ารายการหนึ่งในชุดข้อมูล
	} // ปิด block การทำงานปัจจุบัน
} // ปิด block การทำงานปัจจุบัน

func isRateLimited(key string) bool { // ประกาศฟังก์ชัน isRateLimited
	loginMu.Lock()         // ทำงานคำสั่ง loginMu.Lock()
	defer loginMu.Unlock() // เลื่อนคำสั่งนี้ไปทำตอนฟังก์ชันจบ: loginMu.Unlock()

	attempt := loginAttempts[key]            // กำหนดค่า attempt จาก loginAttempts[key]
	if time.Now().After(attempt.ResetTime) { // ตรวจเงื่อนไข time.Now().After(attempt.ResetTime)
		delete(loginAttempts, key) // ทำงานคำสั่ง delete(loginAttempts, key)
		return false               // คืนค่า false
	} // ปิด block การทำงานปัจจุบัน
	return attempt.Count >= 5 // คืนค่า attempt.Count >= 5
} // ปิด block การทำงานปัจจุบัน

func recordFailedLogin(key string) { // ประกาศฟังก์ชัน recordFailedLogin
	loginMu.Lock()         // ทำงานคำสั่ง loginMu.Lock()
	defer loginMu.Unlock() // เลื่อนคำสั่งนี้ไปทำตอนฟังก์ชันจบ: loginMu.Unlock()

	attempt := loginAttempts[key]            // กำหนดค่า attempt จาก loginAttempts[key]
	if time.Now().After(attempt.ResetTime) { // ตรวจเงื่อนไข time.Now().After(attempt.ResetTime)
		attempt = loginAttempt{ResetTime: time.Now().Add(15 * time.Minute)} // กำหนดค่า attempt จาก loginAttempt{ResetTime: time.Now().Add(15 * time.Minute)}
	} // ปิด block การทำงานปัจจุบัน
	attempt.Count++              // ทำงานคำสั่ง attempt.Count++
	loginAttempts[key] = attempt // กำหนดค่า login Attempts[key] จาก attempt
} // ปิด block การทำงานปัจจุบัน

func clearLoginAttempts(key string) { // ประกาศฟังก์ชัน clearLoginAttempts
	loginMu.Lock()             // ทำงานคำสั่ง loginMu.Lock()
	defer loginMu.Unlock()     // เลื่อนคำสั่งนี้ไปทำตอนฟังก์ชันจบ: loginMu.Unlock()
	delete(loginAttempts, key) // ทำงานคำสั่ง delete(loginAttempts, key)
} // ปิด block การทำงานปัจจุบัน
