package repositories

import (
	"errors"

	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type ProductRepositoryInterface interface {
	Create(product *models.Product) error
	FindByID(id uint) (*models.Product, error)
	FindAll(page uint, limit uint, minPrice int64, maxPrice int64, searchName string, categoryIDs []int,sizeIDs []string) (productList []models.Product ,pageTotal int64,err error) 
	Update(product *models.Product) error
	Delete(id uint) error
	// DeleteImageByProductID(productID uint) error
}

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) FindByID(id uint) (*models.Product, error){
	var product models.Product

	err := r.db.Preload("Variants").Preload("Images").First(&product,id).Error

	if errors.Is(err, gorm.ErrRecordNotFound){
		return nil,nil
	}

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) FindAll(page uint, limit uint, minPrice int64, maxPrice int64, searchName string, categoryIDs []int,sizeIDs []string) (productList []models.Product ,pageTotal int64,err error) {
	var products []models.Product
	var total int64

	productQuery := r.db.Model(&models.Product{})

	// NOTE - เช็คว่า Category มากกว่า 0 ไหม
	if len(categoryIDs) >0 {
        productQuery = productQuery.Where("category_id IN ?", categoryIDs)
	}

	if len(sizeIDs) > 0 {
		productQuery = productQuery.
			Joins("JOIN product_variants ON product_variants.product_id = products.id").
			Where("product_variants.size IN ?",sizeIDs).
			Group("products.id")
	}

	if searchName != "" {
		productQuery = productQuery.Where("name ILIKE ?", "%"+searchName+"%")
	}

	if err := productQuery.Count(&total).Error; err != nil{
		return nil,0,err
	}
	
	offset := (page -1 ) * limit
	pageTotal = (total + int64(limit) - 1) / int64(limit)

	err = productQuery.Preload("Category").Preload("Variants").Preload("Images").Offset(int(offset)).Limit(int(limit)).Order("id").Find(&products).Error
	return products, pageTotal,err
}

func (r *ProductRepository) Update(product *models.Product) error {
	// NOTE - ใช้ transaction 
	tx := r.db.Begin()

	if tx.Error != nil{
		return tx.Error
	}

	// NOTE - ลบ variants เก่าออกก่อน
	if err := tx.Unscoped().Where("product_id = ?", product.ID).Delete(&models.ProductVariant{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// NOTE - ลบ image เก่าออกก่อน
	if err := tx.Unscoped().Where("product_id = ?", product.ID).Delete(&models.ProductImage{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Save(product).Error; err !=nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *ProductRepository) Delete(id uint) error {
	// NOTE- ใช้ transaction
	tx := r.db.Begin()

	if tx.Error != nil {
		return tx.Error
	}

	// NOTE - where หา productID เพื่อลบ variant
	if err := tx.Where("product_id = ?",id).Delete(&models.ProductVariant{}).Error; err !=nil{
		tx.Rollback()
		return err
	}

	// NOTE- ลบ Product
	if err := tx.Delete(&models.Product{}, id).Error; err !=nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// func (r*ProductRepository) DeleteImageByProductID(productID uint) error {
// 	return r.db.Where("product_id",productID).Delete(&models.ProductImage{}).Error
// }