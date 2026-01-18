include .env
export

# ----- DB MIGRATIONS -----

migrate-create:
	migrate create -ext sql -dir db/migrations $(name)
	
migrate-up:
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

.PHONY: migrate-create migrate-up migrate-down dev-http
