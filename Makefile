.PHONY: start-server
start-server:
	go run cmd/app/main.go || true

.PHONY: start-db
start-db:
	docker-compose -f pkg/db/db.docker-compose.yml up

.PHONE: reset-db
reset-db:
	docker container stop openai-go-db
	docker container rm openai-go-db
