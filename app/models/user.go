package models

import (
	"api/app/security"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

type User struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Name      string    `gorm:"not null" json:"name,omitempty"`
	Email     string    `gorm:"not null; unique" json:"email"`
	Password  string    `gorm:"not null" json:"password"`
	Type      string    `gorm:"not null" json:"type,omitempty"`
	Tasks     []Task    `json:"tasks,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

func (user *User) Prepare() error {

	if err := user.validate(); err != nil {
		return err
	}

	if err := user.format(); err != nil {
		return err
	}

	return nil
}

func (user *User) validate() error {
	errs := ""
	if len(user.Name) == 0 {
		errs = "the field name is required can't be empty"
	}

	if len(user.Email) == 0 {
		if len(errs) > 0 {
			errs = fmt.Sprintf("%s\n", errs)
		}

		errs = fmt.Sprintf("%sthe field email is required can't be empty", errs)
	} else if err := checkmail.ValidateFormat(user.Email); err != nil {
		if len(errs) > 0 {
			errs = fmt.Sprintf("%s\n", errs)
		}

		errs = fmt.Sprintf("%sthe email is invalid", errs)
	}

	if len(user.Password) == 0 {
		if len(errs) > 0 {
			errs = fmt.Sprintf("%s\n", errs)
		}

		errs = fmt.Sprintf("%sthe field password is required can't be empty", errs)
	}

	if len(user.Type) == 0 {
		if len(errs) > 0 {
			errs = fmt.Sprintf("%s\n", errs)
		}

		errs = fmt.Sprintf("%sthe field type is required can't be empty", errs)
	} else if user.Type != "Technician" && user.Type != "Manager" {
		if len(errs) > 0 {
			errs = fmt.Sprintf("%s\n", errs)
		}
		errs = fmt.Sprintf("%sthe field type should be Technician or Manager", errs)
	}

	if len(errs) > 0 {
		return errors.New(errs)
	}

	return nil
}

func (user *User) format() error {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	hashedPassword, err := security.Hash(user.Password)

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return nil
}
