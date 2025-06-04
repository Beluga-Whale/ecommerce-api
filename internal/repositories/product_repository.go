package repositories

import (
	"errors"
	"strconv"

	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type ProductRepositoryInterface interface {
	Create(product *models.Product) error
	FindByID(id uint) (*models.Product, error)
	FindAll(page uint, limit uint, minPrice int64, maxPrice int64, searchName string, category string) (productList []models.Product ,pageTotal int64,err error) 
	Update(product *models.Product) error
	Delete(id uint) error
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

	err := r.db.First(&product,id).Error

	if errors.Is(err, gorm.ErrRecordNotFound){
		return nil,nil
	}

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) FindAll(page uint, limit uint, minPrice int64, maxPrice int64, searchName string, category string) (productList []models.Product ,pageTotal int64,err error) {
	var products []models.Product
	var total int64

	productQuery := r.db.Model(&models.Product{})

	// NOTE - เช็คว่า Category เป็นค่าว่างมาไหม
	if category != "" {
		// NOTE - ถ้าไม่เป็นค่าว่างต้องเปลี่ยนให้เป็น int
    	catID, err := strconv.Atoi(category)
    	if err == nil {
        	productQuery = productQuery.Where("category_id = ?", catID)
    	}
	}

	// NOTE - Where price
	productQuery = productQuery.Where("price >= ? AND price <= ?",minPrice,maxPrice)

	if searchName != "" {
		productQuery = productQuery.Where("name ILIKE ?", "%"+searchName+"%")
	}

	if err := productQuery.Count(&total).Error; err != nil{
		return nil,0,err
	}
	
	offset := (page -1 ) * limit
	pageTotal = (total + int64(limit) - 1) / int64(limit)

	err = productQuery.Preload("Category").Offset(int(offset)).Limit(int(limit)).Find(&products).Error
	return products, pageTotal,err
}

func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}