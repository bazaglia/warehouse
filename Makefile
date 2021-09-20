# define standard colors
ifneq (,$(findstring xterm,${TERM}))
	RED := $(shell tput -Txterm setaf 1)
	RESET := $(shell tput -Txterm sgr0)
else
	RED := ""
	RESET := ""
endif

.PHONY: build start test
.SILENT: start

build:
	go build -ldflags="-s -w" -o bin/importer cmd/importer/main.go
	go build -ldflags="-s -w" -o bin/server cmd/server/main.go

start:
ifneq ("$(wildcard .env)","")
	docker-compose up -d
else
	echo "${RED}ERROR${RESET}: Missing .env file. You can copy env.sample to .env to get started"
	exit 1
endif

test:
	go test ./...