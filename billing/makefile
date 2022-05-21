db_conn = "postgres:password@localhost:5432/billing?sslmode=disable"

migrations.up:
	migrate -path db/migrations -database postgres://$(db_conn) -verbose up

migrations.create:
	migrate create -ext sql -dir db/migrations $(name)

swagger.gen:
	swag init --parseDependency --parseInternal --parseDepth 3 -g app/api/server.go
