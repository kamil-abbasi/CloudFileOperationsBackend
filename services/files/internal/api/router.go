package api

import (
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"

	v1 "github.com/kamil-abbasi/CloudFileOperationsBackend/internal/api/v1"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
)

func NewRouter(config *config.Config, cache *cache.Cache) *gin.Engine {
	r := gin.Default()

	r.GET("/healthcheck")

	v1.RegisterRoutes(r.Group("v1"), config, cache)

	return r
}
