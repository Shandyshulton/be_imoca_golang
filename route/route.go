package route

import (
	"be_imoca_golang/controller"
	"be_imoca_golang/middleware"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.RouterGroup, db *sql.DB) {
	
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "IMOCA API - Online")
	})

	// Group AUTH /api/auth
	auth := r.Group("/auth")
	{
		auth.POST("/login", controller.LoginHandler(db))
		auth.POST("/register", controller.RegisterHandler(db))
	}

	// Group PUBLIC /api/public
	public := r.Group("/public")
	{
		public.GET("/hero", controller.GetHeroHandler(db))
		public.GET("/hero/images", controller.GetHeroImagesHandler(db))
		public.GET("/partners", controller.GetPartnersHandler(db))
		public.GET("/partners/:id", controller.GetPartnerByIDHandler(db))
		public.GET("/members", controller.GetMembersHandler(db))
        public.GET("/members/:id", controller.GetMemberByIDHandler(db))
		public.GET("/news", controller.GetNewsHandler(db))
		public.GET("/news/:id", controller.GetNewsByIDHandler(db))
		public.GET("/services", controller.GetServicesHandler(db))
		public.GET("/vision", controller.GetVisionHandler(db))
		public.GET("/mission", controller.GetMissionHandler(db))
	}

	// Group ADMIN /api/admin
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware()) 
	{
		// PARTNER
		admin.GET("/partners", controller.GetPartnersHandler(db))
		admin.POST("/partners", controller.CreatePartnerHandler(db))
		admin.PUT("/partners/:id", controller.UpdatePartnerHandler(db))
		admin.DELETE("/partners/:id", controller.DeletePartnerHandler(db))

		// MEMBER
		admin.GET("/members", controller.GetMembersHandler(db))
		admin.POST("/members", controller.CreateMemberHandler(db))
		admin.PUT("/members/:id", controller.UpdateMemberHandler(db))
		admin.DELETE("/members/:id", controller.DeleteMemberHandler(db))

		// NEWS
		admin.GET("/news", controller.GetNewsHandler(db))
		admin.POST("/news", controller.CreateNewsHandler(db))
		admin.PUT("/news/:id", controller.UpdateNewsHandler(db))
		admin.DELETE("/news/:id", controller.DeleteNewsHandler(db))

		// EDITABLE CONTENT
		admin.PUT("/hero", controller.UpdateHeroHandler(db))
		admin.POST("/hero/images", controller.AddHeroImageHandler(db))
		admin.DELETE("/hero/images/:id", controller.DeleteHeroImageHandler(db))

		admin.PUT("/services", controller.UpdateServicesHandler(db))
		admin.PUT("/vision", controller.UpdateVisionHandler(db))
		admin.PUT("/mission", controller.UpdateMissionHandler(db))
	}
}