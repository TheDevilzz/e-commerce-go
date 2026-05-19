package handler // ประกาศ package handler

import ( // เริ่มส่วน import package ที่ไฟล์นี้ต้องใช้
	"net/http" // นำเข้า package net/http
	"strconv"  // นำเข้า package strconv
	"strings"  // นำเข้า package strings

	"backend-go/internal/model"   // นำเข้า package backend-go/internal/model
	"backend-go/internal/service" // นำเข้า package backend-go/internal/service

	"github.com/gin-gonic/gin" // นำเข้า package github.com/gin-gonic/gin
) // ปิด block หรือกลุ่มรายการก่อนหน้า

type ProductHandler struct { // ประกาศ struct ProductHandler
	productService service.ProductService // ทำงานคำสั่ง productService service.ProductService
} // ปิด block การทำงานปัจจุบัน

func NewProductHandler(productService service.ProductService) *ProductHandler { // ประกาศฟังก์ชัน NewProductHandler
	return &ProductHandler{productService: productService} // คืนค่า &ProductHandler{productService: productService}
} // ปิด block การทำงานปัจจุบัน

func (h *ProductHandler) GetProducts(c *gin.Context) { // ประกาศฟังก์ชัน GetProducts
	query := strings.TrimSpace(c.Query("q"))                // กำหนดค่า query จาก strings.TrimSpace(c.Query("q"))
	categoryID := strings.TrimSpace(c.Query("category_id")) // กำหนดค่า category ID จาก strings.TrimSpace(c.Query("category_id"))
	products, err := h.productService.GetProducts()         // กำหนดค่า products, err จาก h.productService.GetProducts()
	if err != nil {                                         // ตรวจเงื่อนไข err != nil
		Error(c, http.StatusInternalServerError, "failed to get products") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to get products")
		return                                                             // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	if query != "" || categoryID != "" { // ตรวจเงื่อนไข query != "" || categoryID != ""
		filtered := make([]model.Product, 0) // กำหนดค่า filtered จาก make([]model.Product, 0)
		for _, product := range products {   // วนลูปตาม _, product := range products
			matchesQuery := query == "" || strings.Contains(strings.ToLower(product.Name), strings.ToLower(query)) || strings.Contains(strings.ToLower(product.Description), strings.ToLower(query)) // กำหนดค่า matches Query จาก query == "" || strings.Contains(strings.ToLower(product.Name), strings.ToLower(query)) || strings.Contains(strings.ToLower(product.Description), strings.ToLower(query))
			matchesCategory := categoryID == "" || strconv.Itoa(product.CategoryID) == categoryID                                                                                                    // กำหนดค่า matches Category จาก categoryID == "" || strconv.Itoa(product.CategoryID) == categoryID
			if matchesQuery && matchesCategory {                                                                                                                                                     // ตรวจเงื่อนไข matchesQuery && matchesCategory
				filtered = append(filtered, product) // กำหนดค่า filtered จาก append(filtered, product)
			} // ปิด block การทำงานปัจจุบัน
		} // ปิด block การทำงานปัจจุบัน
		products = filtered // กำหนดค่า products จาก filtered
	} // ปิด block การทำงานปัจจุบัน

	Success(c, http.StatusOK, "success", products) // ทำงานคำสั่ง Success(c, http.StatusOK, "success", products)
} // ปิด block การทำงานปัจจุบัน

func (h *ProductHandler) SearchProducts(c *gin.Context) { // ประกาศฟังก์ชัน SearchProducts
	h.GetProducts(c) // ทำงานคำสั่ง h.GetProducts(c)
} // ปิด block การทำงานปัจจุบัน

func (h *ProductHandler) GetProductByID(c *gin.Context) { // ประกาศฟังก์ชัน GetProductByID
	id, err := parsePositiveID(c.Param("id")) // กำหนดค่า id, err จาก parsePositiveID(c.Param("id"))
	if err != nil {                           // ตรวจเงื่อนไข err != nil
		BadRequest(c, "invalid product id") // ทำงานคำสั่ง BadRequest(c, "invalid product id")
		return                              // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	product, err := h.productService.GetProductByID(id) // กำหนดค่า product, err จาก h.productService.GetProductByID(id)
	if err != nil {                                     // ตรวจเงื่อนไข err != nil
		Error(c, http.StatusNotFound, "product not found") // ทำงานคำสั่ง Error(c, http.StatusNotFound, "product not found")
		return                                             // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	Success(c, http.StatusOK, "success", product) // ทำงานคำสั่ง Success(c, http.StatusOK, "success", product)
} // ปิด block การทำงานปัจจุบัน

type ProductRequest struct { // ประกาศ struct ProductRequest
	Name        string  `json:"name" binding:"required,min=2,max=150"`         // ประกาศ field Name พร้อม tag JSON/database
	Description string  `json:"description" binding:"required,min=2,max=1000"` // ประกาศ field Description พร้อม tag JSON/database
	Price       float64 `json:"price" binding:"required,gt=0"`                 // ประกาศ field Price พร้อม tag JSON/database
	Stock       int     `json:"stock" binding:"gte=0"`                         // ประกาศ field Stock พร้อม tag JSON/database
	CategoryID  int     `json:"category_id" binding:"required,gt=0"`           // ประกาศ field CategoryID พร้อม tag JSON/database
	Image       string  `json:"image" binding:"required,max=255"`              // ประกาศ field Image พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน

func (h *ProductHandler) CreateProduct(c *gin.Context) { // ประกาศฟังก์ชัน CreateProduct
	var req ProductRequest // ประกาศตัวแปร req ProductRequest

	if err := c.ShouldBindJSON(&req); err != nil { // ตรวจเงื่อนไข err := c.ShouldBindJSON(&req); err != nil
		BadRequest(c, "invalid request") // ทำงานคำสั่ง BadRequest(c, "invalid request")
		return                           // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	product, err := h.productService.CreateProduct(model.Product{ // กำหนดค่า product, err จาก h.productService.CreateProduct(model.Product
		Name:        strings.TrimSpace(req.Name),        // ระบุค่ารายการหนึ่งในชุดข้อมูล
		Description: strings.TrimSpace(req.Description), // ระบุค่ารายการหนึ่งในชุดข้อมูล
		Price:       req.Price,                          // ระบุค่ารายการหนึ่งในชุดข้อมูล
		Stock:       req.Stock,                          // ระบุค่ารายการหนึ่งในชุดข้อมูล
		CategoryID:  req.CategoryID,                     // ระบุค่ารายการหนึ่งในชุดข้อมูล
		Image:       strings.TrimSpace(req.Image),       // ระบุค่ารายการหนึ่งในชุดข้อมูล
	}) // ทำงานคำสั่ง })
	if err != nil { // ตรวจเงื่อนไข err != nil
		Error(c, http.StatusInternalServerError, "failed to create product") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to create product")
		return                                                               // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	Success(c, http.StatusCreated, "created", product) // ทำงานคำสั่ง Success(c, http.StatusCreated, "created", product)
} // ปิด block การทำงานปัจจุบัน

func (h *ProductHandler) UpdateProduct(c *gin.Context) { // ประกาศฟังก์ชัน UpdateProduct
	id, err := parsePositiveID(c.Param("id")) // กำหนดค่า id, err จาก parsePositiveID(c.Param("id"))
	if err != nil {                           // ตรวจเงื่อนไข err != nil
		BadRequest(c, "invalid product id") // ทำงานคำสั่ง BadRequest(c, "invalid product id")
		return                              // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	var req ProductRequest                         // ประกาศตัวแปร req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil { // ตรวจเงื่อนไข err := c.ShouldBindJSON(&req); err != nil
		BadRequest(c, "invalid request") // ทำงานคำสั่ง BadRequest(c, "invalid request")
		return                           // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	product, err := h.productService.UpdateProduct(id, model.Product{ // กำหนดค่า product, err จาก h.productService.UpdateProduct(id, model.Product
		Name:        strings.TrimSpace(req.Name),        // ระบุค่ารายการหนึ่งในชุดข้อมูล
		Description: strings.TrimSpace(req.Description), // ระบุค่ารายการหนึ่งในชุดข้อมูล
		Price:       req.Price,                          // ระบุค่ารายการหนึ่งในชุดข้อมูล
		Stock:       req.Stock,                          // ระบุค่ารายการหนึ่งในชุดข้อมูล
		CategoryID:  req.CategoryID,                     // ระบุค่ารายการหนึ่งในชุดข้อมูล
		Image:       strings.TrimSpace(req.Image),       // ระบุค่ารายการหนึ่งในชุดข้อมูล
	}) // ทำงานคำสั่ง })
	if err != nil { // ตรวจเงื่อนไข err != nil
		Error(c, http.StatusNotFound, "product not found") // ทำงานคำสั่ง Error(c, http.StatusNotFound, "product not found")
		return                                             // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	Success(c, http.StatusOK, "updated", product) // ทำงานคำสั่ง Success(c, http.StatusOK, "updated", product)
} // ปิด block การทำงานปัจจุบัน

func (h *ProductHandler) DeleteProduct(c *gin.Context) { // ประกาศฟังก์ชัน DeleteProduct
	id, err := parsePositiveID(c.Param("id")) // กำหนดค่า id, err จาก parsePositiveID(c.Param("id"))
	if err != nil {                           // ตรวจเงื่อนไข err != nil
		BadRequest(c, "invalid product id") // ทำงานคำสั่ง BadRequest(c, "invalid product id")
		return                              // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	err = h.productService.DeleteProduct(id) // กำหนดค่า error จาก h.productService.DeleteProduct(id)
	if err != nil {                          // ตรวจเงื่อนไข err != nil
		Error(c, http.StatusInternalServerError, "failed to delete product") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to delete product")
		return                                                               // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	Success(c, http.StatusOK, "deleted", nil) // ทำงานคำสั่ง Success(c, http.StatusOK, "deleted", nil)
} // ปิด block การทำงานปัจจุบัน

type CategoryRequest struct { // ประกาศ struct CategoryRequest
	Name string `json:"name" binding:"required,min=2,max=100"` // ประกาศ field Name พร้อม tag JSON/database
} // ปิด block การทำงานปัจจุบัน

func (h *ProductHandler) CreateCategory(c *gin.Context) { // ประกาศฟังก์ชัน CreateCategory
	var req CategoryRequest                        // ประกาศตัวแปร req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil { // ตรวจเงื่อนไข err := c.ShouldBindJSON(&req); err != nil
		BadRequest(c, "invalid request") // ทำงานคำสั่ง BadRequest(c, "invalid request")
		return                           // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	category, err := h.productService.CreateCategory(model.Category{Name: strings.TrimSpace(req.Name)}) // กำหนดค่า category, err จาก h.productService.CreateCategory(model.Category{Name: strings.TrimSpace(req.Name)})
	if err != nil {                                                                                     // ตรวจเงื่อนไข err != nil
		Error(c, http.StatusInternalServerError, "failed to create category") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to create category")
		return                                                                // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	Success(c, http.StatusCreated, "created", category) // ทำงานคำสั่ง Success(c, http.StatusCreated, "created", category)
} // ปิด block การทำงานปัจจุบัน

func (h *ProductHandler) GetCategories(c *gin.Context) { // ประกาศฟังก์ชัน GetCategories
	categories, err := h.productService.GetCategories() // กำหนดค่า categories, err จาก h.productService.GetCategories()
	if err != nil {                                     // ตรวจเงื่อนไข err != nil
		Error(c, http.StatusInternalServerError, "failed to get categories") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to get categories")
		return                                                               // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน
	Success(c, http.StatusOK, "success", categories) // ทำงานคำสั่ง Success(c, http.StatusOK, "success", categories)
} // ปิด block การทำงานปัจจุบัน

func (h *ProductHandler) GetProductsByCategory(c *gin.Context) { // ประกาศฟังก์ชัน GetProductsByCategory
	id, err := parsePositiveID(c.Param("id")) // กำหนดค่า id, err จาก parsePositiveID(c.Param("id"))
	if err != nil {                           // ตรวจเงื่อนไข err != nil
		BadRequest(c, "invalid category id") // ทำงานคำสั่ง BadRequest(c, "invalid category id")
		return                               // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	products, err := h.productService.GetProductsByCategory(id) // กำหนดค่า products, err จาก h.productService.GetProductsByCategory(id)
	if err != nil {                                             // ตรวจเงื่อนไข err != nil
		Error(c, http.StatusInternalServerError, "failed to get products") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to get products")
		return                                                             // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	Success(c, http.StatusOK, "success", products) // ทำงานคำสั่ง Success(c, http.StatusOK, "success", products)
} // ปิด block การทำงานปัจจุบัน

func (h *ProductHandler) UpdateCategory(c *gin.Context) { // ประกาศฟังก์ชัน UpdateCategory
	id, err := parsePositiveID(c.Param("id")) // กำหนดค่า id, err จาก parsePositiveID(c.Param("id"))
	if err != nil {                           // ตรวจเงื่อนไข err != nil
		BadRequest(c, "invalid category id") // ทำงานคำสั่ง BadRequest(c, "invalid category id")
		return                               // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	var req CategoryRequest                        // ประกาศตัวแปร req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil { // ตรวจเงื่อนไข err := c.ShouldBindJSON(&req); err != nil
		BadRequest(c, "invalid request") // ทำงานคำสั่ง BadRequest(c, "invalid request")
		return                           // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	category, err := h.productService.UpdateCategory(id, model.Category{Name: strings.TrimSpace(req.Name)}) // กำหนดค่า category, err จาก h.productService.UpdateCategory(id, model.Category{Name: strings.TrimSpace(req.Name)})
	if err != nil {                                                                                         // ตรวจเงื่อนไข err != nil
		Error(c, http.StatusNotFound, "category not found") // ทำงานคำสั่ง Error(c, http.StatusNotFound, "category not found")
		return                                              // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	Success(c, http.StatusOK, "updated", category) // ทำงานคำสั่ง Success(c, http.StatusOK, "updated", category)
} // ปิด block การทำงานปัจจุบัน

func (h *ProductHandler) DeleteCategory(c *gin.Context) { // ประกาศฟังก์ชัน DeleteCategory
	id, err := parsePositiveID(c.Param("id")) // กำหนดค่า id, err จาก parsePositiveID(c.Param("id"))
	if err != nil {                           // ตรวจเงื่อนไข err != nil
		BadRequest(c, "invalid category id") // ทำงานคำสั่ง BadRequest(c, "invalid category id")
		return                               // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	if err := h.productService.DeleteCategory(id); err != nil { // ตรวจเงื่อนไข err := h.productService.DeleteCategory(id); err != nil
		Error(c, http.StatusInternalServerError, "failed to delete category") // ทำงานคำสั่ง Error(c, http.StatusInternalServerError, "failed to delete category")
		return                                                                // หยุดการทำงานของฟังก์ชันนี้ทันที
	} // ปิด block การทำงานปัจจุบัน

	Success(c, http.StatusOK, "deleted", nil) // ทำงานคำสั่ง Success(c, http.StatusOK, "deleted", nil)
} // ปิด block การทำงานปัจจุบัน

func parsePositiveID(value string) (uint, error) { // ประกาศฟังก์ชัน parsePositiveID
	id, err := strconv.Atoi(value) // กำหนดค่า id, err จาก strconv.Atoi(value)
	if err != nil || id <= 0 {     // ตรวจเงื่อนไข err != nil || id <= 0
		return 0, strconv.ErrSyntax // คืนค่า 0, strconv.ErrSyntax
	} // ปิด block การทำงานปัจจุบัน
	return uint(id), nil // คืนค่า uint(id), nil
} // ปิด block การทำงานปัจจุบัน
