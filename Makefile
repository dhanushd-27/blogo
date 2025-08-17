server-up:
	go run cmd/server/main.go

server-down:
	kill -9 $(pgrep -f "go run cmd/server/main.go")

docker-up:
	docker compose -f docker/docker-compose.yaml up -d

docker-down:
	docker compose -f docker/docker-compose.yaml down

docker-restart:
	docker compose -f docker/docker-compose.yaml down
	docker compose -f docker/docker-compose.yaml up -d

docker-logs:
	docker compose -f docker/docker-compose.yaml logs -f