package error_formats

import (
	"api/app/utils/error_utils"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func ParseError(err error) error_utils.MessageErr {
	sqlErr, ok := err.(*mysql.MySQLError)

	if !ok {
		if strings.Contains(err.Error(), "record not found") {
			return error_utils.NewNotFoundError("no record matching given the identification")
		}

		return error_utils.NewInternalServerError(fmt.Sprintf("error when trying to save: %s", err.Error()))
	}

	switch sqlErr.Number {
	case 1062:
		error_utils.NewInternalServerError("title already taken")
	}
	return error_utils.NewInternalServerError(fmt.Sprintf("error when processing request %s", err.Error()))
}
