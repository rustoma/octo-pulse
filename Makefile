include ./bin/.env

build:
	@cd ./cmd/api && go build -o ../../bin/api 

run: build
	@./bin/api

test:
	@go test -v ./...

migration_up:
	migrate -path internal/db/migrations/ -database "postgresql://${dbuser}@${host}:${dbport}/${dbname}?sslmode=disable" -verbose up

#migration_up:
	migrate -path internal/db/migrations/ -database "postgresql://${dbuser}:${dbpass}@${host}:${dbport}/${dbname}?sslmode=disable" -verbose up

migration_down:
	migrate -path internal/db/migrations/ -database "postgresql://${dbuser}@${host}:${dbport}/${dbname}?sslmode=disable" -verbose down

migration_fix:
	migrate -path internal/db/migrations/ -database "postgresql://${dbuser}@${host}:${dbport}/${dbname}?sslmode=disable" force VERSION

migration_go_to:
	migrate -path internal/db/migrations/ -database "postgresql://${dbuser}@${host}:${dbport}/${dbname}?sslmode=disable" -verbose goto VERSION

migration_create:
	migrate create -ext sql -dir internal/db/migrations -seq add_domain_to_the_author

task_monit:
	./asynqmon --port=9090 --redis-password=${REDIS_PASSWORD}


build_workers:
	@cd ./cmd/workers && go build -o ../../bin/workers 

build_seed:
	@cd ./internal/scripts && go build -o ../../bin/seed

build_insert_images_script:
	@cd ./internal/scripts/insertImages && go build -o ../../../bin/insert_images

run_workers: build_workers
	@./bin/workers

seed: build_seed
	@./bin/seed

insert_images: build_insert_images_script
	@./bin/insert_images ./assets/images/roofs
