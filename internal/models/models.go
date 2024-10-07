package models

import (
	"gorm.io/gorm"
)

type Song struct {
	gorm.Model         // Встроенные поля ID, CreatedAt, UpdatedAt, DeletedAt
	Group       string `json:"group"`        // Название группы или исполнителя
	Song        string `json:"song"`         // Название песни
	ReleaseDate string `json:"release_date"` // Дата релиза
	Text        string `json:"text"`         // Текст песни
	Link        string `json:"link"`         // Ссылка на песню или источник
}
