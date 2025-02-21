package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"belcamp/internal/models"

	"gorm.io/gorm"
)

// ProductService interface defines the methods for product operations
type ProductService interface {
	// Product operations
	GetProductByID(id uint) (*models.Product, error)
	GetProductBySlug(slug string) (*models.Product, error)
	CreateProduct(product *models.Product) error
	UpdateProduct(product *models.Product) error
	DeleteProduct(id uint) error
	ListProducts(page, pageSize int, filters map[string]interface{}) ([]models.Product, int64, error)

	// Variant operations
	GetVariantByID(id uint) (*models.ProductVariant, error)
	CreateVariant(variant *models.ProductVariant) error
	UpdateVariant(variant *models.ProductVariant) error
	DeleteVariant(id uint) error
	ListVariants(productID uint) ([]models.ProductVariant, error)
	UpdateVariantStock(variantID uint, quantity int) error

	// Category operations
	ListProductsByCategory(categoryID uint, page, pageSize int) ([]models.Product, int64, error)

	// Price operations
	UpdateProductPrices(productID uint, prices map[string]float64) error
	UpdateVariantPrices(variantID uint, prices map[string]float64) error

	// Media operations
	UpdateProductPhotos(productID uint, photos []string) error
	UpdateProductColorPhotos(productID uint, colorPhotos map[string][]string) error

	// Search operations
	SearchProducts(query string, page, pageSize int) ([]models.Product, int64, error)
}

type productService struct {
	db *gorm.DB
}

// NewProductService creates a new ProductService instance
func NewProductService(db *gorm.DB) ProductService {
	return &productService{
		db: db,
	}
}

// GetProductByID retrieves a product by its ID
func (s *productService) GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	result := s.db.
		Preload("Category").
		Preload("Variants").
		First(&product, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, result.Error
	}

	return &product, nil
}

// GetProductBySlug retrieves a product by its slug
func (s *productService) GetProductBySlug(slug string) (*models.Product, error) {
	var product models.Product
	result := s.db.
		Where("slug = ?", slug).
		Preload("Category").
		Preload("Variants").
		First(&product)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, result.Error
	}

	return &product, nil
}

// CreateProduct creates a new product
func (s *productService) CreateProduct(product *models.Product) error {
	// Generate slug if not provided
	if product.Slug == "" {
		product.Slug = generateSlug(product.Name)
	}

	// Validate slug uniqueness
	var count int64
	s.db.Model(&models.Product{}).Where("slug = ?", product.Slug).Count(&count)
	if count > 0 {
		return errors.New("slug must be unique")
	}

	return s.db.Create(product).Error
}

// UpdateProduct updates an existing product
func (s *productService) UpdateProduct(product *models.Product) error {
	// Check if product exists
	existing := &models.Product{}
	if err := s.db.First(existing, product.ID).Error; err != nil {
		return err
	}

	// If slug is being changed, validate uniqueness
	if product.Slug != existing.Slug {
		var count int64
		s.db.Model(&models.Product{}).
			Where("slug = ? AND id != ?", product.Slug, product.ID).
			Count(&count)
		if count > 0 {
			return errors.New("slug must be unique")
		}
	}

	return s.db.Model(product).Updates(product).Error
}

// DeleteProduct soft deletes a product
func (s *productService) DeleteProduct(id uint) error {
	return s.db.Delete(&models.Product{}, id).Error
}

// ListProducts returns a paginated list of products with optional filters
func (s *productService) ListProducts(page, pageSize int, filters map[string]interface{}) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := s.db.Model(&models.Product{}).
		Preload("Category").
		Preload("Variants")

	// Apply filters
	for key, value := range filters {
		query = query.Where(key, value)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated products
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// GetVariantByID retrieves a product variant by its ID
func (s *productService) GetVariantByID(id uint) (*models.ProductVariant, error) {
	var variant models.ProductVariant
	result := s.db.First(&variant, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("variant not found")
		}
		return nil, result.Error
	}

	return &variant, nil
}

// CreateVariant creates a new product variant
func (s *productService) CreateVariant(variant *models.ProductVariant) error {
	// Validate SKU uniqueness
	var count int64
	s.db.Model(&models.ProductVariant{}).Where("sku = ?", variant.SKU).Count(&count)
	if count > 0 {
		return errors.New("SKU must be unique")
	}

	return s.db.Create(variant).Error
}

// UpdateVariant updates an existing product variant
func (s *productService) UpdateVariant(variant *models.ProductVariant) error {
	// Check if variant exists
	existing := &models.ProductVariant{}
	if err := s.db.First(existing, variant.ID).Error; err != nil {
		return err
	}

	// If SKU is being changed, validate uniqueness
	if variant.SKU != existing.SKU {
		var count int64
		s.db.Model(&models.ProductVariant{}).
			Where("sku = ? AND id != ?", variant.SKU, variant.ID).
			Count(&count)
		if count > 0 {
			return errors.New("SKU must be unique")
		}
	}

	return s.db.Model(variant).Updates(variant).Error
}

// DeleteVariant soft deletes a product variant
func (s *productService) DeleteVariant(id uint) error {
	return s.db.Delete(&models.ProductVariant{}, id).Error
}

// ListVariants returns all variants for a product
func (s *productService) ListVariants(productID uint) ([]models.ProductVariant, error) {
	var variants []models.ProductVariant
	err := s.db.Where("product_id = ?", productID).Find(&variants).Error
	return variants, err
}

// UpdateVariantStock updates the stock quantity for a variant
func (s *productService) UpdateVariantStock(variantID uint, quantity int) error {
	result := s.db.Model(&models.ProductVariant{}).
		Where("id = ?", variantID).
		Update("availability", quantity)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("variant not found")
	}

	return nil
}

// ListProductsByCategory returns products in a specific category
func (s *productService) ListProductsByCategory(categoryID uint, page, pageSize int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := s.db.Model(&models.Product{}).
		Where("category_id = ?", categoryID).
		Preload("Category").
		Preload("Variants")

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated products
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// UpdateProductPrices updates the prices JSON field for a product
func (s *productService) UpdateProductPrices(productID uint, prices map[string]float64) error {
	pricesJSON, err := json.Marshal(prices)
	if err != nil {
		return err
	}

	return s.db.Model(&models.Product{}).
		Where("id = ?", productID).
		Update("prices", string(pricesJSON)).Error
}

// UpdateVariantPrices updates the prices JSON field for a variant
func (s *productService) UpdateVariantPrices(variantID uint, prices map[string]float64) error {
	pricesJSON, err := json.Marshal(prices)
	if err != nil {
		return err
	}

	return s.db.Model(&models.ProductVariant{}).
		Where("id = ?", variantID).
		Update("prices", string(pricesJSON)).Error
}

// UpdateProductPhotos updates the photos JSON field for a product
func (s *productService) UpdateProductPhotos(productID uint, photos []string) error {
	photosJSON, err := json.Marshal(photos)
	if err != nil {
		return err
	}

	return s.db.Model(&models.Product{}).
		Where("id = ?", productID).
		Update("photos", string(photosJSON)).Error
}

// UpdateProductColorPhotos updates the color photos JSON field for a product
func (s *productService) UpdateProductColorPhotos(productID uint, colorPhotos map[string][]string) error {
	colorPhotosJSON, err := json.Marshal(colorPhotos)
	if err != nil {
		return err
	}

	return s.db.Model(&models.Product{}).
		Where("id = ?", productID).
		Update("color_photos", string(colorPhotosJSON)).Error
}

// SearchProducts searches products by name, description, or SKU
func (s *productService) SearchProducts(query string, page, pageSize int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	searchQuery := fmt.Sprintf("%%%s%%", strings.ToLower(query))

	baseQuery := s.db.Model(&models.Product{}).
		Preload("Category").
		Preload("Variants").
		Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchQuery, searchQuery)

	// Get total count
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := baseQuery.Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// Helper function to generate slug from name
func generateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)

	return slug
}
