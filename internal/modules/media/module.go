package media

import (
	"thyngo/internal/middleware"

	"github.com/gin-gonic/gin"
)

type MediaModule struct {
	service *Service
}

func New() *MediaModule {
	return &MediaModule{
		service: NewService(),
	}
}

func (m *MediaModule) Name() string {
	return "media"
}

func (m *MediaModule) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("", m.listMediaHandler)
	router.GET("/:slug", m.getMediaHandler)

	router.GET("/movies", m.listMoviesHandler)
	router.GET("/series", m.listSeriesHandler)
	router.GET("/albums", m.listAlbumsHandler)

	protected := router.Group("")
	protected.Use(middleware.RequireAuth())
	{
		protected.POST("", m.createMediaHandler)
		protected.PUT("/:slug", m.updateMediaHandler)
		protected.DELETE("/:slug", m.deleteMediaHandler)

		protected.POST("/movies", m.createMovieHandler)
		protected.PUT("/movies/:id", m.updateMovieHandler)
		protected.POST("/series", m.createSeriesHandler)
		protected.PUT("/series/:id", m.updateSeriesHandler)
		protected.POST("/albums", m.createAlbumHandler)
		protected.PUT("/albums/:id", m.updateAlbumHandler)
		
		ext := protected.Group("/external")
		{
			ext.GET("/movies/search", m.searchExternalMoviesHandler)
			ext.GET("/series/search", m.searchExternalSeriesHandler)
			ext.GET("/music/search", m.searchExternalAlbumsHandler)
			ext.GET("/music/albums/:mbid/artwork", m.getExternalAlbumArtworkHandler)
			ext.GET("/movies/:id/artwork", m.getExternalMovieArtworkHandler)
			ext.GET("/series/:id/artwork", m.getExternalSeriesArtworkHandler)
			ext.POST("/import-external", m.importExternalMediaHandler)
		}
	}
}
