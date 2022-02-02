package repositories

import (
	"api/app/database"
	"api/app/models"
	"api/app/utils/error_formats"
	"api/app/utils/error_utils"

	"gorm.io/gorm"
)

var UserRepo UserRepoInterface = &userRepo{}

type UserRepoInterface interface {
	Get(uint64) (*models.User, error_utils.MessageErr)
	GetByEmail(email string) (*models.User, error_utils.MessageErr)
	Create(*models.User) (*models.User, error_utils.MessageErr)
	Init()
}

type userRepo struct {
	db *gorm.DB
}

func (userRepo *userRepo) Init() {
	userRepo.db = database.Database
}

func NewUserRepository(db *gorm.DB) UserRepoInterface {
	return &userRepo{db: db}
}

func (userRepo *userRepo) Create(user *models.User) (*models.User, error_utils.MessageErr) {
	result := userRepo.db.Debug().Create(&user)

	if result.Error != nil {
		return nil, error_formats.ParseError(result.Error)
	}

	return user, nil
}

func (userRepo *userRepo) Get(userId uint64) (*models.User, error_utils.MessageErr) {
	var user *models.User = &models.User{}
	result := userRepo.db.First(&user, userId)

	if result.Error != nil {
		return nil, error_formats.ParseError(result.Error)
	}

	return user, nil
}

func (userRepo *userRepo) GetByEmail(email string) (*models.User, error_utils.MessageErr) {
	var user *models.User = &models.User{}
	result := userRepo.db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return nil, error_formats.ParseError(result.Error)
	}

	return user, nil
}
