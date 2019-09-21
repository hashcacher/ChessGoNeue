generate-mocks:
	mockery -dir=./core -output=./core/mocks -name=Users

generate-mocks-two:
	mockgen -destination=./core/mocks/mock_users.go -package=mocks github.com/hashcacher/ChessGoNeue/Server/v2/core Users
