package models

import (
	"time"


	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

type Task struct {
	ID          uuid.UUID  `form:"id" db:"id"`
	Name        string     `form:"name" db:"name"`
	Description string     `form:"description" db:"description"`
	FinishedAt  nulls.Time `form:"finish_at" db:"finished_at"`
	CreatedAt   time.Time  `form:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `form:"uptdated_at" db:"updated_at"`
}

//Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Task) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(&validators.StringIsPresent{
		Field:   t.Name,
		Name:    "name",
		Message: "Name can not be empty",
	}), nil

}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Task) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {

	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Task) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

type Tasks []Task

func (t Tasks) InfoTask(id uuid.UUID) Task{

	for _, v := range t {
		if v.ID == id {
			
			return v
		}
		
	}

		
 return Task{}
}
