DSN="postgres://postgres:password@localhost:5434/postgres?sslmode=disable"
SCHEMA_PATH=./schema
up:
	migrate -path $(SCHEMA_PATH) -database $(DSN) up
down:
	migrate -path $(SCHEMA_PATH) -database $(DSN) down