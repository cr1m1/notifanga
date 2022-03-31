POSTGRESQL_URL='postgres://postgres:postgres@localhost:5432/notifanga?sslmode=disable'

dev:
	go run cmd/apiserver/main.go .dev.env

build:

migrate_create:
	migrate create -ext sql -dir ./migrations -seq ${TABLE_NAME}

migrate_up:
	migrate -database ${POSTGRESQL_URL} -path ./migrations/ up 

migrate_down:
	migrate -database ${POSTGRESQL_URL} -path ./migrations/ down

migrate_force_fix:
	migrate -path ./migrations/ -database ${POSTGRESQL_URL} force ${VERSION}

setup:
	go mod tidy
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

clean:
	rm -r ./bin/