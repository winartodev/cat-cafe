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
