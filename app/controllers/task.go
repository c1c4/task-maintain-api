package controllers

import (
	"api/app/authentication"
	"api/app/config"
	"api/app/message"
	"api/app/models"
	"api/app/repositories"
	"api/app/utils/error_utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {
	userID, err := authentication.ExtractUserId(c)
	if err != nil {
		errUnauthorized := error_utils.NewUnauthorizedError(err.Error())
		c.JSON(errUnauthorized.Status(), errUnauthorized)
		return
	}

	if canCreate, err := checkIsMethodAllowed("create", c); err != nil || !canCreate {
		if err != nil {
			c.JSON(err.Status(), err)
		} else {
			errForbidden := error_utils.NewForbiddenError("The user doesn't have the right permission to create a task")
			c.JSON(errForbidden.Status(), errForbidden)
		}
		return
	}

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		errUnprocessibleEntity := error_utils.NewUnprocessibleEntityError("it's not possible to convert the JSON into an object")
		c.JSON(errUnprocessibleEntity.Status(), errUnprocessibleEntity)
		return
	}

	task.UserID = userID

	if err := task.Prepare(); err != nil {
		errPrepare := error_utils.NewBadRequestError(err.Error())
		c.JSON(errPrepare.Status(), errPrepare)
		return
	}

	dbTask, errCreateTask := repositories.TaskRepo.Create(&task)

	if errCreateTask != nil {
		c.JSON(errCreateTask.Status(), errCreateTask)
		return
	}

	//publish the message here
	msg := fmt.Sprintf("The tech %d performed the task %d on date %d-%02d-%02d",
		userID,
		task.ID,
		task.CreatedAt.Year(),
		task.CreatedAt.Month(),
		task.CreatedAt.Day(),
	)
	message.Publish(c.Writer, config.GOOGLE_PROJECT_ID, config.GOOGLE_TOPIC_ID, msg)

	c.JSON(http.StatusCreated, dbTask)
}

func UpdateTask(c *gin.Context) {
	userID, err := authentication.ExtractUserId(c)
	if err != nil {
		errUnauthorized := error_utils.NewUnauthorizedError(err.Error())
		c.JSON(errUnauthorized.Status(), errUnauthorized)
		return
	}

	if canUpdate, err := checkIsMethodAllowed("update", c); err != nil || !canUpdate {
		if err != nil {
			c.JSON(err.Status(), err)
		} else {
			errUnauthorized := error_utils.NewForbiddenError("The user doesn't have the right permission to update a task")
			c.JSON(errUnauthorized.Status(), errUnauthorized)
		}
		return
	}

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errBadRequest := error_utils.NewBadRequestError(fmt.Sprintf("not possible to convert %s into a number", c.Param("id")))
		c.JSON(errBadRequest.Status(), errBadRequest)
		return
	}

	dbTask, errFindTask := repositories.TaskRepo.Get(taskID)

	if errFindTask != nil {
		c.JSON(errFindTask.Status(), errFindTask)
		return
	}

	fmt.Println(dbTask.UserID)
	fmt.Println(userID)

	if dbTask.UserID != userID {
		errForbidden := error_utils.NewForbiddenError("Not possible to update a task that does not belong to you")
		c.JSON(errForbidden.Status(), errForbidden)
		return
	}

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		errUnprocessibleEntity := error_utils.NewUnprocessibleEntityError("it's not possible to convert the JSON into an object")
		c.JSON(errUnprocessibleEntity.Status(), errUnprocessibleEntity)
		return
	}

	fmt.Println(task)

	if err := task.Prepare(); err != nil {
		errPrepare := error_utils.NewBadRequestError(err.Error())
		c.JSON(errPrepare.Status(), errPrepare)
		return
	}

	_, errUpdateTask := repositories.TaskRepo.Update(&task)
	if errUpdateTask != nil {
		c.JSON(errUpdateTask.Status(), errUpdateTask)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func GetTask(c *gin.Context) {
	userID, err := authentication.ExtractUserId(c)
	if err != nil {
		errUnauthorized := error_utils.NewUnauthorizedError(err.Error())
		c.JSON(errUnauthorized.Status(), errUnauthorized)
		return
	}

	if canGetOne, err := checkIsMethodAllowed("get_one", c); err != nil || !canGetOne {
		if err != nil {
			c.JSON(err.Status(), err)
		} else {
			errUnauthorized := error_utils.NewForbiddenError("The user doesn't have the right permission to get a specific task")
			c.JSON(errUnauthorized.Status(), errUnauthorized)
		}
		return
	}

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errBadRequest := error_utils.NewBadRequestError(fmt.Sprintf("not possible to convert %s into a number", c.Param("id")))
		c.JSON(errBadRequest.Status(), errBadRequest)
		return
	}

	dbTask, errFindTask := repositories.TaskRepo.Get(taskID)

	if errFindTask != nil {
		c.JSON(errFindTask.Status(), errFindTask)
		return
	}

	if dbTask.UserID != userID {
		errForbidden := error_utils.NewForbiddenError("Not possible to see a task that does not belong to you")
		c.JSON(errForbidden.Status(), errForbidden)
		return
	}

	c.JSON(http.StatusOK, dbTask)
}

func GetTasksByUser(c *gin.Context) {
	userID, err := authentication.ExtractUserId(c)
	if err != nil {
		errUnauthorized := error_utils.NewUnauthorizedError(err.Error())
		c.JSON(errUnauthorized.Status(), errUnauthorized)
		return
	}

	if canGetYourOwnTasks, err := checkIsMethodAllowed("list_own_tasks", c); err != nil || !canGetYourOwnTasks {
		if err != nil {
			c.JSON(err.Status(), err)
		} else {
			errUnauthorized := error_utils.NewForbiddenError("The user doesn't have the right permission to list his tasks")
			c.JSON(errUnauthorized.Status(), errUnauthorized)
		}
		return
	}

	dbTasks := repositories.TaskRepo.GetAllByUserID(userID)

	c.JSON(http.StatusOK, dbTasks)
}

func GetAllTasks(c *gin.Context) {
	if canGetList, err := checkIsMethodAllowed("list", c); err != nil || !canGetList {
		if err != nil {
			c.JSON(err.Status(), err)
		} else {
			errUnauthorized := error_utils.NewForbiddenError("The user doesn't have the right permission to get all tasks")
			c.JSON(errUnauthorized.Status(), errUnauthorized)
		}
		return
	}

	dbTasks := repositories.TaskRepo.GetAll()

	c.JSON(http.StatusOK, dbTasks)
}

func DeleteTasks(c *gin.Context) {
	if canDelete, err := checkIsMethodAllowed("delete", c); err != nil || !canDelete {
		if err != nil {
			c.JSON(err.Status(), err)
		} else {
			errUnauthorized := error_utils.NewForbiddenError("The user doesn't have the right permission to delete a specific task")
			c.JSON(errUnauthorized.Status(), errUnauthorized)
		}
		return
	}

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		errBadRequest := error_utils.NewBadRequestError(fmt.Sprintf("not possible to convert %s into a number", c.Param("id")))
		c.JSON(errBadRequest.Status(), errBadRequest)
		return
	}

	errDeleteTask := repositories.TaskRepo.Delete(taskID)

	if errDeleteTask != nil {
		c.JSON(errDeleteTask.Status(), errDeleteTask)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func checkIsMethodAllowed(permission string, ctx *gin.Context) (bool, error_utils.MessageErr) {
	isAllowed := false

	userPermissions, err := authentication.ExtractPermissions(ctx)

	if err != nil {
		return isAllowed, error_utils.NewInternalServerError(err.Error())
	}

	for _, userPermission := range userPermissions {
		if permission == userPermission {
			isAllowed = true
			break
		}
	}

	return isAllowed, nil
}
