package db

import (
	"fmt"
	"strings"
	config "taskel/config"
	"taskel/constants"
	model "taskel/models"

	"github.com/go-faker/faker/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func Connect() {
	password := fmt.Sprintf("password=%s", config.Config.DBPassword)
	if config.Config.DBPassword == "" {
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

func AutoMigrate() {
	DB.AutoMigrate(&model.Task{}, &model.User{}, &model.Comment{}, &model.Permission{}, &model.Role{})
}

func Reset() {
	DB.Exec("DELETE FROM permissions")
	DB.Exec("DELETE FROM roles")
	DB.Exec("DELETE FROM comments")
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
	taskSeed()
	roleAndPermissionSeed();
}

func taskSeed() {
	var randomUser model.User
	DB.First(&randomUser)
	for i := 0; i < 15; i++ {
		Key := fmt.Sprintf("TASK-%d", i)
		task := model.Task{Title: faker.Word(), Status: "todo", Key: Key}
		task.Comments = []*model.Comment{}

		for j := 0; j < 5; j++ {
			comment := model.Comment{
				AuthorID: randomUser.ID,
				Comment:  faker.Paragraph(),
			}
			task.Comments = append(task.Comments, &comment)
		}

		DB.Create(&task)
	}
}

func roleAndPermissionSeed() {
	permissios := []model.Permission{
		{Name: constants.RBAC_Task_Write},
	}

	roles := []model.Role{
		{Name: "admin", Permissions: permissios},
		{Name: "user"},
	}

	DB.Create(&roles)

	var users []model.User
	DB.Limit(2).Find(&users)

	users[0].RoleID = &roles[0].ID
	users[1].RoleID = &roles[1].ID
	DB.Save(&users)
}