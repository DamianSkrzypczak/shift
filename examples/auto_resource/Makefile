generate-schemas:
	@go generate ./api/schemas

cover:
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out
