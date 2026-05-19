package repository // ประกาศ package repository

import ( // เริ่มส่วน import package ที่ไฟล์นี้ต้องใช้
	"backend-go/internal/model" // นำเข้า package backend-go/internal/model

	"gorm.io/gorm" // นำเข้า package gorm.io/gorm
) // ปิด block หรือกลุ่มรายการก่อนหน้า

type UserRepository interface { // ประกาศ interface UserRepository
	FindAll() ([]model.User, error)               // ทำงานคำสั่ง FindAll() ([]model.User, error)
	FindByID(id uint) (model.User, error)         // ทำงานคำสั่ง FindByID(id uint) (model.User, error)
	FindByEmail(email string) (model.User, error) // ทำงานคำสั่ง FindByEmail(email string) (model.User, error)
	Create(user model.User) (model.User, error)   // ทำงานคำสั่ง Create(user model.User) (model.User, error)
	Update(user model.User) (model.User, error)   // ทำงานคำสั่ง Update(user model.User) (model.User, error)
	Delete(id uint) error                         // ทำงานคำสั่ง Delete(id uint) error
} // ปิด block การทำงานปัจจุบัน

type userRepository struct { // ประกาศ struct userRepository
	db *gorm.DB // ทำงานคำสั่ง db *gorm.DB
} // ปิด block การทำงานปัจจุบัน

func NewUserRepository(db *gorm.DB) UserRepository { // ประกาศฟังก์ชัน NewUserRepository
	return &userRepository{db: db} // คืนค่า &userRepository{db: db}
} // ปิด block การทำงานปัจจุบัน

func (r *userRepository) FindAll() ([]model.User, error) { // ประกาศฟังก์ชัน FindAll
	var users []model.User         // ประกาศตัวแปร users []model.User
	err := r.db.Find(&users).Error // กำหนดค่า error จาก r.db.Find(&users).Error
	return users, err              // คืนค่า users, err
} // ปิด block การทำงานปัจจุบัน

func (r *userRepository) FindByID(id uint) (model.User, error) { // ประกาศฟังก์ชัน FindByID
	var user model.User                // ประกาศตัวแปร user model.User
	err := r.db.First(&user, id).Error // กำหนดค่า error จาก r.db.First(&user, id).Error
	return user, err                   // คืนค่า user, err
} // ปิด block การทำงานปัจจุบัน

func (r *userRepository) FindByEmail(email string) (model.User, error) { // ประกาศฟังก์ชัน FindByEmail
	var user model.User                                      // ประกาศตัวแปร user model.User
	err := r.db.Where("email = ?", email).First(&user).Error // กำหนดค่า error จาก r.db.Where("email = ?", email).First(&user).Error
	return user, err                                         // คืนค่า user, err
} // ปิด block การทำงานปัจจุบัน

func (r *userRepository) Create(user model.User) (model.User, error) { // ประกาศฟังก์ชัน Create
	err := r.db.Create(&user).Error // กำหนดค่า error จาก r.db.Create(&user).Error
	return user, err                // คืนค่า user, err
} // ปิด block การทำงานปัจจุบัน

func (r *userRepository) Update(user model.User) (model.User, error) { // ประกาศฟังก์ชัน Update
	err := r.db.Save(&user).Error // กำหนดค่า error จาก r.db.Save(&user).Error
	return user, err              // คืนค่า user, err
} // ปิด block การทำงานปัจจุบัน

func (r *userRepository) Delete(id uint) error { // ประกาศฟังก์ชัน Delete
	err := r.db.Delete(&model.User{}, id).Error // กำหนดค่า error จาก r.db.Delete(&model.User{}, id).Error
	return err                                  // คืนค่า err
} // ปิด block การทำงานปัจจุบัน
