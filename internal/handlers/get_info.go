package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"os"
	"strings"
)

func GetSongInfoFromGenius(artist, song string) (*GeniusSearchResponse, error) {
	token := os.Getenv("GENIUS_ACCESS_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("Genius access token not found")
	}

	url := fmt.Sprintf("https://api.genius.com/search?q=%s %s", artist, song)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Добавляем заголовок с токеном доступа
	req.Header.Add("Authorization", "Bearer "+token)

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch song info, status: %d", resp.StatusCode)
	}

	var searchResponse GeniusSearchResponse
	if err = json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}

	return &searchResponse, nil
}

func GetLyricsFromGeniusPage(url string) (string, string, error) {
	// Делаем GET-запрос к URL
	resp, err := http.Get(url)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch song page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("failed to fetch song page, status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse song page: %v", err)
	}

	lyrics := doc.Find(".lyrics").Text()

	if lyrics == "" {
		lyrics = doc.Find("div[class^='Lyrics__Container']").Text()
	}

	if lyrics == "" {
		return "", "", fmt.Errorf("div[class^='Lyrics__Container']")
	}

	releaseDate := doc.Find(".releaseDate").Text() // Дату релиза не получилось пока достать.

	verses := strings.Split(lyrics, "\n\n") // Разделяем на куплеты

	return strings.Join(verses, "\n\n"), releaseDate, nil
}
