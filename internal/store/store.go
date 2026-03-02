package store

import (
	"site/internal/models"
	"sync"
)

// Store предоставляет доступ к данным приложения
type Store struct {
	mu       sync.RWMutex
	projects []models.Project
	users    []models.User
}

// New создаёт новое хранилище
func New() *Store {
	return &Store{
		projects: make([]models.Project, 0),
		users:    make([]models.User, 0),
	}
}

// === Методы для проектов ===

// GetProjects возвращает все проекты
func (s *Store) GetProjects() []models.Project {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// Возвращаем копию чтобы избежать гонки данных
	result := make([]models.Project, len(s.projects))
	copy(result, s.projects)
	return result
}

// SetProjects устанавливает новые проекты
func (s *Store) SetProjects(projects []models.Project) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.projects = projects
}

// === Методы для пользователей ===

// GetUsers возвращает всех пользователей
func (s *Store) GetUsers() []models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]models.User, len(s.users))
	copy(result, s.users)
	return result
}

// SetUsers устанавливает новых пользователей
func (s *Store) SetUsers(users []models.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users = users
}

// GetUserByLogin находит пользователя по логину
func (s *Store) GetUserByLogin(login string) (*models.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.users {
		if u.Login == login {
			return &u, true
		}
	}
	return nil, false
}

// GetUserByID находит пользователя по ID
func (s *Store) GetUserByID(id int) (*models.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.users {
		if u.ID == id {
			return &u, true
		}
	}
	return nil, false
}
