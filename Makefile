include .env

build:
	@cd ./cmd/api && go build -o ../../bin/api 

run: build
	@./bin/api

seed:
	@go run internal/scripts/seed.go

test:
	@go test -v ./...

migration_up:
	migrate -path internal/db/migrations/ -database "postgresql://${dbuser}@${host}:${dbport}/${dbname}?sslmode=disable" -verbose up

migration_down:
	migrate -path internal/db/migrations/ -database "postgresql://${dbuser}@${host}:${dbport}/${dbname}?sslmode=disable" -verbose down

migration_fix:
	migrate -path internal/db/migrations/ -database "postgresql://${dbuser}@${host}:${dbport}/${dbname}?sslmode=disable" force VERSION

migration_go_to:
	migrate -path internal/db/migrations/ -database "postgresql://${dbuser}@${host}:${dbport}/${dbname}?sslmode=disable" -verbose goto VERSION

migration_create:
	migrate create -ext sql -dir internal/db/migrations -seq init

