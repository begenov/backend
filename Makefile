postgres:
	sudo docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

createdb: 
	sudo docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb:
	sudo docker exec -it postgres dropdb  simple_bank
migrateup:
	migrate -path migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc: 
	sqlc generate

test: 
	go test -v -cover ./...

migrateup-test:
	migrate -path migration -database "postgresql://root:secret@localhost:5432/test?sslmode=disable" -verbose up

.PHONY: postgres createbd dropdb migrateup migratedown