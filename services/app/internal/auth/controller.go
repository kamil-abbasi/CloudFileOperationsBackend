package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/auth/dtos"
)

type AuthController struct {
	authService *AuthService
}

func NewController(authService *AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Register(ctx *gin.Context) {
	data, _ := ctx.Get("validatedBody")
	dto := data.(dtos.RegisterUserDto)

	response, err := c.authService.Register(dto)

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *AuthController) Login(ctx *gin.Context) {
	data, _ := ctx.Get("validatedBody")
	dto := data.(dtos.LoginUserDto)

	response, err := c.authService.Login(dto)

	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusCreated, response)
}
