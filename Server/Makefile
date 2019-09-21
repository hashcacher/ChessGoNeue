generate-users-mock:
	mockgen -source=./core/user.go -destination=./core/mocks/user.go -package=mocks github.com/hashcacher/ChessGoNeue/Server/v2/core Users
generate-game-mock:
	mockgen -source=./core/game.go -destination=./core/mocks/game.go -package=mocks github.com/hashcacher/ChessGoNeue/Server/v2/core Games
generate-matchrequest-mock:
	mockgen -source=./core/matchRequest.go -destination=./core/mocks/matchRequest.go -package=mocks github.com/hashcacher/ChessGoNeue/Server/v2/core MatchRequests

generate-mocks: generate-users-mock generate-game-mock generate-matchrequest-mock
