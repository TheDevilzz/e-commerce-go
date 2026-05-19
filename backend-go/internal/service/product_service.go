package service // ประกาศ package service

import ( // เริ่มส่วน import package ที่ไฟล์นี้ต้องใช้
	"backend-go/internal/model"      // นำเข้า package backend-go/internal/model
	"backend-go/internal/repository" // นำเข้า package backend-go/internal/repository
) // ปิด block หรือกลุ่มรายการก่อนหน้า

type ProductService interface { // ประกาศ interface ProductService
	GetProducts() ([]model.Product, error)                                   // ทำงานคำสั่ง GetProducts() ([]model.Product, error)
	GetProductByID(id uint) (model.Product, error)                           // ทำงานคำสั่ง GetProductByID(id uint) (model.Product, error)
	CreateProduct(product model.Product) (model.Product, error)              // ทำงานคำสั่ง CreateProduct(product model.Product) (model.Product, error)
	UpdateProduct(id uint, product model.Product) (model.Product, error)     // ทำงานคำสั่ง UpdateProduct(id uint, product model.Product) (model.Product, error)
	DeleteProduct(id uint) error                                             // ทำงานคำสั่ง DeleteProduct(id uint) error
	GetProductsByCategory(categoryID uint) ([]model.Product, error)          // ทำงานคำสั่ง GetProductsByCategory(categoryID uint) ([]model.Product, error)
	GetCategories() ([]model.Category, error)                                // ทำงานคำสั่ง GetCategories() ([]model.Category, error)
	GetCategoryByID(id uint) (model.Category, error)                         // ทำงานคำสั่ง GetCategoryByID(id uint) (model.Category, error)
	CreateCategory(category model.Category) (model.Category, error)          // ทำงานคำสั่ง CreateCategory(category model.Category) (model.Category, error)
	UpdateCategory(id uint, category model.Category) (model.Category, error) // ทำงานคำสั่ง UpdateCategory(id uint, category model.Category) (model.Category, error)
	DeleteCategory(id uint) error                                            // ทำงานคำสั่ง DeleteCategory(id uint) error
} // ปิด block การทำงานปัจจุบัน

type productService struct { // ประกาศ struct productService
	productRepo repository.ProductRepository // ทำงานคำสั่ง productRepo repository.ProductRepository
} // ปิด block การทำงานปัจจุบัน

func NewProductService(productRepo repository.ProductRepository) ProductService { // ประกาศฟังก์ชัน NewProductService
	return &productService{productRepo: productRepo} // คืนค่า &productService{productRepo: productRepo}
} // ปิด block การทำงานปัจจุบัน

func (s *productService) GetProducts() ([]model.Product, error) { // ประกาศฟังก์ชัน GetProducts
	return s.productRepo.FindAll() // คืนค่า s.productRepo.FindAll()
} // ปิด block การทำงานปัจจุบัน

func (s *productService) GetProductByID(id uint) (model.Product, error) { // ประกาศฟังก์ชัน GetProductByID
	return s.productRepo.FindByID(id) // คืนค่า s.productRepo.FindByID(id)
} // ปิด block การทำงานปัจจุบัน

func (s *productService) CreateProduct(product model.Product) (model.Product, error) { // ประกาศฟังก์ชัน CreateProduct
	return s.productRepo.Create(product) // คืนค่า s.productRepo.Create(product)
} // ปิด block การทำงานปัจจุบัน

func (s *productService) UpdateProduct(id uint, product model.Product) (model.Product, error) { // ประกาศฟังก์ชัน UpdateProduct
	oldProduct, err := s.productRepo.FindByID(id) // กำหนดค่า old Product, err จาก s.productRepo.FindByID(id)
	if err != nil {                               // ตรวจเงื่อนไข err != nil
		return oldProduct, err // คืนค่า oldProduct, err
	} // ปิด block การทำงานปัจจุบัน

	oldProduct.Name = product.Name               // กำหนดค่า Name ของ productเก่า จาก Name ของ productใหม่
	oldProduct.Description = product.Description // กำหนดค่า Description ของ productเก่า จาก Description ของ productใหม่
	oldProduct.Price = product.Price             // กำหนดค่า Price ของ productเก่า จาก Price ของ productใหม่
	oldProduct.Stock = product.Stock             // กำหนดค่า Stock ของ productเก่า จาก Stock ของ productใหม่
	oldProduct.CategoryID = product.CategoryID   // กำหนดค่า Category ID ของ productเก่า จาก Category ID ของ productใหม่
	oldProduct.Image = product.Image             // กำหนดค่า Image ของ productเก่า จาก Image ของ productใหม่

	return s.productRepo.Update(oldProduct) // คืนค่า s.productRepo.Update(oldProduct)
} // ปิด block การทำงานปัจจุบัน

func (s *productService) DeleteProduct(id uint) error { // ประกาศฟังก์ชัน DeleteProduct
	return s.productRepo.Delete(id) // คืนค่า s.productRepo.Delete(id)
} // ปิด block การทำงานปัจจุบัน

func (s *productService) GetProductsByCategory(categoryID uint) ([]model.Product, error) { // ประกาศฟังก์ชัน GetProductsByCategory
	return s.productRepo.FindByCategoryID(categoryID) // คืนค่า s.productRepo.FindByCategoryID(categoryID)
} // ปิด block การทำงานปัจจุบัน

func (s *productService) GetCategories() ([]model.Category, error) { // ประกาศฟังก์ชัน GetCategories
	return s.productRepo.FindAllCategories() // คืนค่า s.productRepo.FindAllCategories()
} // ปิด block การทำงานปัจจุบัน

func (s *productService) GetCategoryByID(id uint) (model.Category, error) { // ประกาศฟังก์ชัน GetCategoryByID
	return s.productRepo.FindCategoryByID(id) // คืนค่า s.productRepo.FindCategoryByID(id)
} // ปิด block การทำงานปัจจุบัน

func (s *productService) CreateCategory(category model.Category) (model.Category, error) { // ประกาศฟังก์ชัน CreateCategory
	return s.productRepo.CreateCategory(category) // คืนค่า s.productRepo.CreateCategory(category)
} // ปิด block การทำงานปัจจุบัน

func (s *productService) UpdateCategory(id uint, category model.Category) (model.Category, error) { // ประกาศฟังก์ชัน UpdateCategory
	oldCategory, err := s.productRepo.FindCategoryByID(id) // กำหนดค่า old Category, err จาก s.productRepo.FindCategoryByID(id)
	if err != nil {                                        // ตรวจเงื่อนไข err != nil
		return oldCategory, err // คืนค่า oldCategory, err
	} // ปิด block การทำงานปัจจุบัน
	oldCategory.Name = category.Name                 // กำหนดค่า Name ของ old Category จาก Name ของ category
	return s.productRepo.UpdateCategory(oldCategory) // คืนค่า s.productRepo.UpdateCategory(oldCategory)
} // ปิด block การทำงานปัจจุบัน

func (s *productService) DeleteCategory(id uint) error { // ประกาศฟังก์ชัน DeleteCategory
	return s.productRepo.DeleteCategory(id) // คืนค่า s.productRepo.DeleteCategory(id)
} // ปิด block การทำงานปัจจุบัน
