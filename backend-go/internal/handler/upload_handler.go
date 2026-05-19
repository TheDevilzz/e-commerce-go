package handler // ประกาศ package handler

import ( // เริ่มส่วน import package ที่ไฟล์นี้ต้องใช้
	"fmt"            // นำเข้า package fmt
	"mime/multipart" // นำเข้า package mime/multipart
	"net/http"       // นำเข้า package net/http
	"os"             // นำเข้า package os
	"path/filepath"  // นำเข้า package path/filepath
	"strings"        // นำเข้า package strings
	"time"           // นำเข้า package time

	"github.com/gin-gonic/gin" // นำเข้า package github.com/gin-gonic/gin
) // ปิด block หรือกลุ่มรายการก่อนหน้า

const maxImageUploadSize = 5 << 20 // ประกาศค่าคงที่ maxImageUploadSize = 5 << 20

var allowedImageExtensions = map[string]bool{ // ประกาศตัวแปร allowedImageExtensions = map[string]bool
	".jpg":  true, // ระบุค่ารายการหนึ่งในชุดข้อมูล
	".jpeg": true, // ระบุค่ารายการหนึ่งในชุดข้อมูล
	".png":  true, // ระบุค่ารายการหนึ่งในชุดข้อมูล
	".webp": true, // ระบุค่ารายการหนึ่งในชุดข้อมูล
} // ปิด block การทำงานปัจจุบัน

func UploadImage(c *gin.Context) { // ประกาศฟังก์ชัน UploadImage
	file, err := c.FormFile("image") // กำหนดค่า file, err จาก c.FormFile("image")
	if err != nil {                  // ตรวจเงื่อนไข err != nil
		BadRequest(c, "image is required") // ทำงานคำสั่ง BadRequest(c, "image is required")
		return                             // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	if err := validateImageFile(file); err != nil { // ตรวจเงื่อนไข err := validateImageFile(file); err != nil
		BadRequest(c, err.Error()) // ทำงานคำสั่ง BadRequest(c, err.Error())
		return                     // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	if err := os.MkdirAll("uploads", 0755); err != nil { // ตรวจเงื่อนไข err := os.MkdirAll("uploads", 0755); err != nil
		Error(c, http.StatusInternalServerError, "failed to prepare upload directory") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to prepare upload directory")
		return                                                                         // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	ext := strings.ToLower(filepath.Ext(file.Filename))         // กำหนดค่า ext จาก strings.ToLower(filepath.Ext(file.Filename))
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext) // กำหนดค่า filename จาก fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	path := "uploads/" + filename                               // กำหนดค่า path จาก "uploads/" + filename

	if err := c.SaveUploadedFile(file, path); err != nil { // ตรวจเงื่อนไข err := c.SaveUploadedFile(file, path); err != nil
		Error(c, http.StatusInternalServerError, "failed to upload image") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to upload image")
		return                                                             // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	Success(c, http.StatusOK, "upload success", gin.H{ // ทำงานคำสั่ง Success(c, http.StatusOK, "upload success", gin.H{
		"url":     "/uploads/" + filename,          // ระบุค่ารายการหนึ่งในชุดข้อมูล
		"content": file.Header.Get("Content-Type"), // ระบุค่ารายการหนึ่งในชุดข้อมูล
	}) // ทำงานคำสั่ง })
} // ปิด block การทำงานปัจจุบัน

func validateImageFile(file *multipart.FileHeader) error { // ประกาศฟังก์ชัน validateImageFile
	if file.Size <= 0 { // ตรวจเงื่อนไข file.Size <= 0
		return fmt.Errorf("image is empty") // คืนค่า fmt.Errorf("image is empty")
	} // ปิด block การทำงานปัจจุบัน
	if file.Size > maxImageUploadSize { // ตรวจเงื่อนไข file.Size > maxImageUploadSize
		return fmt.Errorf("image size must be 5MB or less") // คืนค่า fmt.Errorf("image size must be 5MB or less")
	} // ปิด block การทำงานปัจจุบัน

	ext := strings.ToLower(filepath.Ext(file.Filename)) // กำหนดค่า ext จาก strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageExtensions[ext] {                   // ตรวจเงื่อนไข !allowedImageExtensions[ext]
		return fmt.Errorf("only jpg, jpeg, png, and webp images are allowed") // คืนค่า fmt.Errorf("only jpg, jpeg, png, and webp images are allowed")
	} // ปิด block การทำงานปัจจุบัน

	contentType := file.Header.Get("Content-Type") // กำหนดค่า content Type จาก file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") { // ตรวจเงื่อนไข !strings.HasPrefix(contentType, "image/")
		return fmt.Errorf("file must be an image") // คืนค่า fmt.Errorf("file must be an image")
	} // ปิด block การทำงานปัจจุบัน

	return nil // คืนค่า nil
} // ปิด block การทำงานปัจจุบัน
