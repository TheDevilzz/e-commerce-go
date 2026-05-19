package model // ประกาศ package model

import ( // เริ่มส่วน import package ที่ไฟล์นี้ต้องใช้
	"time" // นำเข้า package time

	"gorm.io/gorm" // นำเข้า package gorm.io/gorm
) // ปิด block หรือกลุ่มรายการก่อนหน้า

type Product struct { // ประกาศ struct Product
	ID          int            `json:"id" gorm:"primaryKey"`                      // ประกาศ field ID พร้อม tag JSON/database
	Name        string         `json:"name" gorm:"size:150;uniqueIndex;not null"` // ประกาศ field Name พร้อม tag JSON/database
	Description string         `json:"description" gorm:"size:1000"`              // ประกาศ field Description พร้อม tag JSON/database
	Price       float64        `json:"price" gorm:"not null"`                     // ประกาศ field Price พร้อม tag JSON/database
	Stock       int            `json:"stock" gorm:"not null"`                     // ประกาศ field Stock พร้อม tag JSON/database
	CategoryID  int            `json:"category_id" gorm:"index;not null"`         // ประกาศ field CategoryID พร้อม tag JSON/database
	Category    Category       `json:"category" gorm:"foreignKey:CategoryID"`     // ประกาศ field Category พร้อม tag JSON/database
	Image       string         `json:"image" gorm:"size:255"`                     // ประกาศ field Image พร้อม tag JSON/database
	CreatedAt   time.Time      `json:"created_at"`                                // ประกาศ field CreatedAt พร้อม tag JSON/database
	UpdatedAt   time.Time      `json:"updated_at"`                                // ประกาศ field UpdatedAt พร้อม tag JSON/database
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`                            // ประกาศ field DeletedAt พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน

type Category struct { // ประกาศ struct Category
	ID        int       `json:"id" gorm:"primaryKey"`                      // ประกาศ field ID พร้อม tag JSON/database
	Name      string    `json:"name" gorm:"size:100;uniqueIndex;not null"` // ประกาศ field Name พร้อม tag JSON/database
	CreatedAt time.Time `json:"created_at"`                                // ประกาศ field CreatedAt พร้อม tag JSON/database
	UpdatedAt time.Time `json:"updated_at"`                                // ประกาศ field UpdatedAt พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน
