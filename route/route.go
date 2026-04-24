package route

import (
	controllers "be_imoca_golang/controller"
	"be_imoca_golang/middleware"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm" //
)

func InitRoutes(r *gin.RouterGroup, db *sql.DB, gormDB *gorm.DB) {

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "IMOCA API - Online")
	})

	heroCtrl := &controllers.HeroController{DB: gormDB}
	orgCtrl := &controllers.OrganizationController{DB: gormDB}
	memberCtrl := &controllers.MemberController{DB: gormDB}
	partnerCtrl := &controllers.PartnerController{DB: gormDB}
	serviceCtrl := &controllers.ServiceController{DB: gormDB}
	missionCtrl := &controllers.MissionController{DB: gormDB}
	visionCtrl := &controllers.VisionController{DB: gormDB}
	newsCtrl := &controllers.NewsController{DB: gormDB}

	// Group AUTH /api/auth
	auth := r.Group("/auth")
	{
		auth.POST("/login", controllers.LoginHandler(db))
		auth.POST("/register", controllers.RegisterHandler(db))
	}

	// Group PUBLIC /api/public
	public := r.Group("/public")
	{
		public.GET("/hero", heroCtrl.GetHero)
		public.GET("/hero/images", heroCtrl.GetHeroImages)

		public.GET("/organization", orgCtrl.GetAll)

		public.GET("/members", memberCtrl.GetAll)
		public.GET("/members/:id", memberCtrl.GetByID)

		public.GET("/partners", partnerCtrl.GetAll)
		public.GET("/partners/:id", partnerCtrl.GetByID)

		public.GET("/news", newsCtrl.GetAll)
		public.GET("/news/:id", newsCtrl.GetByID)

		public.GET("/services", serviceCtrl.GetServices)
		public.GET("/vision", visionCtrl.GetVision)
		public.GET("/mission", missionCtrl.GetMission)
	}

	// Group ADMIN /api/admin
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	{
		// MEMBER
		admin.GET("/members", memberCtrl.GetAll)
		admin.POST("/members", memberCtrl.Create)
		admin.PUT("/members/:id", memberCtrl.Update)
		admin.DELETE("/members/:id", memberCtrl.Delete)

		// PARTNER
		admin.GET("/partners", partnerCtrl.GetAll)
		admin.POST("/partners", partnerCtrl.Create)
		admin.PUT("/partners/:id", partnerCtrl.Update)
		admin.DELETE("/partners/:id", partnerCtrl.Delete)

		// Organization
		admin.GET("/organization", orgCtrl.GetAll)
		admin.POST("/organization", orgCtrl.Create)
		admin.PUT("/organization/:id", orgCtrl.Update)
		admin.DELETE("/organization/:id", orgCtrl.Delete)

		// NEWS
		admin.GET("/news", newsCtrl.GetAll)
		admin.POST("/news", newsCtrl.Create)
		admin.PUT("/news/:id", newsCtrl.Update)
		admin.DELETE("/news/:id", newsCtrl.Delete)

		// EDITABLE CONTENT
		admin.PUT("/hero", heroCtrl.UpdateHero)
		admin.POST("/hero/images", heroCtrl.AddHeroImage)
		admin.DELETE("/hero/images/:id", heroCtrl.DeleteHeroImage)

		admin.PUT("/services", serviceCtrl.UpdateServices)
		admin.PUT("/vision", visionCtrl.UpdateVision)
		admin.PUT("/mission", missionCtrl.UpdateMission)

	}
}
