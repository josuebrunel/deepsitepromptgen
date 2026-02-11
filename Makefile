NAME=dspg
MAIN=cmd/main.go
BIN=bin/${NAME}

deps:
	go install github.com/a-h/templ/cmd/templ@latest
	go get github.com/a-h/templ@latest
	go mod tidy

dev.deps: deps
	go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest

test:
	go test -v -failfast -count=1 -cover -covermode=count -coverprofile=coverage.out ./...
	go tool cover -func coverage.out

templ:
	templ generate

debug.build: build
	go build -gcflags "all=-N -l" -ldflags="-compressdwarf=false" -o ${BIN} ${MAIN}

debug: debug.build
	dlv --listen=:4000 --headless=true --log=true --accept-multiclient --api-version=2 exec ${BIN}

tidy:
	go mod tidy

build: tidy templ
	go build -o ${BIN} ${MAIN}

run: build
	./${BIN}

migration.init:
	go install github.com/pressly/goose/v3/cmd/goose@latest

migration.create:
	@read -p "Enter migration name: " name; \
	goose -dir db/migrations create $${name} sql

migration.up:
	goose -dir db/migrations postgres "host=$${DB_HOST} user=$${DB_USER} password=$${DB_PASSWORD} dbname=$${DB_NAME} port=$${DB_PORT} sslmode=disable" up

migration.down:
	goose -dir db/migrations postgres "host=$${DB_HOST} user=$${DB_USER} password=$${DB_PASSWORD} dbname=$${DB_NAME} port=$${DB_PORT} sslmode=disable" down

migration.status:
	goose -dir db/migrations postgres "host=$${DB_HOST} user=$${DB_USER} password=$${DB_PASSWORD} dbname=$${DB_NAME} port=$${DB_PORT} sslmode=disable" status
