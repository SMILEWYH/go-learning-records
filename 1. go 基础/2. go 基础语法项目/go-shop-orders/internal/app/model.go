package app

// 金额单位固定为“分”（第 03、12 章），客户端不能提交订单单价或总额。
type Cents int64
type OrderStatus string

const (
	Confirmed OrderStatus = "confirmed"
	Shipped   OrderStatus = "shipped"
	Canceled  OrderStatus = "canceled"
)

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price Cents  `json:"price_cents"`
	Stock int    `json:"stock"`
}
type ItemInput struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
type CreateOrder struct {
	RequestID string      `json:"request_id"`
	Customer  string      `json:"customer"`
	Items     []ItemInput `json:"items"`
}
type OrderItem struct {
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	UnitPrice Cents  `json:"unit_price_cents"`
	Quantity  int    `json:"quantity"`
}
type Order struct {
	ID        string      `json:"id"`
	RequestID string      `json:"request_id"`
	Customer  string      `json:"customer"`
	Items     []OrderItem `json:"items"`
	Total     Cents       `json:"total_cents"`
	Status    OrderStatus `json:"status"`
	CreatedAt string      `json:"created_at"`
}
type State struct {
	Version  int                `json:"version"`
	NextID   int                `json:"next_id"`
	Products map[string]Product `json:"products"`
	Orders   map[string]Order   `json:"orders"`
}
