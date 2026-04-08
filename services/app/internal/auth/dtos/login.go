package dtos

type LoginUserDto struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserResponseDto struct {
	AccessToken string `json:"accessToken"`
}
