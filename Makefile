run:
	go run cmd/main.go

test:
	go test ./...

swagger:
	swag init --generalInfo cmd/main.go --output docs

help:
	@echo "📦 Makefile команды:"
	@echo "  run      — запустить сервер"
	@echo "  swagger  — сгенерировать Swagger-документацию"
	@echo "  test     — запустить тесты"
	@echo "  build    — собрать бинарник"
