package model

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // In practice, store hashed passwords
	Email    string `json:"email"`
}
