package types

type Usertypes struct {
	Id            int    `json:"id"`
	Username      string `json:"username"`
	Password_hash string `json:"password_hash"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Role          string `json:"role"`
	Is_verified   bool   `json:"is_verified"`
}

type AdminSuper struct {
	Adminid       uint   `json:"adminid" form:"adminid"`
	Username      string `json:"username" form:"username"`
	Password_hash string `json:"password_hash" form:"password_hash"`
	Role          string `json:"role"`
}
