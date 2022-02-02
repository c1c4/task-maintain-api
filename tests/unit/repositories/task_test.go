package repositories

import (
	"api/app/models"
	"api/app/repositories"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type taskSuite struct {
	suite.Suite
	DB             *gorm.DB
	mock           sqlmock.Sqlmock
	taskRepository repositories.TaskRepoInterface
}

func (s *taskSuite) SetupSuite() {
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

	s.taskRepository = repositories.NewTaskRepository(s.DB)
}

func (s *taskSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestTaskInit(t *testing.T) {
	suite.Run(t, new(taskSuite))
}

func (s *taskSuite) TestCreateTask_Success() {
	task := models.Task{
		Summary:   "Creating a summary",
		UserID:    1,
		CreatedAt: tm,
		UpdatedAt: tm,
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec("INSERT INTO `tasks`").
		WithArgs(task.Summary, task.UserID, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	dbTask, err := s.taskRepository.Create(&task)
	require.NoError(s.T(), err)
	require.NotEqual(s.T(), dbTask.ID, uint64(0))
}

func (s *taskSuite) TestGetTask_Success() {
	task := models.Task{
		ID:        1,
		Summary:   "Recovering a summary",
		UserID:    1,
		CreatedAt: tm,
		UpdatedAt: tm,
	}

	s.mock.ExpectQuery("SELECT(.*)").
		WithArgs(task.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "summary", "user_id", "created_at", "updated_at"}).
			AddRow(task.ID, task.Summary, task.UserID, task.CreatedAt, task.UpdatedAt))

	dbTask, err := s.taskRepository.Get(uint64(1))
	require.NoError(s.T(), err)
	require.Equal(s.T(), dbTask.ID, uint64(1))
}

func (s *userSuite) TestGetTask_NotFound() {
	errorString := "no record matching given the identification"

	//no record matching given the identification
	s.mock.ExpectQuery("SELECT(.*)").
		WithArgs(uint64(1)).
		WillReturnError(errors.New(errorString))

	_, err := s.userRepository.Get(uint64(1))
	require.Error(s.T(), err, errorString)
}

func (s *taskSuite) TestUpdateTask_Success() {
	task := models.Task{
		ID:        1,
		Summary:   "Updating a summary",
		UserID:    1,
		CreatedAt: tm,
		UpdatedAt: tm,
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE `tasks`").
		WithArgs(task.Summary, sqlmock.AnyArg(), task.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	dbTask, err := s.taskRepository.Update(&task)
	require.NoError(s.T(), err)
	require.NotEqual(s.T(), dbTask.ID, uint64(0))
	require.Equal(s.T(), dbTask.ID, task.ID)
}

func (s *taskSuite) TestUpdateTask_NotFound() {
	errorString := "no record matching given the identification"
	task := models.Task{
		ID:        1,
		Summary:   "Updating a summary",
		UserID:    1,
		CreatedAt: tm,
		UpdatedAt: tm,
	}

	//no record matching given the identification
	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE `tasks`").
		WithArgs(task.Summary, sqlmock.AnyArg(), task.ID).
		WillReturnError(errors.New(errorString))
	s.mock.ExpectRollback()

	_, err := s.taskRepository.Update(&task)
	require.Error(s.T(), err, errorString)
}

func (s *taskSuite) TestGetAllTasks_Success() {
	tasks := []models.Task{
		{
			ID:        1,
			Summary:   "Recovering a summary 1",
			UserID:    1,
			CreatedAt: tm,
			UpdatedAt: tm,
		},
		{
			ID:        2,
			Summary:   "Recovering a summary 2",
			UserID:    1,
			CreatedAt: tm,
			UpdatedAt: tm,
		},
	}

	s.mock.ExpectQuery("SELECT(.*)").
		WillReturnRows(sqlmock.NewRows([]string{"id", "summary", "user_id", "created_at", "updated_at"}).
			AddRow(tasks[0].ID, tasks[0].Summary, tasks[0].UserID, tasks[0].CreatedAt, tasks[0].UpdatedAt).
			AddRow(tasks[1].ID, tasks[1].Summary, tasks[1].UserID, tasks[1].CreatedAt, tasks[1].UpdatedAt),
		)

	dbTask := s.taskRepository.GetAll()
	require.Equal(s.T(), len(dbTask), 2)
}

func (s *taskSuite) TestGetAllTasks_Empty() {
	s.mock.ExpectQuery("SELECT(.*)").
		WillReturnRows(sqlmock.NewRows([]string{"id", "summary", "user_id", "created_at", "updated_at"}))

	dbTask := s.taskRepository.GetAll()
	require.Equal(s.T(), len(dbTask), 0)
}

func (s *taskSuite) TestGetAllTasksByUser_Success() {
	tasks := []models.Task{
		{
			ID:        1,
			Summary:   "Recovering a summary 1",
			UserID:    1,
			CreatedAt: tm,
			UpdatedAt: tm,
		},
		{
			ID:        2,
			Summary:   "Recovering a summary 2",
			UserID:    1,
			CreatedAt: tm,
			UpdatedAt: tm,
		},
	}

	s.mock.ExpectQuery("SELECT(.*)").
		WithArgs(tasks[0].UserID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "summary", "user_id", "created_at", "updated_at"}).
			AddRow(tasks[0].ID, tasks[0].Summary, tasks[0].UserID, tasks[0].CreatedAt, tasks[0].UpdatedAt).
			AddRow(tasks[1].ID, tasks[1].Summary, tasks[1].UserID, tasks[1].CreatedAt, tasks[1].UpdatedAt),
		)

	dbTask := s.taskRepository.GetAllByUserID(tasks[0].UserID)
	require.Equal(s.T(), len(dbTask), 2)
}

func (s *taskSuite) TestGetAllTasksByUser_Empty() {
	s.mock.ExpectQuery("SELECT(.*)").
		WithArgs(uint64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "summary", "user_id", "created_at", "updated_at"}))

	dbTask := s.taskRepository.GetAllByUserID(uint64(1))
	require.Equal(s.T(), len(dbTask), 0)
}

func (s *taskSuite) TestDeleteTask_Success() {
	id := uint64(1)

	s.mock.ExpectBegin()
	s.mock.ExpectExec("DELETE FROM `tasks`").WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	err := s.taskRepository.Delete(id)
	require.NoError(s.T(), err)
}
