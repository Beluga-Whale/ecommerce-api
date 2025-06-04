package repositories

import (
	"gorm.io/gorm"
)

type OrderRepositoryInterface interface {
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository{
	return &OrderRepository{db:db}
}

// func FindProductByID (products []models.Product)