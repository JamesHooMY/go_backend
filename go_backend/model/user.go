package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"type:varchar(50);not null;unique" json:"email"`
	Mobile   string `gorm:"type:varchar(11)" json:"mobile"`
	Name     string `gorm:"type:varchar(20);not null" json:"name"`
	Age      int    `gorm:"type:int(3)" json:"age"`
	Password string `gorm:"type:varchar(20);not null" json:"-"`
}
