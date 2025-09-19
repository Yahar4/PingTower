package httprequests

import (
	"fmt"
	"net/http"

	"github.com/PingTower/internal/entities"
)

// Функция для проверки конкретных url'ов из джобов
// привязывается к полю Handler
func HttpCheck(job entities.Job) error {
	resp, err := http.Get(job.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("status %d", resp.StatusCode)
	}

	return nil
}
