package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"newanysock/internal/entity"
	"newanysock/internal/entity/user"
)

var DB *gorm.DB

func Connect() {
	connection, err := gorm.Open(mysql.Open("root:@/anysock?parseTime=true"), &gorm.Config{})

	if err != nil {
		panic("could not connect to the database ")
	}
	DB = connection
	connection.AutoMigrate(&user.Account{})
	connection.AutoMigrate(&entity.License{})
}
