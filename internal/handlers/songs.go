package handlers

import (
	"Music-lib/internal/db"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type Song struct {
	ID uint `gorm:"primaryKey"`
	gorm.Model
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type GeniusSearchResponse struct {
	Response struct {
		Hits []struct {
			Result struct {
				FullTitle string `json:"full_title"`
				URL       string `json:"url"`
			} `json:"result"`
		} `json:"hits"`
	} `json:"response"`
}

// GetSongs godoc
// @Summary Получить все песни
// @Description Получить список всех песен с пагинацией
// @Tags songs
// @Produce json
// @Param page query int false "Номер страницы"
// @Param limit query int false "Количество песен на странице"
// @Success 200 {array} Song
// @Failure 400 {object} ErrorResponse
// @Router /songs [get]

func GetSongs(w http.ResponseWriter, r *http.Request) {
	var songs []Song
	query := db.DB.Model(&Song{})

	// Фильтрация по группе
	group := r.URL.Query().Get("group")
	if group != "" {
		query = query.Where("\"group\" = ?", group) // Используем двойные кавычки вокруг group
	}

	// Фильтрация по названию песни
	song := r.URL.Query().Get("song")
	if song != "" {
		query = query.Where("song = ?", song)
	}

	// Пагинация
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	var pageNumber, limitNumber int
	var err error
	if page != "" {
		pageNumber, err = strconv.Atoi(page)
		if err != nil || pageNumber < 1 {
			pageNumber = 1 // Стандартное значение - первая страница
		}
	} else {
		pageNumber = 1
	}

	if limit != "" {
		limitNumber, err = strconv.Atoi(limit)
		if err != nil || limitNumber < 1 {
			limitNumber = 10 // Стандартное значение - 10 записей на странице
		}
	} else {
		limitNumber = 10
	}

	// Применяем пагинацию
	query = query.Offset((pageNumber - 1) * limitNumber).Limit(limitNumber)

	// Выполняем запрос
	if err = query.Find(&songs).Error; err != nil {
		http.Error(w, "Failed to retrieve songs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response, err := json.MarshalIndent(songs, "", "    ")
	if err != nil {
		http.Error(w, "Failed to format JSON", http.StatusInternalServerError)
		return
	}

	w.Write(response)

	w.WriteHeader(http.StatusOK)
}

//@Summary Получить текст песни
//@Description Возвращает текст песни с пагинацией по куплетам
//@Tags Songs
//@Accept json
//@Produce json
//@Param id path int true "ID песни"
//@Param page query int false "Номер страницы"
//@Param limit query int false "Количество куплетов на странице"
//@Success 200 {array} string "Текст песни с пагинацией"
//@Failure 404 {string} string "Песня не найдена"
//@Router /songs/{id}/text [get] */

func GetSongText(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var song Song
	if err := db.DB.First(&song, id).Error; err != nil {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	// Разделяем текст песни на куплеты по символу новой строки
	verses := strings.Split(song.Text, "\n")

	// Пагинация по куплетам
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	var pageNumber, pageSize int
	fmt.Sscanf(page, "%d", &pageNumber)
	fmt.Sscanf(limit, "%d", &pageSize)

	if pageNumber > 0 && pageSize > 0 {
		start := (pageNumber - 1) * pageSize
		end := start + pageSize
		if start > len(verses) {
			start = len(verses)
		}
		if end > len(verses) {
			end = len(verses)
		}
		verses = verses[start:end]
	}

	w.Header().Set("Content-Type", "application/json")
	response, err := json.MarshalIndent(verses, "", "    ") // 4 пробела отступа
	if err != nil {
		http.Error(w, "Failed to format JSON", http.StatusInternalServerError)
		return
	}

	w.Write(response)
	w.WriteHeader(http.StatusOK)
}

// @Summary Удалить песню
// @Description Удаляет песню из библиотеки по её ID
// @Tags Songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Success 200 {string} string "Песня успешно удалена"
// @Failure 404 {string} string "Песня не найдена"
// @Router /songs/{id} [delete] */

func DeleteSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := db.DB.Delete(&Song{}, id).Error; err != nil {
		http.Error(w, "Failed to delete song", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Song deleted successfully"))
}

//@Summary Обновить данные песни
//@Description Обновляет информацию о песне в библиотеке по её ID
//@Tags Songs
//@Accept json
//@Produce json
//@Param id path int true "ID песни"
//@Param song body Song true "Данные для обновления"
//@Success 200 {object} Song "Обновлённая песня"
//@Failure 404 {string} string "Песня не найдена"
//@Failure 400 {string} string "Неверный формат данных"
//@Router /songs/{id} [put] */

func UpdateSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var song Song
	if err := db.DB.First(&song, id).Error; err != nil {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := db.DB.Save(&song).Error; err != nil {
		http.Error(w, "Failed to update song", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(song)
}
