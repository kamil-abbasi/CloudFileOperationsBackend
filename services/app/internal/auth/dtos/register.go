package dtos

type RegisterUserDto struct {
	Name              string `json:"name" binding:"required"`
	Password          string `json:"password" binding:"required"`
	ConfirmedPassword string `json:"confirmedPassword" binding:"required,eqfield=Password"`
}

type RegisterUserResponseDto struct {
	Id         string `json:"id"`
	MaxStorage uint64 `json:"maxStorage"`
}
