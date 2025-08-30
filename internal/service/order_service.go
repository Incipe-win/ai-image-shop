package service

import (
	"errors"

	"github.com/Incipe-win/ai-tshirt-shop/internal/model"
	"github.com/Incipe-win/ai-tshirt-shop/internal/repository"
	"gorm.io/gorm"
)

type OrderService struct {
	orderRepo *repository.OrderRepository
	cartRepo  *repository.CartRepository
}

type CreateOrderRequest struct {
	CartItemIDs []uint `json:"cart_item_ids" binding:"required,min=1"`
}

type OrderItemResponse struct {
	ID             uint    `json:"id"`
	ProductName    string  `json:"product_name"`
	DesignImageURL string  `json:"design_image_url"`
	Size           string  `json:"size"`
	Color          string  `json:"color"`
	Price          float64 `json:"price"`
	Quantity       int     `json:"quantity"`
}

type OrderResponse struct {
	ID          uint                `json:"id"`
	OrderSN     string              `json:"order_sn"`
	TotalAmount float64             `json:"total_amount"`
	Status      model.OrderStatus   `json:"status"`
	CreatedAt   string              `json:"created_at"`
	OrderItems  []OrderItemResponse `json:"order_items"`
}

func NewOrderService(orderRepo *repository.OrderRepository, cartRepo *repository.CartRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		cartRepo:  cartRepo,
	}
}

func (s *OrderService) CreateOrder(userID uint, req *CreateOrderRequest) (*OrderResponse, error) {
	// 获取购物车项
	var cartItems []model.CartItem
	var totalAmount float64

	for _, cartItemID := range req.CartItemIDs {
		cartItem, err := s.cartRepo.FindByID(cartItemID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("cart item not found")
			}
			return nil, err
		}

		if cartItem.UserID != userID {
			return nil, errors.New("cart item does not belong to user")
		}

		cartItems = append(cartItems, *cartItem)
		totalAmount += cartItem.Product.BasePrice * float64(cartItem.Quantity)
	}

	if len(cartItems) == 0 {
		return nil, errors.New("no valid cart items found")
	}

	// 创建订单
	order := &model.Order{
		UserID:      userID,
		OrderSN:     s.orderRepo.GenerateOrderSN(),
		TotalAmount: totalAmount,
		Status:      model.OrderStatusPending,
	}

	// 创建订单项
	var orderItems []model.OrderItem
	for _, cartItem := range cartItems {
		orderItem := model.OrderItem{
			ProductName:    cartItem.Product.Name,
			DesignImageURL: cartItem.Design.ImageURL,
			Size:           cartItem.Size,
			Color:          cartItem.Color,
			Price:          cartItem.Product.BasePrice,
			Quantity:       cartItem.Quantity,
		}
		orderItems = append(orderItems, orderItem)
	}

	// 在事务中创建订单和订单项
	if err := s.orderRepo.CreateWithItems(order, orderItems); err != nil {
		return nil, err
	}

	// 从购物车中删除已下单的商品
	for _, cartItemID := range req.CartItemIDs {
		s.cartRepo.RemoveItem(cartItemID, userID)
	}

	// 返回创建的订单信息
	return s.GetOrderByID(userID, order.ID)
}

func (s *OrderService) GetOrdersByUserID(userID uint) ([]OrderResponse, error) {
	orders, err := s.orderRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []OrderResponse
	for _, order := range orders {
		responses = append(responses, *s.toOrderResponse(&order))
	}

	return responses, nil
}

func (s *OrderService) GetOrderByID(userID uint, orderID uint) (*OrderResponse, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	if order.UserID != userID {
		return nil, errors.New("order does not belong to user")
	}

	return s.toOrderResponse(order), nil
}

func (s *OrderService) GetOrderByOrderSN(userID uint, orderSN string) (*OrderResponse, error) {
	order, err := s.orderRepo.FindByOrderSN(orderSN)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	if order.UserID != userID {
		return nil, errors.New("order does not belong to user")
	}

	return s.toOrderResponse(order), nil
}

func (s *OrderService) UpdateOrderStatus(userID uint, orderID uint, status model.OrderStatus) error {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("order not found")
		}
		return err
	}

	if order.UserID != userID {
		return errors.New("order does not belong to user")
	}

	return s.orderRepo.UpdateStatus(orderID, status)
}

func (s *OrderService) toOrderResponse(order *model.Order) *OrderResponse {
	var orderItems []OrderItemResponse
	for _, item := range order.OrderItems {
		orderItems = append(orderItems, OrderItemResponse{
			ID:             item.ID,
			ProductName:    item.ProductName,
			DesignImageURL: item.DesignImageURL,
			Size:           item.Size,
			Color:          item.Color,
			Price:          item.Price,
			Quantity:       item.Quantity,
		})
	}

	return &OrderResponse{
		ID:          order.ID,
		OrderSN:     order.OrderSN,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
		OrderItems:  orderItems,
	}
}