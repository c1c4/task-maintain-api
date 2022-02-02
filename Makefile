build.application:
	docker-compose up --build app task-maintain-db

start.application:
	docker-compose up app task-maintain-db

stop.application:
	docker-compose down

test.unit:
	go test ./tests/unit/... -v

test.integration:
	docker-compose up -d task-maintain-db-test
	go test ./tests/integration/... -v
