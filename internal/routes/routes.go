package routes

import (
	"github.com/gin-gonic/gin"
	handlers "url-shortener/internal/handler"
)

func PublicRoutes(r *gin.RouterGroup, h *handlers.Handler) {
	if r == nil || h == nil {
		panic("nil pointer")
	}

	r.GET("/:id", h.GetLinkHandler)
	r.POST("/", h.CreateLinkHandler)
	r.POST("/api/shorten", h.APICreateLinkHandler)
}
