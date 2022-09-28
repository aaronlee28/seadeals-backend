package model

const UserRoleName = "user"
const SellerRoleName = "seller"

type Role struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}
