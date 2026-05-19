package model // ประกาศ package model

import "time" // ทำงานคำสั่ง import "time"

type Order struct { // ประกาศ struct Order
	ID         uint        `json:"id" gorm:"primaryKey"`                                         // ประกาศ field ID พร้อม tag JSON/database
	UserID     uint        `json:"user_id" gorm:"index;not null"`                                // ประกาศ field UserID พร้อม tag JSON/database
	TotalPrice float64     `json:"total_price" gorm:"not null"`                                  // ประกาศ field TotalPrice พร้อม tag JSON/database
	Status     string      `json:"status" gorm:"size:30;default:pending"`                        // ประกาศ field Status พร้อม tag JSON/database
	Items      []OrderItem `json:"items" gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;"` // ประกาศ field Items พร้อม tag JSON/database
	CreatedAt  time.Time   // ทำงานคำสั่ง CreatedAt time.Time
	UpdatedAt  time.Time   // ทำงานคำสั่ง UpdatedAt time.Time
} // ปิด block การทำงานปัจจุบัน

type OrderItem struct { // ประกาศ struct OrderItem
	ID          uint    `json:"id" gorm:"primaryKey"`                  // ประกาศ field ID พร้อม tag JSON/database
	OrderID     uint    `json:"order_id" gorm:"not null"`              // ประกาศ field OrderID พร้อม tag JSON/database
	ProductID   uint    `json:"product_id" gorm:"index;not null"`      // ประกาศ field ProductID พร้อม tag JSON/database
	ProductName string  `json:"product_name" gorm:"size:150;not null"` // ประกาศ field ProductName พร้อม tag JSON/database
	Quantity    int     `json:"quantity" gorm:"not null"`              // ประกาศ field Quantity พร้อม tag JSON/database
	Price       float64 `json:"price" gorm:"not null"`                 // ประกาศ field Price พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน
