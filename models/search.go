package models

import "time"

type Search struct {
	Name string    `form:"name"`
	Date time.Time `form:"date"`
}
