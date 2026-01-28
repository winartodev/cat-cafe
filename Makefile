include .env
export

# ----- DB MIGRATIONS -----

migrate-create:
	migrate create -ext sql -dir db/migrations $(name)
	
migrate-up:
	@echo "Running migrations with DB_URL: $(DB_URL)"
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down 1

migrate-force:
	migrate -path db/migrations -database "$(DB_URL)" force $(v)

seed-create:
	migrate create -ext sql -dir db/seeds $(name)

seed-up:
	@echo "Running seeds using separate table..."
	migrate -path db/seeds -database "$(DB_URL)&x-migrations-table=seed_migrations" -verbose up

seed-down:
	migrate -path db/seeds -database "$(DB_URL)&x-migrations-table=seed_migrations" -verbose down 1

seed-force:
	migrate -path db/seeds -database "$(DB_URL)&x-migrations-table=seed_migrations" force $(v)

#--------------------------

# ----- HOT RELOAD -----
dev-http:
	air -c .air.http.toml
	
# ----------------------

# ----- CONTAINER -----
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f app

docker-restart:
	docker-compose down
	docker-compose up -d --build

docker-clean:
	docker-compose down -v

docker-build:
	docker build -t cat-cafe-api:latest .

# ----------------------

.PHONY: migrate-create migrate-up migrate-down dev-http
