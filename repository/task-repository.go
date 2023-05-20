package repository

import (
	"taskel/db"
	model "taskel/models"
)

func GetTaskByIdOrKey(idOrKey string) (*model.Task, error) {
	var task model.Task
	result := db.DB.Where("key = ?", idOrKey).First(&task)
	if result.Error != nil {
		return nil, result.Error
	}
	return &task, nil
}

func TaskWatch(taskKey string, userID uint) error {
	var task model.Task
	db.DB.Where("key = ?", taskKey).First(&task)
	err := db.DB.Model(&task).Where("key = ?", taskKey).Association("Watchers").Append(
		&model.User{
			ID: userID,
		},
	)
	db.DB.Save(&task)
	return err
}
