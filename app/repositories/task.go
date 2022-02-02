package repositories

import (
	"api/app/database"
	"api/app/models"
	"api/app/utils/error_formats"
	"api/app/utils/error_utils"
	"errors"

	"gorm.io/gorm"
)

var TaskRepo TaskRepoInterface = &taskRepo{}

type TaskRepoInterface interface {
	Get(uint64) (*models.Task, error_utils.MessageErr)
	Create(*models.Task) (*models.Task, error_utils.MessageErr)
	Update(*models.Task) (*models.Task, error_utils.MessageErr)
	GetAll() []models.Task
	GetAllByUserID(userID uint64) []models.Task
	Delete(uint64) error_utils.MessageErr
	Init()
}

type taskRepo struct {
	db *gorm.DB
}

func (taskRepo *taskRepo) Init() {
	taskRepo.db = database.Database
}

func NewTaskRepository(db *gorm.DB) TaskRepoInterface {
	return &taskRepo{db: db}
}

func (taskRepo *taskRepo) Create(task *models.Task) (*models.Task, error_utils.MessageErr) {
	result := taskRepo.db.Create(&task)

	if result.Error != nil {
		return nil, error_formats.ParseError(result.Error)
	}

	return task, nil
}

func (taskRepo *taskRepo) Get(taskId uint64) (*models.Task, error_utils.MessageErr) {
	var task *models.Task = &models.Task{}
	result := taskRepo.db.First(&task, taskId)

	if result.Error != nil {
		return nil, error_formats.ParseError(result.Error)
	}

	return task, nil
}

func (taskRepo *taskRepo) Update(task *models.Task) (*models.Task, error_utils.MessageErr) {
	result := taskRepo.db.Model(&task).Updates(models.Task{Summary: task.Summary})

	if result.Error != nil {
		return nil, error_formats.ParseError(result.Error)
	}

	return task, nil
}

func (taskRepo *taskRepo) GetAll() []models.Task {
	var tasks []models.Task
	taskRepo.db.Find(&tasks)

	return tasks
}

func (taskRepo *taskRepo) GetAllByUserID(userID uint64) []models.Task {
	var tasks []models.Task
	taskRepo.db.Where(&models.Task{UserID: userID}).Find(&tasks)

	return tasks
}

func (taskRepo *taskRepo) Delete(taskId uint64) error_utils.MessageErr {
	result := taskRepo.db.Delete(&models.Task{}, taskId)

	if result.Error != nil || result.RowsAffected == 0 {
		if result.RowsAffected == 0 {
			return error_formats.ParseError(errors.New("record not found"))
		}

		return error_formats.ParseError(result.Error)
	}

	return nil
}
