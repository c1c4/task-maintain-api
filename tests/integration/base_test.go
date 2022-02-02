package integration

import (
	"api/app"
	"api/app/database"
	"api/app/database/migration"
	"api/app/models"
	"api/app/repositories"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SuiteTest struct {
	suite.Suite
}

var (
	technician_token = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjk5OTk5NjU1ODI4LCJwZXJtaXNzaW9ucyI6WyJjcmVhdGUiLCJ1cGRhdGUiLCJnZXRfb25lIiwibGlzdF9vd25fdGFza3MiXSwidXNlcl9pZCI6MX0.aFj60IBurKQoSPQ25jHLu2yzBrsRLQwtlvncgL2_7IY"
	manager_token    = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjk5OTQzNjU1NzM1LCJwZXJtaXNzaW9ucyI6WyJsaXN0IiwiZGVsZXRlIiwibm90aWZpZWQiXSwidXNlcl9pZCI6Mn0.x1WUovq-J4R0RhaSoXbQTvJ9pkZvO-q_oJT9HY-Pkk8"
	baseURL          = "http://localhost:8080/v1"
)

func TestSuite(t *testing.T) {
	suite.Run(t, &SuiteTest{})
}

func (s *SuiteTest) SetupSuite() {
	serverReady := make(chan bool)

	app := app.App{
		ServerReady: serverReady,
	}

	go app.StartApp()
	<-serverReady
}

func (s *SuiteTest) TearDownSuite() {
	p, _ := os.FindProcess(syscall.Getpid())
	p.Signal(syscall.SIGINT)
	database.Database.Migrator().DropTable(&models.User{}, &models.Task{})
}

func (s *SuiteTest) SetupTest() {
	migration.AutoMigration()
}

func (s *SuiteTest) TearDownTest() {
	s.NoError(database.Database.Migrator().DropTable(&models.User{}, &models.Task{}))
}

func (s *SuiteTest) seedOneUserTech() {
	user := models.User{
		Name:     "Tech user",
		Email:    "tech@tech.com",
		Password: "$2a$10$rcF.fKve09NwH29MHP5V4ulfW5/6AAE9QbZTEA6aGEuMsw3ADydPq",
		Type:     "Technician",
	}
	repository := repositories.NewUserRepository(database.Database)
	repository.Create(&user)
}

func (s *SuiteTest) seedOneUserManager() {
	user := models.User{
		Name:     "Manager user",
		Email:    "manager@manager.com",
		Password: "$2a$10$rcF.fKve09NwH29MHP5V4ulfW5/6AAE9QbZTEA6aGEuMsw3ADydPq",
		Type:     "Manager",
	}
	repository := repositories.NewUserRepository(database.Database)
	repository.Create(&user)
}

func (s *SuiteTest) seedOneTask() {
	task := models.Task{
		ID:      1,
		Summary: "This is a summary test",
		UserID:  1,
	}
	repository := repositories.NewTaskRepository(database.Database)
	repository.Create(&task)
}

func (s *SuiteTest) seedMultipleTasks() {
	tasks := []models.Task{
		{
			ID:      1,
			Summary: "This is a summary test",
			UserID:  1,
		},
		{
			Summary: "This is a summary test 2",
			UserID:  1,
		},
		{
			Summary: "This is a summary test 3",
			UserID:  2,
		},
	}
	repository := repositories.NewTaskRepository(database.Database)

	for _, task := range tasks {
		repository.Create(&task)
	}
}
