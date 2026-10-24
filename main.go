package main

import (
	"evermos-mini/config"
	"evermos-mini/controllers"
	"evermos-mini/models"
	"evermos-mini/repositories"
	"evermos-mini/routes"
	"evermos-mini/usecases"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Inisialisasi Database
	db := config.InitDB()

	// Eksekusi AutoMigrate skema database
	err := db.AutoMigrate(
		&models.User{},
		&models.Toko{},
		&models.Alamat{},
		&models.Category{},
		&models.Produk{},
		&models.LogProduk{},
		&models.Transaksi{},
		&models.DetailTransaksi{},
	)
	if err != nil {
		log.Fatalf("Auto migration database gagal: %v", err)
	}
	log.Println("Migrasi tabel database berhasil selesai.")

	// Eksekusi seeder kategori default
	config.SeedData(db)

	// Pastikan folder uploads tersedia
	_ = os.MkdirAll("./uploads", os.ModePerm)

	// Inisialisasi Layer Repository
	userRepo := repositories.NewUserRepository(db)
	tokoRepo := repositories.NewTokoRepository(db)
	alamatRepo := repositories.NewAlamatRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	productRepo := repositories.NewProductRepository(db)
	trxRepo := repositories.NewTrxRepository(db)

	// Inisialisasi Layer Usecase
	authUsecase := usecases.NewAuthUsecase(userRepo, tokoRepo)
	userUsecase := usecases.NewUserUsecase(userRepo)
	tokoUsecase := usecases.NewTokoUsecase(tokoRepo)
	alamatUsecase := usecases.NewAlamatUsecase(alamatRepo)
	categoryUsecase := usecases.NewCategoryUsecase(categoryRepo)
	productUsecase := usecases.NewProductUsecase(productRepo, tokoRepo, categoryRepo)
	trxUsecase := usecases.NewTrxUsecase(trxRepo, alamatRepo)
	provCityUsecase := usecases.NewProvCityUsecase()

	// Inisialisasi Layer Controller
	authCtrl := controllers.NewAuthController(authUsecase)
	userCtrl := controllers.NewUserController(userUsecase)
	tokoCtrl := controllers.NewTokoController(tokoUsecase)
	alamatCtrl := controllers.NewAlamatController(alamatUsecase)
	categoryCtrl := controllers.NewCategoryController(categoryUsecase)
	productCtrl := controllers.NewProductController(productUsecase)
	trxCtrl := controllers.NewTrxController(trxUsecase)
	provCityCtrl := controllers.NewProvCityController(provCityUsecase)

	// Inisialisasi Fiber
	app := fiber.New(fiber.Config{
		AppName: "Evermos Mini Project - Rakamin Backend",
	})

	// Static route untuk file upload
	app.Static("/uploads", "./uploads")

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, token",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	app.Use(logger.New())
	app.Use(recover.New())

	// Registrasi Routes
	routes.SetupRoutes(routes.RouteConfig{
		App:          app,
		AuthCtrl:     authCtrl,
		UserCtrl:     userCtrl,
		TokoCtrl:     tokoCtrl,
		AlamatCtrl:   alamatCtrl,
		CategoryCtrl: categoryCtrl,
		ProductCtrl:  productCtrl,
		TrxCtrl:      trxCtrl,
		ProvCityCtrl: provCityCtrl,
	})

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = ":8080"
	}
	if port[0] != ':' {
		port = ":" + port
	}

	fmt.Printf("Aplikasi berjalan di http://localhost%s\n", port)
	if err := app.Listen(port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
