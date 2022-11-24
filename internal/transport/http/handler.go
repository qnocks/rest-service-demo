package http

import (
	"github.com/gin-gonic/gin"
	"team-task/internal/config"
	"team-task/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Init(cfg *config.Config) (*gin.Engine, *gin.Engine) {
	setRouter := h.setupSetAPI(cfg)
	getRouter := h.setupGetAPI()

	h.initBackup([]*gin.Engine{setRouter, getRouter})

	return setRouter, getRouter
}

func (h *Handler) setupSetAPI(cfg *config.Config) *gin.Engine {
	setRouter := gin.New()
	setRouter.Use(gin.Recovery(), gin.Logger())

	authorized := setRouter.Group("/", gin.BasicAuth(gin.Accounts{cfg.Auth.Username: cfg.Auth.Password}))
	authorized.POST("/", h.set)

	return setRouter
}

func (h *Handler) setupGetAPI() *gin.Engine {
	getRouter := gin.New()
	getRouter.Use(gin.Recovery(), gin.Logger())

	getRouter.GET("/", h.get)

	return getRouter
}

func (h *Handler) initBackup(routers []*gin.Engine) {
	for _, router := range routers {
		router.POST("/backup", h.backup)
	}
}
