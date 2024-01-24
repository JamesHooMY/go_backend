package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(20);not null"`
	Age      int    `gorm:"type:int(3);not null"`
	Mobile   string `gorm:"type:varchar(11);not null;unique"`
	Email    string `gorm:"type:varchar(50);not null;unique"`
	Password string `gorm:"type:varchar(20);not null"`
}
