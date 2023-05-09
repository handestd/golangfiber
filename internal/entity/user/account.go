package user

import "github.com/golang-jwt/jwt/v5"

type Account struct {
	Id       int    `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
	Email    string `json:"email" gorm:"column:email"`
	Balance  int64  `json:"balance" gorm:"column:balance"`
	LastIp   string `json:"lastip" gorm:"column:lastip"`
	Avatar   string `json:"avatar" gorm:"column:avatar"`
	Role     string `json:"role" gorm:"column:role"`
	GroupId  string `json:"groupid" gorm:"column:groupid"`
}

func (Account) TableName() string {
	return "user_account"
}

type UserClaim struct {
	jwt.RegisteredClaims
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
