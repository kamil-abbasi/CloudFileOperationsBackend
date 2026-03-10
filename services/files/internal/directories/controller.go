package directories

import (
	"github.com/gin-gonic/gin"
	"github.com/kamil-abbasi/CloudFileOperationsBackend/internal/config"
	"github.com/patrickmn/go-cache"
)

type DirectoriesController struct {
	config             *config.Config
	cache              *cache.Cache
	directoriesService *DirectoriesService
}

func NewController(config *config.Config, cache *cache.Cache, service *DirectoriesService) DirectoriesController {
	return DirectoriesController{
		cache:              cache,
		config:             config,
		directoriesService: service,
	}
}

func (c *DirectoriesController) Download(ctx *gin.Context) {}

func (c *DirectoriesController) Create(ctx *gin.Context) {}

func (c *DirectoriesController) Remove(ctx *gin.Context) {}
