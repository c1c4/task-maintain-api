package models

import (
	"errors"
	"time"
)

type Task struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Summary   string    `gorm:"size:2500;not null" json:"summary,omitempty"`
	UserID    uint64    `json:"userId,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"modifiedAt,omitempty"`
}

func (task *Task) Prepare() error {

	if err := task.validate(); err != nil {
		return err
	}

	return nil
}

func (task *Task) validate() error {
	if len(task.Summary) == 0 {
		return errors.New("the field summary is required can't be empty")
	} else if len(task.Summary) > 2500 {
		return errors.New("the summary is too long need to be less or equal to 2500 characters")
	}

	return nil
}
