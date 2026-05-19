package service // ประกาศ package service

import ( // เริ่มส่วน import package ที่ไฟล์นี้ต้องใช้
	"backend-go/internal/model"      // นำเข้า package backend-go/internal/model
	"backend-go/internal/repository" // นำเข้า package backend-go/internal/repository
) // ปิด block หรือกลุ่มรายการก่อนหน้า

type UserService interface { // ประกาศ interface UserService
	FindAllUsers() ([]model.User, error)              // ทำงานคำสั่ง FindAllUsers() ([]model.User, error)
	FindUserByID(id uint) (model.User, error)         // ทำงานคำสั่ง FindUserByID(id uint) (model.User, error)
	FindUserByEmail(email string) (model.User, error) // ทำงานคำสั่ง FindUserByEmail(email string) (model.User, error)
	CreateUser(user model.User) (model.User, error)   // ทำงานคำสั่ง CreateUser(user model.User) (model.User, error)
	UpdateUser(user model.User) (model.User, error)   // ทำงานคำสั่ง UpdateUser(user model.User) (model.User, error)
	DeleteUser(id uint) error                         // ทำงานคำสั่ง DeleteUser(id uint) error
} // ปิด block การทำงานปัจจุบัน

type userService struct { // ประกาศ struct userService
	userRepo repository.UserRepository // ทำงานคำสั่ง userRepo repository.UserRepository
} // ปิด block การทำงานปัจจุบัน

func NewUserService(userRepo repository.UserRepository) UserService { // ประกาศฟังก์ชัน NewUserService
	return &userService{userRepo: userRepo} // คืนค่า &userService{userRepo: userRepo}
} // ปิด block การทำงานปัจจุบัน

func (s *userService) FindAllUsers() ([]model.User, error) { // ประกาศฟังก์ชัน FindAllUsers
	return s.userRepo.FindAll() // คืนค่า s.userRepo.FindAll()
} // ปิด block การทำงานปัจจุบัน

func (s *userService) FindUserByID(id uint) (model.User, error) { // ประกาศฟังก์ชัน FindUserByID
	return s.userRepo.FindByID(id) // คืนค่า s.userRepo.FindByID(id)
} // ปิด block การทำงานปัจจุบัน

func (s *userService) FindUserByEmail(email string) (model.User, error) { // ประกาศฟังก์ชัน FindUserByEmail
	return s.userRepo.FindByEmail(email) // คืนค่า s.userRepo.FindByEmail(email)
} // ปิด block การทำงานปัจจุบัน

func (s *userService) CreateUser(user model.User) (model.User, error) { // ประกาศฟังก์ชัน CreateUser
	return s.userRepo.Create(user) // คืนค่า s.userRepo.Create(user)
} // ปิด block การทำงานปัจจุบัน

func (s *userService) UpdateUser(user model.User) (model.User, error) { // ประกาศฟังก์ชัน UpdateUser
	return s.userRepo.Update(user) // คืนค่า s.userRepo.Update(user)
} // ปิด block การทำงานปัจจุบัน

func (s *userService) DeleteUser(id uint) error { // ประกาศฟังก์ชัน DeleteUser
	return s.userRepo.Delete(id) // คืนค่า s.userRepo.Delete(id)
} // ปิด block การทำงานปัจจุบัน
