include .env 

export 

MIGRATION_PATH = ./migrations
DATABASE_URL = postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL)

# DATABASE_URL = postgres://postgres:market123@localhost:5432/chiapp?sslmode=disable

migrate-up:
	migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) up

migrate-up-one:
	migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) up $(version)

migrate-down-one:
	migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) down $(version)

migrate-goto:
	migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) goto $(version)

migrate-down:
	migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) down

migrate-force:
	migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) force $(version)

migrate-version:
	migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) version

create-migration:
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(name)

swagger:
	swag init -g cmd/server/main.go --parseDependency --parseInternal


