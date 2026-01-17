run:
	go run ./main.go

mockgen:
	mockgen -source=internal/repository/user.go -destination=internal/repository/mocks/user_mock.go -package=mocks	
	