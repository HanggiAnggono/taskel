package db

import (
	"fmt"
	"strings"
	model "taskel/models"

	"github.com/go-faker/faker/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func Connect() {
	var dsn string = "host=localhost user=postgres password=postgres dbname=taskel_dev port=5432 client_encoding=utf8"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Failed connecting to database %s\n", dsn))
	}
}

func Reset() {
	DB.Exec("DELETE FROM users")
	DB.Exec("DELETE FROM tasks")
}

func Seed() {
	// generate 10 fake Users
	for i := 0; i < 10; i++ {
		email := faker.Email()
		firstName := faker.FirstName()
		password, _ := model.UserHashPassword("123456")
		user := model.User{Name: firstName, Username: strings.ToLower(firstName), Email: &email, Password: password}
		DB.Create(&user)
	}

	// generate 15 fake tasks
	for i := 0; i < 15; i++ {
		task := model.Task{Name: faker.Sentence(), Status: "todo"}
		DB.Create(&task)
	}
}
