package model // ประกาศ package model

import "time" // ทำงานคำสั่ง import "time"

type CartItem struct { // ประกาศ struct CartItem
	ID        uint      `json:"id" gorm:"primaryKey"`                // ประกาศ field ID พร้อม tag JSON/database
	UserID    uint      `json:"user_id" gorm:"index;not null"`       // ประกาศ field UserID พร้อม tag JSON/database
	ProductID uint      `json:"product_id" gorm:"index;not null"`    // ประกาศ field ProductID พร้อม tag JSON/database
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"` // ประกาศ field Product พร้อม tag JSON/database
	Quantity  int       `json:"quantity" gorm:"not null"`            // ประกาศ field Quantity พร้อม tag JSON/database
	CreatedAt time.Time `json:"created_at"`                          // ประกาศ field CreatedAt พร้อม tag JSON/database
	UpdatedAt time.Time `json:"updated_at"`                          // ประกาศ field UpdatedAt พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน

type WishlistItem struct { // ประกาศ struct WishlistItem
	ID        uint      `json:"id" gorm:"primaryKey"`                                     // ประกาศ field ID พร้อม tag JSON/database
	UserID    uint      `json:"user_id" gorm:"index:idx_user_product,unique;not null"`    // ประกาศ field UserID พร้อม tag JSON/database
	ProductID uint      `json:"product_id" gorm:"index:idx_user_product,unique;not null"` // ประกาศ field ProductID พร้อม tag JSON/database
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`                      // ประกาศ field Product พร้อม tag JSON/database
	CreatedAt time.Time `json:"created_at"`                                               // ประกาศ field CreatedAt พร้อม tag JSON/database
	UpdatedAt time.Time `json:"updated_at"`                                               // ประกาศ field UpdatedAt พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน

type Address struct { // ประกาศ struct Address
	ID        uint      `json:"id" gorm:"primaryKey"`            // ประกาศ field ID พร้อม tag JSON/database
	UserID    uint      `json:"user_id" gorm:"index;not null"`   // ประกาศ field UserID พร้อม tag JSON/database
	Type      string    `json:"type" gorm:"size:30;not null"`    // ประกาศ field Type พร้อม tag JSON/database
	Name      string    `json:"name" gorm:"size:100;not null"`   // ประกาศ field Name พร้อม tag JSON/database
	Street    string    `json:"street" gorm:"size:255;not null"` // ประกาศ field Street พร้อม tag JSON/database
	City      string    `json:"city" gorm:"size:100;not null"`   // ประกาศ field City พร้อม tag JSON/database
	State     string    `json:"state" gorm:"size:100;not null"`  // ประกาศ field State พร้อม tag JSON/database
	Zip       string    `json:"zip" gorm:"size:20;not null"`     // ประกาศ field Zip พร้อม tag JSON/database
	IsDefault bool      `json:"is_default" gorm:"default:false"` // ประกาศ field IsDefault พร้อม tag JSON/database
	CreatedAt time.Time `json:"created_at"`                      // ประกาศ field CreatedAt พร้อม tag JSON/database
	UpdatedAt time.Time `json:"updated_at"`                      // ประกาศ field UpdatedAt พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน

type PaymentMethod struct { // ประกาศ struct PaymentMethod
	ID        uint      `json:"id" gorm:"primaryKey"`            // ประกาศ field ID พร้อม tag JSON/database
	UserID    uint      `json:"user_id" gorm:"index;not null"`   // ประกาศ field UserID พร้อม tag JSON/database
	Type      string    `json:"type" gorm:"size:30;not null"`    // ประกาศ field Type พร้อม tag JSON/database
	Last4     string    `json:"last4" gorm:"size:4;not null"`    // ประกาศ field Last4 พร้อม tag JSON/database
	Expiry    string    `json:"expiry" gorm:"size:7;not null"`   // ประกาศ field Expiry พร้อม tag JSON/database
	IsDefault bool      `json:"is_default" gorm:"default:false"` // ประกาศ field IsDefault พร้อม tag JSON/database
	CreatedAt time.Time `json:"created_at"`                      // ประกาศ field CreatedAt พร้อม tag JSON/database
	UpdatedAt time.Time `json:"updated_at"`                      // ประกาศ field UpdatedAt พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน

type Promotion struct { // ประกาศ struct Promotion
	ID          uint      `json:"id" gorm:"primaryKey"`                     // ประกาศ field ID พร้อม tag JSON/database
	Code        string    `json:"code" gorm:"size:50;uniqueIndex;not null"` // ประกาศ field Code พร้อม tag JSON/database
	Description string    `json:"description" gorm:"size:255;not null"`     // ประกาศ field Description พร้อม tag JSON/database
	Discount    float64   `json:"discount" gorm:"not null"`                 // ประกาศ field Discount พร้อม tag JSON/database
	Type        string    `json:"type" gorm:"size:30;not null"`             // ประกาศ field Type พร้อม tag JSON/database
	StartDate   string    `json:"start_date" gorm:"size:10;not null"`       // ประกาศ field StartDate พร้อม tag JSON/database
	EndDate     string    `json:"end_date" gorm:"size:10;not null"`         // ประกาศ field EndDate พร้อม tag JSON/database
	Status      string    `json:"status" gorm:"size:30;index;not null"`     // ประกาศ field Status พร้อม tag JSON/database
	UsageCount  int       `json:"usage_count" gorm:"default:0"`             // ประกาศ field UsageCount พร้อม tag JSON/database
	CreatedAt   time.Time `json:"created_at"`                               // ประกาศ field CreatedAt พร้อม tag JSON/database
	UpdatedAt   time.Time `json:"updated_at"`                               // ประกาศ field UpdatedAt พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน
