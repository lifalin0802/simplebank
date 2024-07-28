postgres:
	docker run  --name postgres16 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p:5432:5432 -d postgres:16

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres16 dropdb simple_bank

migrationup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrationdown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
		
rmpostgres:
	docker rm postgres16

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: createdb postgres dropdb rmpostgres migrationup migrationdown sqlc 