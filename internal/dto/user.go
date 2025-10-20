package dto

type UserUpdate struct {
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}

type UserPatch struct {
	Email     *string    `json:"email"`
	Password  *string    `json:"password"`
}
