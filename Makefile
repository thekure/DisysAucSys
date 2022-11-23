gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpc/interface.proto

server0:
	go run server/server.go 0
	
server1:
	go run server/server.go 1

server2:
	go run server/server.go 2
