package handler // ประกาศ package handler

import ( // เริ่มส่วน import package ที่ไฟล์นี้ต้องใช้
	"backend-go/internal/model" // นำเข้า package backend-go/internal/model
	"bytes"                     // นำเข้า package bytes
	"encoding/json"             // นำเข้า package encoding/json
	"fmt"                       // นำเข้า package fmt
	"image"                     // นำเข้า package image
	_ "image/jpeg"              // ทำงานคำสั่ง _ "image/jpeg"
	_ "image/png"               // ทำงานคำสั่ง _ "image/png"
	"io"                        // นำเข้า package io
	"mime/multipart"            // นำเข้า package mime/multipart
	"net/http"                  // นำเข้า package net/http
	"os"                        // นำเข้า package os
	"path/filepath"             // นำเข้า package path/filepath
	"time"                      // นำเข้า package time

	"github.com/gin-gonic/gin"             // นำเข้า package github.com/gin-gonic/gin
	"github.com/makiuchi-d/gozxing"        // นำเข้า package github.com/makiuchi-d/gozxing
	"github.com/makiuchi-d/gozxing/qrcode" // นำเข้า package github.com/makiuchi-d/gozxing/qrcode
	"gorm.io/gorm"                         // นำเข้า package gorm.io/gorm
) // ปิด block หรือกลุ่มรายการก่อนหน้า

type PaymentHandler struct { // ประกาศ struct PaymentHandler
	db *gorm.DB // ทำงานคำสั่ง db *gorm.DB
} // ปิด block การทำงานปัจจุบัน

func NewPaymentHandler(db *gorm.DB) *PaymentHandler { // ประกาศฟังก์ชัน NewPaymentHandler
	return &PaymentHandler{db: db} // คืนค่า &PaymentHandler{db: db}
} // ปิด block การทำงานปัจจุบัน

type UploadSlipRequest struct { // ประกาศ struct UploadSlipRequest
	OrderID uint    `form:"order_id" binding:"required"`    // ทำงานคำสั่ง OrderID uint `form:"order_id" binding:"required"`
	Amount  float64 `form:"amount" binding:"required,gt=0"` // ทำงานคำสั่ง Amount float64 `form:"amount" binding:"required,gt=0"`
} // ปิด block การทำงานปัจจุบัน

type VerifyResponse struct { // ประกาศ struct VerifyResponse
	Status string     `json:"status"` // ประกาศ field Status พร้อม tag JSON/database
	Data   VerifyData `json:"data"`   // ประกาศ field Data พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน

type VerifyData struct { // ประกาศ struct VerifyData
	TransID  string       `json:"trans_id"`  // ประกาศ field TransID พร้อม tag JSON/database
	Type     string       `json:"type"`      // ประกาศ field Type พร้อม tag JSON/database
	Amount   float64      `json:"amount"`    // ประกาศ field Amount พร้อม tag JSON/database
	DateTime string       `json:"date_time"` // ประกาศ field DateTime พร้อม tag JSON/database
	Sender   VerifySender `json:"sender"`    // ประกาศ field Sender พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน

type VerifySender struct { // ประกาศ struct VerifySender
	Name      string `json:"name"`       // ประกาศ field Name พร้อม tag JSON/database
	Bank      string `json:"bank"`       // ประกาศ field Bank พร้อม tag JSON/database
	AccountNo string `json:"account_no"` // ประกาศ field AccountNo พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน

func (h *PaymentHandler) UploadSlip(c *gin.Context) { // ประกาศฟังก์ชัน UploadSlip
	var req UploadSlipRequest                  // ประกาศตัวแปร req UploadSlipRequest
	userId := c.GetUint("userId")              // กำหนดค่า user Id จาก c.GetUint("userId")
	if err := c.ShouldBind(&req); err != nil { // ตรวจเงื่อนไข err := c.ShouldBind(&req); err != nil
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()}) // ส่ง JSON response กลับไปยัง client
		return                                                       // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	var order model.Order                                                                                 // ประกาศตัวแปร order model.Order
	if err := h.db.Where("id = ? AND user_id = ?", req.OrderID, userId).First(&order).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Where("id = ? AND user_id = ?", req.OrderID, userId).First(&order).Error; err != nil
		c.JSON(http.StatusNotFound, gin.H{"message": "order not found"}) // ส่ง JSON response กลับไปยัง client
		return                                                           // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	if order.Status == "paid" { // ตรวจเงื่อนไข order.Status == "paid"
		c.JSON(http.StatusBadRequest, gin.H{"message": "order is already paid"}) // ส่ง JSON response กลับไปยัง client
		return                                                                   // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	if fmt.Sprintf("%.2f", order.TotalPrice) != fmt.Sprintf("%.2f", req.Amount) { // ตรวจเงื่อนไข fmt.Sprintf("%.2f", order.TotalPrice) != fmt.Sprintf("%.2f", req.Amount)
		c.JSON(http.StatusBadRequest, gin.H{"message": "amount does not match order total"}) // ส่ง JSON response กลับไปยัง client
		return                                                                               // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	file, err := c.FormFile("image") // กำหนดค่า file, err จาก c.FormFile("image")
	if err != nil {                  // ตรวจเงื่อนไข err != nil
		c.JSON(http.StatusBadRequest, gin.H{"message": "image is required"}) // ส่ง JSON response กลับไปยัง client
		return                                                               // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	path, data, err := h.ReadQRCode(file, req.Amount) // กำหนดค่า path, data, err จาก h.ReadQRCode(file, req.Amount)
	if err != nil {                                   // ตรวจเงื่อนไข err != nil
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()}) // ส่ง JSON response กลับไปยัง client
		return                                                       // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	var existingPayment model.Payment                                                              // ประกาศตัวแปร existingPayment model.Payment
	if err := h.db.Where("ref = ?", data.Data.TransID).First(&existingPayment).Error; err == nil { // ตรวจเงื่อนไข err := h.db.Where("ref = ?", data.Data.TransID).First(&existingPayment).Error; err == nil
		c.JSON(http.StatusBadRequest, gin.H{"message": "this slip has already been used"}) // ส่ง JSON response กลับไปยัง client
		return                                                                             // หยุดการทำงานของฟังก์ชันนี้ทันที
	} else if err != gorm.ErrRecordNotFound { // ทำงานคำสั่ง } else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to check payment"}) // ส่ง JSON response กลับไปยัง client
		return                                                                              // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	payment := model.Payment{ // กำหนดค่า payment จาก model.Payment
		UserID:  userId,            // ระบุค่ารายการหนึ่งในชุดข้อมูล
		OrderID: req.OrderID,       // ระบุค่ารายการหนึ่งในชุดข้อมูล
		Amount:  req.Amount,        // ระบุค่ารายการหนึ่งในชุดข้อมูล
		SlipURL: "/" + path,        // ระบุค่ารายการหนึ่งในชุดข้อมูล
		Status:  "verified",        // ระบุค่ารายการหนึ่งในชุดข้อมูล
		Ref:     data.Data.TransID, // ระบุค่ารายการหนึ่งในชุดข้อมูล
	} // ปิด block การทำงานปัจจุบัน

	tx := h.db.Begin()   // กำหนดค่า tx จาก h.db.Begin()
	if tx.Error != nil { // ตรวจเงื่อนไข tx.Error != nil
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to start transaction"}) // ส่ง JSON response กลับไปยัง client
		return                                                                                  // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	if err := tx.Create(&payment).Error; err != nil { // ตรวจเงื่อนไข err := tx.Create(&payment).Error; err != nil
		tx.Rollback()                                                                        // ทำงานคำสั่ง tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create payment"}) // ส่ง JSON response กลับไปยัง client
		return                                                                               // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	order.Status = "paid"                         // กำหนดค่า Status ของ order จาก "paid"
	if err := tx.Save(&order).Error; err != nil { // ตรวจเงื่อนไข err := tx.Save(&order).Error; err != nil
		tx.Rollback()                                                                      // ทำงานคำสั่ง tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update order"}) // ส่ง JSON response กลับไปยัง client
		return                                                                             // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	if err := tx.Commit().Error; err != nil { // ตรวจเงื่อนไข err := tx.Commit().Error; err != nil
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to commit payment"}) // ส่ง JSON response กลับไปยัง client
		return                                                                               // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	c.JSON(http.StatusOK, gin.H{ // ส่ง JSON response กลับไปยัง client
		"message": "payment verified", // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"url":     "/" + path,         // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"data":    data,               // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"payment": payment,            // ระบุค่ารายการหนึ่งในชุดข้อมูล
	}) // ทำงานคำสั่ง })
} // ปิด block การทำงานปัจจุบัน

func (h *PaymentHandler) ReadQRCode(file *multipart.FileHeader, amount float64) (string, *VerifyResponse, error) { // ประกาศฟังก์ชัน ReadQRCode
	ext := filepath.Ext(file.Filename)                          // กำหนดค่า ext จาก filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext) // กำหนดค่า file Name จาก fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	path := "slips/" + fileName                                 // กำหนดค่า path จาก "slips/" + fileName

	if err := os.MkdirAll("slips", os.ModePerm); err != nil { // ตรวจเงื่อนไข err := os.MkdirAll("slips", os.ModePerm); err != nil
		return "", nil, err // คืนค่า "", nil, err
	} // ปิด block การทำงานปัจจุบัน

	src, err := file.Open() // กำหนดค่า src, err จาก file.Open()
	if err != nil {         // ตรวจเงื่อนไข err != nil
		return "", nil, err // คืนค่า "", nil, err
	} // ปิด block การทำงานปัจจุบัน
	defer src.Close() // เลื่อนคำสั่งนี้ไปทำตอนฟังก์ชันจบ: src.Close()

	dst, err := os.Create(path) // กำหนดค่า dst, err จาก os.Create(path)
	if err != nil {             // ตรวจเงื่อนไข err != nil
		return "", nil, err // คืนค่า "", nil, err
	} // ปิด block การทำงานปัจจุบัน
	defer dst.Close() // เลื่อนคำสั่งนี้ไปทำตอนฟังก์ชันจบ: dst.Close()

	if _, err = dst.ReadFrom(src); err != nil { // ตรวจเงื่อนไข _, err = dst.ReadFrom(src); err != nil
		return "", nil, err // คืนค่า "", nil, err
	} // ปิด block การทำงานปัจจุบัน

	QR, err := os.Open(path) // กำหนดค่า QR, err จาก os.Open(path)
	if err != nil {          // ตรวจเงื่อนไข err != nil
		return "", nil, err // คืนค่า "", nil, err
	} // ปิด block การทำงานปัจจุบัน
	defer QR.Close() // เลื่อนคำสั่งนี้ไปทำตอนฟังก์ชันจบ: QR.Close()

	img, _, err := image.Decode(QR) // กำหนดค่า img,  , err จาก image.Decode(QR)
	if err != nil {                 // ตรวจเงื่อนไข err != nil
		return "", nil, err // คืนค่า "", nil, err
	} // ปิด block การทำงานปัจจุบัน

	bitmap, err := gozxing.NewBinaryBitmapFromImage(img) // กำหนดค่า bitmap, err จาก gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {                                      // ตรวจเงื่อนไข err != nil
		return "", nil, err // คืนค่า "", nil, err
	} // ปิด block การทำงานปัจจุบัน

	qrReader := qrcode.NewQRCodeReader()        // กำหนดค่า qr Reader จาก qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bitmap, nil) // กำหนดค่า result, err จาก qrReader.Decode(bitmap, nil)
	if err != nil {                             // ตรวจเงื่อนไข err != nil
		return "", nil, fmt.Errorf("failed to read QR code from image") // คืนค่า "", nil, fmt.Errorf("failed to read QR code from image")
	} // ปิด block การทำงานปัจจุบัน

	data, err := h.CheckPayment(result.GetText(), amount) // กำหนดค่า data, err จาก h.CheckPayment(result.GetText(), amount)
	if err != nil {                                       // ตรวจเงื่อนไข err != nil
		return "", nil, err // คืนค่า "", nil, err
	} // ปิด block การทำงานปัจจุบัน

	return path, data, nil // คืนค่า path, data, nil
} // ปิด block การทำงานปัจจุบัน

func (h *PaymentHandler) CheckPayment(qrPayload string, amount float64) (*VerifyResponse, error) { // ประกาศฟังก์ชัน CheckPayment
	url := "https://api.nearbyshop.xyz/slipVerify/noSlip" // กำหนดค่า url จาก "https://api.nearbyshop.xyz/slipVerify/noSlip"
	token := os.Getenv("NEARBYSHOP_TOKEN")                // กำหนดค่า token จาก os.Getenv("NEARBYSHOP_TOKEN")
	if token == "" {                                      // ตรวจเงื่อนไข token == ""
		return nil, fmt.Errorf("NEARBYSHOP_TOKEN is not configured") // คืนค่า nil, fmt.Errorf("NEARBYSHOP_TOKEN is not configured")
	} // ปิด block การทำงานปัจจุบัน

	body := &bytes.Buffer{}                                      // กำหนดค่า body จาก &bytes.Buffer{}
	writer := multipart.NewWriter(body)                          // กำหนดค่า writer จาก multipart.NewWriter(body)
	_ = writer.WriteField("qr_payload", qrPayload)               // กำหนดค่า   จาก writer.WriteField("qr_payload", qrPayload)
	_ = writer.WriteField("amount", fmt.Sprintf("%.2f", amount)) // กำหนดค่า   จาก writer.WriteField("amount", fmt.Sprintf("%.2f", amount))
	writer.Close()                                               // ทำงานคำสั่ง writer.Close()

	req, err := http.NewRequest(http.MethodPost, url, body) // กำหนดค่า req, err จาก http.NewRequest(http.MethodPost, url, body)
	if err != nil {                                         // ตรวจเงื่อนไข err != nil
		return nil, err // คืนค่า nil, err
	} // ปิด block การทำงานปัจจุบัน
	req.Header.Set("Authorization", "Bearer "+token)             // ทำงานคำสั่ง req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", writer.FormDataContentType()) // ทำงานคำสั่ง req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 15 * time.Second} // กำหนดค่า client จาก &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)                       // กำหนดค่า resp, err จาก client.Do(req)
	if err != nil {                                   // ตรวจเงื่อนไข err != nil
		return nil, err // คืนค่า nil, err
	} // ปิด block การทำงานปัจจุบัน
	defer resp.Body.Close() // เลื่อนคำสั่งนี้ไปทำตอนฟังก์ชันจบ: resp.Body.Close()

	data, err := io.ReadAll(resp.Body) // กำหนดค่า data, err จาก io.ReadAll(resp.Body)
	if err != nil {                    // ตรวจเงื่อนไข err != nil
		return nil, err // คืนค่า nil, err
	} // ปิด block การทำงานปัจจุบัน

	if resp.StatusCode != http.StatusOK { // ตรวจเงื่อนไข resp.StatusCode != http.StatusOK
		return nil, fmt.Errorf("payment verification failed: status %d body: %s", resp.StatusCode, string(data)) // คืนค่า nil, fmt.Errorf("payment verification failed: status %d body: %s", resp.StatusCode, string(data))
	} // ปิด block การทำงานปัจจุบัน

	var verifyResp VerifyResponse                             // ประกาศตัวแปร verifyResp VerifyResponse
	if err := json.Unmarshal(data, &verifyResp); err != nil { // ตรวจเงื่อนไข err := json.Unmarshal(data, &verifyResp); err != nil
		return nil, err // คืนค่า nil, err
	} // ปิด block การทำงานปัจจุบัน

	if verifyResp.Status != "success" { // ตรวจเงื่อนไข verifyResp.Status != "success"
		return nil, fmt.Errorf("payment status is not success") // คืนค่า nil, fmt.Errorf("payment status is not success")
	} // ปิด block การทำงานปัจจุบัน

	return &verifyResp, nil // คืนค่า &verifyResp, nil
} // ปิด block การทำงานปัจจุบัน
