package route

import (
	"be_imoca_golang/controller"
	"be_imoca_golang/middleware" 
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "IMOCA API - Online")
	})

	// Group AUTH 
	auth := r.Group("/auth")
	{
		auth.POST("/login", controller.LoginHandler(db))
		auth.POST("/register", controller.RegisterHandler(db)) // Untuk buat akun admin
	}

	// Group PUBLIC 
	public := r.Group("/public")
	{
		public.GET("/partners", controller.GetPartnersHandler(db))
		public.GET("/partners/:id", controller.GetPartnerByIDHandler(db))
		public.GET("/news", controller.GetNewsHandler(db))
		public.GET("/news/:id", controller.GetNewsByIDHandler(db))
		public.GET("/members", controller.GetMembersHandler(db))
		
	}

	// Group ADMIN 
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware()) 
	{
		// PARTNER
		admin.POST("/partners", controller.CreatePartnerHandler(db))
		admin.PUT("/partners/:id", controller.UpdatePartnerHandler(db))
		admin.DELETE("/partners/:id", controller.DeletePartnerHandler(db))

		// MEMBER
		admin.POST("/members", controller.CreateMemberHandler(db))
		admin.PUT("/members/:id", controller.UpdateMemberHandler(db))
		admin.DELETE("/members/:id", controller.DeleteMemberHandler(db))

		// NEWS
		admin.POST("/news", controller.CreateNewsHandler(db))
		admin.PUT("/news/:id", controller.UpdateNewsHandler(db))
		admin.DELETE("/news/:id", controller.DeleteNewsHandler(db))
	}
}