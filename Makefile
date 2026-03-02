.PHONY: run build test clean deps

# Запуск в режиме разработки
run:
	go run cmd/server/main.go

# Сборка бинарника
build:
	go build -o bin/server cmd/server/main.go

# Запуск тестов
test:
	go test ./internal/... -v

# Очистка
clean:
	rm -f bin/server
	rm -f data/temp_*.xlsx

# Установка зависимостей
deps:
	go mod tidy

# Линтинг (если установлен golangci-lint)
lint:
	golangci-lint run ./...

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  make run    - Запустить сервер в режиме разработки"
	@echo "  make build  - Собрать бинарный файл"
	@echo "  make test   - Запустить тесты"
	@echo "  make clean  - Очистить артефакты сборки"
	@echo "  make deps   - Обновить зависимости"
	@echo "  make lint   - Запустить линтер"