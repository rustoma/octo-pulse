include .env

build:
	@cd ./cmd/api && go build -o ../../bin/api 

run: build
	@./bin/api

seed:
	@go run internal/cripts/seed.go

test:
	@go test -v ./...

migration_up:
	migrate -path database/migration/ -database "postgresql://${dbuser}@${host}:${dbport}/${dbname}?sslmode=disable" -verbose up

migration_down:
	migrate -path database/migration/ -database "postgresql://${dbuser}@${host}:${dbport}/${dbname}?sslmode=disable" -verbose down

migration_fix:
	migrate -path database/migration/ -database "postgresql://${dbuser}@${host}:${dbport}/${dbname}?sslmode=disable" force VERSION

migration_go_to:
	migrate -path database/migration/ -database "postgresql://${dbuser}@${host}:${dbport}/${dbname}?sslmode=disable" -verbose goto VERSION
