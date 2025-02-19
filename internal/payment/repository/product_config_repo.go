package repository

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	. "github.com/mavrk-mose/pay/internal/payment/models"
)

type ProductConfigRepo struct {
	DB *sqlx.DB
}

func NewProductConfigRepo(db *sqlx.DB) *ProductConfigRepo {
	return &ProductConfigRepo{DB: db}
}

func (r *ProductConfigRepo) GetProductConfig(ctx *gin.Context, productName string) (*ProductConfiguration, error) {
	var config ProductConfiguration
	query := "SELECT id, product_name, fee_percentage FROM product_configurations WHERE product_name = $1"
	err := r.DB.GetContext(ctx, &config, query, productName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product configuration: %v", err)
	}
	return &config, nil
}
