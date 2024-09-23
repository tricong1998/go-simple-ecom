package dto

type CreateUserPoint struct {
	OrderId uint `json:"order_id" gorm:"unique"`
	UserId  uint `json:"user_id" `
	Amount  uint `json:"amount"`
}
