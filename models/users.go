package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Users struct {
	ID          uuid.UUID `db:"id"`
	FirstName   string    `db:"first_name"`
	LastName    string    `db:"last_name"`
	DNI         string    `db:"dni"`
	PhoneNumber int       `db:"phone_number"`
	Email       string    `db:"email"`
	JobTitle    string    `db:"job_title"`
	Age         string    `db:"age"`
	Avg         string    `db:"avg"`
	Birthdate   time.Time `db:"birthdate"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type SliceUsers []Users

func (t SliceUsers) AvgAge(avg SliceUsers) Users {
	for _, v := range avg {
		return v
	}
	return Users{}
}
