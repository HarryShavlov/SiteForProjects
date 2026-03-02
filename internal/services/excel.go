package services

import (
	"fmt"
	"os"
	"strings"

	"site/internal/models"

	"github.com/xuri/excelize/v2"
)

// LoadProjectsFromFile загружает проекты из Excel файла
func LoadProjectsFromFile(filepath string) ([]models.Project, error) {
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer f.Close()

	return parseProjectsSheet(f, "Sheet1")
}

// LoadUsersFromFile загружает пользователей из Excel файла
func LoadUsersFromFile(filepath string) ([]models.User, error) {
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer f.Close()

	return parseUsersSheet(f, "Sheet1")
}

func parseProjectsSheet(f *excelize.File, sheet string) ([]models.Project, error) {
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения листа: %w", err)
	}

	projects := make([]models.Project, 0)
	for i, row := range rows {
		if i == 0 || len(row) < 2 {
			continue
		}
		project := models.Project{
			ID:           i,
			Title:        getCell(row, 1),
			Subtitle:     getCell(row, 2),
			Level:        getCell(row, 3),
			Tags:         getCell(row, 4),
			Orders:       getCell(row, 5),
			Provider:     getCell(row, 6),
			ProviderCol:  getCell(row, 7),
			Duration:     getCell(row, 8),
			Participants: getCell(row, 9),
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func parseUsersSheet(f *excelize.File, sheet string) ([]models.User, error) {
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения листа: %w", err)
	}

	users := make([]models.User, 0)
	for i, row := range rows {
		if i == 0 || len(row) < 3 {
			continue
		}
		password := getCell(row, 2)
		// Хешируем пароль если он ещё не захеширован
		if !strings.HasPrefix(password, "$2a$") {
			password, _ = HashPassword(password)
		}
		user := models.User{
			ID:       i,
			Login:    getCell(row, 1),
			Password: password,
			Role:     getCell(row, 3),
			Name:     getCell(row, 4),
			Class:    getCell(row, 5),
		}
		users = append(users, user)
	}
	return users, nil
}

// ParseProjectsFromTempFile парсит временный Excel файл после загрузки
func ParseProjectsFromTempFile(tempPath string) ([]models.Project, error) {
	defer os.Remove(tempPath) // Удаляем после обработки
	return LoadProjectsFromFile(tempPath)
}

// getCell безопасно получает значение ячейки
func getCell(row []string, index int) string {
	if index < len(row) {
		return row[index]
	}
	return ""
}
