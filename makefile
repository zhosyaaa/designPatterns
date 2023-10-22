postgres:
	docker run --rm -d --name postgres15 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=1079 -p 5432:5432 postgres:15

createdb:
	docker exec -it postgres15 createdb --username=postgres --owner=postgres bank

migrateup:
	migrate -path internal/db/migration -database "postgresql://postgres:1079@localhost:5432/bank?sslmode=disable" -verbose up

migratedown:
	migrate -path  internal/db/migration -database "postgresql://postgres:1079@localhost:5432/bank?sslmode=disable" -verbose down

dropdb:
	docker exec -it postgres15 dropdb -U postgres bank

.PHONY: postgres createdb dropdb