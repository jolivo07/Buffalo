package models

import (
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"
)

type Tasks struct {
	ID          uuid.UUID  `form:"id" db:"id"`
	Name_task   string     `form:"name_task" db:"name_task"`
	Description string     `form:"description" db:"description"`
	Finish_at   nulls.Time `form:"finish_at" db:"finish_at"`
	CreatedAt   time.Time  `form:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `form:"uptdated_at" db:"updated_at"`
}
