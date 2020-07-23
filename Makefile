BIN_NAME=web-backend

dep:
	@echo "Fetching dependencies..."
	@go mod tidy

test: dep
	@echo "Testing source code..."
	@go test -race -short -cover -coverprofile cover.out
	@go tool cover -func cover.out

build: gen-swagger dep
	@echo "Building binary..."
	@go build -o ./${BIN_NAME} main.go

gen-swagger:
	@echo "Updating API documentation..."
	@swag init
