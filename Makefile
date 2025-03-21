run:
	@echo "Запуск сервиса:"
	@docker compose up -d

stop:
	@echo "Остановка сервиса:"
	@docker compose down

fmt:
	@go fmt ./...

vet: fmt
	@go vet ./...

unit-tests: vet
	@echo "Запуск unit-тестов для основной логики shortener-сервиса:"
	@go test -v ./internal/services/...

	@echo "Запуск unit-тестов для обработчиков:"
	@go test -v ./internal/handlers/...
	
	@echo "Запуск unit-тестов для in-memory режима:"
	@go test -v ./internal/storage/memory/...

up-postgres: vet
	@echo "Запуск тестового контейнера с PostrgeSQL:"
	@docker build -t psql_test:test internal/storage/database/. && docker run --rm -p 5432:5432 --name psql_test -d psql_test:test

integration-tests: up-postgres 
	@echo "Запуск integration-тестов для postgres режима:"
	@sleep 5
	@go test -v ./internal/storage/database/...

test-cover:
	@go test -cover ./...

clean:
	@if [ $$(docker ps -q -f name=psql_test) ]; then docker stop psql_test; fi
	@go clean -testcache