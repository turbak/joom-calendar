all: start-env migrate

start-env:
	docker-compose up -d

migrate:
	goose -dir migrations postgres "user=postgres dbname=joom_calendar password=postgrespw sslmode=disable" up