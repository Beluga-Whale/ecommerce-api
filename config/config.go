package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Beluga-Whale/ecommerce-api/internal/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

var TestDB *gorm.DB

func LoadEnv() {
	env := os.Getenv("APP_ENV")

	if env == "production" {
		fmt.Println("✅ Running in production mode: using ENV variables only")
		return
	}

	if env == "" {
		env = "development"
	}

	envFileMap := map[string]string{
		"development":    ".env",
		"test":           ".env.test",
		"test.localhost": ".env.test.localhost",
		"production":     ".env.production",
	}

	envFile, ok := envFileMap[env]
	if !ok {
		log.Fatalf("❌ Invalid APP_ENV: %s", env)
	}

	// ✅ ใช้ runtime.Caller เพื่อให้ไม่หลุด directory เวลา go test
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("❌ Failed to get current file path")
	}

	// currentFile → /path/to/project/server/config/config.go
	serverDir := filepath.Join(filepath.Dir(currentFile), "..") // เดินขึ้นจาก /config → /server
	envPath := filepath.Join(serverDir, envFile)

	fmt.Println("🔧 Loading env from:", envPath)

	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("❌ Failed to load env: %v", err)
	}
}

func ConnectDB() {
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
	os.Getenv("HOST"),
	os.Getenv("USER_NAME"),
	os.Getenv("PASSWORD"),
	os.Getenv("DATABASE_NAME"),
	os.Getenv("PORT"),
	os.Getenv("SSL_MODE"),
	)

	fmt.Println("🔍 ENV:", os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("DATABASE_NAME"), os.Getenv("USER_NAME"))


	// New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
		SlowThreshold: time.Second, // Slow SQL threshold
		LogLevel:      logger.Info, // Log level
		Colorful:      true,        // Enable color
		},
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger, // add Logger
	})

	if err != nil {
		log.Fatal("Fail to connect DB : ",err)
	}

	fmt.Println("Connect DB Success!")

	DB.Exec(`
	DO $$ BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status') THEN
		CREATE TYPE status AS ENUM ('pending', 'paid', 'shipped', 'cancel');
	END IF;
	END$$;
	`)

	DB.Exec(`
	DO $$ BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role') THEN
		CREATE TYPE role AS ENUM ('user', 'admin');
	END IF;
	END$$;
	`)

	DB.Exec(`
	DO $$ BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment_status') THEN
		CREATE TYPE payment_status AS ENUM ('payed', 'failed');
	END IF;
	END$$;
	`)


	// NOTE - AutoMigrate จะตรวจสอบและอัปเดตฐานข้อมูล
	err = DB.AutoMigrate(
		&models.CartItem{},   // NOTE - ให้ตรวจสอบตาราง CartItem
		&models.Category{},   // NOTE - ให้ตรวจสอบตาราง Category
		&models.Coupon{},   // NOTE - ให้ตรวจสอบตาราง Coupon
		&models.Order{},   // NOTE - ให้ตรวจสอบตาราง Order
		&models.OrderItem{},   // NOTE - ให้ตรวจสอบตาราง OrderItem
		&models.Payment{},   // NOTE - ให้ตรวจสอบตาราง Payment
		&models.Product{},   // NOTE - ให้ตรวจสอบตาราง Product
		&models.ProductVariant{}, // NOTE - ให้ตรวจสอบตาราง ProductVariant
		&models.Review{},   // NOTE - ให้ตรวจสอบตาราง Review
		&models.User{},   // NOTE - ให้ตรวจสอบตาราง User
		&models.ProductImage{}, // NOTE - ให้ตรวจสอบตาราง ProductImage
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

}

func ConnectTestDB() {
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
	os.Getenv("HOST"),
	os.Getenv("USER_NAME"),
	os.Getenv("PASSWORD"),
	os.Getenv("DATABASE_NAME"),
	os.Getenv("PORT"),
	os.Getenv("SSL_MODE"),
	)

	fmt.Println("🔍 ENV:", os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("DATABASE_NAME"), os.Getenv("USER_NAME"))


	// New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
		SlowThreshold: time.Second, // Slow SQL threshold
		LogLevel:      logger.Silent, // Log level
		Colorful:      true,        // Enable color
		},
	)

	TestDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger, // add Logger
	})

	if err != nil {
		log.Fatal("Fail to connect TestDB : ",err)
	}

	fmt.Println("Connect TestDB Success!")

	TestDB.Exec(`
	DO $$ BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status') THEN
		CREATE TYPE status AS ENUM ('pending', 'paid', 'shipped', 'cancel');
	END IF;
	END$$;
	`)

	TestDB.Exec(`
	DO $$ BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role') THEN
		CREATE TYPE role AS ENUM ('user', 'admin');
	END IF;
	END$$;
	`)

	TestDB.Exec(`
	DO $$ BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment_status') THEN
		CREATE TYPE payment_status AS ENUM ('payed', 'failed');
	END IF;
	END$$;
	`)


	// NOTE - AutoMigrate จะตรวจสอบและอัปเดตฐานข้อมูล
	err = TestDB.AutoMigrate(
		&models.CartItem{},   // NOTE - ให้ตรวจสอบตาราง CartItem
		&models.Category{},   // NOTE - ให้ตรวจสอบตาราง Category
		&models.Coupon{},   // NOTE - ให้ตรวจสอบตาราง Coupon
		&models.Order{},   // NOTE - ให้ตรวจสอบตาราง Order
		&models.OrderItem{},   // NOTE - ให้ตรวจสอบตาราง OrderItem
		&models.Payment{},   // NOTE - ให้ตรวจสอบตาราง Payment
		&models.Product{},   // NOTE - ให้ตรวจสอบตาราง Product
		&models.ProductVariant{}, // NOTE - ให้ตรวจสอบตาราง ProductVariant
		&models.Review{},   // NOTE - ให้ตรวจสอบตาราง Review
		&models.User{},   // NOTE - ให้ตรวจสอบตาราง User
		&models.ProductImage{}, // NOTE - ให้ตรวจสอบตาราง ProductImage
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

}