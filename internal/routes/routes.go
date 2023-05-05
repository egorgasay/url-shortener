package routes

import (
	"github.com/gin-gonic/gin"
	"net/http/pprof"
	handlers "url-shortener/internal/handler"
)

// PublicRoutes routes for unregistered users.
func PublicRoutes(r *gin.RouterGroup, h *handlers.Handler) {
	if r == nil || h == nil {
		panic("nil pointer")
	}

	r.GET("/:id", h.GetLinkHandler)
	r.GET("/api/user/urls", h.GetAllLinksHandler)
	r.GET("/ping", h.Ping)
	r.GET("/api/internal/stats", h.GetStatsHandler)

	r.Any("/debug/pprof/", gin.WrapF(pprof.Index))
	r.Any("/debug/pprof/cmdline", gin.WrapF(pprof.Cmdline))
	r.Any("/debug/pprof/profile", gin.WrapF(pprof.Profile))
	r.Any("/debug/pprof/symbol", gin.WrapF(pprof.Symbol))
	r.Any("/debug/pprof/trace", gin.WrapF(pprof.Trace))

	r.POST("/api/shorten/batch", h.BatchHandler)
	r.POST("/", h.CreateLinkHandler)
	r.POST("/api/shorten", h.APICreateLinkHandler)

	r.DELETE("/api/user/urls", h.APIDeleteLinksHandler)
}
