package db

import (
	"fmt"
	"strings"
	config "taskel/config"
	model "taskel/models"

	"github.com/go-faker/faker/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func Connect() {
	password := fmt.Sprintf("password=%s", config.Config.DBPassword)
	if (config.Config.DBPassword == "") {
		password = ""
	}

	var dsn string = fmt.Sprintf("host=%s port=%s user=%s %s dbname=%s client_encoding=UTF8 sslmode=disable",
		config.Config.DBHost,
		config.Config.DBPort,
		config.Config.DBUser,
		password,
		config.Config.DBName,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Failed connecting to database %s\n", dsn))
	}
}

func Reset() {
	DB.Exec("DELETE FROM task_users")
	DB.Exec("DELETE FROM tasks")
	DB.Exec("DELETE FROM users")
}

func Seed() {
	// generate 10 fake Users
	for i := 0; i < 10; i++ {
		email := fmt.Sprintf("hanggi_anggono+%d@yahoo.com", i)
		firstName := faker.FirstName()
		password, _ := model.UserHashPassword("123456")
		user := model.User{Name: firstName, Username: strings.ToLower(firstName), Email: &email, Password: password}
		DB.Create(&user)
	}

	// generate 15 fake tasks
	for i := 0; i < 15; i++ {
		task := model.Task{Title: faker.Sentence(), Status: "todo"}
		DB.Create(&task)
	}
}
