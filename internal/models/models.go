package models

// User представляет пользователя системы
type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login" excel:"Логин"`
	Password string `json:"-" excel:"Пароль"`  // Не отдаём в JSON
	Role     string `json:"role" excel:"Роль"` // "admin" или "student"
	Name     string `json:"name" excel:"Имя"`
	Class    string `json:"class" excel:"Класс"`
}

// Project представляет тему проекта
type Project struct {
	ID           int    `json:"id"`
	Title        string `json:"title" excel:"Название"`
	Subtitle     string `json:"subtitle" excel:"Описание"`
	Level        string `json:"level" excel:"Уровень"`
	Tags         string `json:"tags" excel:"Теги"`
	Orders       string `json:"orders" excel:"Заказов"`
	Provider     string `json:"provider" excel:"Провайдер"`
	ProviderCol  string `json:"providerColor" excel:"Цвет"`
	Duration     string `json:"duration" excel:"Срок"`
	Participants string `json:"participants" excel:"Участников"`
}

// LoginRequest — запрос на авторизацию
type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse — ответ после успешного входа
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// UploadResponse — ответ после загрузки файла
type UploadResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Count   int    `json:"count"`
}
