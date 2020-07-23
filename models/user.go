package models

import basemodel "web-backend-patal/config"

type User struct {
	basemodel.BaseModel
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
}
