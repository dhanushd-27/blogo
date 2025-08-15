server-run: 
	go run cmd/server/main.go

server-close:
	kill -9 $(pgrep -f "go run cmd/server/main.go")