package models

type Credential struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type RegisterInfo struct {
	Username  string `form:"username" json:"username" binding:"required"`
	Email     string `form:"email" json:"email" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
	FirstName string `form:"firstname" json:"firstname" binding:"required"`
	LastName  string `form:"lastname" json:"lastname" binding:"required"`
	Pid string `form:"pid" json:"pid"`
}

type User struct {
	Username string `form:"username" json:"username" binding:"required"`
}
