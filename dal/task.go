package dal

import (
	"errors"

	"gorm.io/gorm"
)

func GetAllTask(projectId int64) ([]Task, error) {
	var tasks []Task
	result := DB.Model(&Task{}).Where("project_id = ?", projectId).Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}

func GetTask(taskId int64) (Task, error) {
	var task Task
	result := DB.Model(&Task{}).Where("id = ?", taskId).First(&task)
	if result.Error != nil {
		return Task{}, result.Error
	}
	return task, nil
}

func AddTask(task Task) (int64, error) {
	result := DB.Model(&Task{}).Create(&task)
	if result.Error != nil {
		return 0, result.Error
	}
	return task.Id, nil
}

func CheckIfProjectHasRunningTask(projectId int64) (bool, error) {
	result := DB.Model(&Task{}).Where("project_id = ? && status = ?", projectId, TaskStatusRunning).First(&Task{})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}

func SetProjectRunningTaskStatus(projectId int64, status string) error {
	result := DB.Model(&Task{}).
		Where("project_id = ? && status = ?", projectId, TaskStatusRunning).
		Update("status", status).Limit(1)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
