package entity

type User struct {
	Username          string `json:"username"`
	Password          string `json:"password"`
	PasswordConfirmed string `json:"password_confirmed"`
}

type UserResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Error    string `json:"error"`
}
