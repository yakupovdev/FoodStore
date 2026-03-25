include .env
export

export PROJECT_ROOT=$(shell pwd)

app-build:
	@docker compose up --build -d

app-up:
	@docker compose up -d
app-down:
	@docker compose down

go-up:
	@docker compose up -d foodstore-backend
go-down:
	@docker compose stop foodstore-backend

postgres-up:
	@docker compose up -d foodstore-postgres
postgres-down:
	@docker compose stop foodstore-postgres

pgadmin-up:
	@docker compose up -d foodstore-pgadmin
pgadmin-down:
	@docker compose stop foodstore-pgadmin
