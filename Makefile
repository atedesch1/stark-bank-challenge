build:
	go build -o bin/stark

run: build
	ENVIRONMENT=dev ./bin/stark
