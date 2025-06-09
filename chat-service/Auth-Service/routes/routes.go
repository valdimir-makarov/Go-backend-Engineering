package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/Auth-Service/controllers" // Your actual handlers
)

// SetupRoutes configures the API routes
func SetupRoutes(r *gin.Engine, ctrl *controllers.Controller) {
	r.POST("/login", ctrl.Login)
	r.POST("/register", ctrl.Register)
	r.GET("/profile", ctrl.Profile)
	r.POST("/logout", ctrl.Logout)
}
