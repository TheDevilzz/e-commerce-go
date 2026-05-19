package model // ประกาศ package model

import "time" // ทำงานคำสั่ง import "time"

type User struct { // ประกาศ struct User
	ID          uint      `json:"id" gorm:"primaryKey"`                         // ประกาศ field ID พร้อม tag JSON/database
	Username    string    `json:"username" gorm:"size:50;uniqueIndex;not null"` // ประกาศ field Username พร้อม tag JSON/database
	Email       string    `json:"email" gorm:"size:191;uniqueIndex"`            // ประกาศ field Email พร้อม tag JSON/database
	Name        string    `json:"name" gorm:"size:100;not null"`                // ประกาศ field Name พร้อม tag JSON/database
	Password    string    `json:"-" gorm:"size:255;not null"`                   // ประกาศ field Password พร้อม tag JSON/database
	Phone       string    `json:"phone" gorm:"size:30;not null"`                // ประกาศ field Phone พร้อม tag JSON/database
	DateOfBirth string    `json:"date_of_birth" gorm:"not null"`                // ประกาศ field DateOfBirth พร้อม tag JSON/database
	Role        string    `json:"role" gorm:"size:30;index;not null"`           // ประกาศ field Role พร้อม tag JSON/database
	CreatedAt   time.Time `json:"created_at"`                                   // ประกาศ field CreatedAt พร้อม tag JSON/database
	UpdatedAt   time.Time `json:"updated_at"`                                   // ประกาศ field UpdatedAt พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน

type RefreshToken struct { // ประกาศ struct RefreshToken
	ID        uint       `json:"id" gorm:"primaryKey"`                  // ประกาศ field ID พร้อม tag JSON/database
	UserID    uint       `json:"user_id" gorm:"index;not null"`         // ประกาศ field UserID พร้อม tag JSON/database
	User      User       `json:"-" gorm:"foreignKey:UserID"`            // ประกาศ field User พร้อม tag JSON/database
	TokenHash string     `json:"-" gorm:"size:64;uniqueIndex;not null"` // ประกาศ field TokenHash พร้อม tag JSON/database
	ExpiresAt time.Time  `json:"expires_at" gorm:"index;not null"`      // ประกาศ field ExpiresAt พร้อม tag JSON/database
	RevokedAt *time.Time `json:"revoked_at" gorm:"index"`               // ประกาศ field RevokedAt พร้อม tag JSON/database
	CreatedAt time.Time  `json:"created_at"`                            // ประกาศ field CreatedAt พร้อม tag JSON/database
	UpdatedAt time.Time  `json:"updated_at"`                            // ประกาศ field UpdatedAt พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน
