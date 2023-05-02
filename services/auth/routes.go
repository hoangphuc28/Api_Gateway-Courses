package auth

import (
	"github.com/Zhoangp/Api_Gateway-Courses/config"
	"github.com/Zhoangp/Api_Gateway-Courses/services/auth/internal/delivery/http"
	"github.com/Zhoangp/Api_Gateway-Courses/services/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine, c *config.Config, middleware *middleware.MiddleareManager) {
	//define auth client
	client := NewAuthServiceClient(c)
	authHandler := http.NewAuthHandler(c, client)

	//create auth routes
	routes := r.Group("/auth")
	routes.POST("/login", authHandler.Login())
	routes.POST("/register", authHandler.Register())
	routes.GET("/token", authHandler.NewToken())
	routes.POST("/account", authHandler.GetTokenVerifyAccount())
	routes.POST("/password", authHandler.GetTokenResetPassword())

	routes.Use(middleware.RequireVerifyToken())
	routes.PATCH("/account/:token", authHandler.VerifyAccount())
}
