POSTGRESQL_URL='postgres://postgres:postgres@localhost:5432/notifanga?sslmode=disable'

dev:
	go run cmd/apiserver/main.go .dev.env

build:
	go build -o ./bin/apiserver -v ./cmd/apiserver/

swagger:	
	swag init -g internal/app/delivery/http/server.go

migrate_create:
	migrate create -ext sql -dir ./migrations -seq ${TABLE_NAME}

migrate_up:
	migrate -database ${POSTGRESQL_URL} -path ./migrations/ up 

migrate_down:
	migrate -database ${POSTGRESQL_URL} -path ./migrations/ down

migrate_force_fix:
	migrate -path ./migrations/ -database ${POSTGRESQL_URL} force ${VERSION}

mock_repository:
	mockgen -destination=internal/app/services/mocks/${NAME}.go \
	-package=mocks 	-source=internal/app/services/${NAME}.go \
	github.com/sea-auca/auca-issue-collector/internal/app/services/${NAME}.go ${NAME}Repository

setup:
	go mod tidy
	go install github.com/swaggo/swag/cmd/swag@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

clean:
	rm -r ./bin/