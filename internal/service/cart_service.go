package service

import (
	"errors"

	"github.com/Incipe-win/ai-tshirt-shop/internal/model"
	"github.com/Incipe-win/ai-tshirt-shop/internal/repository"
	"gorm.io/gorm"
)

type CartService struct {
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
	designRepo  *repository.DesignRepository
}

type AddToCartRequest struct {
	ProductID uint  `json:"product_id" binding:"required"`
	DesignID  *uint `json:"design_id"`
	Quantity  int   `json:"quantity" binding:"required,gt=0"`
}

type UpdateCartRequest struct {
	Quantity int `json:"quantity" binding:"required,gt=0"`
}

type CartItemResponse struct {
	ID       uint                     `json:"id"`
	Product  *ProductResponse         `json:"product"`
	Design   *CartDesignResponse      `json:"design"`
	Quantity int                      `json:"quantity"`
}

type CartDesignResponse struct {
	ID       uint   `json:"id"`
	Prompt   string `json:"prompt"`
	ImageURL string `json:"image_url"`
	Style    string `json:"style"`
}

type CartResponse struct {
	Items      []CartItemResponse `json:"items"`
	TotalItems int                `json:"total_items"`
	TotalValue float64            `json:"total_value"`
}

func NewCartService(cartRepo *repository.CartRepository, productRepo *repository.ProductRepository, designRepo *repository.DesignRepository) *CartService {
	return &CartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
		designRepo:  designRepo,
	}
}

func (s *CartService) AddToCart(userID uint, req *AddToCartRequest) error {
	// 验证商品是否存在
	product, err := s.productRepo.FindByID(req.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	if !product.IsActive {
		return errors.New("product is not available")
	}

	var designID uint
	
	// 如果提供了设计ID，验证设计是否存在且属于该用户
	if req.DesignID != nil {
		design, err := s.designRepo.FindByID(*req.DesignID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("design not found")
			}
			return err
		}

		if design.UserID != userID {
			return errors.New("design does not belong to user")
		}
		
		designID = *req.DesignID
	} else {
		// 如果没有提供设计ID，获取用户的设计作品
		designs, err := s.designRepo.FindByUserID(userID)
		if err != nil {
			return err
		}
		
		// 如果用户没有设计作品，创建一个默认设计
		if len(designs) == 0 {
			defaultDesign := &model.Design{
				UserID:    userID,
				Prompt:    "默认创意设计",
				Category:  "general",
				ImageURL:  "/static/images/ai.png",
				Style:     "",
			}
			
			if err := s.designRepo.Create(defaultDesign); err != nil {
				return err
			}
			designID = defaultDesign.ID
		} else {
			// 使用用户最新的设计
			designID = designs[0].ID
		}
	}

	cartItem := &model.CartItem{
		UserID:    userID,
		ProductID: req.ProductID,
		DesignID:  designID,
		Quantity:  req.Quantity,
	}

	return s.cartRepo.AddItem(cartItem)
}

func (s *CartService) GetCart(userID uint) (*CartResponse, error) {
	items, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	var cartItems []CartItemResponse
	var totalItems int
	var totalValue float64

	for _, item := range items {
		cartItem := CartItemResponse{
			ID:       item.ID,
			Quantity: item.Quantity,
		}

		// 构造产品信息
		if item.Product.ID != 0 {
			cartItem.Product = &ProductResponse{
				ID:          item.Product.ID,
				Name:        item.Product.Name,
				Description: item.Product.Description,
				BasePrice:   item.Product.BasePrice,
				Category:    item.Product.Category,
				Brand:       item.Product.Brand,
				IsActive:    item.Product.IsActive,
				ImageURL:    item.Product.ImageURL,
				CreatedAt:   item.Product.CreatedAt.Format("2006-01-02 15:04:05"),
			}
		}

		// 构造设计信息
		if item.Design.ID != 0 {
			cartItem.Design = &CartDesignResponse{
				ID:       item.Design.ID,
				Prompt:   item.Design.Prompt,
				ImageURL: item.Design.ImageURL,
				Style:    item.Design.Style,
			}
		}

		cartItems = append(cartItems, cartItem)
		totalItems += item.Quantity
		if item.Product.ID != 0 {
			totalValue += item.Product.BasePrice * float64(item.Quantity)
		}
	}

	return &CartResponse{
		Items:      cartItems,
		TotalItems: totalItems,
		TotalValue: totalValue,
	}, nil
}

func (s *CartService) UpdateCartItem(userID uint, itemID uint, req *UpdateCartRequest) error {
	// 验证购物车项是否属于该用户
	item, err := s.cartRepo.FindByID(itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("cart item not found")
		}
		return err
	}

	if item.UserID != userID {
		return errors.New("cart item does not belong to user")
	}

	return s.cartRepo.UpdateQuantity(itemID, req.Quantity)
}

func (s *CartService) RemoveFromCart(userID uint, itemID uint) error {
	return s.cartRepo.RemoveItem(itemID, userID)
}

func (s *CartService) ClearCart(userID uint) error {
	return s.cartRepo.ClearCartByUserID(userID)
}