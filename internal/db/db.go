package db

import (
	"fmt"
	"log"
	"os"

	"Music-lib/internal/models" // Подключаем модели для миграций
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Загружаем переменные окружения из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Формируем строку подключения на основе переменных окружения
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))

	// Подключаемся к базе данных
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Вызываем функцию миграций
	migrateDB()
}

// Миграции для создания/обновления таблиц
func migrateDB() {
	err := DB.AutoMigrate(&models.Song{}) // Добавляем модели, которые нужно мигрировать
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	fmt.Println("Database migrated successfully")
}
