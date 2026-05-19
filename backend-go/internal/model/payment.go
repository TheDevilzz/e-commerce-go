package model // ประกาศ package model

import "time" // ทำงานคำสั่ง import "time"

type Payment struct { // ประกาศ struct Payment
	ID        uint      `json:"id" gorm:"primaryKey"`                     // ประกาศ field ID พร้อม tag JSON/database
	UserID    uint      `json:"user_id" gorm:"index;not null"`            // ประกาศ field UserID พร้อม tag JSON/database
	OrderID   uint      `json:"order_id" gorm:"index;not null"`           // ประกาศ field OrderID พร้อม tag JSON/database
	Ref       string    `json:"ref" gorm:"size:191;uniqueIndex;not null"` // ประกาศ field Ref พร้อม tag JSON/database
	Amount    float64   `json:"amount" gorm:"not null"`                   // ประกาศ field Amount พร้อม tag JSON/database
	SlipURL   string    `json:"slip_url" gorm:"size:255;not null"`        // ประกาศ field SlipURL พร้อม tag JSON/database
	Status    string    `json:"status" gorm:"size:30;default:verified"`   // ประกาศ field Status พร้อม tag JSON/database
	CreatedAt time.Time `json:"created_at"`                               // ประกาศ field CreatedAt พร้อม tag JSON/database
	UpdatedAt time.Time `json:"updated_at"`                               // ประกาศ field UpdatedAt พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน
