package routes

import (
	"test-task/internal/modules/auth"

	"github.com/gin-gonic/gin"
)

type AppRouter struct {
	Prefix  string
	Version string
	Routes  *gin.RouterGroup
}

func NewAppRouter(engine *gin.Engine, prefix string, version string) *AppRouter {
	router := engine.Group(prefix + version)
	return &AppRouter{
		Prefix:  prefix,
		Version: version,
		Routes:  router,
	}
}

func (r *AppRouter) RegisterAuthRoutes(handler *auth.Handler) {
	router := r.Routes.Group("/auth")

	router.POST("/login", handler.LoginUserHandler)
	router.POST("/signup", handler.RegisterUserHandler)
	router.POST("/issue-tokens/:id", handler.IssueTokensHandler)
	router.POST("/refresh-tokens", handler.RefreshTokensHandler)
}
