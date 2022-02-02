package controllers

import (
	"api/app/authentication"
	"api/app/models"
	"api/app/repositories"
	"api/app/security"
	"api/app/utils/error_utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		errUnprocessibleEntity := error_utils.NewUnprocessibleEntityError("it's not possible to convert the JSON into an object")
		c.JSON(errUnprocessibleEntity.Status(), errUnprocessibleEntity)
		return
	}

	dbUser, errGetByEmail := repositories.UserRepo.GetByEmail(user.Email)

	if errGetByEmail != nil {
		c.JSON(errGetByEmail.Status(), errGetByEmail)
		return
	}

	if err := security.VerifyPassword(dbUser.Password, user.Password); err != nil {
		errUnauthorized := error_utils.NewUnauthorizedError(err.Error())
		c.JSON(errUnauthorized.Status(), errUnauthorized)
		return
	}

	token, err := authentication.CreateToken(dbUser.ID, dbUser.Type)
	if err != nil {
		errInternalServer := error_utils.NewInternalServerError(err.Error())
		c.JSON(errInternalServer.Status(), errInternalServer)
		return
	}

	userID := strconv.FormatUint(dbUser.ID, 10)

	c.JSON(http.StatusOK, models.AuthenticationData{ID: userID, Token: token})
}
