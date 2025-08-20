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

migrate-up:
	export POSTGRES_URL='postgres://postgres:postgres@localhost:5434/blogo?sslmode=disable' && migrate -database ${POSTGRES_URL} -path internal/db/migration up

migrate-down:
	export POSTGRES_URL='postgres://postgres:postgres@localhost:5434/blogo?sslmode=disable' && migrate -database ${POSTGRES_URL} -path internal/db/migration down

sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate