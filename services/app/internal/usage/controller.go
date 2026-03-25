package usage

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/shared"
)

type UsageController struct {
	usageService *UsageService
}

func NewController(usageService *UsageService) UsageController {
	return UsageController{
		usageService: usageService,
	}
}

func (c *UsageController) CalculateForLocation(ctx *gin.Context) {
	location := "/" + ctx.Query("location")

	location = filepath.Clean(location)

	result, err := c.usageService.CalculateForLocation(location)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &shared.HttpError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})

		return
	}

	ctx.JSON(http.StatusOK, result)
}
