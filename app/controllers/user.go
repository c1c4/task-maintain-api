package controllers

import (
	"api/app/models"
	"api/app/repositories"
	"api/app/utils/error_utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		errUnprocessibleEntity := error_utils.NewUnprocessibleEntityError("it's not possible to convert the JSON into an object")
		c.JSON(errUnprocessibleEntity.Status(), errUnprocessibleEntity)
		return
	}

	if err := user.Prepare(); err != nil {
		errPrepare := error_utils.NewBadRequestError(err.Error())
		c.JSON(errPrepare.Status(), errPrepare)
		return
	}

	dbUser, errCreateUser := repositories.UserRepo.Create(&user)

	if errCreateUser != nil {
		c.JSON(errCreateUser.Status(), errCreateUser)
		return
	}

	c.JSON(http.StatusCreated, dbUser)
}
