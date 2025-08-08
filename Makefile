NAME=dspg
MAIN=cmd/main.go
BIN=bin/${NAME}

deps:
	go get github.com/labstack/echo/v5@v5.0.0-20230722203903-ec5b858dab61
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
