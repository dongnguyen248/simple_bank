postgres:
	docker run --name postgres12 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -p 5432:5432 -d postgres:12-alpine
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres12 dropdb --username=root simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose up
migrateup1:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
migratedown:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/dongnguyen248/simple_bank/db/sqlc Store
# The following targets are used to manage the database and run migrations.
.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock
# This Makefile is used to create and drop the database for the simple_bank application.