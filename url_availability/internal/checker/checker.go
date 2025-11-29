package checker

import (
	"net/http"
	"strings"
	"time"
	"url_availability/internal/models"
)

func normalize(url string) string {
	if !strings.HasPrefix(url, "http") {
		return "http://" + url
	}
	return url
}

// Проверка доступности страницы согласно условиям

func CheckLink(url string) models.LinkStatus {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(normalize(url))
	if err != nil {
		return models.NotAvailable
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return models.Available
	}

	return models.NotAvailable
}
