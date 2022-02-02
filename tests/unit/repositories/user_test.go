package repositories

import (
	"api/app/models"
	"api/app/repositories"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var tm = time.Now()

type userSuite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	userRepository repositories.UserRepoInterface
}

func (s *userSuite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
		DriverName:                "mysql",
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	require.NoError(s.T(), err)

	s.userRepository = repositories.NewUserRepository(s.DB)
}

func (s *userSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestUserInit(t *testing.T) {
	suite.Run(t, new(userSuite))
}

func (s *userSuite) TestCreateUser_Success() {
	user := models.User{
		Name:      "name",
		Email:     "email@email.com",
		Password:  "pass",
		Type:      "Manager",
		CreatedAt: tm,
		UpdatedAt: tm,
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec("INSERT INTO `users`").
		WithArgs(user.Name, user.Email, user.Password, user.Type, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	dbUser, err := s.userRepository.Create(&user)
	require.NoError(s.T(), err)
	require.NotEqual(s.T(), dbUser.ID, uint64(0))
}

func (s *userSuite) TestCreateUser_UniqueConstraint() {
	errorString := "Error 1062: Duplicate entry 'ceci-tech@gmail.com' for key 'users.email'"
	user := models.User{
		Name:      "name",
		Email:     "email@email.com",
		Password:  "pass",
		Type:      "Manager",
		CreatedAt: tm,
		UpdatedAt: tm,
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec("INSERT INTO `users`").
		WithArgs(user.Name, user.Email, user.Password, user.Type, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New(errorString))
	s.mock.ExpectRollback()

	_, err := s.userRepository.Create(&user)
	require.Error(s.T(), err, errorString)

}

func (s *userSuite) TestGetUser_Success() {
	user := models.User{
		Name:      "name",
		Email:     "email@email.com",
		Password:  "pass",
		Type:      "Manager",
		CreatedAt: tm,
		UpdatedAt: tm,
	}

	s.mock.ExpectQuery("SELECT(.*)").
		WithArgs(uint64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password", "type", "created_at", "updated_at"}).
			AddRow(uint64(1), user.Name, user.Password, user.Type, user.CreatedAt, user.UpdatedAt))

	res, err := s.userRepository.Get(uint64(1))
	require.NoError(s.T(), err)
	require.Equal(s.T(), res.ID, uint64(1))
}

func (s *userSuite) TestGetUser_NotFound() {
	errorString := "no record matching given the identification"

	//no record matching given the identification
	s.mock.ExpectQuery("SELECT(.*)").
		WithArgs(uint64(1)).
		WillReturnError(errors.New(errorString))

	_, err := s.userRepository.Get(uint64(1))
	require.Error(s.T(), err, errorString)
}

func (s *userSuite) TestGetUserByEmail_Success() {
	user := models.User{
		Name:      "name",
		Email:     "email@email.com",
		Password:  "pass",
		Type:      "Manager",
		CreatedAt: tm,
		UpdatedAt: tm,
	}

	s.mock.ExpectQuery("SELECT(.*)").
		WithArgs(user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password", "type", "created_at", "updated_at"}).
			AddRow(uint64(1), user.Name, user.Password, user.Type, user.CreatedAt, user.UpdatedAt))

	res, err := s.userRepository.GetByEmail(user.Email)
	require.NoError(s.T(), err)
	require.Equal(s.T(), res.ID, uint64(1))
}

func (s *userSuite) TestGetUserByEmail_NotFound() {
	var (
		errorString = "no record matching given the identification"
		email       = "email@email.com"
	)

	//no record matching given the identification
	s.mock.ExpectQuery("SELECT(.*)").
		WithArgs(email).
		WillReturnError(errors.New(errorString))

	_, err := s.userRepository.GetByEmail(email)
	require.Error(s.T(), err, errorString)
}
