package repository

import (
	"taskel/db"
	model "taskel/models"
)

func TaskWatch(taskID uint, userID uint) error {
	var task model.Task
	db.DB.First(&task, taskID)
	err := db.DB.Model(&task).Where("id = ?", taskID).Association("Watchers").Append(
		&model.User{
			ID: userID,
		},
	)
	db.DB.Save(&task)
	return err
}
