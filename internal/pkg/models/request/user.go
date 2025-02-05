package request

type User struct {
	Id            int    `json:"id"`
	Username      string `json:"username"`
	Password_hash string `json:"password_hash"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Role          string `json:"role"`
	Is_verified   bool   `json:"is_verified"`
}
