DB_URL=postgresql://root:secret@localhost:5432/scoreit?sslmode=disable

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

db_docs:
	dbdocs build docs/db.dbml

db_schema:
	dbml2sql --postgres -o docs/schema.sql docs/db.dbml

sqlc:
	sqlc generate

test:
	go test -json -coverpkg=./... -cover ./... -coverprofile=coverage.out -short ./tools 2>&1 | gotestfmt

test-with-report:
	go test -json -coverpkg=./... -cover ./... -coverprofile=coverage.out -short ./tools 2>&1 | tee /tmp/gotest.log | gotestfmt

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/kwalter26/scoreit-api-go/db/sqlc Store

gin:
	GIN_MODE=release gin -i run main.go --all --port 8080

.PHONY: db_docs db_schema test test_report