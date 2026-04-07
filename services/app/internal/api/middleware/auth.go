package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/auth"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
)

func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.Error(auth.ErrMissingAuthHeader)
			ctx.Abort()
			return
		}

		arr := strings.Split(authHeader, " ")

		if len(arr) < 2 {
			ctx.Error(auth.ErrMissingToken)
			ctx.Abort()
			return
		}

		tokenString := arr[1]

		bytes, err := os.ReadFile(cfg.JwtPublicKeyPath)

		if err != nil {
			ctx.Error(fmt.Errorf("failed to read jwt private key"))
			ctx.Abort()
			return
		}

		key, err := jwt.ParseRSAPublicKeyFromPEM(bytes)

		if err != nil {
			ctx.Error(fmt.Errorf("failed to parse jwt public key"))
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			return key, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}))

		if err != nil {
			ctx.Error(fmt.Errorf("failed to parse jwt token"))
			ctx.Abort()
			return
		}

		if !token.Valid {
			ctx.Error(auth.ErrInvalidToken)
			ctx.Abort()
			return
		}

		sub, err := token.Claims.GetSubject()

		if err != nil {
			ctx.Error(fmt.Errorf("Missing sub field in token"))
			ctx.Abort()
			return
		}

		ctx.Set("userId", sub)
		ctx.Next()
	}
}
