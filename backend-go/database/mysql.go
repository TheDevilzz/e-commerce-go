package database // ประกาศ package database

import ( // เริ่มส่วน import package ที่ไฟล์นี้ต้องใช้
	"fmt" // นำเข้า package fmt
	"log" // นำเข้า package log

	"backend-go/config" // นำเข้า package backend-go/config

	"gorm.io/driver/mysql" // นำเข้า package gorm.io/driver/mysql
	"gorm.io/gorm"         // นำเข้า package gorm.io/gorm
) // ปิด block หรือกลุ่มรายการก่อนหน้า

func ConnectDB(cfg config.Config) *gorm.DB { // ประกาศฟังก์ชัน ConnectDB
	dsn := fmt.Sprintf( // กำหนดค่า dsn จาก fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", // ระบุค่ารายการหนึ่งในชุดข้อมูล
		cfg.DBUser,     // ระบุค่ารายการหนึ่งในชุดข้อมูล
		cfg.DBPassword, // ระบุค่ารายการหนึ่งในชุดข้อมูล
		cfg.DBHost,     // ระบุค่ารายการหนึ่งในชุดข้อมูล
		cfg.DBPort,     // ระบุค่ารายการหนึ่งในชุดข้อมูล
		cfg.DBName,     // ระบุค่ารายการหนึ่งในชุดข้อมูล
	) // ปิด block หรือกลุ่มรายการก่อนหน้า

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) // กำหนดค่า db, err จาก gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {                                       // ตรวจเงื่อนไข err != nil
		log.Fatal("Database connection failed:", err) // ทำงานคำสั่ง log.Fatal("Database connection failed:", err)
	} // ปิด block การทำงานปัจจุบัน

	return db // คืนค่า db
} // ปิด block การทำงานปัจจุบัน
