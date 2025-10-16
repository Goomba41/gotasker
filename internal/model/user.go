package model

import (
	"time"
)

type User struct {
	ID        uint  
	Email     string
	Password  string
	CreatedAt time.Time
}
