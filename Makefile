run:
	go run ./main.go

mockgen:
	mockgen -source=internal/repository/user.go -destination=internal/repository/mocks/user_mock.go -package=mocks

test:
	go test -v ./... -coverprofile=coverage.out

coverage:
	go tool cover -html=coverage.out