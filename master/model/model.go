package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name       string
	GroupNames string
}

func InitSchema(orm *gorm.DB) {
	models := []interface{}{
		&User{},
	}

	for _, model := range models {
		orm.AutoMigrate(model) // create table and migrate
	}
}
