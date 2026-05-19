package handler // ประกาศ package handler

import ( // เริ่มส่วน import package ที่ไฟล์นี้ต้องใช้
	"net/http" // นำเข้า package net/http
	"strings"  // นำเข้า package strings

	"backend-go/internal/model" // นำเข้า package backend-go/internal/model

	"github.com/gin-gonic/gin" // นำเข้า package github.com/gin-gonic/gin
	"gorm.io/gorm"             // นำเข้า package gorm.io/gorm
) // ปิด block หรือกลุ่มรายการก่อนหน้า

type CustomerHandler struct { // ประกาศ struct CustomerHandler
	db *gorm.DB // ทำงานคำสั่ง db *gorm.DB
} // ปิด block การทำงานปัจจุบัน

func NewCustomerHandler(db *gorm.DB) *CustomerHandler { // ประกาศฟังก์ชัน NewCustomerHandler
	return &CustomerHandler{db: db} // คืนค่า &CustomerHandler{db: db}
} // ปิด block การทำงานปัจจุบัน

type CartItemRequest struct { // ประกาศ struct CartItemRequest
	ProductID uint `json:"product_id" binding:"required,gt=0"`      // ประกาศ field ProductID พร้อม tag JSON/database
	Quantity  int  `json:"quantity" binding:"required,gt=0,lte=99"` // ประกาศ field Quantity พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน

type AddressRequest struct { // ประกาศ struct AddressRequest
	Type      string `json:"type" binding:"required,max=30"`          // ประกาศ field Type พร้อม tag JSON/database
	Name      string `json:"name" binding:"required,min=2,max=100"`   // ประกาศ field Name พร้อม tag JSON/database
	Street    string `json:"street" binding:"required,min=2,max=255"` // ประกาศ field Street พร้อม tag JSON/database
	City      string `json:"city" binding:"required,min=2,max=100"`   // ประกาศ field City พร้อม tag JSON/database
	State     string `json:"state" binding:"required,min=2,max=100"`  // ประกาศ field State พร้อม tag JSON/database
	Zip       string `json:"zip" binding:"required,min=2,max=20"`     // ประกาศ field Zip พร้อม tag JSON/database
	IsDefault bool   `json:"is_default"`                              // ประกาศ field IsDefault พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน

type PaymentMethodRequest struct { // ประกาศ struct PaymentMethodRequest
	Type      string `json:"type" binding:"required,max=30"`        // ประกาศ field Type พร้อม tag JSON/database
	Last4     string `json:"last4" binding:"required,len=4"`        // ประกาศ field Last4 พร้อม tag JSON/database
	Expiry    string `json:"expiry" binding:"required,min=4,max=7"` // ประกาศ field Expiry พร้อม tag JSON/database
	IsDefault bool   `json:"is_default"`                            // ประกาศ field IsDefault พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน

func (h *CustomerHandler) GetCart(c *gin.Context) { // ประกาศฟังก์ชัน GetCart
	userID := c.GetUint("userID")                                                                            // กำหนดค่า user ID จาก c.GetUint("userID")
	var items []model.CartItem                                                                               // ประกาศตัวแปร items []model.CartItem
	if err := h.db.Preload("Product.Category").Where("user_id = ?", userID).Find(&items).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Preload("Product.Category").Where("user_id = ?", userID).Find(&items).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to get cart") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to get cart")
		return                                                         // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	Success(c, http.StatusOK, "success", items) // ทำงานคำสั่ง Success(c, http.StatusOK, "success", items)
} // ปิด block การทำงานปัจจุบัน

func (h *CustomerHandler) AddCartItem(c *gin.Context) { // ประกาศฟังก์ชัน AddCartItem
	userID := c.GetUint("userID")                  // กำหนดค่า user ID จาก c.GetUint("userID")
	var req CartItemRequest                        // ประกาศตัวแปร req CartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil { // ตรวจเงื่อนไข err := c.ShouldBindJSON(&req); err != nil
		BadRequest(c, "invalid request") // ทำงานคำสั่ง BadRequest(c, "invalid request")
		return                           // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	var product model.Product                                         // ประกาศตัวแปร product model.Product
	if err := h.db.First(&product, req.ProductID).Error; err != nil { // ตรวจเงื่อนไข err := h.db.First(&product, req.ProductID).Error; err != nil
		Error(c, http.StatusNotFound, "product not found") // ทำงานคำสั่ง Error(c, http.StatusNotFound, "product not found")
		return                                             // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	var item model.CartItem                                                                       // ประกาศตัวแปร item model.CartItem
	err := h.db.Where("user_id = ? AND product_id = ?", userID, req.ProductID).First(&item).Error // กำหนดค่า error จาก h.db.Where("user_id = ? AND product_id = ?", userID, req.ProductID).First(&item).Error
	if err == nil {                                                                               // ตรวจเงื่อนไข err == nil
		item.Quantity += req.Quantity                  // ทำงานคำสั่ง item.Quantity += req.Quantity
		if err := h.db.Save(&item).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Save(&item).Error; err != nil
			Error(c, http.StatusInternalServerError, "failed to update cart") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to update cart")
			return                                                            // หยุดการทำงานของฟังก์ชันนี้ทันที
		} // ปิด block การทำงานปัจจุบัน
		Success(c, http.StatusOK, "updated", item) // ทำงานคำสั่ง Success(c, http.StatusOK, "updated", item)
		return                                     // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	if err != gorm.ErrRecordNotFound { // ตรวจเงื่อนไข err != gorm.ErrRecordNotFound
		Error(c, http.StatusInternalServerError, "failed to check cart") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to check cart")
		return                                                           // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	item = model.CartItem{UserID: userID, ProductID: req.ProductID, Quantity: req.Quantity} // กำหนดค่า item จาก model.CartItem{UserID: userID, ProductID: req.ProductID, Quantity: req.Quantity}
	if err := h.db.Create(&item).Error; err != nil {                                        // ตรวจเงื่อนไข err := h.db.Create(&item).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to add cart item") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to add cart item")
		return                                                              // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	Success(c, http.StatusCreated, "created", item) // ทำงานคำสั่ง Success(c, http.StatusCreated, "created", item)
} // ปิด block การทำงานปัจจุบัน

func (h *CustomerHandler) UpdateCartItem(c *gin.Context) { // ประกาศฟังก์ชัน UpdateCartItem
	userID := c.GetUint("userID")             // กำหนดค่า user ID จาก c.GetUint("userID")
	id, err := parsePositiveID(c.Param("id")) // กำหนดค่า id, err จาก parsePositiveID(c.Param("id"))
	if err != nil {                           // ตรวจเงื่อนไข err != nil
		BadRequest(c, "invalid cart item id") // ทำงานคำสั่ง BadRequest(c, "invalid cart item id")
		return                                // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	var req CartItemRequest                        // ประกาศตัวแปร req CartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil { // ตรวจเงื่อนไข err := c.ShouldBindJSON(&req); err != nil
		BadRequest(c, "invalid request") // ทำงานคำสั่ง BadRequest(c, "invalid request")
		return                           // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	var item model.CartItem                                                                     // ประกาศตัวแปร item model.CartItem
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&item).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&item).Error; err != nil
		Error(c, http.StatusNotFound, "cart item not found") // ทำงานคำสั่ง Error(c, http.StatusNotFound, "cart item not found")
		return                                               // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	item.Quantity = req.Quantity                   // กำหนดค่า Quantity ของ item จาก Quantity ของ request
	if err := h.db.Save(&item).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Save(&item).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to update cart item") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to update cart item")
		return                                                                 // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	Success(c, http.StatusOK, "updated", item) // ทำงานคำสั่ง Success(c, http.StatusOK, "updated", item)
} // ปิด block การทำงานปัจจุบัน

func (h *CustomerHandler) DeleteCartItem(c *gin.Context) { // ประกาศฟังก์ชัน DeleteCartItem
	userID := c.GetUint("userID")             // กำหนดค่า user ID จาก c.GetUint("userID")
	id, err := parsePositiveID(c.Param("id")) // กำหนดค่า id, err จาก parsePositiveID(c.Param("id"))
	if err != nil {                           // ตรวจเงื่อนไข err != nil
		BadRequest(c, "invalid cart item id") // ทำงานคำสั่ง BadRequest(c, "invalid cart item id")
		return                                // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.CartItem{}).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.CartItem{}).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to remove cart item") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to remove cart item")
		return                                                                 // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	Success(c, http.StatusOK, "deleted", nil) // ทำงานคำสั่ง Success(c, http.StatusOK, "deleted", nil)
} // ปิด block การทำงานปัจจุบัน

func (h *CustomerHandler) GetWishlist(c *gin.Context) { // ประกาศฟังก์ชัน GetWishlist
	userID := c.GetUint("userID")                                                                            // กำหนดค่า user ID จาก c.GetUint("userID")
	var items []model.WishlistItem                                                                           // ประกาศตัวแปร items []model.WishlistItem
	if err := h.db.Preload("Product.Category").Where("user_id = ?", userID).Find(&items).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Preload("Product.Category").Where("user_id = ?", userID).Find(&items).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to get wishlist") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to get wishlist")
		return                                                             // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	Success(c, http.StatusOK, "success", items) // ทำงานคำสั่ง Success(c, http.StatusOK, "success", items)
} // ปิด block การทำงานปัจจุบัน

func (h *CustomerHandler) AddWishlistItem(c *gin.Context) { // ประกาศฟังก์ชัน AddWishlistItem
	userID := c.GetUint("userID") // กำหนดค่า user ID จาก c.GetUint("userID")
	var req struct {              // ประกาศตัวแปร req struct
		ProductID uint `json:"product_id" binding:"required,gt=0"` // ประกาศ field ProductID พร้อม tag JSON/database
	} // ปิด block การทำงานปัจจุบัน
	if err := c.ShouldBindJSON(&req); err != nil { // ตรวจเงื่อนไข err := c.ShouldBindJSON(&req); err != nil
		BadRequest(c, "invalid request") // ทำงานคำสั่ง BadRequest(c, "invalid request")
		return                           // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	item := model.WishlistItem{UserID: userID, ProductID: req.ProductID}                                                  // กำหนดค่า item จาก model.WishlistItem{UserID: userID, ProductID: req.ProductID}
	if err := h.db.FirstOrCreate(&item, model.WishlistItem{UserID: userID, ProductID: req.ProductID}).Error; err != nil { // ตรวจเงื่อนไข err := h.db.FirstOrCreate(&item, model.WishlistItem{UserID: userID, ProductID: req.ProductID}).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to add wishlist item") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to add wishlist item")
		return                                                                  // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	Success(c, http.StatusCreated, "created", item) // ทำงานคำสั่ง Success(c, http.StatusCreated, "created", item)
} // ปิด block การทำงานปัจจุบัน

func (h *CustomerHandler) DeleteWishlistItem(c *gin.Context) { // ประกาศฟังก์ชัน DeleteWishlistItem
	userID := c.GetUint("userID")             // กำหนดค่า user ID จาก c.GetUint("userID")
	id, err := parsePositiveID(c.Param("id")) // กำหนดค่า id, err จาก parsePositiveID(c.Param("id"))
	if err != nil {                           // ตรวจเงื่อนไข err != nil
		BadRequest(c, "invalid wishlist item id") // ทำงานคำสั่ง BadRequest(c, "invalid wishlist item id")
		return                                    // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.WishlistItem{}).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.WishlistItem{}).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to remove wishlist item") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to remove wishlist item")
		return                                                                     // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	Success(c, http.StatusOK, "deleted", nil) // ทำงานคำสั่ง Success(c, http.StatusOK, "deleted", nil)
} // ปิด block การทำงานปัจจุบัน

func (h *CustomerHandler) GetAddresses(c *gin.Context) { // ประกาศฟังก์ชัน GetAddresses
	userID := c.GetUint("userID")                                                    // กำหนดค่า user ID จาก c.GetUint("userID")
	var addresses []model.Address                                                    // ประกาศตัวแปร addresses []model.Address
	if err := h.db.Where("user_id = ?", userID).Find(&addresses).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Where("user_id = ?", userID).Find(&addresses).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to get addresses") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to get addresses")
		return                                                              // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	Success(c, http.StatusOK, "success", addresses) // ทำงานคำสั่ง Success(c, http.StatusOK, "success", addresses)
} // ปิด block การทำงานปัจจุบัน

func (h *CustomerHandler) CreateAddress(c *gin.Context) { // ประกาศฟังก์ชัน CreateAddress
	h.saveAddress(c, 0) // ทำงานคำสั่ง h.saveAddress(c, 0)
} // ปิด block การทำงานปัจจุบัน

func (h *CustomerHandler) UpdateAddress(c *gin.Context) { // ประกาศฟังก์ชัน UpdateAddress
	id, err := parsePositiveID(c.Param("id")) // กำหนดค่า id, err จาก parsePositiveID(c.Param("id"))
	if err != nil {                           // ตรวจเงื่อนไข err != nil
		BadRequest(c, "invalid address id") // ทำงานคำสั่ง BadRequest(c, "invalid address id")
		return                              // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	h.saveAddress(c, id) // ทำงานคำสั่ง h.saveAddress(c, id)
} // ปิด block การทำงานปัจจุบัน

func (h *CustomerHandler) saveAddress(c *gin.Context, id uint) { // ประกาศฟังก์ชัน saveAddress
	userID := c.GetUint("userID")                  // กำหนดค่า user ID จาก c.GetUint("userID")
	var req AddressRequest                         // ประกาศตัวแปร req AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil { // ตรวจเงื่อนไข err := c.ShouldBindJSON(&req); err != nil
		BadRequest(c, "invalid request") // ทำงานคำสั่ง BadRequest(c, "invalid request")
		return                           // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	address := model.Address{ // กำหนดค่า address จาก model.Address
		UserID: userID, Type: strings.TrimSpace(req.Type), Name: strings.TrimSpace(req.Name), // ระบุค่ารายการหนึ่งในชุดข้อมูล
		Street: strings.TrimSpace(req.Street), City: strings.TrimSpace(req.City), State: strings.TrimSpace(req.State), // ระบุค่ารายการหนึ่งในชุดข้อมูล
		Zip: strings.TrimSpace(req.Zip), IsDefault: req.IsDefault, // ระบุค่ารายการหนึ่งในชุดข้อมูล
	} // ปิด block การทำงานปัจจุบัน
	if id > 0 { // ตรวจเงื่อนไข id > 0
		if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&address).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&address).Error; err != nil
			Error(c, http.StatusNotFound, "address not found") // ทำงานคำสั่ง Error(c, http.StatusNotFound, "address not found")
			return                                             // หยุดการทำงานของฟังก์ชันนี้ทันที
		} // ปิด block การทำงานปัจจุบัน
		address.Type, address.Name, address.Street = req.Type, req.Name, req.Street                               // กำหนดค่า address.Type, address.Name, address.Street จาก req.Type, req.Name, req.Street
		address.City, address.State, address.Zip, address.IsDefault = req.City, req.State, req.Zip, req.IsDefault // กำหนดค่า address.City, address.State, address.Zip, address.Is Default จาก req.City, req.State, req.Zip, req.IsDefault
	} // ปิด block การทำงานปัจจุบัน
	if req.IsDefault { // ตรวจเงื่อนไข req.IsDefault
		h.db.Model(&model.Address{}).Where("user_id = ?", userID).Update("is_default", false) // กำหนดค่า h.db.Model(&model.Address{}).Where("user id จาก ?", userID).Update("is_default", false)
	} // ปิด block การทำงานปัจจุบัน
	if err := h.db.Save(&address).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Save(&address).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to save address") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to save address")
		return                                                             // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	Success(c, map[bool]int{true: http.StatusOK, false: http.StatusCreated}[id > 0], "saved", address) // ทำงานคำสั่ง Success(c, map[bool]int{true: http.StatusOK, false: http.StatusCreated}[id > 0], "saved", address)
} // ปิด block การทำงานปัจจุบัน

func (h *CustomerHandler) DeleteAddress(c *gin.Context) { // ประกาศฟังก์ชัน DeleteAddress
	h.deleteOwned(c, "address", &model.Address{}) // ทำงานคำสั่ง h.deleteOwned(c, "address", &model.Address{})
} // ปิด block การทำงานปัจจุบัน

func (h *CustomerHandler) GetPaymentMethods(c *gin.Context) { // ประกาศฟังก์ชัน GetPaymentMethods
	userID := c.GetUint("userID")                                                  // กำหนดค่า user ID จาก c.GetUint("userID")
	var methods []model.PaymentMethod                                              // ประกาศตัวแปร methods []model.PaymentMethod
	if err := h.db.Where("user_id = ?", userID).Find(&methods).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Where("user_id = ?", userID).Find(&methods).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to get payment methods") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to get payment methods")
		return                                                                    // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	Success(c, http.StatusOK, "success", methods) // ทำงานคำสั่ง Success(c, http.StatusOK, "success", methods)
} // ปิด block การทำงานปัจจุบัน

func (h *CustomerHandler) CreatePaymentMethod(c *gin.Context) { // ประกาศฟังก์ชัน CreatePaymentMethod
	h.savePaymentMethod(c, 0) // ทำงานคำสั่ง h.savePaymentMethod(c, 0)
} // ปิด block การทำงานปัจจุบัน

func (h *CustomerHandler) UpdatePaymentMethod(c *gin.Context) { // ประกาศฟังก์ชัน UpdatePaymentMethod
	id, err := parsePositiveID(c.Param("id")) // กำหนดค่า id, err จาก parsePositiveID(c.Param("id"))
	if err != nil {                           // ตรวจเงื่อนไข err != nil
		BadRequest(c, "invalid payment method id") // ทำงานคำสั่ง BadRequest(c, "invalid payment method id")
		return                                     // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	h.savePaymentMethod(c, id) // ทำงานคำสั่ง h.savePaymentMethod(c, id)
} // ปิด block การทำงานปัจจุบัน

func (h *CustomerHandler) savePaymentMethod(c *gin.Context, id uint) { // ประกาศฟังก์ชัน savePaymentMethod
	userID := c.GetUint("userID")                  // กำหนดค่า user ID จาก c.GetUint("userID")
	var req PaymentMethodRequest                   // ประกาศตัวแปร req PaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil { // ตรวจเงื่อนไข err := c.ShouldBindJSON(&req); err != nil
		BadRequest(c, "invalid request") // ทำงานคำสั่ง BadRequest(c, "invalid request")
		return                           // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	method := model.PaymentMethod{UserID: userID, Type: req.Type, Last4: req.Last4, Expiry: req.Expiry, IsDefault: req.IsDefault} // กำหนดค่า method จาก model.PaymentMethod{UserID: userID, Type: req.Type, Last4: req.Last4, Expiry: req.Expiry, IsDefault: req.IsDefault}
	if id > 0 {                                                                                                                   // ตรวจเงื่อนไข id > 0
		if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&method).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&method).Error; err != nil
			Error(c, http.StatusNotFound, "payment method not found") // ทำงานคำสั่ง Error(c, http.StatusNotFound, "payment method not found")
			return                                                    // หยุดการทำงานของฟังก์ชันนี้ทันที
		} // ปิด block การทำงานปัจจุบัน
		method.Type, method.Last4, method.Expiry, method.IsDefault = req.Type, req.Last4, req.Expiry, req.IsDefault // กำหนดค่า method.Type, method.Last4, method.Expiry, method.Is Default จาก req.Type, req.Last4, req.Expiry, req.IsDefault
	} // ปิด block การทำงานปัจจุบัน
	if req.IsDefault { // ตรวจเงื่อนไข req.IsDefault
		h.db.Model(&model.PaymentMethod{}).Where("user_id = ?", userID).Update("is_default", false) // กำหนดค่า h.db.Model(&model.Payment Method{}).Where("user id จาก ?", userID).Update("is_default", false)
	} // ปิด block การทำงานปัจจุบัน
	if err := h.db.Save(&method).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Save(&method).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to save payment method") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to save payment method")
		return                                                                    // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	Success(c, map[bool]int{true: http.StatusOK, false: http.StatusCreated}[id > 0], "saved", method) // ทำงานคำสั่ง Success(c, map[bool]int{true: http.StatusOK, false: http.StatusCreated}[id > 0], "saved", method)
} // ปิด block การทำงานปัจจุบัน

func (h *CustomerHandler) DeletePaymentMethod(c *gin.Context) { // ประกาศฟังก์ชัน DeletePaymentMethod
	h.deleteOwned(c, "payment method", &model.PaymentMethod{}) // ทำงานคำสั่ง h.deleteOwned(c, "payment method", &model.PaymentMethod{})
} // ปิด block การทำงานปัจจุบัน

func (h *CustomerHandler) deleteOwned(c *gin.Context, label string, modelValue any) { // ประกาศฟังก์ชัน deleteOwned
	userID := c.GetUint("userID")             // กำหนดค่า user ID จาก c.GetUint("userID")
	id, err := parsePositiveID(c.Param("id")) // กำหนดค่า id, err จาก parsePositiveID(c.Param("id"))
	if err != nil {                           // ตรวจเงื่อนไข err != nil
		BadRequest(c, "invalid "+label+" id") // ทำงานคำสั่ง BadRequest(c, "invalid "+label+" id")
		return                                // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).Delete(modelValue).Error; err != nil { // ตรวจเงื่อนไข err := h.db.Where("id = ? AND user_id = ?", id, userID).Delete(modelValue).Error; err != nil
		Error(c, http.StatusInternalServerError, "failed to delete "+label) // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to delete "+label)
		return                                                              // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	Success(c, http.StatusOK, "deleted", nil) // ทำงานคำสั่ง Success(c, http.StatusOK, "deleted", nil)
} // ปิด block การทำงานปัจจุบัน
