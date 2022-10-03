package model

const UserRoleName = "user"
const SellerRoleName = "seller"
const Level1RoleName = "level1"

type Role struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}
