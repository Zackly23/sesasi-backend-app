package config

import (
	"fmt"
	"log"
	"os"
	"sesasi-backend-app/models"
	"sesasi-backend-app/seeders"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct{
	DB *gorm.DB

}

func (d Database) getDatabaseConfig() map[string]interface{} {
	return map[string]interface{}{
		"host":     os.Getenv("DB_HOST"),
		"port":     os.Getenv("DB_PORT"),
		"user":     os.Getenv("DB_USER"),
		"password": os.Getenv("DB_PASSWORD"),
		"dbname":   os.Getenv("DB_NAME"),
		"sslmode":  os.Getenv("SSL_MODE"), // Use "require" for production
	}
}


func (d *Database) ConnectDatabase() (*gorm.DB, error) {
	dbConfig := d.getDatabaseConfig()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig["host"], dbConfig["port"], dbConfig["user"], dbConfig["password"], dbConfig["dbname"], dbConfig["sslmode"])

	var db *gorm.DB
	var err error

	// Retry loop
	maxRetries := 10
	for i := 1; i <= maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			// Try ping
			sqlDB, _ := db.DB()
			if pingErr := sqlDB.Ping(); pingErr == nil {
				fmt.Println("✅ Connected to PostgreSQL")
				break
			} else {
				fmt.Printf("❌ Ping failed: %v (attempt %d/%d)\n", pingErr, i, maxRetries)
			}
		} else {
			fmt.Printf("❌ Failed to connect to DB: %v (attempt %d/%d)\n", err, i, maxRetries)
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("💥 final failure after %d retries: %w", maxRetries, err)
	}

	// Auto migrate models
	if err := db.AutoMigrate(models.GetModels()...); err != nil {
		fmt.Println("Failed to auto migrate models:", err)
		return nil, err
	}

	// Seed Database
	if err := seeders.PermissionSeeder(db); err != nil {
		fmt.Println("Failed Seeding Pengajuan Permission")
	}

	if err := seeders.RolePermissionSeeder(db); err != nil {
		fmt.Println("Failed Seeding Role Permission ")
	}

	if err := seeders.UserSeeder(db); err != nil {
		fmt.Println("Failed Seeding User")
	}

	if err := seeders.PengajuanIzinSeed(db); err != nil {
		fmt.Println("Failed Seeding Pengajuan Izin")
	}

	if err := seeders.PengajuanIzinCommentSeed(db); err != nil {
		fmt.Println("Failed Seeding Pengajuan Izin Comment")
	}

	log.Println("Database connected successfully")

	return db, nil
}

func SetupDatabase() *gorm.DB {
	d := &Database{}              // bikin instance
	db, err := d.ConnectDatabase() // panggil method
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	log.Println("✅ Database connected successfully")
	
	return db
}
