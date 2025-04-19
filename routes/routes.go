package routes

import (
	"github.com/gin-gonic/gin"
	//"github.com/tnqbao/gau_phim_backend/api/admin"
	//"github.com/tnqbao/gau_phim_backend/api/admin/movie"
	//"github.com/tnqbao/gau_phim_backend/api/authed"
	//"github.com/tnqbao/gau_phim_backend/api/public"
	"github.com/tnqbao/gau_truyen_backend/middlewares"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.CORSMiddleware())
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	//apiRoutes := r.Group("/api/gautruyen")
	//{
	//	publicRouter := apiRoutes.Group("/")
	//	{
	//		publicRouter.GET("/home-page", public.GetHomePageData)
	//		publicRouter.GET("/category/:slug", public.GetListMovieByCategory)
	//		publicRouter.GET("/type/:slug", public.GetListMovieByType)
	//		publicRouter.GET("/nation/:slug", public.GetListMovieByNation)
	//
	//	}
	//	adminRoutes := apiRoutes.Group("/")
	//	{
	//		adminRoutes.Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	//		adminRoutes.PUT("/crawl", movie.CrawlMovieFromUrl)
	//		adminRoutes.POST("/movie", movie.CreateMovie)
	//		adminRoutes.PUT("/home-page/hero", admin.UpdateHeroHomePage)
	//		adminRoutes.PUT("/home-page/release", admin.UpdateReleaseHomePage)
	//		adminRoutes.PUT("/home-page/featured", admin.UpdateFeaturedHomePage)
	//	}
	//
	//	authedRoutes := apiRoutes.Group("/")
	//	{
	//		authedRoutes.Use(middlewares.AuthMiddleware())
	//		authedRoutes.POST("/like", authed.AddMovieLiked)
	//		authedRoutes.GET("/likes", authed.GetListMovieLiked)
	//		authedRoutes.DELETE("/like", authed.RemoveMovieLiked)
	//
	//		authedRoutes.GET("/history", authed.GetHistoryView)
	//		authedRoutes.POST("/history", authed.UpdateHistoryView)
	//		authedRoutes.DELETE("/history/:slug", authed.DeleteHistoryViewForSlug)
	//		authedRoutes.DELETE("/history", authed.DeleteAllHistoryView)
	//	}
	//
	//}
	return r
}
