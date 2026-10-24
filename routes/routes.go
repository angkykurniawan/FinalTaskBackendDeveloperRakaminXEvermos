package routes

import (
	"evermos-mini/controllers"
	"evermos-mini/middleware"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App          *fiber.App
	AuthCtrl     *controllers.AuthController
	UserCtrl     *controllers.UserController
	TokoCtrl     *controllers.TokoController
	AlamatCtrl   *controllers.AlamatController
	CategoryCtrl *controllers.CategoryController
	ProductCtrl  *controllers.ProductController
	TrxCtrl      *controllers.TrxController
	ProvCityCtrl *controllers.ProvCityController
}

func SetupRoutes(cfg RouteConfig) {
	app := cfg.App

	// ==================== AUTH ROUTES ====================
	auth := app.Group("/auth")
	auth.Post("/register", cfg.AuthCtrl.Register)
	auth.Post("/login", cfg.AuthCtrl.Login)

	// ==================== USER & ALAMAT ROUTES ====================
	user := app.Group("/user", middleware.JWTProtected)
	user.Get("/", cfg.UserCtrl.GetMyProfile)
	user.Put("/", cfg.UserCtrl.UpdateMyProfile)

	// Alamat Kirim
	user.Get("/alamat", cfg.AlamatCtrl.GetAllAlamat)
	user.Get("/alamat/:id", cfg.AlamatCtrl.GetAlamatByID)
	user.Post("/alamat", cfg.AlamatCtrl.CreateAlamat)
	user.Put("/alamat/:id", cfg.AlamatCtrl.UpdateAlamat)
	user.Delete("/alamat/:id", cfg.AlamatCtrl.DeleteAlamat)

	// ==================== TOKO ROUTES ====================
	// Endpoint publik & terproteksi toko
	app.Get("/toko", cfg.TokoCtrl.GetAllToko)
	app.Get("/toko/my", middleware.JWTProtected, cfg.TokoCtrl.GetMyToko)
	app.Get("/toko/:id_toko", cfg.TokoCtrl.GetTokoByID)
	app.Put("/toko/:id_toko", middleware.JWTProtected, cfg.TokoCtrl.UpdateToko)

	// ==================== CATEGORY ROUTES ====================
	app.Get("/category", cfg.CategoryCtrl.GetAllCategory)
	app.Get("/category/:id", cfg.CategoryCtrl.GetCategoryByID)

	categoryAdmin := app.Group("/category", middleware.JWTProtected, middleware.AdminOnly)
	categoryAdmin.Post("/", cfg.CategoryCtrl.CreateCategory)
	categoryAdmin.Put("/:id", cfg.CategoryCtrl.UpdateCategory)
	categoryAdmin.Delete("/:id", cfg.CategoryCtrl.DeleteCategory)

	// ==================== PRODUCT ROUTES ====================
	app.Get("/product", cfg.ProductCtrl.GetAllProduct)
	app.Get("/product/:id", cfg.ProductCtrl.GetProductByID)

	productProtected := app.Group("/product", middleware.JWTProtected)
	productProtected.Post("/", cfg.ProductCtrl.CreateProduct)
	productProtected.Put("/:id", cfg.ProductCtrl.UpdateProduct)
	productProtected.Delete("/:id", cfg.ProductCtrl.DeleteProduct)

	// ==================== TRANSAKSI ROUTES ====================
	trx := app.Group("/trx", middleware.JWTProtected)
	trx.Get("/", cfg.TrxCtrl.GetAllTrx)
	trx.Get("/:id", cfg.TrxCtrl.GetTrxByID)
	trx.Post("/", cfg.TrxCtrl.CreateTrx)
	trx.Put("/:id/status", cfg.TrxCtrl.UpdateStatusTrx)

	// ==================== PROVINSI & KOTA ROUTES ====================
	provcity := app.Group("/provcity")
	provcity.Get("/listprovincies", cfg.ProvCityCtrl.GetListProvincies)
	provcity.Get("/listcities/:prov_id", cfg.ProvCityCtrl.GetListCities)
	provcity.Get("/detailprovince/:prov_id", cfg.ProvCityCtrl.GetDetailProvince)
	provcity.Get("/detailcity/:city_id", cfg.ProvCityCtrl.GetDetailCity)
}
