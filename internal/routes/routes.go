package routes

import (
	"github.com/gin-gonic/gin"
	handlers "url-shortener/internal/handler"
)

func PublicRoutes(r *gin.RouterGroup, h handlers.Handler) {
	r.GET("/:id", h.GetLinkHandler)
	r.POST("/", h.CreateLinkHandler)
	r.POST("/api/shorten", h.APICreateLinkHandler)
}
