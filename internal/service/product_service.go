package service

import (
	"github.com/Incipe-win/ai-tshirt-shop/internal/model"
	"github.com/Incipe-win/ai-tshirt-shop/internal/repository"
)

type ProductService struct {
	productRepo *repository.ProductRepository
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	BasePrice   float64 `json:"base_price" binding:"required,gt=0"`
	Category    string  `json:"category"`
	Brand       string  `json:"brand"`
}

type ProductResponse struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	BasePrice    float64 `json:"base_price"`
	Category     string  `json:"category"`
	Brand        string  `json:"brand"`
	IsActive     bool    `json:"is_active"`
	ImageURL     string  `json:"image_url"`
	DesignPrompt string  `json:"design_prompt"`
	DesignStyle  string  `json:"design_style"`
	CreatorName  string  `json:"creator_name"`
	CreatorID    *uint   `json:"creator_id"`
	CreatedAt    string  `json:"created_at"`
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (s *ProductService) CreateProduct(req *CreateProductRequest) (*ProductResponse, error) {
	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		BasePrice:   req.BasePrice,
		Category:    req.Category,
		Brand:       req.Brand,
		IsActive:    true,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

func (s *ProductService) GetAllProducts() ([]ProductResponse, error) {
	products, err := s.productRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []ProductResponse
	for _, product := range products {
		responses = append(responses, *s.toProductResponse(&product))
	}

	return responses, nil
}

func (s *ProductService) GetProductByID(id uint) (*ProductResponse, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

func (s *ProductService) GetProductsByCategory(category string) ([]ProductResponse, error) {
	products, err := s.productRepo.FindByCategory(category)
	if err != nil {
		return nil, err
	}

	var responses []ProductResponse
	for _, product := range products {
		responses = append(responses, *s.toProductResponse(&product))
	}

	return responses, nil
}

func (s *ProductService) toProductResponse(product *model.Product) *ProductResponse {
	return &ProductResponse{
		ID:           product.ID,
		Name:         product.Name,
		Description:  product.Description,
		BasePrice:    product.BasePrice,
		Category:     product.Category,
		Brand:        product.Brand,
		IsActive:     product.IsActive,
		ImageURL:     product.ImageURL,
		DesignPrompt: product.DesignPrompt,
		DesignStyle:  product.DesignStyle,
		CreatorName:  product.CreatorName,
		CreatorID:    product.CreatorID,
		CreatedAt:    product.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}