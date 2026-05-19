package repository // ประกาศ package repository

import ( // เริ่มส่วน import package ที่ไฟล์นี้ต้องใช้
	"backend-go/internal/model" // นำเข้า package backend-go/internal/model

	"gorm.io/gorm" // นำเข้า package gorm.io/gorm
) // ปิด block หรือกลุ่มรายการก่อนหน้า

type ProductRepository interface { // ประกาศ interface ProductRepository
	FindAll() ([]model.Product, error)                              // ทำงานคำสั่ง FindAll() ([]model.Product, error)
	FindByID(id uint) (model.Product, error)                        // ทำงานคำสั่ง FindByID(id uint) (model.Product, error)
	Create(product model.Product) (model.Product, error)            // ทำงานคำสั่ง Create(product model.Product) (model.Product, error)
	Update(product model.Product) (model.Product, error)            // ทำงานคำสั่ง Update(product model.Product) (model.Product, error)
	Delete(id uint) error                                           // ทำงานคำสั่ง Delete(id uint) error
	FindByCategoryID(categoryID uint) ([]model.Product, error)      // ทำงานคำสั่ง FindByCategoryID(categoryID uint) ([]model.Product, error)
	FindAllCategories() ([]model.Category, error)                   // ทำงานคำสั่ง FindAllCategories() ([]model.Category, error)
	FindCategoryByID(id uint) (model.Category, error)               // ทำงานคำสั่ง FindCategoryByID(id uint) (model.Category, error)
	CreateCategory(category model.Category) (model.Category, error) // ทำงานคำสั่ง CreateCategory(category model.Category) (model.Category, error)
	UpdateCategory(category model.Category) (model.Category, error) // ทำงานคำสั่ง UpdateCategory(category model.Category) (model.Category, error)
	DeleteCategory(id uint) error                                   // ทำงานคำสั่ง DeleteCategory(id uint) error
} // ปิด block การทำงานปัจจุบัน

type productRepository struct { // ประกาศ struct productRepository
	db *gorm.DB // ทำงานคำสั่ง db *gorm.DB
} // ปิด block การทำงานปัจจุบัน

func NewProductRepository(db *gorm.DB) ProductRepository { // ประกาศฟังก์ชัน NewProductRepository
	return &productRepository{db: db} // คืนค่า &productRepository{db: db}
} // ปิด block การทำงานปัจจุบัน

func (r *productRepository) FindAll() ([]model.Product, error) { // ประกาศฟังก์ชัน FindAll
	var products []model.Product                          // ประกาศตัวแปร products []model.Product
	err := r.db.Preload("Category").Find(&products).Error // กำหนดค่า error จาก r.db.Preload("Category").Find(&products).Error
	return products, err                                  // คืนค่า products, err
} // ปิด block การทำงานปัจจุบัน

func (r *productRepository) FindByID(id uint) (model.Product, error) { // ประกาศฟังก์ชัน FindByID
	var product model.Product                                 // ประกาศตัวแปร product model.Product
	err := r.db.Preload("Category").First(&product, id).Error // กำหนดค่า error จาก r.db.Preload("Category").First(&product, id).Error
	return product, err                                       // คืนค่า product, err
} // ปิด block การทำงานปัจจุบัน

func (r *productRepository) Create(product model.Product) (model.Product, error) { // ประกาศฟังก์ชัน Create
	err := r.db.Create(&product).Error // กำหนดค่า error จาก r.db.Create(&product).Error
	return product, err                // คืนค่า product, err
} // ปิด block การทำงานปัจจุบัน

func (r *productRepository) Update(product model.Product) (model.Product, error) { // ประกาศฟังก์ชัน Update
	err := r.db.Save(&product).Error // กำหนดค่า error จาก r.db.Save(&product).Error
	return product, err              // คืนค่า product, err
} // ปิด block การทำงานปัจจุบัน

func (r *productRepository) Delete(id uint) error { // ประกาศฟังก์ชัน Delete
	return r.db.Transaction(func(tx *gorm.DB) error { // คืนค่า r.db.Transaction(func(tx *gorm.DB) error
		if err := tx.Where("product_id = ?", id).Delete(&model.WishlistItem{}).Error; err != nil { // ตรวจเงื่อนไข err := tx.Where("product_id = ?", id).Delete(&model.WishlistItem{}).Error; err != nil
			return err // คืนค่า err
		} // ปิด block การทำงานปัจจุบัน
		if err := tx.Where("product_id = ?", id).Delete(&model.CartItem{}).Error; err != nil { // ตรวจเงื่อนไข err := tx.Where("product_id = ?", id).Delete(&model.CartItem{}).Error; err != nil
			return err // คืนค่า err
		} // ปิด block การทำงานปัจจุบัน

		return tx.Delete(&model.Product{}, id).Error // คืนค่า tx.Delete(&model.Product{}, id).Error
	}) // ทำงานคำสั่ง })
} // ปิด block การทำงานปัจจุบัน

func (r *productRepository) FindByCategoryID(categoryID uint) ([]model.Product, error) { // ประกาศฟังก์ชัน FindByCategoryID
	var products []model.Product                                                               // ประกาศตัวแปร products []model.Product
	err := r.db.Preload("Category").Where("category_id = ?", categoryID).Find(&products).Error // กำหนดค่า error จาก r.db.Preload("Category").Where("category_id = ?", categoryID).Find(&products).Error
	return products, err                                                                       // คืนค่า products, err
} // ปิด block การทำงานปัจจุบัน

func (r *productRepository) FindAllCategories() ([]model.Category, error) { // ประกาศฟังก์ชัน FindAllCategories
	var categories []model.Category     // ประกาศตัวแปร categories []model.Category
	err := r.db.Find(&categories).Error // กำหนดค่า error จาก r.db.Find(&categories).Error
	return categories, err              // คืนค่า categories, err
} // ปิด block การทำงานปัจจุบัน

func (r *productRepository) FindCategoryByID(id uint) (model.Category, error) { // ประกาศฟังก์ชัน FindCategoryByID
	var category model.Category            // ประกาศตัวแปร category model.Category
	err := r.db.First(&category, id).Error // กำหนดค่า error จาก r.db.First(&category, id).Error
	return category, err                   // คืนค่า category, err
} // ปิด block การทำงานปัจจุบัน

func (r *productRepository) CreateCategory(category model.Category) (model.Category, error) { // ประกาศฟังก์ชัน CreateCategory
	err := r.db.Create(&category).Error // กำหนดค่า error จาก r.db.Create(&category).Error
	return category, err                // คืนค่า category, err
} // ปิด block การทำงานปัจจุบัน

func (r *productRepository) UpdateCategory(category model.Category) (model.Category, error) { // ประกาศฟังก์ชัน UpdateCategory
	err := r.db.Save(&category).Error // กำหนดค่า error จาก r.db.Save(&category).Error
	return category, err              // คืนค่า category, err
} // ปิด block การทำงานปัจจุบัน

func (r *productRepository) DeleteCategory(id uint) error { // ประกาศฟังก์ชัน DeleteCategory
	return r.db.Transaction(func(tx *gorm.DB) error {
		productIDs := tx.Model(&model.Product{}).
			Select("id").
			Where("category_id = ?", id)

		if err := tx.Where("product_id IN (?)", productIDs).
			Delete(&model.WishlistItem{}).Error; err != nil {
			return err
		}

		if err := tx.Where("product_id IN (?)", productIDs).
			Delete(&model.CartItem{}).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().
			Where("category_id = ?", id).
			Delete(&model.Product{}).Error; err != nil {
			return err
		}

		return tx.Delete(&model.Category{}, id).Error
	})
} // ปิด block การทำงานปัจจุบัน
