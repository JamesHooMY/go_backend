package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(20);not null" json:"name"`
	Age      int    `gorm:"type:int(3)" json:"age"`
	Mobile   string `gorm:"type:varchar(11);unique" json:"mobile"`
	Email    string `gorm:"type:varchar(50);not null;unique" json:"email"`
	Password string `gorm:"type:varchar(20);not null" json:"-"`
}
