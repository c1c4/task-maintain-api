package authentication

import (
	"api/app/config"
	"api/app/constants"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreateToken(userID uint64, userType string) (string, error) {
	permissionsByType := constants.GetSupportedPermissionsByUserType()
	permissions := jwt.MapClaims{}
	permissions["authorized"] = true
	permissions["exp"] = time.Now().Add(time.Hour * 6).Unix()
	permissions["user_id"] = userID
	permissions["permissions"] = permissionsByType[userType]
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)
	return token.SignedString([]byte(config.SECRETKEY))
}

func ValidateToken(c *gin.Context) error {
	tokenString := extractToken(c)
	token, err := jwt.Parse(tokenString, returnVerificationKey)

	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("invalid token")
}

func ExtractUserId(c *gin.Context) (uint64, error) {
	tokenString := extractToken(c)
	token, err := jwt.Parse(tokenString, returnVerificationKey)

	if err != nil {
		return 0, err
	}

	if permissions, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		usuarioID, erro := strconv.ParseUint(fmt.Sprintf("%.0f", permissions["user_id"]), 10, 64)
		if erro != nil {
			return 0, erro
		}

		return usuarioID, nil
	}

	return 0, errors.New("invalid token")
}

func ExtractPermissions(c *gin.Context) ([]string, error) {
	tokenString := extractToken(c)
	token, err := jwt.Parse(tokenString, returnVerificationKey)

	if err != nil {
		return []string{}, nil
	}

	if permissions, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userPermissionsInterface := permissions["permissions"]
		userPermissions := []string{}

		// cast the userPermissionsInterface into a slice of interfaces for each value cast to string and store it in a new slice of strings
		for _, value := range userPermissionsInterface.([]interface{}) {
			userPermissions = append(userPermissions, value.(string))
		}

		return userPermissions, nil
	}

	return []string{}, errors.New("invalid token")
}

func extractToken(c *gin.Context) string {
	token := c.Request.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}

	return ""
}

func returnVerificationKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signature method! %v", token.Header["alg"])
	}

	return []byte(config.SECRETKEY), nil
}
