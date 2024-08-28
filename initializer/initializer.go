package initializer

import (
	"test-task/internal/config"
	db "test-task/internal/database"
	"test-task/internal/modules/auth"
	"test-task/internal/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type AppWrapper struct {
	Engine   *gin.Engine
	Config   *config.Config
	Database *db.DBHandler
	Router   *routes.AppRouter
}

func InitializeApp() (*AppWrapper, error) {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		return nil, err
	}

	dbHandler := db.InitDB(cfg.DBSource)
	app := setupGin(cfg)

	router := routes.NewAppRouter(app, "/api", "/v1")
	err = InitializeModule(dbHandler, cfg, auth.InitAuthService, auth.NewHandler, router.RegisterAuthRoutes)
	if err != nil {
		return nil, err
	}

	return &AppWrapper{
		Engine:   app,
		Config:   cfg,
		Database: &dbHandler,
		Router:   router,
	}, nil
}

// generic service initializer
type Service interface{}
type Handler interface{}

func InitializeModule[T Service, H Handler](
	dbHandler db.DBHandler,
	cfg *config.Config,
	initService func(dbHandler db.DBHandler, cfg *config.Config) (T, error),
	createHandler func(T, *config.Config) H,
	registerRoutes func(H)) error {

	service, err := initService(dbHandler, cfg)
	if err != nil {
		return err
	}

	handler := createHandler(service, cfg)

	registerRoutes(handler)
	return nil
}

// cors configuration
func CorsConfig(cfg *config.Config) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// gin setup
func setupGin(cfg *config.Config) *gin.Engine {
	app := gin.New()
	app.Use(CorsConfig(cfg))

	if gin.Mode() == gin.DebugMode {
		app.Use(gin.Logger())
	}

	app.Use(gin.Recovery())

	return app
}
