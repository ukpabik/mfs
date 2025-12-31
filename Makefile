test:
	go test -race ./... 

run:
	go run main.go

client:
	go run ./internal/client/main.go