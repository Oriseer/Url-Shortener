build: 
	@go build -o bin/fs cmd/main/main.go

run: build
	@./bin/fs

