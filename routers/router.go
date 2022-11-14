package routers

import (
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/controllers"
	"bitbucket.org/isbtotogroup/isbpanel_api_backend/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Init() *fiber.App {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		// Set some security headers:
		// c.Set("Content-Security-Policy", "frame-ancestors 'none'")
		// c.Set("X-XSS-Protection", "1; mode=block")
		// c.Set("X-Content-Type-Options", "nosniff")
		// c.Set("X-Download-Options", "noopen")
		// c.Set("Strict-Transport-Security", "max-age=5184000")
		// c.Set("X-Frame-Options", "SAMEORIGIN")
		// c.Set("X-DNS-Prefetch-Control", "off")

		// Go to next middleware:
		return c.Next()
	})
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(compress.New())
	app.Get("/ipaddress", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":      fiber.StatusOK,
			"message":     "Success",
			"record":      "data",
			"BASEURL":     c.BaseURL(),
			"HOSTNAME":    c.Hostname(),
			"IP":          c.IP(),
			"IPS":         c.IPs(),
			"OriginalURL": c.OriginalURL(),
			"Path":        c.Path(),
			"Protocol":    c.Protocol(),
			"Subdomain":   c.Subdomains(),
		})
	})
	app.Get("/dashboard", monitor.New())

	app.Post("/api/login", controllers.CheckLogin)
	app.Post("/api/valid", middleware.JWTProtected(), controllers.Home)
	app.Post("/api/alladmin", middleware.JWTProtected(), controllers.Adminhome)
	app.Post("/api/detailadmin", middleware.JWTProtected(), controllers.AdminDetail)
	app.Post("/api/saveadmin", middleware.JWTProtected(), controllers.AdminSave)

	app.Post("/api/alladminrule", middleware.JWTProtected(), controllers.Adminrulehome)
	app.Post("/api/saveadminrule", middleware.JWTProtected(), controllers.AdminruleSave)

	app.Post("/api/pasaran", middleware.JWTProtected(), controllers.Pasaranhome)
	app.Post("/api/pasaransave", middleware.JWTProtected(), controllers.Pasaransave)
	app.Post("/api/keluaran", middleware.JWTProtected(), controllers.Keluaranhome)
	app.Post("/api/keluaransave", middleware.JWTProtected(), controllers.Keluaransave)
	app.Post("/api/keluarandelete", middleware.JWTProtected(), controllers.Keluarandelete)

	app.Post("/api/prediksi", middleware.JWTProtected(), controllers.Prediksihome)
	app.Post("/api/prediksisave", middleware.JWTProtected(), controllers.Prediksisave)
	app.Post("/api/prediksidelete", middleware.JWTProtected(), controllers.Prediksidelete)

	app.Post("/api/tafsirmimpi", middleware.JWTProtected(), controllers.Tafsirmimpihome)
	app.Post("/api/tafsirmimpisave", middleware.JWTProtected(), controllers.Tafsirmimpisave)

	app.Post("/api/news", middleware.JWTProtected(), controllers.Newshome)
	app.Post("/api/newssave", middleware.JWTProtected(), controllers.Newssave)
	app.Post("/api/newsdelete", middleware.JWTProtected(), controllers.Newsdelete)
	app.Post("/api/categorynews", middleware.JWTProtected(), controllers.Categoryhome)
	app.Post("/api/categorynewssave", middleware.JWTProtected(), controllers.Categorysave)
	app.Post("/api/categorynewsdelete", middleware.JWTProtected(), controllers.Categorydelete)

	app.Post("/api/movie", middleware.JWTProtected(), controllers.Moviehome)
	app.Post("/api/movienotcdn", middleware.JWTProtected(), controllers.Moviehomenotcdn)
	app.Post("/api/moviebanner", middleware.JWTProtected(), controllers.Moviehomebanner)
	app.Post("/api/movietrouble", middleware.JWTProtected(), controllers.Movietroublehome)
	app.Post("/api/moviemini", middleware.JWTProtected(), controllers.Movieminihome)
	app.Post("/api/moviesave", middleware.JWTProtected(), controllers.Moviesave)
	app.Post("/api/moviedelete", middleware.JWTProtected(), controllers.Moviedelete)
	app.Post("/api/movieseries", middleware.JWTProtected(), controllers.Moviehomeseries)
	app.Post("/api/movieseriestrouble", middleware.JWTProtected(), controllers.Moviehomeseriestrouble)
	app.Post("/api/movieseriessave", middleware.JWTProtected(), controllers.Movieseriessave)
	app.Post("/api/movieseriesseason", middleware.JWTProtected(), controllers.Seasonhome)
	app.Post("/api/movieseriesseasonsave", middleware.JWTProtected(), controllers.Seasonsave)
	app.Post("/api/movieseriesseasondelete", middleware.JWTProtected(), controllers.Seasondelete)
	app.Post("/api/movieseriesepisode", middleware.JWTProtected(), controllers.Episodehome)
	app.Post("/api/movieseriesepisodesave", middleware.JWTProtected(), controllers.Episodesave)
	app.Post("/api/movieseriesepisodedelete", middleware.JWTProtected(), controllers.Episodedelete)
	app.Post("/api/moviecloudualbum", middleware.JWTProtected(), controllers.Moviecloud)
	app.Post("/api/moviecloudupdate", middleware.JWTProtected(), controllers.Movieupdatecloud)
	app.Post("/api/movieclouddelete", middleware.JWTProtected(), controllers.Moviedeletecloud)
	app.Post("/api/moviecloudupload", middleware.JWTProtected(), controllers.Movieuploadcloud)
	app.Post("/api/genremovie", middleware.JWTProtected(), controllers.Genrehome)
	app.Post("/api/genremoviesave", middleware.JWTProtected(), controllers.Genresave)
	app.Post("/api/genremoviedelete", middleware.JWTProtected(), controllers.Genredelete)
	app.Post("/api/moviebannersave", middleware.JWTProtected(), controllers.MoviebannerSave)
	app.Post("/api/moviebannerdelete", middleware.JWTProtected(), controllers.MoviebannerDelete)

	app.Post("/api/slider", middleware.JWTProtected(), controllers.Sliderhome)
	app.Post("/api/slidersave", middleware.JWTProtected(), controllers.Slidersave)
	app.Post("/api/sliderdelete", middleware.JWTProtected(), controllers.Sliderdelete)

	app.Post("/api/domain", middleware.JWTProtected(), controllers.Domainhome)
	app.Post("/api/domainsave", middleware.JWTProtected(), controllers.DomainSave)

	app.Post("/api/webagen", middleware.JWTProtected(), controllers.Websiteagenhome)
	app.Post("/api/webagensave", middleware.JWTProtected(), controllers.Websiteagensave)
	app.Post("/api/game", middleware.JWTProtected(), controllers.Gamehome)
	app.Post("/api/gamesave", middleware.JWTProtected(), controllers.Gamesave)

	app.Post("/api/cloudflare", middleware.JWTProtected(), controllers.Moviecloud2)
	app.Post("/api/album", middleware.JWTProtected(), controllers.Albumhome)
	app.Post("/api/albumsave", middleware.JWTProtected(), controllers.Albumsave)

	app.Post("/api/crm", middleware.JWTProtected(), controllers.Crmhome)
	app.Post("/api/crmsales", middleware.JWTProtected(), controllers.Crmsales)
	app.Post("/api/crmdeposit", middleware.JWTProtected(), controllers.Crmdeposit)
	app.Post("/api/crmisbtv", middleware.JWTProtected(), controllers.Crmisbtvhome)
	app.Post("/api/crmduniafilm", middleware.JWTProtected(), controllers.Crmduniafilm)
	app.Post("/api/crmsave", middleware.JWTProtected(), controllers.CrmSave)
	app.Post("/api/crmsavestatus", middleware.JWTProtected(), controllers.CrmSavestatus)
	app.Post("/api/crmsalessave", middleware.JWTProtected(), controllers.CrmSalesSave)
	app.Post("/api/crmsalesdelete", middleware.JWTProtected(), controllers.CrmSalesdelete)
	app.Post("/api/crmsavesource", middleware.JWTProtected(), controllers.CrmSavesource)
	app.Post("/api/crmsavedatabase", middleware.JWTProtected(), controllers.CrmSavedatabase)
	app.Post("/api/crmsavemaintenance", middleware.JWTProtected(), controllers.CrmSavemaintenance)

	app.Post("/api/slotprovider", middleware.JWTProtected(), controllers.Providerslothome)
	app.Post("/api/slotprovidersave", middleware.JWTProtected(), controllers.ProviderslotSave)
	app.Post("/api/prediksislot", middleware.JWTProtected(), controllers.Prediksislothome)
	app.Post("/api/prediksislotsave", middleware.JWTProtected(), controllers.PrediksislotSave)
	app.Post("/api/prediksislotdelete", middleware.JWTProtected(), controllers.PrediksislotDelete)
	app.Post("/api/prediksislotgenerator", middleware.JWTProtected(), controllers.PrediksislotGenerator)

	app.Post("/api/banner", middleware.JWTProtected(), controllers.Bannerhome)
	app.Post("/api/bannersave", middleware.JWTProtected(), controllers.Bannersave)

	app.Post("/api/departement", middleware.JWTProtected(), controllers.Departementhome)
	app.Post("/api/departementsave", middleware.JWTProtected(), controllers.DepartementSave)
	app.Post("/api/employee", middleware.JWTProtected(), controllers.Employeehome)
	app.Post("/api/employeebydepart", middleware.JWTProtected(), controllers.EmployeeByDepart)
	app.Post("/api/employeebysalesperformance", middleware.JWTProtected(), controllers.EmployeeBySalesPerformance)
	app.Post("/api/employeesave", middleware.JWTProtected(), controllers.EmployeeSave)
	return app
}
