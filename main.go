package main

import (
	_ "Music-lib/docs"
	data "Music-lib/internal/db"
	"Music-lib/internal/handlers"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

//	@Summary		Добавить новую песню
//	@Description	Добавляет новую песню в библиотеку
//	@Accept			json
//	@Produce		json
//	@Param			song	body		Song	true	"Песня"
//	@Success		201		{object}	Song
//	@Failure		400		{object}	ErrorResponse
//	@Router			/songs [post]

func AddSong(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting to process AddSong request")

	var song handlers.Song
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		log.Println("Error decoding song:", err)
		http.Error(w, "Failed to decode the song", http.StatusBadRequest)
		return
	}

	geniusResponse, err := handlers.GetSongInfoFromGenius(song.Group, song.Song)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve song info: %v", err), http.StatusInternalServerError)
		return
	}

	// Если API вернул данные, обогащаем нашу запись
	if len(geniusResponse.Response.Hits) > 0 {
		song.Link = geniusResponse.Response.Hits[0].Result.URL

		lyrics, releaseDate, err := handlers.GetLyricsFromGeniusPage(song.Link)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve lyrics: %v", err), http.StatusInternalServerError)
		}
		song.Text = lyrics
		song.ReleaseDate = releaseDate
	}

	// Сохраняем песню в базе данных
	if err = data.DB.Create(&song).Error; err != nil {
		log.Println("Error saving song:", err)
		http.Error(w, "Failed to save song", http.StatusInternalServerError)
		return
	}

	log.Println("Successfully saved song")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(song)
}

func main() {

	data.ConnectDB()

	// Добавь сюда маршруты и сервер, которые мы создадим позже
	r := mux.NewRouter()

	r.HandleFunc("/songs", AddSong).Methods("POST")
	r.HandleFunc("/songs", handlers.GetSongs).Methods("GET")
	r.HandleFunc("/info", handlers.GetSongs).Methods("GET")
	r.HandleFunc("/songs/{id:[0-9]+}/text", handlers.GetSongText).Methods("GET")
	r.HandleFunc("/songs/{id:[0-9]+}", handlers.DeleteSong).Methods("DELETE")
	r.HandleFunc("/songs/{id:[0-9]+}", handlers.UpdateSong).Methods("PUT")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", r))
}
