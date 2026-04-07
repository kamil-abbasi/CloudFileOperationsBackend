package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/auth"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/directories"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/files"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/shared"
)

func Error() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last()

			if err == nil || err.Err == nil {
				ctx.JSON(http.StatusInternalServerError, &shared.HttpError{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error",
				})

				return
			}

			actualErr := err.Err

			log.Printf("%+v\n", actualErr)

			if errors.Is(actualErr, files.ErrAlreadyExists) {
				ctx.JSON(http.StatusConflict, &shared.HttpError{
					Code:    http.StatusConflict,
					Message: "File already exists",
				})
			} else if errors.Is(actualErr, files.ErrNotFound) {
				ctx.JSON(http.StatusNotFound, &shared.HttpError{
					Code:    http.StatusNotFound,
					Message: "File not found",
				})
			} else if errors.Is(actualErr, directories.ErrAlreadyExists) {
				ctx.JSON(http.StatusConflict, &shared.HttpError{
					Code:    http.StatusConflict,
					Message: "Directory already exists",
				})
			} else if errors.Is(actualErr, directories.ErrNotFound) {
				ctx.JSON(http.StatusNotFound, &shared.HttpError{
					Code:    http.StatusNotFound,
					Message: "Directory not found",
				})
			} else if errors.Is(actualErr, auth.ErrUserNotFound) {
				ctx.JSON(http.StatusNotFound, &shared.HttpError{
					Code:    http.StatusNotFound,
					Message: "User not found",
				})
			} else if errors.Is(actualErr, auth.ErrUserAlreadyExists) {
				ctx.JSON(http.StatusConflict, &shared.HttpError{
					Code:    http.StatusConflict,
					Message: "User already exists",
				})
			} else if errors.Is(actualErr, auth.ErrPasswordsDoNotMatch) {
				ctx.JSON(http.StatusBadRequest, &shared.HttpError{
					Code:    http.StatusBadRequest,
					Message: "Passwords do not match",
				})
			} else if errors.Is(actualErr, auth.ErrInvalidCredentials) {
				ctx.JSON(http.StatusUnauthorized, &shared.HttpError{
					Code:    http.StatusUnauthorized,
					Message: "Invalid credentials",
				})
			} else if errors.Is(actualErr, auth.ErrMissingAuthHeader) {
				ctx.JSON(http.StatusUnauthorized, &shared.HttpError{
					Code:    http.StatusUnauthorized,
					Message: "Missing auth header",
				})
			} else if errors.Is(actualErr, auth.ErrMissingToken) {
				ctx.JSON(http.StatusUnauthorized, &shared.HttpError{
					Code:    http.StatusUnauthorized,
					Message: "Missing token",
				})
			} else if errors.Is(actualErr, auth.ErrInvalidToken) {
				ctx.JSON(http.StatusUnauthorized, &shared.HttpError{
					Code:    http.StatusUnauthorized,
					Message: "Invalid token",
				})
			} else {
				ctx.JSON(http.StatusInternalServerError, &shared.HttpError{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error",
				})
			}
		}
	}
}
