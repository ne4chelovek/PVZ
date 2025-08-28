package routes

import (
	"PVZ/internal/middleware"
	"PVZ/internal/utils"
	"PVZ/pkg/handler"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(handler *handler.PVZHandler, tokenUtils utils.TokenUtils) (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(
		middleware.LoggingMiddleware(),
		middleware.PrometheusMiddleware(),
		gin.Recovery(),
	)

	// Открытые эндпоинты
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
	r.POST("/dummyLogin", handler.LoginDummy)

	// Защищённые: все требуют токен
	protected := r.Group("/")

	// ПВЗ — только для moderator
	protected.POST("/pvz", middleware.AuthMiddleware(tokenUtils, "moderator"), handler.CreatePVZ)
	protected.GET("/pvz", handler.GetAllPVZ) // employee и moderator

	// Приёмки — только employee
	protected.POST("/receptions", middleware.AuthMiddleware(tokenUtils, "employee"), handler.CreateReception)
	protected.POST("/products", middleware.AuthMiddleware(tokenUtils, "employee"), handler.AddProduct)
	protected.POST("/pvz/:pvzId/close_last_reception", middleware.AuthMiddleware(tokenUtils, "employee"), handler.CloseLastReception)
	protected.POST("/pvz/:pvzId/delete_last_product", middleware.AuthMiddleware(tokenUtils, "employee"), handler.DeleteLastProduct)

	return r, nil
}
