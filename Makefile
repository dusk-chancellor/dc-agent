

local-run:
	go run main.go

local-build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/agent ./main.go

build:
	docker build --tag agent .

run:
	docker run -d agent 
